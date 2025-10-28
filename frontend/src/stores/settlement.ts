// Pinia store for settlement cost calculations
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { SettlementRequest, SettlementResponse, ProfilePoint, Totals, HourlyBreakdown, Assumptions } from '@/types/settlement'
import { runSettlement as runSettlementAPI } from '@/services/dataClient'

export const useSettlementStore = defineStore('settlement', () => {
  // State
  const data = ref<SettlementResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const totals = computed((): Totals | null => {
    return data.value?.totals || null
  })

  const byHour = computed((): HourlyBreakdown[] => {
    return data.value?.by_hour || []
  })

  const assumptions = computed((): Assumptions | null => {
    return data.value?.assumptions || null
  })

  const source = computed(() => {
    return data.value?.source_prices || null
  })

  const period = computed(() => {
    return data.value?.period || null
  })

  // Formatted values for display
  const formattedTotalCost = computed((): string => {
    if (!totals.value) return '—'
    return totals.value.cost_yen.toLocaleString('ja-JP', {
      minimumFractionDigits: 1,
      maximumFractionDigits: 1
    })
  })

  const formattedTotalKWh = computed((): string => {
    if (!totals.value) return '—'
    return totals.value.kwh.toLocaleString('ja-JP', {
      minimumFractionDigits: 1,
      maximumFractionDigits: 1
    })
  })

  // Actions
  async function runSettlement(profile: ProfilePoint[], area: string, date: string, pvOffsetPct: number = 0.15) {
    loading.value = true
    error.value = null

    try {
      const request: SettlementRequest = {
        profile,
        prices: { area, date },
        pv_offset_pct: pvOffsetPct
      }

      data.value = await runSettlementAPI(request)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to run settlement'
      console.error('[settlement-store] Settlement error:', err)
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
    totals,
    byHour,
    assumptions,
    source,
    period,
    formattedTotalCost,
    formattedTotalKWh,
    // Actions
    runSettlement,
    reset
  }
})
