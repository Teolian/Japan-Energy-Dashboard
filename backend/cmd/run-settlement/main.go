// Package main provides a CLI tool to run settlement calculations.
// Usage: go run main.go -profile profile.json -area tokyo -date 2025-10-23 -pv 0.15
// Output: settlement-result.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/teo/aversome/backend/internal/jepx"
	"github.com/teo/aversome/backend/internal/settlement"
)

func main() {
	var profilePath, area, date string
	var pvOffset float64

	flag.StringVar(&profilePath, "profile", "", "Path to consumption profile JSON file")
	flag.StringVar(&area, "area", "tokyo", "Area for JEPX prices: tokyo or kansai")
	flag.StringVar(&date, "date", "", "Date for JEPX prices (YYYY-MM-DD)")
	flag.Float64Var(&pvOffset, "pv", 0.0, "PV offset percentage (0.0-1.0, e.g., 0.15 for 15%)")
	flag.Parse()

	if profilePath == "" {
		log.Fatal("Error: -profile is required")
	}
	if date == "" {
		log.Fatal("Error: -date is required")
	}

	// Validate area
	if area != "tokyo" && area != "kansai" {
		log.Fatalf("Error: invalid area %q (must be 'tokyo' or 'kansai')", area)
	}

	// Validate PV offset
	if pvOffset < 0 || pvOffset > 1 {
		log.Fatalf("Error: pv offset must be between 0.0 and 1.0, got %v", pvOffset)
	}

	log.Printf("Running settlement calculation...")
	log.Printf("  Profile: %s", profilePath)
	log.Printf("  Area: %s", area)
	log.Printf("  Date: %s", date)
	log.Printf("  PV Offset: %.1f%%", pvOffset*100)

	// Load consumption profile
	profile, err := loadProfile(profilePath)
	if err != nil {
		log.Fatalf("Failed to load profile: %v", err)
	}
	log.Printf("Loaded %d hourly profile points", len(profile))

	// Load JEPX prices from generated JSON
	pricesPath := fmt.Sprintf("../public/data/jp/jepx/spot-%s-%s.json", area, date)
	pricesResp, err := loadJEPXPrices(pricesPath)
	if err != nil {
		log.Fatalf("Failed to load JEPX prices: %v\nHint: Run 'go run cmd/fetch-jepx/main.go --date %s --area %s' first", err, date, area)
	}
	log.Printf("Loaded %d hourly price points", len(pricesResp.PriceYenPerKwh))

	// Build settlement request
	req := &settlement.Request{
		Profile: profile,
		Prices: settlement.PricesRequest{
			Area: area,
			Date: date,
		},
		PVOffsetPct: pvOffset,
	}

	// Calculate settlement
	resp, err := settlement.Calculate(req, pricesResp.PriceYenPerKwh, pricesResp.Source)
	if err != nil {
		log.Fatalf("Settlement calculation failed: %v", err)
	}

	log.Printf("✓ Settlement calculated successfully")
	log.Printf("  Period: %s to %s", resp.Period.From, resp.Period.To)
	log.Printf("  Total kWh: %.1f", resp.Totals.KWh)
	log.Printf("  Total Cost: ¥%.1f", resp.Totals.CostYen)

	// Write result to JSON
	outputPath := "settlement-result.json"
	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	log.Printf("✓ Result written to %s (%d bytes)", outputPath, len(jsonData))
}

// loadProfile loads consumption profile from JSON file.
func loadProfile(path string) ([]settlement.ProfilePoint, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var profile []settlement.ProfilePoint
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return profile, nil
}

// loadJEPXPrices loads JEPX price data from generated JSON artifact.
func loadJEPXPrices(path string) (*jepx.Response, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var resp jepx.Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &resp, nil
}
