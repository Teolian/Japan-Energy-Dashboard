// Package generation provides estimation logic for generation mix when real data unavailable.
package generation

import (
	"fmt"
	"math"
	"time"

	"github.com/teo/aversome/backend/internal/demand"
	"github.com/teo/aversome/backend/internal/jepx"
)

// Estimator estimates generation mix from demand and price data.
type Estimator struct{}

// NewEstimator creates a new generation mix estimator.
func NewEstimator() *Estimator {
	return &Estimator{}
}

// EstimateFromDemandAndPrice estimates generation mix using demand and JEPX price patterns.
// Logic:
// - Solar peaks 11:00-14:00, correlates with low prices (duck curve)
// - Nuclear provides base load (~25-30% constant)
// - LNG/Coal fill remaining demand
// - Wind/Hydro small constant percentage
func (e *Estimator) EstimateFromDemandAndPrice(demandResp *demand.Response, jepxResp *jepx.Response) (*Response, error) {
	if len(demandResp.Series) == 0 || len(jepxResp.PriceYenPerKwh) == 0 {
		return nil, fmt.Errorf("empty demand or price data")
	}

	// Calculate price statistics for solar correlation
	_, minPrice, maxPrice := e.calculatePriceStats(jepxResp)

	resp := NewResponse(string(demandResp.Area), demandResp.Date)
	resp.Source = Source{
		Name: "Estimated (demand + price correlation)",
		URL:  "Internal calculation",
	}

	// Generate hourly generation points
	for i, demandPoint := range demandResp.Series {
		if i >= len(jepxResp.PriceYenPerKwh) {
			break
		}

		price := jepxResp.PriceYenPerKwh[i].Price
		totalDemand := demandPoint.DemandMW
		hour := demandPoint.Timestamp.Hour()

		// Estimate solar generation (peaks at noon, correlates with low prices)
		solarMW := e.estimateSolar(hour, price, minPrice, maxPrice, totalDemand)

		// Nuclear base load (25-30% in Japan)
		nuclearMW := totalDemand * 0.27

		// Wind (low penetration in Japan, ~3%)
		windMW := totalDemand * 0.03

		// Hydro (pumped storage + run-of-river, ~8%)
		hydroMW := totalDemand * 0.08

		// Remaining demand filled by fossil fuels
		fossilMW := totalDemand - solarMW - nuclearMW - windMW - hydroMW
		if fossilMW < 0 {
			fossilMW = 0
		}

		// Split fossil: 60% LNG, 30% Coal, 10% Other (Japan's typical mix)
		lngMW := fossilMW * 0.60
		coalMW := fossilMW * 0.30
		otherMW := fossilMW * 0.10

		point := GenerationPoint{
			Timestamp: demandPoint.Timestamp,
			SolarMW:   solarMW,
			WindMW:    windMW,
			HydroMW:   hydroMW,
			NuclearMW: nuclearMW,
			LNGMW:     lngMW,
			CoalMW:    coalMW,
			OtherMW:   otherMW,
			TotalMW:   totalDemand,
		}

		resp.Series = append(resp.Series, point)
	}

	// Calculate metadata
	resp.CalculateMeta()

	return resp, nil
}

// estimateSolar estimates solar generation based on hour and price.
// Solar peaks at midday (11:00-14:00) and correlates with low prices.
func (e *Estimator) estimateSolar(hour int, price, minPrice, maxPrice, totalDemand float64) float64 {
	// Time-based solar curve (0 at night, peak at noon)
	timeFactor := e.solarTimeCurve(hour)

	// Price-based adjustment (low prices → more solar)
	// When price is at minimum, solar is at maximum potential
	priceFactor := 1.0
	if maxPrice > minPrice {
		// Normalize price to 0-1 range, then invert (low price = high solar)
		normalizedPrice := (price - minPrice) / (maxPrice - minPrice)
		priceFactor = 1.0 - (normalizedPrice * 0.3) // Up to 30% variation based on price
	}

	// Maximum solar penetration in Japan: ~15-20% of total demand at peak
	maxSolarPct := 0.18

	solarMW := totalDemand * maxSolarPct * timeFactor * priceFactor

	// Ensure non-negative
	if solarMW < 0 {
		solarMW = 0
	}

	return solarMW
}

// solarTimeCurve returns solar generation factor (0-1) based on hour.
// Models typical solar generation curve: ramps up 6am, peaks 11am-2pm, down by 6pm.
func (e *Estimator) solarTimeCurve(hour int) float64 {
	// No solar at night (10pm - 5am)
	if hour >= 22 || hour <= 5 {
		return 0
	}

	// Sunrise ramp (6am - 10am): quadratic growth
	if hour >= 6 && hour <= 10 {
		t := float64(hour-6) / 4.0 // 0 to 1
		return t * t               // Quadratic ramp
	}

	// Peak hours (11am - 2pm): maximum output
	if hour >= 11 && hour <= 14 {
		return 1.0
	}

	// Sunset decline (3pm - 6pm): quadratic decline
	if hour >= 15 && hour <= 18 {
		t := float64(18-hour) / 3.0 // 1 to 0
		return t * t                 // Quadratic decline
	}

	// Evening tail (7pm - 9pm): minimal
	if hour >= 19 && hour <= 21 {
		return 0.05
	}

	return 0
}

// calculatePriceStats computes average, min, and max prices.
func (e *Estimator) calculatePriceStats(jepxResp *jepx.Response) (avg, min, max float64) {
	if len(jepxResp.PriceYenPerKwh) == 0 {
		return 0, 0, 0
	}

	sum := 0.0
	min = jepxResp.PriceYenPerKwh[0].Price
	max = jepxResp.PriceYenPerKwh[0].Price

	for _, p := range jepxResp.PriceYenPerKwh {
		sum += p.Price
		if p.Price < min {
			min = p.Price
		}
		if p.Price > max {
			max = p.Price
		}
	}

	avg = sum / float64(len(jepxResp.PriceYenPerKwh))
	return avg, min, max
}

// EstimateWithSeasonalAdjustment applies seasonal adjustments to base estimation.
// Winter: higher nuclear (heating demand), Summer: higher solar (longer days).
func (e *Estimator) EstimateWithSeasonalAdjustment(baseResp *Response, date time.Time) *Response {
	// Determine season
	month := date.Month()

	// Summer months (June-August): boost solar by 10%
	isSummer := month >= 6 && month <= 8

	// Winter months (December-February): reduce solar by 20%, boost nuclear
	isWinter := month == 12 || month <= 2

	for i := range baseResp.Series {
		if isSummer {
			baseResp.Series[i].SolarMW *= 1.10
			// Rebalance: reduce fossil slightly
			reduction := baseResp.Series[i].SolarMW * 0.10
			baseResp.Series[i].LNGMW -= reduction * 0.6
			baseResp.Series[i].CoalMW -= reduction * 0.4
		}

		if isWinter {
			baseResp.Series[i].SolarMW *= 0.80
			// Increase nuclear base load
			baseResp.Series[i].NuclearMW *= 1.05
			// Rebalance
			increase := baseResp.Series[i].NuclearMW * 0.05
			baseResp.Series[i].LNGMW -= increase * 0.6
			baseResp.Series[i].CoalMW -= increase * 0.4
		}

		// Ensure no negatives
		baseResp.Series[i].LNGMW = math.Max(baseResp.Series[i].LNGMW, 0)
		baseResp.Series[i].CoalMW = math.Max(baseResp.Series[i].CoalMW, 0)
	}

	// Recalculate metadata
	baseResp.CalculateMeta()

	return baseResp
}
