// Settlement-lite types matching backend/internal/settlement/types.go
// Based on AGENT_TECH_SPEC.md §3.5

export interface ProfilePoint {
  ts: string // ISO8601 with Asia/Tokyo offset
  kwh: number // Consumption in kWh
}

export interface PricesRequest {
  area: string // "tokyo" | "kansai"
  date: string // YYYY-MM-DD
}

export interface SettlementRequest {
  profile: ProfilePoint[] // Hourly consumption profile
  prices: PricesRequest // JEPX price reference
  pv_offset_pct: number // PV offset percentage (0.0-1.0)
}

export interface Period {
  from: string // ISO8601 timestamp
  to: string // ISO8601 timestamp
}

export interface Totals {
  kwh: number // Total consumption in kWh
  cost_yen: number // Total cost in JPY
}

export interface HourlyBreakdown {
  ts: string // ISO8601 with Asia/Tokyo offset
  kwh: number // Consumption in kWh
  price: number // Price in JPY/kWh
  cost: number // Cost in JPY (kwh × price × (1 - pv%))
}

export interface Assumptions {
  pv_offset_pct: number // PV offset percentage
  area: string // Price area
}

export interface SettlementSource {
  name: string // e.g., "JEPX"
  url: string // Original data source URL
}

export interface SettlementResponse {
  period: Period // Time range
  totals: Totals // Aggregated results
  by_hour: HourlyBreakdown[] // Per-hour breakdown
  assumptions: Assumptions // Calculation parameters
  source_prices: SettlementSource // Price data attribution
}

// Helper to convert demand data (MW) to profile (kWh)
// Assumes hourly intervals: 1 hour × MW = MWh, then × 1000 = kWh
export function demandToProfile(demandMw: number[], timestamps: string[]): ProfilePoint[] {
  return timestamps.map((ts, i) => ({
    ts,
    kwh: (demandMw[i] || 0) * 1000 // MW × 1 hour = MWh, then MWh × 1000 = kWh
  }))
}
