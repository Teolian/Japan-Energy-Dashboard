// Weather and solar radiation types for Open-Meteo API
export interface WeatherResponse {
  latitude: number
  longitude: number
  timezone: string
  hourly: {
    time: string[] // ISO 8601 timestamps
    shortwave_radiation: number[] // W/m² - Global Horizontal Irradiance (GHI)
    direct_radiation: number[] // W/m² - Direct Normal Irradiance (DNI)
    diffuse_radiation: number[] // W/m² - Diffuse Horizontal Irradiance (DHI)
    cloud_cover: number[] // % - Cloud cover percentage
    temperature_2m: number[] // °C - Temperature at 2m height
  }
}

export interface SolarDataPoint {
  ts: string // ISO 8601 timestamp
  hour: number // 0-23
  ghi: number // W/m² - Global Horizontal Irradiance
  dni: number // W/m² - Direct Normal Irradiance
  dhi: number // W/m² - Diffuse Horizontal Irradiance
  cloud_cover: number // % - Cloud coverage
  temperature: number // °C
  estimated_pv_generation_mw?: number // Estimated PV generation (if capacity known)
}

export interface SolarForecast {
  location: string
  latitude: number
  longitude: number
  date: string
  timezone: string
  data: SolarDataPoint[]
  peak_radiation_hour: number
  avg_radiation: number // Daily average W/m²
  total_radiation_kwh_m2: number // Daily total kWh/m²
}

// Location coordinates for Japanese cities
export const CITY_COORDINATES = {
  tokyo: { lat: 35.6762, lon: 139.6503, name: 'Tokyo' },
  kansai: { lat: 34.6937, lon: 135.5023, name: 'Osaka (Kansai)' },
  hokkaido: { lat: 43.0642, lon: 141.3469, name: 'Sapporo' },
  kyushu: { lat: 33.5904, lon: 130.4017, name: 'Fukuoka' }
} as const

export type Area = keyof typeof CITY_COORDINATES

// PV system assumptions for estimation
export const PV_SYSTEM_PARAMS = {
  efficiency: 0.18, // 18% panel efficiency (typical crystalline silicon)
  performance_ratio: 0.75, // 75% performance ratio (accounts for losses)
  capacity_mw_per_km2: 50 // Typical utility-scale solar farm density
}
