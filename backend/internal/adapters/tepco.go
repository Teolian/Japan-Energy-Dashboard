// Package adapters provides data source adapters for normalization.
// Adapters are pure functions: input artifact → validated typed structs.
// Per AGENT_TECH_SPEC.md §4: handle CSV variability, auto-detect headers.
package adapters

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/teo/aversome/backend/internal/demand"
	"github.com/teo/aversome/backend/pkg/timeutil"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// TEPCOAdapter normalizes Tokyo Electric Power Company demand data.
type TEPCOAdapter struct {
	sourceURL string
}

// NewTEPCOAdapter creates a new TEPCO data adapter.
func NewTEPCOAdapter() *TEPCOAdapter {
	return &TEPCOAdapter{
		sourceURL: "https://www.tepco.co.jp/forecast/html/download-j.html",
	}
}

// ParseCSV parses TEPCO CSV data into normalized demand.Response.
// CSV format (example):
//   DATE,TIME,実績(万kW),予測(万kW)
//   2025-10-24,0:00,2665.4,2701.0
//   2025-10-24,1:00,2589.2,2630.5
//
// Notes:
// - Header detection by column names (order may vary)
// - Units: 万kW (10,000 kW) = 10 MW → convert to MW
// - Forecast may be empty for some rows
func (a *TEPCOAdapter) ParseCSV(reader io.Reader, date string) (*demand.Response, error) {
	// Convert from Shift-JIS to UTF-8
	// TEPCO CSV files are encoded in Shift-JIS (Japanese encoding)
	utf8Reader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())

	csvReader := csv.NewReader(utf8Reader)
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields (for metadata rows)

	// Find the header row (skip metadata rows)
	// The actual data starts with a row beginning with "DATE"
	var header []string
	var err error
	for {
		header, err = csvReader.Read()
		if err != nil {
			return nil, fmt.Errorf("failed to find CSV header: %w", err)
		}
		// Check if this row starts with "DATE" (the actual header)
		if len(header) > 0 && strings.ToUpper(strings.TrimSpace(header[0])) == "DATE" {
			break
		}
	}

	// Auto-detect column indices
	colIndices := a.detectColumns(header)
	if colIndices["date"] == -1 || colIndices["time"] == -1 || colIndices["actual"] == -1 {
		return nil, fmt.Errorf("required columns not found in header: %v", header)
	}

	resp := demand.NewResponse(demand.AreaTokyo, date)
	resp.Source = demand.Source{
		Name: "TEPCO",
		URL:  a.sourceURL,
	}

	// Parse date for validation
	baseDate, err := timeutil.ParseDate(date)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}

	hasForecast := false
	lineNum := 1
	seenHours := make(map[int]bool) // Track hours to avoid duplicates

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

		// Extract date and time
		rowDate := strings.TrimSpace(record[colIndices["date"]])
		rowTime := strings.TrimSpace(record[colIndices["time"]])

		// Normalize date format (2025/11/1 → 2025-11-01)
		normalizedRowDate := a.normalizeDate(rowDate)
		normalizedDate := a.normalizeDate(date)

		// Only include rows matching the requested date
		if normalizedRowDate != normalizedDate {
			continue
		}

		// Parse hour from time string (e.g., "0:00", "13:00")
		// Filter only hourly data (skip 5-minute intervals like "0:05", "0:10")
		hour, minutes, err := a.parseTime(rowTime)
		if err != nil {
			return nil, fmt.Errorf("invalid time format at line %d: %s", lineNum, rowTime)
		}

		// Skip non-hourly data points (only include :00 minutes)
		if minutes != 0 {
			continue
		}

		// Skip duplicate hours (CSV contains multiple blocks with same hours)
		if seenHours[hour] {
			continue
		}
		seenHours[hour] = true

		// Create timestamp
		ts := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), hour, 0, 0, 0, timeutil.TokyoLocation)

		// Parse actual demand (万kW → MW)
		actualStr := strings.TrimSpace(record[colIndices["actual"]])
		actual, err := strconv.ParseFloat(actualStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid actual value at line %d: %s", lineNum, actualStr)
		}
		demandMW := actual * 10.0 // 万kW to MW

		// Parse forecast if column exists and value is present
		var forecastMW *float64
		if colIndices["forecast"] != -1 && colIndices["forecast"] < len(record) {
			forecastStr := strings.TrimSpace(record[colIndices["forecast"]])
			if forecastStr != "" {
				forecast, err := strconv.ParseFloat(forecastStr, 64)
				if err == nil {
					fVal := forecast * 10.0
					forecastMW = &fVal
					hasForecast = true
				}
			}
		}

		point := demand.SeriesPoint{
			Timestamp:  ts,
			DemandMW:   demandMW,
			ForecastMW: forecastMW,
		}
		resp.Series = append(resp.Series, point)
	}

	// Validate we have data
	if len(resp.Series) == 0 {
		return nil, fmt.Errorf("no data found for date %s", date)
	}

	// Add warning if no forecast data
	if !hasForecast {
		resp.Meta = &demand.Meta{
			Warning: "Forecast data not available for this date",
		}
	}

	return resp, nil
}

// detectColumns finds column indices by header names.
// Returns map with keys: date, time, actual, forecast.
func (a *TEPCOAdapter) detectColumns(header []string) map[string]int {
	indices := map[string]int{
		"date":     -1,
		"time":     -1,
		"actual":   -1,
		"forecast": -1,
	}

	for i, col := range header {
		col = strings.ToLower(strings.TrimSpace(col))
		switch {
		case strings.Contains(col, "date") || col == "日付":
			indices["date"] = i
		case strings.Contains(col, "time") || col == "時刻":
			indices["time"] = i
		case strings.Contains(col, "実績") || strings.Contains(col, "actual"):
			indices["actual"] = i
		case strings.Contains(col, "予測") || strings.Contains(col, "予想") || strings.Contains(col, "forecast"):
			indices["forecast"] = i
		}
	}

	return indices
}

// parseTime extracts hour and minutes from time string like "0:00", "13:00", "0:05".
func (a *TEPCOAdapter) parseTime(timeStr string) (hour int, minutes int, err error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format: %s", timeStr)
	}
	hour, err = strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("invalid hour: %s", timeStr)
	}
	minutes, err = strconv.Atoi(parts[1])
	if err != nil || minutes < 0 || minutes > 59 {
		return 0, 0, fmt.Errorf("invalid minutes: %s", timeStr)
	}
	return hour, minutes, nil
}

// normalizeDate converts date from various formats to YYYY-MM-DD.
// Supports: "2025/11/1", "2025-11-01", "2025/11/01"
func (a *TEPCOAdapter) normalizeDate(dateStr string) string {
	// Replace / with -
	dateStr = strings.ReplaceAll(dateStr, "/", "-")

	// Parse and reformat to ensure zero-padding
	// Try parsing as 2006-1-2 (without zero padding)
	t, err := time.Parse("2006-1-2", dateStr)
	if err == nil {
		return t.Format("2006-01-02")
	}

	// Try parsing as 2006-01-02 (with zero padding)
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t.Format("2006-01-02")
	}

	// Return as-is if parsing fails
	return dateStr
}
