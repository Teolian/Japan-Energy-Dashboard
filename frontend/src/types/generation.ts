// Generation mix types - electricity generation by fuel type
// Estimated from demand + JEPX price correlation

export interface GenerationPoint {
  ts: string  // ISO8601 timestamp
  solar_mw: number
  wind_mw: number
  hydro_mw: number
  nuclear_mw: number
  lng_mw: number
  coal_mw: number
  other_mw: number
  total_mw: number
}

export interface GenerationSource {
  name: string
  url: string
}

export interface GenerationMeta {
  avg_renewable_pct: number     // (Solar + Wind + Hydro) / Total * 100
  avg_carbon_gco2_kwh: number   // Average carbon intensity
  peak_solar_mw: number
  peak_wind_mw: number
}

export interface GenerationResponse {
  date: string          // YYYY-MM-DD
  area: string          // tokyo | kansai
  timezone: string      // Asia/Tokyo
  timescale: string     // hourly
  series: GenerationPoint[]
  source: GenerationSource
  meta?: GenerationMeta
}

// Chart-friendly data structure
export interface GenerationChartData {
  time: string  // HH:00 format
  solar: number
  wind: number
  hydro: number
  nuclear: number
  lng: number
  coal: number
  other: number
  total: number
  renewable_pct: number
  carbon_gco2_kwh: number
}

// Carbon intensity classification
export type CarbonLevel = 'low' | 'medium' | 'high' | 'very-high'

export interface CarbonIntensity {
  value: number
  level: CarbonLevel
  label: string
  color: string
}

// Helper: Get carbon intensity level
export function getCarbonLevel(gco2kwh: number): CarbonIntensity {
  if (gco2kwh < 200) {
    return {
      value: gco2kwh,
      level: 'low',
      label: 'Very Clean',
      color: 'text-green-600 dark:text-green-400'
    }
  } else if (gco2kwh < 300) {
    return {
      value: gco2kwh,
      level: 'medium',
      label: 'Clean',
      color: 'text-blue-600 dark:text-blue-400'
    }
  } else if (gco2kwh < 350) {
    return {
      value: gco2kwh,
      level: 'high',
      label: 'Moderate',
      color: 'text-orange-600 dark:text-orange-400'
    }
  } else {
    return {
      value: gco2kwh,
      level: 'very-high',
      label: 'High Carbon',
      color: 'text-red-600 dark:text-red-400'
    }
  }
}

// Helper: Calculate renewable percentage for a point
export function calculateRenewablePct(point: GenerationPoint): number {
  if (point.total_mw === 0) return 0
  const renewable = point.solar_mw + point.wind_mw + point.hydro_mw
  return (renewable / point.total_mw) * 100
}

// Helper: Calculate carbon intensity for a point
export function calculateCarbonIntensity(point: GenerationPoint): number {
  if (point.total_mw === 0) return 0
  // Emission factors (gCO2/kWh): LNG 350, Coal 850, Other 500
  const carbon = (point.lng_mw * 350 + point.coal_mw * 850 + point.other_mw * 500) / point.total_mw
  return carbon
}

// Helper: Transform response to chart data
export function toChartData(response: GenerationResponse): GenerationChartData[] {
  return response.series.map(point => {
    const renewable_pct = calculateRenewablePct(point)
    const carbon_gco2_kwh = calculateCarbonIntensity(point)

    return {
      time: new Date(point.ts).toLocaleTimeString('ja-JP', {
        hour: '2-digit',
        minute: '2-digit',
        timeZone: 'Asia/Tokyo'
      }),
      solar: point.solar_mw,
      wind: point.wind_mw,
      hydro: point.hydro_mw,
      nuclear: point.nuclear_mw,
      lng: point.lng_mw,
      coal: point.coal_mw,
      other: point.other_mw,
      total: point.total_mw,
      renewable_pct,
      carbon_gco2_kwh
    }
  })
}
