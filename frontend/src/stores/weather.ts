import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { SolarForecast, Area } from '@/types/weather'
import { fetchSolarForecast, generateMockSolarForecast } from '@/services/weatherClient'

export const useWeatherStore = defineStore('weather', () => {
  // State
  const tokyoForecast = ref<SolarForecast | null>(null)
  const kansaiForecast = ref<SolarForecast | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const useMockData = ref(false) // Toggle for offline development

  // Getters
  const hasForecastData = computed(() => tokyoForecast.value !== null || kansaiForecast.value !== null)

  const forecastForArea = computed(() => (area: Area) => {
    if (area === 'tokyo') return tokyoForecast.value
    if (area === 'kansai') return kansaiForecast.value
    return null
  })

  // Get radiation values as array for charting
  const radiationValues = computed(() => (area: Area) => {
    const forecast = area === 'tokyo' ? tokyoForecast.value : kansaiForecast.value
    return forecast?.data.map(d => d.ghi) || []
  })

  // Get estimated PV generation values
  const pvGenerationValues = computed(() => (area: Area) => {
    const forecast = area === 'tokyo' ? tokyoForecast.value : kansaiForecast.value
    return forecast?.data.map(d => d.estimated_pv_generation_mw || 0) || []
  })

  // Get cloud cover values
  const cloudCoverValues = computed(() => (area: Area) => {
    const forecast = area === 'tokyo' ? tokyoForecast.value : kansaiForecast.value
    return forecast?.data.map(d => d.cloud_cover) || []
  })

  // Actions
  async function fetchForecast(area: Area, date: string) {
    loading.value = true
    error.value = null

    try {
      let forecast: SolarForecast

      if (useMockData.value) {
        // Use mock data for offline development
        forecast = generateMockSolarForecast(area, date)
      } else {
        // Try real API, fall back to mock on error
        try {
          forecast = await fetchSolarForecast(area, date)
        } catch (apiError) {
          console.warn(`[WeatherStore] API failed for ${area}, using mock data:`, apiError)
          forecast = generateMockSolarForecast(area, date)
        }
      }

      // Store based on area
      if (area === 'tokyo') {
        tokyoForecast.value = forecast
      } else if (area === 'kansai') {
        kansaiForecast.value = forecast
      }

      return forecast
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch weather data'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchBothAreas(date: string) {
    await Promise.all([
      fetchForecast('tokyo', date),
      fetchForecast('kansai', date)
    ])
  }

  function toggleMockData() {
    useMockData.value = !useMockData.value
  }

  function clearData() {
    tokyoForecast.value = null
    kansaiForecast.value = null
    error.value = null
  }

  return {
    // State
    tokyoForecast,
    kansaiForecast,
    loading,
    error,
    useMockData,

    // Getters
    hasForecastData,
    forecastForArea,
    radiationValues,
    pvGenerationValues,
    cloudCoverValues,

    // Actions
    fetchForecast,
    fetchBothAreas,
    toggleMockData,
    clearData
  }
})
