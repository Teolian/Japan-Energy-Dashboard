// Feature flags composable
// Controls visibility of features via environment variables

import { computed } from 'vue'

export interface FeatureFlags {
  isReserveEnabled: boolean
  isJEPXEnabled: boolean
  isSettlementEnabled: boolean
  isComparisonEnabled: boolean
}

export function useFeatureFlags(): FeatureFlags {
  // Helper to parse env flag (defaults to true if not set or not explicitly "false")
  const isEnabled = (envVar: string | undefined): boolean => {
    return envVar !== 'false'
  }

  const isReserveEnabled = computed(() =>
    isEnabled(import.meta.env.VITE_FEAT_RESERVE)
  )

  const isJEPXEnabled = computed(() =>
    isEnabled(import.meta.env.VITE_FEAT_JEPX)
  )

  const isSettlementEnabled = computed(() =>
    isEnabled(import.meta.env.VITE_FEAT_SETTLEMENT)
  )

  const isComparisonEnabled = computed(() =>
    isEnabled(import.meta.env.VITE_FEAT_COMPARISON)
  )

  return {
    isReserveEnabled: isReserveEnabled.value,
    isJEPXEnabled: isJEPXEnabled.value,
    isSettlementEnabled: isSettlementEnabled.value,
    isComparisonEnabled: isComparisonEnabled.value
  }
}
