// Package jepx provides types and business logic for JEPX spot price data.
// Follows AGENT_TECH_SPEC.md §3.3 API contract.
package jepx

// PricePoint represents a single hourly spot price.
type PricePoint struct {
	Timestamp string  `json:"ts"`    // ISO8601 with Asia/Tokyo offset (e.g., "2025-10-23T00:00:00+09:00")
	Price     float64 `json:"price"` // JPY/kWh
}

// Source contains attribution for the data source.
type Source struct {
	Name string `json:"name"` // e.g., "JEPX"
	URL  string `json:"url"`  // Original data source URL
}

// Meta contains optional metadata and warnings.
type Meta struct {
	Warning string `json:"warning,omitempty"` // Non-blocking warning message
}

// Response is the top-level response structure for JEPX spot price endpoint.
// GET /api/jp/jepx/spot?date=YYYY-MM-DD&area=tokyo
type Response struct {
	Date             string       `json:"date"`               // YYYY-MM-DD format
	Area             string       `json:"area"`               // e.g., "tokyo", "kansai"
	Timescale        string       `json:"timescale"`          // Always "hourly"
	PriceYenPerKwh   []PricePoint `json:"price_yen_per_kwh"`  // 24 hourly price points
	Source           Source       `json:"source"`             // Data attribution
	Meta             *Meta        `json:"meta,omitempty"`     // Optional metadata/warnings
}

// NewResponse creates a properly initialized Response with defaults.
func NewResponse(date, area string) *Response {
	return &Response{
		Date:           date,
		Area:           area,
		Timescale:      "hourly",
		PriceYenPerKwh: make([]PricePoint, 0, 24), // Pre-allocate for 24 hours
	}
}
