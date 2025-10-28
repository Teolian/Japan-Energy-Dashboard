// Package settlement provides cost calculation logic.
package settlement

import (
	"fmt"
	"math"

	"github.com/teo/aversome/backend/internal/jepx"
)

// Calculate computes settlement cost from consumption profile and JEPX prices.
// Formula: cost = Σ(kWh × JPY/kWh × (1 - pv_offset_pct))
// Rounding: 0.1 JPY for costs, 0.1 kWh for consumption
func Calculate(req *Request, prices []jepx.PricePoint, priceSource jepx.Source) (*Response, error) {
	if len(req.Profile) == 0 {
		return nil, fmt.Errorf("profile is empty")
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("prices are empty")
	}
	if req.PVOffsetPct < 0 || req.PVOffsetPct > 1 {
		return nil, fmt.Errorf("pv_offset_pct must be between 0 and 1, got %v", req.PVOffsetPct)
	}

	// Build price lookup map by timestamp
	priceMap := make(map[string]float64)
	for _, pp := range prices {
		priceMap[pp.Timestamp] = pp.Price
	}

	resp := NewResponse()
	resp.Assumptions = Assumptions{
		PVOffsetPct: req.PVOffsetPct,
		Area:        req.Prices.Area,
	}
	resp.SourcePrices = Source{
		Name: priceSource.Name,
		URL:  priceSource.URL,
	}

	var totalKWh, totalCost float64
	var firstTS, lastTS string

	// Calculate per-hour costs
	for i, profilePoint := range req.Profile {
		ts := profilePoint.Timestamp
		kwh := profilePoint.KWh

		// Find matching price
		price, ok := priceMap[ts]
		if !ok {
			return nil, fmt.Errorf("no price found for timestamp %s", ts)
		}

		// Apply PV offset: effective consumption = kwh × (1 - pv_offset_pct)
		effectiveKWh := kwh * (1 - req.PVOffsetPct)

		// Calculate cost: effective_kwh × price (before rounding)
		cost := effectiveKWh * price

		// Accumulate totals with unrounded values to avoid rounding errors
		totalKWh += kwh
		totalCost += cost

		// Round for display in breakdown
		kwhRounded := roundTo(kwh, 0.1)
		costRounded := roundTo(cost, 0.1)

		// Add to breakdown
		breakdown := HourlyBreakdown{
			Timestamp: ts,
			KWh:       kwhRounded,
			Price:     price,
			Cost:      costRounded,
		}
		resp.ByHour = append(resp.ByHour, breakdown)

		// Track period
		if i == 0 {
			firstTS = ts
		}
		lastTS = ts
	}

	// Set totals with rounding
	resp.Totals = Totals{
		KWh:     roundTo(totalKWh, 0.1),
		CostYen: roundTo(totalCost, 0.1),
	}

	// Set period
	resp.Period = Period{
		From: firstTS,
		To:   lastTS,
	}

	return resp, nil
}

// roundTo rounds a float64 to the nearest multiple of precision.
// Examples:
//   roundTo(12345.67, 0.1) = 12345.7
//   roundTo(302100.456, 0.1) = 302100.5
func roundTo(value, precision float64) float64 {
	if precision == 0 {
		return value
	}
	return math.Round(value/precision) * precision
}
