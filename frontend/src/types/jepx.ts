// JEPX (Japan Electric Power Exchange) spot price types
// Based on backend/internal/jepx/types.go and AGENT_TECH_SPEC §3.3

export interface PricePoint {
  ts: string // ISO8601 timestamp with Asia/Tokyo offset
  price: number // JPY/kWh
}

export interface JEPXSource {
  name: string
  url: string
  accessed_at: string // ISO8601
}

export interface JEPXMeta {
  min_price: number
  max_price: number
  avg_price: number
}

export interface JEPXResponse {
  date: string // YYYY-MM-DD
  area: string // "tokyo" | "kansai"
  timescale: string // "hourly"
  price_yen_per_kwh: PricePoint[]
  source: JEPXSource
  meta?: JEPXMeta
}

// Helper to extract price values for charting
export function extractPriceValues(response: JEPXResponse): number[] {
  return response.price_yen_per_kwh.map(p => p.price)
}

// Helper to get price at specific hour (0-23)
export function getPriceAtHour(response: JEPXResponse, hour: number): number | null {
  const point = response.price_yen_per_kwh[hour]
  return point ? point.price : null
}
