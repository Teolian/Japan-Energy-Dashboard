// Package main fetches solar radiation forecast from Open-Meteo API.
// Usage: go run main.go -area tokyo -date 2025-11-07 -output forecast.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Hourly    struct {
		Time                []string  `json:"time"`
		ShortwaveRadiation  []float64 `json:"shortwave_radiation"`
		DirectRadiation     []float64 `json:"direct_radiation"`
		DiffuseRadiation    []float64 `json:"diffuse_radiation"`
		CloudCover          []float64 `json:"cloud_cover"`
		Temperature2m       []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

type SolarDataPoint struct {
	Timestamp              string  `json:"ts"`
	Hour                   int     `json:"hour"`
	GHI                    float64 `json:"ghi"`
	DNI                    float64 `json:"dni"`
	DHI                    float64 `json:"dhi"`
	CloudCover             float64 `json:"cloud_cover"`
	Temperature            float64 `json:"temperature"`
	EstimatedPVGenerationMW float64 `json:"estimated_pv_generation_mw"`
}

type SolarForecast struct {
	Location            string           `json:"location"`
	Latitude            float64          `json:"latitude"`
	Longitude           float64          `json:"longitude"`
	Date                string           `json:"date"`
	Timezone            string           `json:"timezone"`
	PeakRadiationHour   int              `json:"peak_radiation_hour"`
	AvgRadiation        float64          `json:"avg_radiation"`
	TotalRadiationKWhM2 float64          `json:"total_radiation_kwh_m2"`
	Data                []SolarDataPoint `json:"data"`
}

var cityCoordinates = map[string]struct {
	name string
	lat  float64
	lon  float64
}{
	"tokyo":  {"Tokyo", 35.6895, 139.6917},
	"kansai": {"Osaka", 34.6937, 135.5023},
}

const (
	pvEfficiency       = 0.20 // 20% panel efficiency
	performanceRatio   = 0.75 // 75% system performance ratio
	openMeteoBaseURL   = "https://api.open-meteo.com/v1/forecast"
)

func main() {
	var area, date, outputPath string
	flag.StringVar(&area, "area", "tokyo", "Area: tokyo or kansai")
	flag.StringVar(&date, "date", "", "Date in YYYY-MM-DD format")
	flag.StringVar(&outputPath, "output", "", "Output file path")
	flag.Parse()

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	coords, ok := cityCoordinates[area]
	if !ok {
		log.Fatalf("Invalid area: %s (must be tokyo or kansai)", area)
	}

	log.Printf("Fetching solar forecast for %s on %s...", coords.name, date)

	// Build API URL
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&start_date=%s&end_date=%s&hourly=shortwave_radiation,direct_radiation,diffuse_radiation,cloud_cover,temperature_2m&timezone=Asia/Tokyo",
		openMeteoBaseURL, coords.lat, coords.lon, date, date)

	// Fetch data
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch from Open-Meteo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Open-Meteo API error: %s", resp.Status)
	}

	var apiResp OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	// Transform to our format
	solarData := make([]SolarDataPoint, 0, len(apiResp.Hourly.Time))
	var totalGHI float64
	peakIdx := 0
	maxGHI := 0.0

	for i, timestamp := range apiResp.Hourly.Time {
		t, _ := time.Parse("2006-01-02T15:04", timestamp)
		hour := t.Hour()

		ghi := apiResp.Hourly.ShortwaveRadiation[i]
		dni := apiResp.Hourly.DirectRadiation[i]
		dhi := apiResp.Hourly.DiffuseRadiation[i]
		cloudCover := apiResp.Hourly.CloudCover[i]
		temp := apiResp.Hourly.Temperature2m[i]

		// Estimate PV generation (simplified)
		pvGenMW := (ghi * pvEfficiency * performanceRatio) / 1000

		solarData = append(solarData, SolarDataPoint{
			Timestamp:              timestamp,
			Hour:                   hour,
			GHI:                    ghi,
			DNI:                    dni,
			DHI:                    dhi,
			CloudCover:             cloudCover,
			Temperature:            temp,
			EstimatedPVGenerationMW: pvGenMW,
		})

		totalGHI += ghi
		if ghi > maxGHI {
			maxGHI = ghi
			peakIdx = hour
		}
	}

	// Calculate aggregates
	avgRadiation := 0.0
	if len(solarData) > 0 {
		avgRadiation = totalGHI / float64(len(solarData))
	}
	totalRadiationKWhM2 := totalGHI / 1000

	forecast := SolarForecast{
		Location:            coords.name,
		Latitude:            apiResp.Latitude,
		Longitude:           apiResp.Longitude,
		Date:                date,
		Timezone:            apiResp.Timezone,
		PeakRadiationHour:   peakIdx,
		AvgRadiation:        avgRadiation,
		TotalRadiationKWhM2: totalRadiationKWhM2,
		Data:                solarData,
	}

	// Write output
	output, err := json.MarshalIndent(forecast, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if outputPath == "" {
		outputPath = fmt.Sprintf("weather-%s-%s.json", area, date)
	}

	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}

	log.Printf("✅ Solar forecast saved to %s", outputPath)
	log.Printf("   Peak radiation: %.0f W/m² at %d:00", maxGHI, peakIdx)
	log.Printf("   Daily total: %.2f kWh/m²", totalRadiationKWhM2)
}
