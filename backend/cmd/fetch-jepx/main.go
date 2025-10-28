// Package main provides a pipeline job to fetch JEPX spot price data and normalize to JSON.
// Usage: go run main.go -date 2025-10-23 -area tokyo
// Output: /public/data/jp/jepx/spot-{area}-YYYY-MM-DD.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/teo/aversome/backend/internal/adapters"
	"github.com/teo/aversome/backend/pkg/timeutil"
)

func main() {
	var date, area string
	flag.StringVar(&date, "date", "", "Date in YYYY-MM-DD format (defaults to today)")
	flag.StringVar(&area, "area", "tokyo", "Area: tokyo or kansai")
	flag.Parse()

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

	log.Printf("Fetching JEPX spot prices for %s area on %s...", area, date)

	// Open test CSV (in production, this would be fetched via HTTP)
	csvPath := filepath.Join("internal", "adapters", "testdata", "jepx-sample.csv")
	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Failed to open JEPX CSV: %v", err)
	}
	defer f.Close()

	// Parse CSV using JEPX adapter
	adapter := adapters.NewJEPXAdapter()
	resp, err := adapter.ParseCSV(f, date, area)
	if err != nil {
		log.Fatalf("Failed to parse JEPX CSV: %v", err)
	}

	log.Printf("Parsed %d hourly price points", len(resp.PriceYenPerKwh))
	if resp.Meta != nil && resp.Meta.Warning != "" {
		log.Printf("Warning: %s", resp.Meta.Warning)
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

	log.Printf("✓ Successfully wrote %s (%d bytes)", outputPath, len(jsonData))
}
