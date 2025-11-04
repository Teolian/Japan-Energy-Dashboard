// Package main provides HTTP-based JEPX spot price data fetching with fallback to testdata.
// Usage: go run main.go -area tokyo -date 2025-10-24 --use-http
// Output: /public/data/jp/jepx/spot-{area}-YYYY-MM-DD.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/teo/aversome/backend/internal/adapters"
	pkghttp "github.com/teo/aversome/backend/pkg/http"
	"github.com/teo/aversome/backend/pkg/logger"
	"github.com/teo/aversome/backend/pkg/sources"
	"github.com/teo/aversome/backend/pkg/timeutil"
)

func main() {
	var date, area string
	var useHTTP, jsonLog bool

	flag.StringVar(&date, "date", "", "Date in YYYY-MM-DD format (defaults to today)")
	flag.StringVar(&area, "area", "tokyo", "Area: tokyo or kansai")
	flag.BoolVar(&useHTTP, "use-http", false, "Use real HTTP fetching (default: testdata)")
	flag.BoolVar(&jsonLog, "json-log", false, "Enable JSON structured logging")
	flag.Parse()

	// Initialize logger
	lgr := logger.New(jsonLog)

	// Default to today if no date provided
	if date == "" {
		date = timeutil.FormatDate(time.Now())
	}

	// Validate date format
	if _, err := timeutil.ParseDate(date); err != nil {
		log.Fatalf("Invalid date format: %v", err)
	}

	// Validate area
	if area != "tokyo" && area != "kansai" {
		log.Fatalf("Invalid area: %s (must be 'tokyo' or 'kansai')", area)
	}

	lgr.Info(fmt.Sprintf("Fetching JEPX spot prices for %s area on %s (HTTP: %v)", area, date, useHTTP))

	// Load source configuration
	cfg := sources.LoadConfig()

	// Get data source
	var reader io.ReadCloser
	var sourceName string
	var fetchDuration time.Duration

	if useHTTP {
		// Attempt HTTP fetch
		start := time.Now()

		fetcher := pkghttp.NewFetcher(pkghttp.DefaultConfig())

		// JEPX provides CSV data via their market data API
		// Note: The actual URL structure may need to be verified/updated
		// Format example: https://www.jepx.jp/market/excel/spot_YYYYMMDD.csv
		parsedDate, _ := timeutil.ParseDate(date)
		dateStr := parsedDate.Format("20060102") // YYYYMMDD

		// JEPX URL already has trailing slash, don't add another one
		baseURL := strings.TrimRight(cfg.JEPX.URL, "/")
		url := fmt.Sprintf("%s/market/excel/spot_%s.csv", baseURL, dateStr)
		sourceName = cfg.JEPX.Name

		lgr.Info(fmt.Sprintf("Attempting HTTP fetch from %s", url))

		data, err := fetcher.Fetch(url)
		fetchDuration = time.Since(start)

		if err != nil {
			lgr.LogFetch(sourceName, "failure", "", "HTTP fetch failed, falling back to testdata", fetchDuration, err)
			useHTTP = false // Fallback to testdata
		} else {
			lgr.LogFetch(sourceName, "success", "", "HTTP fetch successful", fetchDuration, nil)
			reader = data
			defer reader.Close()
		}
	}

	// Fallback to testdata if HTTP failed or not requested
	if !useHTTP {
		csvPath := filepath.Join("internal", "adapters", "testdata", "jepx-sample.csv")
		sourceName = "JEPX (testdata)"

		lgr.Info(fmt.Sprintf("Using testdata: %s", csvPath))

		f, err := os.Open(csvPath)
		if err != nil {
			log.Fatalf("Failed to open testdata CSV: %v", err)
		}
		defer f.Close()
		reader = f
	}

	// Parse CSV using JEPX adapter
	parseStart := time.Now()

	adapter := adapters.NewJEPXAdapter()
	resp, err := adapter.ParseCSV(reader, date, area)

	parseDuration := time.Since(parseStart)

	if err != nil {
		lgr.LogFetch(sourceName, "failure", "", "CSV parsing failed", parseDuration, err)
		log.Fatalf("Failed to parse JEPX CSV: %v", err)
	}

	lgr.Info(fmt.Sprintf("Parsed %d hourly price points", len(resp.PriceYenPerKwh)))
	if resp.Meta != nil && resp.Meta.Warning != "" {
		lgr.Info(fmt.Sprintf("Warning: %s", resp.Meta.Warning))
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create output directory
	outputDir := filepath.Join("..", "public", "data", "jp", "jepx")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write JSON file
	outputPath := filepath.Join(outputDir, fmt.Sprintf("spot-%s-%s.json", area, date))
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}

	mode := "testdata"
	if useHTTP {
		mode = "HTTP"
	}

	lgr.LogFetch(
		sourceName,
		"success",
		outputPath,
		fmt.Sprintf("Successfully wrote JEPX spot price data (%s mode)", mode),
		fetchDuration+parseDuration,
		nil,
	)

	log.Printf("✓ Successfully wrote %s (%d bytes)", outputPath, len(jsonData))
}
