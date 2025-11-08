// Package main provides a pipeline job to fetch generation mix data from OCCTO and normalize to JSON.
// Usage: go run main.go -area tokyo -date 2025-11-08 --use-http
// Output: /public/data/jp/{area}/generation-YYYY-MM-DD.json
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
	var area, date, outputPath string
	var useHTTP bool
	flag.StringVar(&area, "area", "tokyo", "Area: tokyo or kansai")
	flag.StringVar(&date, "date", "", "Date in YYYY-MM-DD format (defaults to today)")
	flag.StringVar(&outputPath, "output", "", "Output file path (defaults to public/data/jp/{area}/generation-{date}.json)")
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

	// Validate area
	if area != "tokyo" && area != "kansai" {
		log.Fatalf("Invalid area: %s (must be tokyo or kansai)", area)
	}

	log.Printf("Fetching OCCTO generation mix data for %s/%s (HTTP: %v)...", area, date, useHTTP)

	// Get data source
	var reader io.ReadCloser
	var sourceName string

	if useHTTP {
		// HTTP fetch from OCCTO public API
		fetcher := pkghttp.NewFetcher(pkghttp.BrowserConfig())

		// OCCTO CSV download URL
		// jhSybt=03: 電源種別供給力 (generation capacity by fuel type)
		// Format: ?jhSybt=03&tgtYmdFrom=YYYY/MM/DD&tgtYmdTo=YYYY/MM/DD
		dateFormatted := parsedDate.Format("2006/01/02") // YYYY/MM/DD
		url := fmt.Sprintf(
			"https://web-kohyo.occto.or.jp/kks-web-public/download/downloadCsv?jhSybt=03&tgtYmdFrom=%s&tgtYmdTo=%s",
			dateFormatted, dateFormatted,
		)
		sourceName = "OCCTO"

		log.Printf("Attempting HTTP fetch from %s", url)

		data, err := fetcher.Fetch(url)
		if err != nil {
			log.Fatalf("HTTP fetch failed: %v", err)
		}

		log.Printf("✓ HTTP fetch successful")
		reader = data
		defer reader.Close()
	} else {
		log.Fatalf("Testdata mode not implemented for generation mix (use --use-http)")
	}

	// Parse CSV using OCCTO adapter
	adapter := adapters.NewOCCTOAdapter()
	resp, err := adapter.ParseGenerationMixCSV(reader, date, area)
	if err != nil {
		log.Fatalf("Failed to parse OCCTO CSV: %v", err)
	}

	resp.Source.Name = sourceName

	log.Printf("Parsed %d generation points", len(resp.Series))
	if resp.Meta != nil {
		log.Printf("Renewable penetration: %.1f%%", resp.Meta.AvgRenewablePct)
		log.Printf("Carbon intensity: %.1f gCO2/kWh", resp.Meta.AvgCarbonGCO2KWh)
		log.Printf("Peak solar: %.1f MW", resp.Meta.PeakSolarMW)
	}

	// Determine output path
	if outputPath == "" {
		outputPath = filepath.Join("public", "data", "jp", area, fmt.Sprintf("generation-%s.json", date))
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write JSON output
	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	log.Printf("✓ Successfully wrote generation mix data to %s", outputPath)
}
