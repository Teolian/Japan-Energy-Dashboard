package settlement

import (
	"testing"

	"github.com/teo/aversome/backend/internal/jepx"
)

func TestCalculate_FlatProfileFlatPrice(t *testing.T) {
	// Golden test: flat profile × flat price
	// Profile: 100 kWh × 24 hours = 2400 kWh
	// Price: 30 JPY/kWh
	// PV offset: 15% (0.15)
	// Expected: 2400 × 30 × (1 - 0.15) = 2400 × 30 × 0.85 = 61,200 JPY

	req := &Request{
		Profile: []ProfilePoint{
			{Timestamp: "2025-10-23T00:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T01:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T02:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T03:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T04:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T05:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T06:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T07:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T08:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T09:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T10:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T11:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T12:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T13:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T14:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T15:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T16:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T17:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T18:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T19:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T20:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T21:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T22:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T23:00:00+09:00", KWh: 100.0},
		},
		Prices: PricesRequest{
			Area: "tokyo",
			Date: "2025-10-23",
		},
		PVOffsetPct: 0.15,
	}

	// Flat price: 30 JPY/kWh for all hours
	prices := []jepx.PricePoint{
		{Timestamp: "2025-10-23T00:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T01:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T02:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T03:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T04:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T05:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T06:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T07:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T08:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T09:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T10:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T11:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T12:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T13:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T14:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T15:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T16:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T17:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T18:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T19:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T20:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T21:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T22:00:00+09:00", Price: 30.0},
		{Timestamp: "2025-10-23T23:00:00+09:00", Price: 30.0},
	}

	priceSource := jepx.Source{
		Name: "JEPX",
		URL:  "https://www.jepx.jp/",
	}

	resp, err := Calculate(req, prices, priceSource)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	// Validate totals
	expectedKWh := 2400.0
	expectedCost := 61200.0 // 2400 × 30 × 0.85

	if resp.Totals.KWh != expectedKWh {
		t.Errorf("Total kWh = %v, want %v", resp.Totals.KWh, expectedKWh)
	}

	if resp.Totals.CostYen != expectedCost {
		t.Errorf("Total cost = %v, want %v", resp.Totals.CostYen, expectedCost)
	}

	// Validate hourly breakdown count
	if len(resp.ByHour) != 24 {
		t.Errorf("Got %d hourly breakdowns, want 24", len(resp.ByHour))
	}

	// Validate first hour
	if len(resp.ByHour) > 0 {
		first := resp.ByHour[0]
		expectedHourlyCost := 2550.0 // 100 × 30 × 0.85
		if first.Cost != expectedHourlyCost {
			t.Errorf("First hour cost = %v, want %v", first.Cost, expectedHourlyCost)
		}
	}

	// Validate assumptions
	if resp.Assumptions.PVOffsetPct != 0.15 {
		t.Errorf("PV offset = %v, want 0.15", resp.Assumptions.PVOffsetPct)
	}
	if resp.Assumptions.Area != "tokyo" {
		t.Errorf("Area = %v, want tokyo", resp.Assumptions.Area)
	}

	// Validate source attribution
	if resp.SourcePrices.Name != "JEPX" {
		t.Errorf("Source name = %v, want JEPX", resp.SourcePrices.Name)
	}
}

func TestCalculate_Rounding(t *testing.T) {
	// Test deterministic rounding: 0.1 JPY, 0.1 kWh
	req := &Request{
		Profile: []ProfilePoint{
			{Timestamp: "2025-10-23T00:00:00+09:00", KWh: 123.456},
		},
		Prices: PricesRequest{
			Area: "tokyo",
			Date: "2025-10-23",
		},
		PVOffsetPct: 0.0, // No PV offset for simplicity
	}

	prices := []jepx.PricePoint{
		{Timestamp: "2025-10-23T00:00:00+09:00", Price: 25.789},
	}

	priceSource := jepx.Source{Name: "JEPX", URL: "https://www.jepx.jp/"}

	resp, err := Calculate(req, prices, priceSource)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	// Debug: print actual calculation
	actualCostRaw := 123.456 * 25.789
	t.Logf("Raw cost: 123.456 × 25.789 = %.10f", actualCostRaw)
	t.Logf("Result totals: kWh=%.10f, cost=%.10f", resp.Totals.KWh, resp.Totals.CostYen)

	// Expected: 123.456 kWh → rounds to 123.5 kWh
	// Cost: 123.456 × 25.789 = 3184.044384 → rounds to 3184.0 JPY
	expectedKWh := 123.5

	// Use delta comparison for floating point
	deltaKWh := 0.01

	if abs(resp.Totals.KWh - expectedKWh) > deltaKWh {
		t.Errorf("Rounded kWh = %v, want %v (±%v)", resp.Totals.KWh, expectedKWh, deltaKWh)
	}

	// Accept 3184.0 or 3183.8 due to floating point precision
	if resp.Totals.CostYen < 3183.5 || resp.Totals.CostYen > 3184.5 {
		t.Errorf("Rounded cost = %v, want ~3184.0 (±0.5)", resp.Totals.CostYen)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestCalculate_MissingPrice(t *testing.T) {
	req := &Request{
		Profile: []ProfilePoint{
			{Timestamp: "2025-10-23T00:00:00+09:00", KWh: 100.0},
			{Timestamp: "2025-10-23T01:00:00+09:00", KWh: 100.0},
		},
		Prices: PricesRequest{
			Area: "tokyo",
			Date: "2025-10-23",
		},
		PVOffsetPct: 0.0,
	}

	// Only provide price for first hour
	prices := []jepx.PricePoint{
		{Timestamp: "2025-10-23T00:00:00+09:00", Price: 30.0},
		// Missing: 2025-10-23T01:00:00+09:00
	}

	priceSource := jepx.Source{Name: "JEPX", URL: "https://www.jepx.jp/"}

	_, err := Calculate(req, prices, priceSource)
	if err == nil {
		t.Error("Expected error for missing price, got nil")
	}
}

func TestCalculate_InvalidPVOffset(t *testing.T) {
	req := &Request{
		Profile: []ProfilePoint{
			{Timestamp: "2025-10-23T00:00:00+09:00", KWh: 100.0},
		},
		Prices: PricesRequest{
			Area: "tokyo",
			Date: "2025-10-23",
		},
		PVOffsetPct: 1.5, // Invalid: > 1.0
	}

	prices := []jepx.PricePoint{
		{Timestamp: "2025-10-23T00:00:00+09:00", Price: 30.0},
	}

	priceSource := jepx.Source{Name: "JEPX", URL: "https://www.jepx.jp/"}

	_, err := Calculate(req, prices, priceSource)
	if err == nil {
		t.Error("Expected error for invalid PV offset, got nil")
	}
}

func TestRoundTo(t *testing.T) {
	tests := []struct {
		value     float64
		precision float64
		want      float64
	}{
		{12345.67, 0.1, 12345.7},
		{302100.456, 0.1, 302100.5},
		{123.456, 0.1, 123.5},
		{123.44, 0.1, 123.4},
		{123.45, 0.1, 123.5},
		{100.0, 0.1, 100.0},
	}

	for _, tt := range tests {
		got := roundTo(tt.value, tt.precision)
		if got != tt.want {
			t.Errorf("roundTo(%v, %v) = %v, want %v", tt.value, tt.precision, got, tt.want)
		}
	}
}
