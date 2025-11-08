// Package generation provides types for electricity generation mix data.
package generation

import (
	"encoding/json"
	"time"
)

// Source represents the data source metadata.
type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// GenerationPoint represents generation capacity by fuel type at a specific time.
type GenerationPoint struct {
	Timestamp  time.Time `json:"ts"`
	SolarMW    float64   `json:"solar_mw"`
	WindMW     float64   `json:"wind_mw"`
	HydroMW    float64   `json:"hydro_mw"`
	NuclearMW  float64   `json:"nuclear_mw"`
	LNGMW      float64   `json:"lng_mw"`
	CoalMW     float64   `json:"coal_mw"`
	OtherMW    float64   `json:"other_mw"`
	TotalMW    float64   `json:"total_mw"`
}

// Response represents the complete generation mix response.
type Response struct {
	Date       string            `json:"date"`        // YYYY-MM-DD
	Area       string            `json:"area"`        // tokyo, kansai, etc.
	Timezone   string            `json:"timezone"`    // Asia/Tokyo
	Timescale  string            `json:"timescale"`   // hourly
	Series     []GenerationPoint `json:"series"`
	Source     Source            `json:"source"`
	Meta       *Meta             `json:"meta,omitempty"`
}

// Meta contains aggregated metrics.
type Meta struct {
	AvgRenewablePct  float64 `json:"avg_renewable_pct"`  // (Solar + Wind + Hydro) / Total
	AvgCarbonGCO2KWh float64 `json:"avg_carbon_gco2_kwh"` // Average carbon intensity
	PeakSolarMW      float64 `json:"peak_solar_mw"`
	PeakWindMW       float64 `json:"peak_wind_mw"`
}

// NewResponse creates a new generation mix response.
func NewResponse(area, date string) *Response {
	return &Response{
		Date:      date,
		Area:      area,
		Timezone:  "Asia/Tokyo",
		Timescale: "hourly",
		Series:    []GenerationPoint{},
	}
}

// CalculateMeta computes aggregated metrics from series data.
func (r *Response) CalculateMeta() {
	if len(r.Series) == 0 {
		return
	}

	var renewablePctSum float64
	var carbonSum float64
	var peakSolar, peakWind float64

	for _, point := range r.Series {
		if point.TotalMW > 0 {
			renewableMW := point.SolarMW + point.WindMW + point.HydroMW
			renewablePct := (renewableMW / point.TotalMW) * 100
			renewablePctSum += renewablePct

			// Calculate carbon intensity (simplified emission factors in gCO2/kWh)
			// Nuclear: 0, Solar: 0, Wind: 0, Hydro: 0
			// LNG: 350 gCO2/kWh, Coal: 850 gCO2/kWh, Other: 500 gCO2/kWh (avg)
			carbonIntensity := (point.LNGMW*350 + point.CoalMW*850 + point.OtherMW*500) / point.TotalMW
			carbonSum += carbonIntensity
		}

		if point.SolarMW > peakSolar {
			peakSolar = point.SolarMW
		}
		if point.WindMW > peakWind {
			peakWind = point.WindMW
		}
	}

	count := float64(len(r.Series))
	r.Meta = &Meta{
		AvgRenewablePct:  renewablePctSum / count,
		AvgCarbonGCO2KWh: carbonSum / count,
		PeakSolarMW:      peakSolar,
		PeakWindMW:       peakWind,
	}
}

// MarshalJSON implements custom JSON serialization.
func (r *Response) MarshalJSON() ([]byte, error) {
	type Alias Response

	// Format timestamps as ISO8601
	type SeriesAlias struct {
		Timestamp  string  `json:"ts"`
		SolarMW    float64 `json:"solar_mw"`
		WindMW     float64 `json:"wind_mw"`
		HydroMW    float64 `json:"hydro_mw"`
		NuclearMW  float64 `json:"nuclear_mw"`
		LNGMW      float64 `json:"lng_mw"`
		CoalMW     float64 `json:"coal_mw"`
		OtherMW    float64 `json:"other_mw"`
		TotalMW    float64 `json:"total_mw"`
	}

	series := make([]SeriesAlias, len(r.Series))
	for i, point := range r.Series {
		series[i] = SeriesAlias{
			Timestamp: point.Timestamp.Format(time.RFC3339),
			SolarMW:   point.SolarMW,
			WindMW:    point.WindMW,
			HydroMW:   point.HydroMW,
			NuclearMW: point.NuclearMW,
			LNGMW:     point.LNGMW,
			CoalMW:    point.CoalMW,
			OtherMW:   point.OtherMW,
			TotalMW:   point.TotalMW,
		}
	}

	return json.Marshal(&struct {
		*Alias
		Series []SeriesAlias `json:"series"`
	}{
		Alias:  (*Alias)(r),
		Series: series,
	})
}
