// Pinia store for OCCTO reserve margin data
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Area } from '@/types/demand'
import type { ReserveResponse, AreaReserve, ReserveStatus, StatusConfig } from '@/types/reserve'
import { STATUS_CONFIGS } from '@/types/reserve'
import { getReserveData } from '@/services/dataClient'

export const useReserveStore = defineStore('reserve', () => {
  // State
  const data = ref<ReserveResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const reserveForArea = computed(() => {
    return (area: Area): AreaReserve | null => {
      if (!data.value) return null
      return data.value.areas.find(a => a.area === area) || null
    }
  })

  const statusConfig = computed(() => {
    return (status: ReserveStatus): StatusConfig => {
      return STATUS_CONFIGS[status]
    }
  })

  const hasWarning = computed(() => {
    return data.value?.meta?.warning !== undefined
  })

  const warning = computed(() => {
    return data.value?.meta?.warning || null
  })

  const source = computed(() => {
    return data.value?.source || null
  })

  // Actions
  async function fetchReserveData(date: string) {
    loading.value = true
    error.value = null

    try {
      data.value = await getReserveData(date)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch reserve data'
      console.error('[reserve-store] Fetch error:', err)
    } finally {
      loading.value = false
    }
  }

  function reset() {
    data.value = null
    loading.value = false
    error.value = null
  }

  return {
    // State
    data,
    loading,
    error,
    // Getters
    reserveForArea,
    statusConfig,
    hasWarning,
    warning,
    source,
    // Actions
    fetchReserveData,
    reset
  }
})
