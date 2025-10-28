// Pinia store for JEPX spot price data
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Area } from '@/types/demand'
import type { JEPXResponse } from '@/types/jepx'
import { getJEPXData } from '@/services/dataClient'

export const useJEPXStore = defineStore('jepx', () => {
  // State - store prices per area
  const tokyoData = ref<JEPXResponse | null>(null)
  const kansaiData = ref<JEPXResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const pricesForArea = computed(() => {
    return (area: Area): JEPXResponse | null => {
      return area === 'tokyo' ? tokyoData.value : kansaiData.value
    }
  })

  const hasWarning = computed(() => {
    return (area: Area): boolean => {
      const data = pricesForArea.value(area)
      return (data?.meta as any)?.warning !== undefined
    }
  })

  const warning = computed(() => {
    return (area: Area): string | null => {
      const data = pricesForArea.value(area)
      return (data?.meta as any)?.warning || null
    }
  })

  const source = computed(() => {
    return (area: Area) => {
      const data = pricesForArea.value(area)
      return data?.source || null
    }
  })

  // Get price values array (for charting)
  const priceValues = computed(() => {
    return (area: Area): number[] => {
      const data = pricesForArea.value(area)
      return data?.price_yen_per_kwh.map(p => p.price) || []
    }
  })

  // Get price at specific hour
  const priceAtHour = computed(() => {
    return (area: Area, hour: number): number | null => {
      const data = pricesForArea.value(area)
      if (!data) return null
      const point = data.price_yen_per_kwh[hour]
      return point ? point.price : null
    }
  })

  // Actions
  async function fetchJEPXData(area: Area, date: string) {
    loading.value = true
    error.value = null

    try {
      const data = await getJEPXData(area, date)
      if (area === 'tokyo') {
        tokyoData.value = data
      } else {
        kansaiData.value = data
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch JEPX data'
      console.error(`[jepx-store] Fetch error for ${area}:`, err)
    } finally {
      loading.value = false
    }
  }

  // Fetch both areas at once
  async function fetchBothAreas(date: string) {
    loading.value = true
    error.value = null

    try {
      const [tokyo, kansai] = await Promise.all([
        getJEPXData('tokyo', date),
        getJEPXData('kansai', date)
      ])
      tokyoData.value = tokyo
      kansaiData.value = kansai
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch JEPX data'
      console.error('[jepx-store] Fetch error:', err)
    } finally {
      loading.value = false
    }
  }

  function reset() {
    tokyoData.value = null
    kansaiData.value = null
    loading.value = false
    error.value = null
  }

  return {
    // State
    tokyoData,
    kansaiData,
    loading,
    error,
    // Getters
    pricesForArea,
    hasWarning,
    warning,
    source,
    priceValues,
    priceAtHour,
    // Actions
    fetchJEPXData,
    fetchBothAreas,
    reset
  }
})
