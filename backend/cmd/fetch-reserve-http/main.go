// Package main provides a pipeline job to fetch reserve margin data from OCCTO and normalize to JSON.
// Usage: go run main.go -date 2025-10-24 --use-http
// Output: /public/data/jp/system/reserve-YYYY-MM-DD.json
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
	pkghttp "github.com/teo/aversome/backend/pkg/http"
	"github.com/teo/aversome/backend/pkg/timeutil"
)

func main() {
	var date, outputPath string
	var useHTTP bool
	flag.StringVar(&date, "date", "", "Date in YYYY-MM-DD format (defaults to today)")
	flag.StringVar(&outputPath, "output", "", "Output file path (defaults to public/data/jp/system/reserve-{date}.json)")
	flag.BoolVar(&useHTTP, "use-http", false, "Use real HTTP fetching (default: testdata)")
	flag.Parse()

	// Default to today if no date provided
	if date == "" {
		date = timeutil.FormatDate(time.Now())
	}

	// Validate date format
	parsedDate, err := timeutil.ParseDate(date)
	if err != nil {
		log.Fatalf("Invalid date format: %v", err)
	}

	log.Printf("Fetching OCCTO reserve margin data for %s (HTTP: %v)...", date, useHTTP)

	// Get data source
	var reader io.ReadCloser
	var sourceName string

	if useHTTP {
		// HTTP fetch from OCCTO public API
		fetcher := pkghttp.NewFetcher(pkghttp.BrowserConfig())

		// OCCTO CSV download URL
		// Format: ?jhSybt=02&tgtYmdFrom=YYYY/MM/DD&tgtYmdTo=YYYY/MM/DD
		dateFormatted := parsedDate.Format("2006/01/02") // YYYY/MM/DD
		url := fmt.Sprintf(
			"https://web-kohyo.occto.or.jp/kks-web-public/download/downloadCsv?jhSybt=02&tgtYmdFrom=%s&tgtYmdTo=%s",
			dateFormatted, dateFormatted,
		)
		sourceName = "OCCTO"

		log.Printf("Attempting HTTP fetch from %s", url)

		data, err := fetcher.Fetch(url)
		if err != nil {
			log.Printf("⚠️  HTTP fetch failed, falling back to testdata: %v", err)
			useHTTP = false // Fallback
		} else {
			log.Printf("✓ HTTP fetch successful")
			reader = data
			defer reader.Close()
		}
	}

	// Fallback to testdata if HTTP failed or not requested
	if !useHTTP {
		csvPath := filepath.Join("internal", "adapters", "testdata", "occto-sample.csv")
		sourceName = "OCCTO (testdata)"

		log.Printf("Using testdata: %s", csvPath)

		f, err := os.Open(csvPath)
		if err != nil {
			log.Fatalf("Failed to open testdata CSV: %v", err)
		}
		defer f.Close()
		reader = f
	}

	// Parse CSV using OCCTO adapter
	adapter := adapters.NewOCCTOAdapter()
	resp, err := adapter.ParseCSV(reader, date)
	if err != nil {
		log.Fatalf("Failed to parse OCCTO CSV: %v", err)
	}

	resp.Source.Name = sourceName

	log.Printf("Parsed %d areas", len(resp.Areas))
	if resp.Meta != nil && resp.Meta.Warning != "" {
		log.Printf("Warning: %s", resp.Meta.Warning)
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Determine output path
	if outputPath == "" {
		// Default path (system-level data)
		outputDir := filepath.Join("public", "data", "jp", "system")
		outputPath = filepath.Join(outputDir, fmt.Sprintf("reserve-%s.json", date))
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

	log.Printf("✓ Successfully wrote %s (%d bytes)", outputPath, len(jsonData))
}
