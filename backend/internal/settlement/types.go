// Package settlement provides types and business logic for settlement cost calculations.
// Follows AGENT_TECH_SPEC.md §3.5 API contract.
package settlement

// ProfilePoint represents a single hourly consumption point.
type ProfilePoint struct {
	Timestamp string  `json:"ts"`  // ISO8601 with Asia/Tokyo offset
	KWh       float64 `json:"kwh"` // Consumption in kWh
}

// PricesRequest specifies which JEPX price data to use.
type PricesRequest struct {
	Area string `json:"area"` // e.g., "tokyo", "kansai"
	Date string `json:"date"` // YYYY-MM-DD format
}

// Request is the input for settlement calculation.
// POST /api/settlements/run
type Request struct {
	Profile     []ProfilePoint `json:"profile"`              // Hourly consumption profile
	Prices      PricesRequest  `json:"prices"`               // JEPX price reference
	PVOffsetPct float64        `json:"pv_offset_pct"`        // PV offset percentage (0.0-1.0)
}

// Period represents the time range of the settlement.
type Period struct {
	From string `json:"from"` // ISO8601 timestamp
	To   string `json:"to"`   // ISO8601 timestamp
}

// Totals contains aggregated settlement results.
type Totals struct {
	KWh     float64 `json:"kwh"`      // Total consumption in kWh
	CostYen float64 `json:"cost_yen"` // Total cost in JPY
}

// HourlyBreakdown represents per-hour settlement details.
type HourlyBreakdown struct {
	Timestamp string  `json:"ts"`    // ISO8601 with Asia/Tokyo offset
	KWh       float64 `json:"kwh"`   // Consumption in kWh
	Price     float64 `json:"price"` // Price in JPY/kWh
	Cost      float64 `json:"cost"`  // Cost in JPY (kwh × price × (1 - pv%))
}

// Assumptions contains the parameters used in the calculation.
type Assumptions struct {
	PVOffsetPct float64 `json:"pv_offset_pct"` // PV offset percentage
	Area        string  `json:"area"`          // Price area
}

// Source contains attribution for price data.
type Source struct {
	Name string `json:"name"` // e.g., "JEPX"
	URL  string `json:"url"`  // Original data source URL
}

// Response is the settlement calculation result.
type Response struct {
	Period       Period            `json:"period"`        // Time range
	Totals       Totals            `json:"totals"`        // Aggregated results
	ByHour       []HourlyBreakdown `json:"by_hour"`       // Per-hour breakdown
	Assumptions  Assumptions       `json:"assumptions"`   // Calculation parameters
	SourcePrices Source            `json:"source_prices"` // Price data attribution
}

// NewResponse creates a properly initialized Response with defaults.
func NewResponse() *Response {
	return &Response{
		ByHour: make([]HourlyBreakdown, 0, 24), // Pre-allocate for 24 hours
	}
}
