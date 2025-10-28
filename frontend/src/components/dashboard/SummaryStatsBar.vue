<script setup lang="ts">
import { computed } from 'vue'
import { Activity, TrendingUp, DollarSign, Shield } from 'lucide-vue-next'
import { useDemandStore } from '@/stores/demand'
import { useJEPXStore } from '@/stores/jepx'
import { useReserveStore } from '@/stores/reserve'
import Tooltip from '@/components/common/Tooltip.vue'

const demandStore = useDemandStore()
const jepxStore = useJEPXStore()
const reserveStore = useReserveStore()

// Calculate combined peak from both areas
const totalPeak = computed(() => {
  const tokyoPeak = demandStore.tokyoMetrics?.peak || 0
  const kansaiPeak = demandStore.kansaiMetrics?.peak || 0
  return (tokyoPeak + kansaiPeak).toLocaleString('en-US')
})

// Calculate average price across both areas
const avgPrice = computed(() => {
  const tokyoPrices = jepxStore.priceValues('tokyo')
  const kansaiPrices = jepxStore.priceValues('kansai')

  const allPrices = [...tokyoPrices, ...kansaiPrices]
  if (allPrices.length === 0) return '—'

  const avg = allPrices.reduce((sum, p) => sum + p, 0) / allPrices.length
  return `¥${avg.toFixed(1)}`
})

// Get worst reserve status
const systemStatus = computed(() => {
  const tokyo = reserveStore.reserveForArea('tokyo')
  const kansai = reserveStore.reserveForArea('kansai')

  if (!tokyo || !kansai) return { label: 'Unknown', color: 'text-gray-500', bgColor: 'bg-gray-100 dark:bg-gray-800' }

  // Take the worse status
  const minReserve = Math.min(tokyo.reserve_margin_pct, kansai.reserve_margin_pct)

  if (minReserve >= 8) {
    return { label: 'Stable', color: 'text-green-600 dark:text-green-400', bgColor: 'bg-green-50 dark:bg-green-900/20' }
  } else if (minReserve >= 5) {
    return { label: 'Watch', color: 'text-yellow-600 dark:text-yellow-400', bgColor: 'bg-yellow-50 dark:bg-yellow-900/20' }
  } else {
    return { label: 'Tight', color: 'text-red-600 dark:text-red-400', bgColor: 'bg-red-50 dark:bg-red-900/20' }
  }
})

// Data freshness
const dataDate = computed(() => {
  return demandStore.currentDate
})
</script>

<template>
  <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-8 py-4">
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <!-- System Status -->
      <div class="flex items-center gap-3">
        <div :class="['p-2 rounded-lg', systemStatus.bgColor]">
          <Shield :size="20" :class="systemStatus.color" />
        </div>
        <div>
          <div class="text-xs text-gray-500 dark:text-gray-400">System Status</div>
          <Tooltip :content="`Power supply reserve: ${systemStatus.label}`">
            <div :class="['text-sm font-semibold cursor-help', systemStatus.color]">
              {{ systemStatus.label }}
            </div>
          </Tooltip>
        </div>
      </div>

      <!-- Combined Peak -->
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-lg bg-blue-50 dark:bg-blue-900/20">
          <Activity :size="20" class="text-blue-600 dark:text-blue-400" />
        </div>
        <div>
          <div class="text-xs text-gray-500 dark:text-gray-400">Combined Peak</div>
          <Tooltip content="Sum of Tokyo + Kansai peak demand">
            <div class="text-sm font-semibold text-gray-900 dark:text-white cursor-help">
              {{ totalPeak }} MW
            </div>
          </Tooltip>
        </div>
      </div>

      <!-- Average Price -->
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-lg bg-orange-50 dark:bg-orange-900/20">
          <DollarSign :size="20" class="text-orange-600 dark:text-orange-400" />
        </div>
        <div>
          <div class="text-xs text-gray-500 dark:text-gray-400">Avg Spot Price</div>
          <Tooltip content="Average JEPX spot price across both areas">
            <div class="text-sm font-semibold text-gray-900 dark:text-white cursor-help">
              {{ avgPrice }}/kWh
            </div>
          </Tooltip>
        </div>
      </div>

      <!-- Data Status -->
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-lg bg-purple-50 dark:bg-purple-900/20">
          <TrendingUp :size="20" class="text-purple-600 dark:text-purple-400" />
        </div>
        <div>
          <div class="text-xs text-gray-500 dark:text-gray-400">Data Date</div>
          <Tooltip content="Currently viewing data for this date">
            <div class="text-sm font-semibold text-gray-900 dark:text-white cursor-help">
              {{ dataDate }}
            </div>
          </Tooltip>
        </div>
      </div>
    </div>
  </div>
</template>
