// Pinia store for generation mix data
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  GenerationResponse,
  GenerationChartData
} from '@/types/generation'
import { toChartData, getCarbonLevel } from '@/types/generation'
import { fetchGenerationMix } from '@/services/dataClient'

export const useGenerationStore = defineStore('generation', () => {
  // State
  const tokyoData = ref<GenerationResponse | null>(null)
  const kansaiData = ref<GenerationResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Computed: Chart data
  const tokyoChartData = computed<GenerationChartData[]>(() => {
    if (!tokyoData.value) return []
    return toChartData(tokyoData.value)
  })

  const kansaiChartData = computed<GenerationChartData[]>(() => {
    if (!kansaiData.value) return []
    return toChartData(kansaiData.value)
  })

  // Computed: Tokyo metrics
  const tokyoMetrics = computed(() => {
    if (!tokyoData.value?.meta) return null

    const meta = tokyoData.value.meta
    const carbonLevel = getCarbonLevel(meta.avg_carbon_gco2_kwh)

    return {
      renewablePct: meta.avg_renewable_pct,
      carbonIntensity: carbonLevel,
      peakSolarMW: meta.peak_solar_mw,
      peakWindMW: meta.peak_wind_mw
    }
  })

  // Computed: Kansai metrics
  const kansaiMetrics = computed(() => {
    if (!kansaiData.value?.meta) return null

    const meta = kansaiData.value.meta
    const carbonLevel = getCarbonLevel(meta.avg_carbon_gco2_kwh)

    return {
      renewablePct: meta.avg_renewable_pct,
      carbonIntensity: carbonLevel,
      peakSolarMW: meta.peak_solar_mw,
      peakWindMW: meta.peak_wind_mw
    }
  })

  // Computed: Greenest hour (highest renewable %)
  const greenestHour = computed(() => {
    if (tokyoChartData.value.length === 0) return null

    let maxRenewable = tokyoChartData.value[0]
    if (!maxRenewable) return null

    tokyoChartData.value.forEach(point => {
      if (maxRenewable && point.renewable_pct > maxRenewable.renewable_pct) {
        maxRenewable = point
      }
    })

    return {
      time: maxRenewable.time,
      renewablePct: maxRenewable.renewable_pct,
      carbonGCO2: maxRenewable.carbon_gco2_kwh
    }
  })

  // Computed: Cleanest hour (lowest carbon)
  const cleanestHour = computed(() => {
    if (tokyoChartData.value.length === 0) return null

    let minCarbon = tokyoChartData.value[0]
    if (!minCarbon) return null

    tokyoChartData.value.forEach(point => {
      if (minCarbon && point.carbon_gco2_kwh < minCarbon.carbon_gco2_kwh) {
        minCarbon = point
      }
    })

    return {
      time: minCarbon.time,
      carbonGCO2: minCarbon.carbon_gco2_kwh,
      renewablePct: minCarbon.renewable_pct
    }
  })

  // Actions
  async function fetchTokyo(date: string) {
    loading.value = true
    error.value = null

    try {
      const data = await fetchGenerationMix('tokyo', date)
      tokyoData.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch Tokyo generation data'
      console.error('Error fetching Tokyo generation:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchKansai(date: string) {
    loading.value = true
    error.value = null

    try {
      const data = await fetchGenerationMix('kansai', date)
      kansaiData.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch Kansai generation data'
      console.error('Error fetching Kansai generation:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchBoth(date: string) {
    await Promise.all([
      fetchTokyo(date),
      fetchKansai(date)
    ])
  }

  function clearData() {
    tokyoData.value = null
    kansaiData.value = null
    error.value = null
  }

  return {
    // State
    tokyoData,
    kansaiData,
    loading,
    error,

    // Computed
    tokyoChartData,
    kansaiChartData,
    tokyoMetrics,
    kansaiMetrics,
    greenestHour,
    cleanestHour,

    // Actions
    fetchTokyo,
    fetchKansai,
    fetchBoth,
    clearData
  }
})
