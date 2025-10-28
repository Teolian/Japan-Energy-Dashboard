// Package main provides HTTP-based demand data fetching with fallback to testdata.
// Usage: go run main.go -area tokyo -date 2025-10-24 --use-http
// Output: /public/data/jp/{area}/demand-YYYY-MM-DD.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/teo/aversome/backend/internal/adapters"
	"github.com/teo/aversome/backend/internal/demand"
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

	lgr.Info(fmt.Sprintf("Fetching %s demand data for %s (HTTP: %v)", area, date, useHTTP))

	// Load source configuration
	cfg := sources.LoadConfig()

	// Get data source
	var reader io.ReadCloser
	var sourceName string
	var fetchDuration time.Duration

	if useHTTP {
		// Attempt HTTP fetch
		start := time.Now()
		var url string
		if area == "tokyo" {
			url = cfg.TEPCO.URL
			sourceName = cfg.TEPCO.Name
		} else {
			url = cfg.Kansai.URL
			sourceName = cfg.Kansai.Name
		}

		lgr.Info(fmt.Sprintf("Attempting HTTP fetch from %s", url))

		fetcher := pkghttp.NewFetcher(pkghttp.DefaultConfig())
		var err error
		reader, err = fetcher.Fetch(url)
		fetchDuration = time.Since(start)

		if err != nil {
			lgr.LogFetch(sourceName, "failure", "", "HTTP fetch failed, falling back to testdata", fetchDuration, err)
			useHTTP = false // Fallback to testdata
		} else {
			lgr.LogFetch(sourceName, "success", "", "HTTP fetch successful", fetchDuration, nil)
			defer reader.Close()
		}
	}

	// Fallback to testdata if HTTP failed or not requested
	if !useHTTP {
		csvPath := ""
		if area == "tokyo" {
			csvPath = filepath.Join("internal", "adapters", "testdata", "tepco-sample.csv")
			sourceName = "TEPCO (testdata)"
		} else {
			csvPath = filepath.Join("internal", "adapters", "testdata", "kansai-sample.csv")
			sourceName = "Kansai (testdata)"
		}

		lgr.Info(fmt.Sprintf("Using testdata: %s", csvPath))

		f, err := os.Open(csvPath)
		if err != nil {
			log.Fatalf("Failed to open testdata CSV: %v", err)
		}
		defer f.Close()
		reader = f
	}

	// Parse CSV using appropriate adapter
	var resp *demand.Response
	var err error

	parseStart := time.Now()

	switch area {
	case "tokyo":
		adapter := adapters.NewTEPCOAdapter()
		resp, err = adapter.ParseCSV(reader, date)
	case "kansai":
		adapter := adapters.NewKansaiAdapter()
		resp, err = adapter.ParseCSV(reader, date)
	}

	parseDuration := time.Since(parseStart)

	if err != nil {
		lgr.LogFetch(sourceName, "failure", "", "CSV parsing failed", parseDuration, err)
		log.Fatalf("Failed to parse CSV: %v", err)
	}

	lgr.Info(fmt.Sprintf("Parsed %d data points", len(resp.Series)))

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create output directory
	outputDir := filepath.Join("..", "public", "data", "jp", area)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write JSON file
	outputPath := filepath.Join(outputDir, fmt.Sprintf("demand-%s.json", date))
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
		fmt.Sprintf("Successfully wrote demand data (%s mode)", mode),
		fetchDuration+parseDuration,
		nil,
	)

	log.Printf("✓ Successfully wrote %s (%d bytes)", outputPath, len(jsonData))
}
