// Package adapters provides data source adapters for normalization.
package adapters

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/teo/aversome/backend/internal/jepx"
)

// JEPXAdapter normalizes JEPX (Japan Electric Power Exchange) day-ahead spot price data.
type JEPXAdapter struct {
	sourceURL string
}

// NewJEPXAdapter creates a new JEPX data adapter.
func NewJEPXAdapter() *JEPXAdapter {
	return &JEPXAdapter{
		sourceURL: "https://www.jepx.jp/",
	}
}

// ParseCSV parses JEPX CSV data into normalized jepx.Response.
// CSV format (expected):
//   Date,Hour,Tokyo_Price,Kansai_Price
//   2025-10-23,0,24.32,23.15
//   2025-10-23,1,22.50,21.80
//   ...
//   2025-10-23,23,26.10,25.30
//
// Or Japanese headers:
//   日付,時,東京価格,関西価格
//
// Notes:
// - Header detection by column names (order may vary)
// - Prices are in JPY/kWh
// - Hour is 0-23 (24-hour format)
// - Multiple areas per date
func (a *JEPXAdapter) ParseCSV(reader io.Reader, date, area string) (*jepx.Response, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Auto-detect column indices
	colIndices := a.detectColumns(header, area)
	if colIndices["date"] == -1 || colIndices["hour"] == -1 || colIndices["price"] == -1 {
		return nil, fmt.Errorf("required columns not found in header: %v", header)
	}

	resp := jepx.NewResponse(date, area)
	resp.Source = jepx.Source{
		Name: "JEPX",
		URL:  a.sourceURL,
	}

	lineNum := 1
	hoursSeen := make(map[int]bool)

	// Read data rows
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV line %d: %w", lineNum, err)
		}
		lineNum++

		// Extract date (format: YYYY-MM-DD or similar)
		rowDate := strings.TrimSpace(record[colIndices["date"]])
		rowDate = a.normalizeDate(rowDate)

		// Only include rows matching the requested date
		if rowDate != date {
			continue
		}

		// Parse hour (0-23)
		hourStr := strings.TrimSpace(record[colIndices["hour"]])
		hour, err := strconv.Atoi(hourStr)
		if err != nil {
			return nil, fmt.Errorf("invalid hour at line %d: %s", lineNum, hourStr)
		}
		if hour < 0 || hour > 23 {
			return nil, fmt.Errorf("hour out of range (0-23) at line %d: %d", lineNum, hour)
		}

		// Parse price (JPY/kWh)
		priceStr := strings.TrimSpace(record[colIndices["price"]])
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price at line %d: %s", lineNum, priceStr)
		}

		// Build timestamp in Asia/Tokyo timezone
		timestamp := a.buildTimestamp(date, hour)

		pricePoint := jepx.PricePoint{
			Timestamp: timestamp,
			Price:     price,
		}
		resp.PriceYenPerKwh = append(resp.PriceYenPerKwh, pricePoint)
		hoursSeen[hour] = true
	}

	// Validate we have data
	if len(resp.PriceYenPerKwh) == 0 {
		return nil, fmt.Errorf("no data found for date %s and area %s", date, area)
	}

	// Add warning if missing hours
	if len(hoursSeen) < 24 {
		resp.Meta = &jepx.Meta{
			Warning: fmt.Sprintf("Data for %d hours available (expected 24)", len(hoursSeen)),
		}
	}

	return resp, nil
}

// detectColumns finds column indices by header names.
// Returns map with keys: date, hour, price.
func (a *JEPXAdapter) detectColumns(header []string, area string) map[string]int {
	indices := map[string]int{
		"date":  -1,
		"hour":  -1,
		"price": -1,
	}

	// Determine price column name based on area
	priceColumn := area + "_price"
	if area == "tokyo" {
		priceColumn = "tokyo_price"
	} else if area == "kansai" {
		priceColumn = "kansai_price"
	}

	for i, col := range header {
		col = strings.ToLower(strings.TrimSpace(col))
		switch {
		case strings.Contains(col, "date") || col == "日付":
			indices["date"] = i
		case strings.Contains(col, "hour") || col == "時" || col == "時刻":
			indices["hour"] = i
		case strings.Contains(col, priceColumn) || col == "東京価格" && area == "tokyo" || col == "関西価格" && area == "kansai":
			indices["price"] = i
		case strings.Contains(col, "price") || strings.Contains(col, "価格"):
			// Fallback: generic price column
			if indices["price"] == -1 {
				indices["price"] = i
			}
		}
	}

	return indices
}

// normalizeDate converts various date formats to YYYY-MM-DD.
// Handles: YYYY-MM-DD, YYYY/MM/DD, YYYYMMDD
func (a *JEPXAdapter) normalizeDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)

	// Try parsing common formats
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"20060102",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}

	// Return as-is if no format matches
	return dateStr
}

// buildTimestamp creates ISO8601 timestamp with Asia/Tokyo offset.
// Example: "2025-10-23T15:00:00+09:00"
func (a *JEPXAdapter) buildTimestamp(date string, hour int) string {
	// Parse date
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		// Fallback to current date if parsing fails
		t = time.Now()
	}

	// Load Asia/Tokyo timezone
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.FixedZone("JST", 9*60*60) // Fallback: UTC+9
	}

	// Create timestamp with specified hour
	timestamp := time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, loc)

	return timestamp.Format(time.RFC3339)
}
