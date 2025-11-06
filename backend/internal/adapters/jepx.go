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
// CSV format (japanesepower.org):
//   datetime,Date,PeriodID,System Price Yen/kWh,Tokyo Yen/kWh,Kansai Yen/kWh,...
//   2025-11-03 00:00:00,2025-11-03,1,9.34,11.23,7.32,...
//   2025-11-03 00:30:00,2025-11-03,2,9.00,10.72,7.00,...
//   ...
//
// Notes:
// - 30-minute intervals (48 periods/day) - we extract only hourly (:00:00)
// - Prices are in JPY/kWh
// - Multiple areas in columns (Tokyo, Kansai, etc.)
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
	if colIndices["datetime"] == -1 && colIndices["date"] == -1 {
		return nil, fmt.Errorf("datetime or date column not found in header: %v", header)
	}
	if colIndices["price"] == -1 {
		return nil, fmt.Errorf("price column for area %s not found in header: %v", area, header)
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

		// Extract datetime or date+hour
		var rowDate string
		var hour int

		if colIndices["datetime"] != -1 {
			// Parse datetime column (format: "2025-11-03 00:00:00")
			datetimeStr := strings.TrimSpace(record[colIndices["datetime"]])

			// Extract date and hour from datetime
			parts := strings.Split(datetimeStr, " ")
			if len(parts) >= 2 {
				rowDate = a.normalizeDate(parts[0])

				// Extract hour from time (HH:MM:SS)
				timeParts := strings.Split(parts[1], ":")
				if len(timeParts) >= 1 {
					hour, err = strconv.Atoi(timeParts[0])
					if err != nil {
						continue // Skip invalid rows
					}

					// Only include hourly data (ignore :30:00 intervals)
					if len(timeParts) >= 2 {
						minute := timeParts[1]
						if minute != "00" {
							continue // Skip 30-minute intervals
						}
					}
				}
			} else {
				continue // Skip malformed datetime
			}
		} else {
			// Fallback: separate date and hour columns
			rowDate = strings.TrimSpace(record[colIndices["date"]])
			rowDate = a.normalizeDate(rowDate)

			hourStr := strings.TrimSpace(record[colIndices["hour"]])
			hour, err = strconv.Atoi(hourStr)
			if err != nil {
				return nil, fmt.Errorf("invalid hour at line %d: %s", lineNum, hourStr)
			}
		}

		// Only include rows matching the requested date
		if rowDate != date {
			continue
		}

		// Validate hour range
		if hour < 0 || hour > 23 {
			continue // Skip out-of-range hours
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
// Returns map with keys: datetime, date, hour, price.
func (a *JEPXAdapter) detectColumns(header []string, area string) map[string]int {
	indices := map[string]int{
		"datetime": -1,
		"date":     -1,
		"hour":     -1,
		"price":    -1,
	}

	// Determine area-specific price column patterns
	// Format: "Tokyo Yen/kWh", "Kansai Yen/kWh"
	areaCapitalized := strings.Title(strings.ToLower(area))

	for i, col := range header {
		colLower := strings.ToLower(strings.TrimSpace(col))

		switch {
		case colLower == "datetime":
			indices["datetime"] = i
		case colLower == "date" || col == "日付":
			indices["date"] = i
		case colLower == "hour" || col == "時" || col == "時刻":
			indices["hour"] = i
		case strings.Contains(colLower, strings.ToLower(areaCapitalized)) && strings.Contains(colLower, "yen/kwh"):
			// Match "Tokyo Yen/kWh", "Kansai Yen/kWh"
			indices["price"] = i
		case col == "東京価格" && area == "tokyo":
			indices["price"] = i
		case col == "関西価格" && area == "kansai":
			indices["price"] = i
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
