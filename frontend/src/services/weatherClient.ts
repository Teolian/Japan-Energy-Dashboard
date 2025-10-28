// Weather API client using Open-Meteo (free, no auth required)
import type { WeatherResponse, SolarForecast, Area, SolarDataPoint } from '@/types/weather'
import { CITY_COORDINATES, PV_SYSTEM_PARAMS } from '@/types/weather'

const OPEN_METEO_BASE_URL = 'https://api.open-meteo.com/v1/forecast'

export async function fetchSolarForecast(area: Area, date: string): Promise<SolarForecast> {
  const coords = CITY_COORDINATES[area]

  // Parse date to get start and end dates for API
  const startDate = date
  const endDate = date // Same day forecast

  // Build API URL with required parameters
  const params = new URLSearchParams({
    latitude: coords.lat.toString(),
    longitude: coords.lon.toString(),
    start_date: startDate,
    end_date: endDate,
    hourly: [
      'shortwave_radiation',
      'direct_radiation',
      'diffuse_radiation',
      'cloud_cover',
      'temperature_2m'
    ].join(','),
    timezone: 'Asia/Tokyo'
  })

  const url = `${OPEN_METEO_BASE_URL}?${params.toString()}`

  try {
    const response = await fetch(url)

    if (!response.ok) {
      throw new Error(`Open-Meteo API error: ${response.statusText}`)
    }

    const data: WeatherResponse = await response.json()

    // Transform to our format
    const solarData: SolarDataPoint[] = data.hourly.time.map((time, i) => {
      const hour = new Date(time).getHours()
      const ghi = data.hourly.shortwave_radiation[i] || 0
      const dni = data.hourly.direct_radiation[i] || 0
      const dhi = data.hourly.diffuse_radiation[i] || 0
      const cloud_cover = data.hourly.cloud_cover[i] || 0
      const temperature = data.hourly.temperature_2m[i] || 0

      // Estimate PV generation (very simplified)
      // Real calculation: PV_MW = GHI * Area * Efficiency * Performance_Ratio / 1000
      // Here we assume 1 MW of installed capacity
      const estimated_pv_generation_mw = (ghi * PV_SYSTEM_PARAMS.efficiency * PV_SYSTEM_PARAMS.performance_ratio) / 1000

      return {
        ts: time,
        hour,
        ghi,
        dni,
        dhi,
        cloud_cover,
        temperature,
        estimated_pv_generation_mw
      }
    })

    // Calculate aggregates
    const validRadiation = solarData.filter(d => d.ghi > 0)
    const peak_radiation_hour = solarData.reduce((maxIdx, curr, idx, arr) =>
      curr.ghi > (arr[maxIdx]?.ghi || 0) ? idx : maxIdx, 0
    )
    const avg_radiation = validRadiation.length > 0
      ? validRadiation.reduce((sum, d) => sum + d.ghi, 0) / validRadiation.length
      : 0
    const total_radiation_kwh_m2 = solarData.reduce((sum, d) => sum + d.ghi, 0) / 1000 // W/m² to kWh/m²

    return {
      location: coords.name,
      latitude: data.latitude,
      longitude: data.longitude,
      date: startDate,
      timezone: data.timezone,
      data: solarData,
      peak_radiation_hour,
      avg_radiation,
      total_radiation_kwh_m2
    }
  } catch (error) {
    console.error('[WeatherClient] Failed to fetch solar forecast:', error)
    throw error
  }
}

// Mock data generator for offline development
export function generateMockSolarForecast(area: Area, date: string): SolarForecast {
  const coords = CITY_COORDINATES[area]

  const data: SolarDataPoint[] = Array.from({ length: 24 }, (_, hour) => {
    // Simulate solar curve: 0 at night, peak at noon
    let ghi = 0
    if (hour >= 6 && hour <= 18) {
      // Bell curve peaking at hour 12
      const hourFromNoon = Math.abs(hour - 12)
      ghi = 800 * Math.exp(-Math.pow(hourFromNoon / 4, 2))
    }

    const dni = ghi * 0.7 // Rough approximation
    const dhi = ghi * 0.3
    const cloud_cover = Math.random() * 40 // 0-40% clouds
    const temperature = 15 + Math.sin((hour - 6) / 12 * Math.PI) * 8 // Temp curve

    const estimated_pv_generation_mw = (ghi * PV_SYSTEM_PARAMS.efficiency * PV_SYSTEM_PARAMS.performance_ratio) / 1000

    return {
      ts: `${date}T${hour.toString().padStart(2, '0')}:00:00+09:00`,
      hour,
      ghi,
      dni,
      dhi,
      cloud_cover,
      temperature,
      estimated_pv_generation_mw
    }
  })

  const peak_radiation_hour = 12
  const validRadiation = data.filter(d => d.ghi > 0)
  const avg_radiation = validRadiation.reduce((sum, d) => sum + d.ghi, 0) / validRadiation.length
  const total_radiation_kwh_m2 = data.reduce((sum, d) => sum + d.ghi, 0) / 1000

  return {
    location: coords.name,
    latitude: coords.lat,
    longitude: coords.lon,
    date,
    timezone: 'Asia/Tokyo',
    data,
    peak_radiation_hour,
    avg_radiation,
    total_radiation_kwh_m2
  }
}
