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
	var date, area, outputPath string
	var useHTTP, jsonLog bool

	flag.StringVar(&date, "date", "", "Date in YYYY-MM-DD format (defaults to today)")
	flag.StringVar(&area, "area", "tokyo", "Area: tokyo or kansai")
	flag.StringVar(&outputPath, "output", "", "Output file path (defaults to public/data/jp/{area}/demand-{date}.json)")
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
		var err error

		fetcher := pkghttp.NewFetcher(pkghttp.DefaultConfig())

		if area == "tokyo" {
			// TEPCO provides data in monthly ZIP archives
			// Format: YYYYMM_power_usage.zip containing YYYYMMDD_power_usage.csv files
			parsedDate, _ := timeutil.ParseDate(date)
			yearMonth := parsedDate.Format("200601") // YYYYMM
			dayFile := parsedDate.Format("20060102") + "_power_usage.csv" // YYYYMMDD_power_usage.csv

			url = fmt.Sprintf("https://www.tepco.co.jp/forecast/html/images/%s_power_usage.zip", yearMonth)
			sourceName = cfg.TEPCO.Name

			lgr.Info(fmt.Sprintf("Attempting to fetch TEPCO ZIP from %s (looking for %s)", url, dayFile))

			reader, err = fetcher.FetchFromZip(url, dayFile)
		} else {
			// Kansai: Use OCCTO API (jhSybt=02 provides 30-minute interval demand data)
			// OCCTO provides demand data for all 10 regions including Kansai
			parsedDate, _ := timeutil.ParseDate(date)
			dateFormatted := parsedDate.Format("2006/01/02") // YYYY/MM/DD
			url = fmt.Sprintf(
				"https://web-kohyo.occto.or.jp/kks-web-public/download/downloadCsv?jhSybt=02&tgtYmdFrom=%s&tgtYmdTo=%s",
				dateFormatted, dateFormatted,
			)
			sourceName = "OCCTO"

			lgr.Info(fmt.Sprintf("Attempting to fetch Kansai demand from OCCTO: %s", url))

			reader, err = fetcher.Fetch(url)
		}

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
		if useHTTP && sourceName == "OCCTO" {
			// Use OCCTO adapter for Kansai when fetching from OCCTO
			adapter := adapters.NewOCCTOAdapter()
			resp, err = adapter.ParseDemandCSV(reader, date, demand.AreaKansai)
		} else {
			// Use Kansai adapter for testdata
			adapter := adapters.NewKansaiAdapter()
			resp, err = adapter.ParseCSV(reader, date)
		}
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

	// Determine output path
	if outputPath == "" {
		// Default path
		outputDir := filepath.Join("public", "data", "jp", area)
		outputPath = filepath.Join(outputDir, fmt.Sprintf("demand-%s.json", date))
	}

	// Create output directory
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write JSON file
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
