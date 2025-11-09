// Data client with dual-mode support (MOCK/LIVE)
// Per AGENT_INTEGRATION_GUIDE.md §2-3

import type { Area, DemandResponse, DataMode } from '@/types/demand'
import type { ReserveResponse, ReserveStatus } from '@/types/reserve'
import type { JEPXResponse, PricePoint } from '@/types/jepx'
import type { SettlementRequest, SettlementResponse } from '@/types/settlement'
import type { GenerationResponse } from '@/types/generation'

const STORAGE_KEY = 'jp-energy-data-mode'

// Get data mode from URL query, localStorage, or env
function getDataMode(): DataMode {
  // 1. Check URL query ?mode=live or ?mode=mock
  const urlParams = new URLSearchParams(window.location.search)
  const urlMode = urlParams.get('mode')
  if (urlMode === 'live' || urlMode === 'mock') {
    // Persist to localStorage
    localStorage.setItem(STORAGE_KEY, urlMode)
    return urlMode
  }

  // 2. Check localStorage
  const storedMode = localStorage.getItem(STORAGE_KEY)
  if (storedMode === 'live' || storedMode === 'mock') {
    return storedMode
  }

  // 3. Check env at build time
  const envMode = import.meta.env.VITE_DATA_MODE
  if (envMode === 'live' || envMode === 'mock') {
    return envMode
  }

  // 4. Default to mock
  return 'mock'
}

// Generate mock data matching backend JSON structure
function generateMockData(area: Area, date: string): DemandResponse {
  const baseValues = area === 'tokyo' ?
    { min: 25000, max: 40000 } :
    { min: 12000, max: 19000 }

  const series = Array.from({ length: 24 }, (_, hour) => {
    const hourFactor = Math.sin((hour - 6) * Math.PI / 12) * 0.3 + 0.7
    const randomness = Math.random() * 0.1 + 0.95
    const actual = Math.round((baseValues.min + (baseValues.max - baseValues.min) * hourFactor) * randomness)
    const forecast = Math.round(actual * (Math.random() * 0.1 + 0.95))

    return {
      ts: `${date}T${hour.toString().padStart(2, '0')}:00:00+09:00`,
      demand_mw: actual,
      forecast_mw: forecast
    }
  })

  return {
    area,
    date,
    timezone: 'Asia/Tokyo',
    timescale: 'hourly',
    series,
    source: {
      name: area === 'tokyo' ? 'TEPCO (Mock)' : 'Kansai (Mock)',
      url: 'https://example.com/mock'
    }
  }
}

// Fetch live data from static JSON files (GitHub Actions approach)
async function fetchLiveData(area: Area, date: string): Promise<DemandResponse> {
  const response = await fetch(`/data/jp/${area}/demand-${date}.json`)

  if (!response.ok) {
    throw new Error(`Failed to fetch ${area} demand for ${date}: ${response.statusText}`)
  }

  const data: DemandResponse = await response.json()
  return data
}

// Main data fetcher with fallback
export async function getDemandData(area: Area, date: string): Promise<DemandResponse> {
  const mode = getDataMode()

  try {
    if (mode === 'live') {
      try {
        return await fetchLiveData(area, date)
      } catch (error) {
        console.warn(`[live-data] Failed to fetch live data for ${area}/${date}, falling back to mock`, error)
        // Graceful fallback to mock
        const mockData = generateMockData(area, date)
        mockData.meta = {
          warning: 'Live data unavailable, showing mock data'
        }
        return mockData
      }
    }
  } catch {
    // Mode check failed, use mock
  }

  return generateMockData(area, date)
}

// Get current data mode (for UI display)
export function getCurrentDataMode(): DataMode {
  return getDataMode()
}

// Set data mode programmatically
export function setDataMode(mode: DataMode): void {
  localStorage.setItem(STORAGE_KEY, mode)
  window.location.reload() // Reload to apply new mode
}

// === RESERVE MARGIN DATA ===

// Generate mock reserve margin data
function generateMockReserveData(date: string): ReserveResponse {
  // Random but stable reserves (using date as seed)
  const seed = date.split('-').reduce((acc, val) => acc + parseInt(val), 0)
  const tokyoReserve = 5 + ((seed % 10) / 10) * 5 // 5-10%
  const kansaiReserve = 6 + ((seed % 12) / 10) * 4 // 6-10%

  const deriveStatus = (pct: number): ReserveStatus => {
    if (pct >= 8) return 'stable'
    if (pct >= 5) return 'watch'
    return 'tight'
  }

  return {
    date,
    areas: [
      {
        area: 'tokyo',
        reserve_margin_pct: Math.round(tokyoReserve * 10) / 10,
        status: deriveStatus(tokyoReserve)
      },
      {
        area: 'kansai',
        reserve_margin_pct: Math.round(kansaiReserve * 10) / 10,
        status: deriveStatus(kansaiReserve)
      }
    ],
    source: {
      name: 'OCCTO (Mock)',
      url: 'https://www.occto.or.jp/'
    }
  }
}

// Fetch live reserve margin data from static JSON files
async function fetchLiveReserveData(date: string): Promise<ReserveResponse> {
  const response = await fetch(`/data/jp/system/reserve-${date}.json`)

  if (!response.ok) {
    throw new Error(`Failed to fetch reserve data for ${date}: ${response.statusText}`)
  }

  const data: ReserveResponse = await response.json()
  return data
}

// Main reserve data fetcher with fallback
export async function getReserveData(date: string): Promise<ReserveResponse> {
  const mode = getDataMode()

  try {
    if (mode === 'live') {
      try {
        return await fetchLiveReserveData(date)
      } catch (error) {
        console.warn(`[live-data] Failed to fetch live reserve data for ${date}, falling back to mock`, error)
        // Graceful fallback to mock
        const mockData = generateMockReserveData(date)
        mockData.meta = {
          warning: 'Live data unavailable, showing mock data'
        }
        return mockData
      }
    }
  } catch {
    // Mode check failed, use mock
  }

  return generateMockReserveData(date)
}

// === JEPX SPOT PRICE DATA ===

// Generate mock JEPX spot price data
function generateMockJEPXData(area: Area, date: string): JEPXResponse {
  // Generate realistic hourly prices (20-45 JPY/kWh)
  // Lower at night (0-6), higher during day (9-20)
  const seed = date.split('-').reduce((acc, val) => acc + parseInt(val), 0)

  const pricePoints: PricePoint[] = Array.from({ length: 24 }, (_, hour) => {
    // Base price varies by area
    const basePrice = area === 'tokyo' ? 30 : 28

    // Time-of-day factor: low at night, high during day
    let todFactor = 1.0
    if (hour >= 0 && hour < 6) todFactor = 0.7 // Night (low)
    else if (hour >= 9 && hour < 20) todFactor = 1.3 // Day (high)
    else todFactor = 0.9 // Morning/evening

    // Add some randomness seeded by date + hour
    const randomness = ((seed + hour) % 10) / 10 * 0.2 + 0.9 // 0.9-1.1

    const price = Math.round(basePrice * todFactor * randomness * 10) / 10

    return {
      ts: `${date}T${hour.toString().padStart(2, '0')}:00:00+09:00`,
      price
    }
  })

  const prices = pricePoints.map(p => p.price)
  const minPrice = Math.min(...prices)
  const maxPrice = Math.max(...prices)
  const avgPrice = Math.round((prices.reduce((a, b) => a + b, 0) / prices.length) * 10) / 10

  return {
    date,
    area,
    timescale: 'hourly',
    price_yen_per_kwh: pricePoints,
    source: {
      name: 'JEPX (Mock)',
      url: 'https://www.jepx.jp/',
      accessed_at: new Date().toISOString()
    },
    meta: {
      min_price: minPrice,
      max_price: maxPrice,
      avg_price: avgPrice
    }
  }
}

// Fetch live JEPX price data from static JSON files
async function fetchLiveJEPXData(area: Area, date: string): Promise<JEPXResponse> {
  const response = await fetch(`/data/jp/jepx/spot-${area}-${date}.json`)

  if (!response.ok) {
    throw new Error(`Failed to fetch JEPX data for ${area}/${date}: ${response.statusText}`)
  }

  const data: JEPXResponse = await response.json()
  return data
}

// Main JEPX data fetcher with fallback
export async function getJEPXData(area: Area, date: string): Promise<JEPXResponse> {
  const mode = getDataMode()

  try {
    if (mode === 'live') {
      try {
        return await fetchLiveJEPXData(area, date)
      } catch (error) {
        console.warn(`[live-data] Failed to fetch live JEPX data for ${area}/${date}, falling back to mock`, error)
        // Graceful fallback to mock
        const mockData = generateMockJEPXData(area, date)
        mockData.meta = {
          ...mockData.meta,
          warning: 'Live data unavailable, showing mock data'
        } as any
        return mockData
      }
    }
  } catch {
    // Mode check failed, use mock
  }

  return generateMockJEPXData(area, date)
}

// === SETTLEMENT COST CALCULATION ===

// Generate mock settlement result
function generateMockSettlementData(req: SettlementRequest): SettlementResponse {
  // Calculate settlement cost based on profile and mock prices
  const { profile, prices, pv_offset_pct } = req

  if (profile.length === 0) {
    throw new Error('Profile cannot be empty')
  }

  // Generate mock prices if needed (or use actual JEPX data)
  const basePrice = prices.area === 'tokyo' ? 30 : 28
  const hourlyPrices = profile.map((_, idx) => {
    const hour = idx % 24
    let todFactor = 1.0
    if (hour >= 0 && hour < 6) todFactor = 0.7
    else if (hour >= 9 && hour < 20) todFactor = 1.3
    else todFactor = 0.9
    return Math.round(basePrice * todFactor * 10) / 10
  })

  // Calculate hourly breakdown
  const byHour = profile.map((point, idx) => {
    const price = hourlyPrices[idx] || 0
    const effectiveKWh = point.kwh * (1 - pv_offset_pct)
    const cost = effectiveKWh * price

    return {
      ts: point.ts,
      kwh: Math.round(point.kwh * 10) / 10,
      price: Math.round(price * 10) / 10,
      cost: Math.round(cost * 10) / 10
    }
  })

  // Calculate totals
  const totalKWh = profile.reduce((sum, p) => sum + p.kwh, 0)
  const totalCost = byHour.reduce((sum, h) => sum + h.cost, 0)

  // Determine period
  const timestamps = profile.map(p => p.ts).sort()
  const from = timestamps[0]!
  const to = timestamps[timestamps.length - 1]!

  return {
    period: { from, to },
    totals: {
      kwh: Math.round(totalKWh * 10) / 10,
      cost_yen: Math.round(totalCost * 10) / 10
    },
    by_hour: byHour,
    assumptions: {
      pv_offset_pct,
      area: prices.area
    },
    source_prices: {
      name: 'JEPX (Mock)',
      url: 'https://www.jepx.jp/'
    }
  }
}

// Fetch live settlement calculation from API
async function fetchLiveSettlement(req: SettlementRequest): Promise<SettlementResponse> {
  const response = await fetch('/api/settlements/run', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(req)
  })

  if (!response.ok) {
    throw new Error(`Failed to run settlement: ${response.statusText}`)
  }

  const data: SettlementResponse = await response.json()
  return data
}

// Main settlement runner with fallback
export async function runSettlement(req: SettlementRequest): Promise<SettlementResponse> {
  const mode = getDataMode()

  try {
    if (mode === 'live') {
      try {
        return await fetchLiveSettlement(req)
      } catch (error) {
        console.warn('[live-data] Failed to run live settlement, falling back to mock', error)
        return generateMockSettlementData(req)
      }
    }
  } catch {
    // Mode check failed, use mock
  }

  return generateMockSettlementData(req)
}

// ==================== GENERATION MIX DATA ====================

// Fetch generation mix data (estimated from demand + JEPX)
export async function fetchGenerationMix(area: Area, date: string): Promise<GenerationResponse> {
  const mode = getDataMode()

  // LIVE mode: fetch from static JSON files
  if (mode === 'live') {
    const response = await fetch(`/data/jp/${area}/generation-${date}.json`)

    if (!response.ok) {
      throw new Error(`Failed to fetch generation data for ${area}/${date}: ${response.statusText}`)
    }

    return await response.json()
  }

  // MOCK mode: no mock data for generation yet, throw error
  throw new Error('Generation mix data only available in LIVE mode')
}
