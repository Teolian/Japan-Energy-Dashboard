// Package main estimates generation mix from existing demand and JEPX price data.
// Usage: go run main.go -area tokyo -date 2025-11-03
// Output: public/data/jp/{area}/generation-{date}.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/teo/aversome/backend/internal/demand"
	"github.com/teo/aversome/backend/internal/generation"
	"github.com/teo/aversome/backend/internal/jepx"
	"github.com/teo/aversome/backend/pkg/timeutil"
)

func main() {
	var area, date, outputPath string
	flag.StringVar(&area, "area", "tokyo", "Area: tokyo or kansai")
	flag.StringVar(&date, "date", "", "Date in YYYY-MM-DD format (defaults to today)")
	flag.StringVar(&outputPath, "output", "", "Output file path (defaults to public/data/jp/{area}/generation-{date}.json)")
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

	log.Printf("Estimating generation mix for %s/%s...", area, date)

	// Load demand data
	demandPath := filepath.Join("public", "data", "jp", area, fmt.Sprintf("demand-%s.json", date))
	demandData, err := os.ReadFile(demandPath)
	if err != nil {
		log.Fatalf("Failed to read demand data: %v (file: %s)", err, demandPath)
	}

	var demandResp demand.Response
	if err := json.Unmarshal(demandData, &demandResp); err != nil {
		log.Fatalf("Failed to parse demand JSON: %v", err)
	}

	log.Printf("✓ Loaded demand data: %d points", len(demandResp.Series))

	// Load JEPX price data
	jepxPath := filepath.Join("public", "data", "jp", "jepx", fmt.Sprintf("spot-%s-%s.json", area, date))
	jepxData, err := os.ReadFile(jepxPath)
	if err != nil {
		log.Fatalf("Failed to read JEPX data: %v (file: %s)", err, jepxPath)
	}

	var jepxResp jepx.Response
	if err := json.Unmarshal(jepxData, &jepxResp); err != nil {
		log.Fatalf("Failed to parse JEPX JSON: %v", err)
	}

	log.Printf("✓ Loaded JEPX data: %d price points", len(jepxResp.PriceYenPerKwh))

	// Estimate generation mix
	estimator := generation.NewEstimator()
	genResp, err := estimator.EstimateFromDemandAndPrice(&demandResp, &jepxResp)
	if err != nil {
		log.Fatalf("Failed to estimate generation mix: %v", err)
	}

	// Apply seasonal adjustment
	genResp = estimator.EstimateWithSeasonalAdjustment(genResp, parsedDate)

	log.Printf("✓ Estimated generation mix: %d points", len(genResp.Series))
	if genResp.Meta != nil {
		log.Printf("  Renewable penetration: %.1f%%", genResp.Meta.AvgRenewablePct)
		log.Printf("  Carbon intensity: %.1f gCO2/kWh", genResp.Meta.AvgCarbonGCO2KWh)
		log.Printf("  Peak solar: %.1f MW", genResp.Meta.PeakSolarMW)
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
	jsonData, err := json.MarshalIndent(genResp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	log.Printf("✓ Successfully wrote estimated generation mix to %s", outputPath)
}
