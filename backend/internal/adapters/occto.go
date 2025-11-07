// Package adapters provides data source adapters for normalization.
package adapters

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/teo/aversome/backend/internal/demand"
	"github.com/teo/aversome/backend/internal/reserve"
)

// OCCTOAdapter normalizes OCCTO (Organization for Cross-regional Coordination of
// Transmission Operators) reserve margin data.
type OCCTOAdapter struct {
	sourceURL string
}

// NewOCCTOAdapter creates a new OCCTO data adapter.
func NewOCCTOAdapter() *OCCTOAdapter {
	return &OCCTOAdapter{
		sourceURL: "https://www.occto.or.jp/",
	}
}

// ParseCSV parses OCCTO CSV data into normalized reserve.Response.
// CSV format (from web-kohyo.occto.or.jp):
//   "2025/11/03 22:59 UPDATE"  ← skip this line
//   "対象年月日","時刻","ブロックNo","エリア名","広域ブロック需要(MW)",...,"エリア需要(MW)","エリア供給力(MW)","エリア予備力(MW)"
//   "2025/11/03","00:30","1","北海道",41823.764,46992.204,...,2854,3243,389
//   "2025/11/03","00:30","1","東北",41823.764,46992.204,...,6589.48,7763.072,1173.592
//
// Notes:
// - First line is timestamp (skip it)
// - 30-minute intervals (we aggregate to daily averages)
// - エリア名 = area name, エリア需要(MW) = demand, エリア供給力(MW) = capacity
// - Reserve margin calculated: (capacity - demand) / capacity * 100
func (a *OCCTOAdapter) ParseCSV(reader io.Reader, date string) (*reserve.Response, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.LazyQuotes = true       // Handle quoted fields more flexibly
	csvReader.FieldsPerRecord = -1    // Allow variable number of fields

	// Skip first line (UPDATE timestamp)
	_, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read first line: %w", err)
	}

	// Read header (second line)
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Auto-detect column indices
	colIndices := a.detectColumns(header)
	if colIndices["date"] == -1 || colIndices["area"] == -1 {
		return nil, fmt.Errorf("required columns (date, area) not found in header: %v", header)
	}
	if colIndices["demand"] == -1 || colIndices["capacity"] == -1 {
		return nil, fmt.Errorf("required columns (demand, capacity) not found in header: %v", header)
	}

	resp := reserve.NewResponse(date)
	resp.Source = reserve.Source{
		Name: "OCCTO",
		URL:  a.sourceURL,
	}

	lineNum := 1
	// Aggregate data by area (multiple 30-min intervals per area)
	areaData := make(map[string]struct {
		demandSum   float64
		capacitySum float64
		count       int
	})

	// Normalize date format (2025/11/03 → 2025-11-03)
	normalizedDate := strings.ReplaceAll(date, "-", "/")

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

		// Extract date (format: 2025/11/03)
		rowDate := strings.TrimSpace(record[colIndices["date"]])

		// Only include rows matching the requested date
		if rowDate != normalizedDate && rowDate != date {
			continue
		}

		// Extract area
		areaStr := strings.TrimSpace(record[colIndices["area"]])
		area := a.normalizeArea(areaStr)
		if area == "" {
			continue // Skip unknown areas
		}

		// Parse demand (MW)
		demandStr := strings.TrimSpace(record[colIndices["demand"]])
		demand, err := strconv.ParseFloat(demandStr, 64)
		if err != nil {
			continue // Skip invalid rows
		}

		// Parse capacity (MW)
		capacityStr := strings.TrimSpace(record[colIndices["capacity"]])
		capacity, err := strconv.ParseFloat(capacityStr, 64)
		if err != nil {
			continue // Skip invalid rows
		}

		// Aggregate by area
		data := areaData[area]
		data.demandSum += demand
		data.capacitySum += capacity
		data.count++
		areaData[area] = data
	}

	// Calculate averages and reserve margins
	for area, data := range areaData {
		if data.count == 0 {
			continue
		}

		avgDemand := data.demandSum / float64(data.count)
		avgCapacity := data.capacitySum / float64(data.count)

		// Calculate reserve margin: (capacity - demand) / capacity * 100
		reservePct := 0.0
		if avgCapacity > 0 {
			reservePct = ((avgCapacity - avgDemand) / avgCapacity) * 100
		}

		// Derive status from percentage
		status := reserve.DeriveStatus(reservePct)

		areaReserve := reserve.AreaReserve{
			Area:             area,
			ReserveMarginPct: reservePct,
			Status:           status,
		}
		resp.Areas = append(resp.Areas, areaReserve)
	}

	// Validate we have data
	if len(resp.Areas) == 0 {
		return nil, fmt.Errorf("no data found for date %s", date)
	}

	return resp, nil
}

// ParseDemandCSV parses OCCTO CSV data into demand.Response for a specific area.
// Uses the same CSV format as ParseCSV but extracts hourly demand time-series.
// Aggregates 30-minute intervals into hourly averages.
func (a *OCCTOAdapter) ParseDemandCSV(reader io.Reader, date string, targetArea demand.Area) (*demand.Response, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1

	// Skip first line (UPDATE timestamp)
	_, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read first line: %w", err)
	}

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Auto-detect column indices (including time)
	colIndices := a.detectDemandColumns(header)
	if colIndices["date"] == -1 || colIndices["time"] == -1 || colIndices["area"] == -1 || colIndices["demand"] == -1 {
		return nil, fmt.Errorf("required columns not found in header: %v", header)
	}

	resp := demand.NewResponse(targetArea, date)
	resp.Source = demand.Source{
		Name: "OCCTO",
		URL:  a.sourceURL,
	}

	// Normalize date format (2025-11-03 → 2025/11/03)
	normalizedDate := strings.ReplaceAll(date, "-", "/")

	// Aggregate 30-minute intervals into hourly data
	hourlyData := make(map[int]struct {
		demandSum float64
		count     int
	})

	lineNum := 1
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV line %d: %w", lineNum, err)
		}
		lineNum++

		// Extract date
		rowDate := strings.TrimSpace(record[colIndices["date"]])
		if rowDate != normalizedDate && rowDate != date {
			continue
		}

		// Extract area
		areaStr := strings.TrimSpace(record[colIndices["area"]])
		area := a.normalizeArea(areaStr)

		// Only process the target area
		if string(targetArea) != area {
			continue
		}

		// Extract time (format: "00:30", "01:00", etc.)
		timeStr := strings.TrimSpace(record[colIndices["time"]])

		// Parse hour from time
		timeParts := strings.Split(timeStr, ":")
		if len(timeParts) < 1 {
			continue
		}
		hour, err := strconv.Atoi(timeParts[0])
		if err != nil {
			continue
		}

		// Parse demand (MW)
		demandStr := strings.TrimSpace(record[colIndices["demand"]])
		demandVal, err := strconv.ParseFloat(demandStr, 64)
		if err != nil {
			continue
		}

		// Aggregate into hourly buckets
		data := hourlyData[hour]
		data.demandSum += demandVal
		data.count++
		hourlyData[hour] = data
	}

	// Convert aggregated data to SeriesPoints
	location, _ := time.LoadLocation("Asia/Tokyo")
	baseDate, _ := time.Parse("2006-01-02", date)

	for hour := 0; hour < 24; hour++ {
		data, exists := hourlyData[hour]
		if !exists || data.count == 0 {
			continue
		}

		avgDemand := data.demandSum / float64(data.count)

		timestamp := time.Date(
			baseDate.Year(), baseDate.Month(), baseDate.Day(),
			hour, 0, 0, 0, location,
		)

		point := demand.SeriesPoint{
			Timestamp:  timestamp,
			DemandMW:   avgDemand,
			ForecastMW: nil, // OCCTO doesn't provide forecasts
		}
		resp.Series = append(resp.Series, point)
	}

	// Validate we have data
	if len(resp.Series) == 0 {
		return nil, fmt.Errorf("no demand data found for area %s on date %s", targetArea, date)
	}

	return resp, nil
}

// detectDemandColumns finds column indices for demand parsing (includes time).
func (a *OCCTOAdapter) detectDemandColumns(header []string) map[string]int {
	indices := map[string]int{
		"date":   -1,
		"time":   -1,
		"area":   -1,
		"demand": -1,
	}

	for i, col := range header {
		colTrimmed := strings.TrimSpace(col)

		switch {
		case strings.Contains(colTrimmed, "対象年月日"):
			indices["date"] = i
		case strings.Contains(colTrimmed, "時刻") || colTrimmed == "対象時刻":
			indices["time"] = i
		case colTrimmed == "エリア名":
			indices["area"] = i
		case colTrimmed == "エリア需要(MW)":
			indices["demand"] = i
		}
	}

	return indices
}

// detectColumns finds column indices by header names.
// Returns map with keys: date, area, demand, capacity.
func (a *OCCTOAdapter) detectColumns(header []string) map[string]int {
	indices := map[string]int{
		"date":     -1,
		"area":     -1,
		"demand":   -1,
		"capacity": -1,
	}

	for i, col := range header {
		colTrimmed := strings.TrimSpace(col)

		switch {
		// Date: "対象年月日"
		case strings.Contains(colTrimmed, "対象年月日") || strings.Contains(colTrimmed, "date"):
			indices["date"] = i

		// Area: "エリア名"
		case colTrimmed == "エリア名" || strings.Contains(colTrimmed, "area"):
			indices["area"] = i

		// Demand: "エリア需要(MW)" (NOT "広域ブロック需要")
		case colTrimmed == "エリア需要(MW)":
			indices["demand"] = i

		// Capacity: "エリア供給力(MW)" (NOT "広域ブロック供給力")
		case colTrimmed == "エリア供給力(MW)":
			indices["capacity"] = i
		}
	}

	return indices
}

// normalizeArea converts various area representations to canonical form.
// Handles both English and Japanese names.
func (a *OCCTOAdapter) normalizeArea(areaStr string) string {
	areaStr = strings.ToLower(strings.TrimSpace(areaStr))

	switch {
	case areaStr == "tokyo" || areaStr == "東京" || strings.Contains(areaStr, "tokyo"):
		return "tokyo"
	case areaStr == "kansai" || areaStr == "関西" || strings.Contains(areaStr, "kansai"):
		return "kansai"
	default:
		// Return as-is for other areas (future expansion)
		return areaStr
	}
}
