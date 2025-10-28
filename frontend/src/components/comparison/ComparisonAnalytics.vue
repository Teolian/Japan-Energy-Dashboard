<script setup lang="ts">
import { computed } from 'vue'
import { useDemandStore } from '@/stores/demand'
import { useJEPXStore } from '@/stores/jepx'
import { useReserveStore } from '@/stores/reserve'
import { Zap, DollarSign, AlertTriangle, Lightbulb } from 'lucide-vue-next'

const demandStore = useDemandStore()
const jepxStore = useJEPXStore()
const reserveStore = useReserveStore()

// Calculate analytics
const analytics = computed(() => {
  const tokyoPrices = jepxStore.priceValues('tokyo')
  const kansaiPrices = jepxStore.priceValues('kansai')

  if (tokyoPrices.length === 0 || kansaiPrices.length === 0) {
    return null
  }

  // Average prices
  const tokyoAvgPrice = tokyoPrices.reduce((a, b) => a + b, 0) / tokyoPrices.length
  const kansaiAvgPrice = kansaiPrices.reduce((a, b) => a + b, 0) / kansaiPrices.length
  const priceDiff = Math.abs(tokyoAvgPrice - kansaiAvgPrice)
  const priceDiffPct = (priceDiff / Math.min(tokyoAvgPrice, kansaiAvgPrice)) * 100
  const cheaperArea = tokyoAvgPrice < kansaiAvgPrice ? 'Tokyo' : 'Kansai'

  // Peak price hours
  const tokyoPeakHour = tokyoPrices.indexOf(Math.max(...tokyoPrices))
  const kansaiPeakHour = kansaiPrices.indexOf(Math.max(...kansaiPrices))
  const peakDiff = Math.abs(tokyoPeakHour - kansaiPeakHour)

  // Arbitrage potential (max price difference in same hour)
  const hourlyDiffs = tokyoPrices.map((tp, i) => Math.abs(tp - (kansaiPrices[i] || 0)))
  const maxArbitrage = Math.max(...hourlyDiffs)
  const arbitrageHour = hourlyDiffs.indexOf(maxArbitrage)

  // Demand comparison
  const tokyoPeak = demandStore.tokyoMetrics?.peak || 0
  const kansaiPeak = demandStore.kansaiMetrics?.peak || 0
  const demandRatio = (tokyoPeak / kansaiPeak).toFixed(2)

  // Reserve comparison
  const tokyoReserve = reserveStore.reserveForArea('tokyo')
  const kansaiReserve = reserveStore.reserveForArea('kansai')
  const reserveDiff = tokyoReserve && kansaiReserve
    ? Math.abs(tokyoReserve.reserve_margin_pct - kansaiReserve.reserve_margin_pct)
    : null
  const saferArea = tokyoReserve && kansaiReserve
    ? (tokyoReserve.reserve_margin_pct > kansaiReserve.reserve_margin_pct ? 'Tokyo' : 'Kansai')
    : null

  // Cost savings calculation (for 1000 kWh consumption)
  const savingsPerMWh = priceDiff * 1000
  const dailySavings = savingsPerMWh // JPY per day for 1 MWh consumption

  return {
    tokyoAvgPrice,
    kansaiAvgPrice,
    priceDiff,
    priceDiffPct,
    cheaperArea,
    tokyoPeakHour,
    kansaiPeakHour,
    peakDiff,
    maxArbitrage,
    arbitrageHour,
    demandRatio,
    reserveDiff,
    saferArea,
    dailySavings,
    tokyoPeak,
    kansaiPeak
  }
})
</script>

<template>
  <div v-if="analytics" class="space-y-4">
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Regional Analytics</h3>

    <!-- Key Metrics Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Price Difference -->
      <div class="p-4 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
        <div class="flex items-center gap-2 mb-2">
          <DollarSign :size="18" class="text-blue-600 dark:text-blue-400" />
          <span class="text-sm font-medium text-blue-900 dark:text-blue-200">Price Spread</span>
        </div>
        <div class="text-2xl font-bold text-blue-900 dark:text-blue-100">
          ¥{{ analytics.priceDiff.toFixed(2) }}
        </div>
        <div class="text-xs text-blue-700 dark:text-blue-300 mt-1">
          {{ analytics.cheaperArea }} is {{ analytics.priceDiffPct.toFixed(1) }}% cheaper
        </div>
      </div>

      <!-- Arbitrage Opportunity -->
      <div class="p-4 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800">
        <div class="flex items-center gap-2 mb-2">
          <Zap :size="18" class="text-green-600 dark:text-green-400" />
          <span class="text-sm font-medium text-green-900 dark:text-green-200">Max Arbitrage</span>
        </div>
        <div class="text-2xl font-bold text-green-900 dark:text-green-100">
          ¥{{ analytics.maxArbitrage.toFixed(2) }}
        </div>
        <div class="text-xs text-green-700 dark:text-green-300 mt-1">
          At {{ analytics.arbitrageHour }}:00
        </div>
      </div>

      <!-- Reserve Safety -->
      <div v-if="analytics.saferArea" class="p-4 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800">
        <div class="flex items-center gap-2 mb-2">
          <AlertTriangle :size="18" class="text-purple-600 dark:text-purple-400" />
          <span class="text-sm font-medium text-purple-900 dark:text-purple-200">Reserve Gap</span>
        </div>
        <div class="text-2xl font-bold text-purple-900 dark:text-purple-100">
          {{ analytics.reserveDiff?.toFixed(1) }}%
        </div>
        <div class="text-xs text-purple-700 dark:text-purple-300 mt-1">
          {{ analytics.saferArea }} has higher margin
        </div>
      </div>
    </div>

    <!-- Detailed Comparison Table -->
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
          <tr>
            <th class="px-4 py-3 text-left font-medium text-gray-700 dark:text-gray-300">Metric</th>
            <th class="px-4 py-3 text-right font-medium text-gray-700 dark:text-gray-300">Tokyo</th>
            <th class="px-4 py-3 text-right font-medium text-gray-700 dark:text-gray-300">Kansai</th>
            <th class="px-4 py-3 text-right font-medium text-gray-700 dark:text-gray-300">Difference</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <!-- Average Price -->
          <tr>
            <td class="px-4 py-3 text-gray-700 dark:text-gray-300">Avg Spot Price</td>
            <td class="px-4 py-3 text-right font-mono text-gray-900 dark:text-gray-100">
              ¥{{ analytics.tokyoAvgPrice.toFixed(2) }}
            </td>
            <td class="px-4 py-3 text-right font-mono text-gray-900 dark:text-gray-100">
              ¥{{ analytics.kansaiAvgPrice.toFixed(2) }}
            </td>
            <td class="px-4 py-3 text-right font-semibold" :class="analytics.tokyoAvgPrice < analytics.kansaiAvgPrice ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'">
              {{ analytics.priceDiffPct.toFixed(1) }}%
            </td>
          </tr>

          <!-- Peak Demand -->
          <tr>
            <td class="px-4 py-3 text-gray-700 dark:text-gray-300">Peak Demand</td>
            <td class="px-4 py-3 text-right font-mono text-gray-900 dark:text-gray-100">
              {{ analytics.tokyoPeak.toFixed(0) }} MW
            </td>
            <td class="px-4 py-3 text-right font-mono text-gray-900 dark:text-gray-100">
              {{ analytics.kansaiPeak.toFixed(0) }} MW
            </td>
            <td class="px-4 py-3 text-right font-mono text-gray-600 dark:text-gray-400">
              {{ analytics.demandRatio }}x
            </td>
          </tr>

          <!-- Peak Hour -->
          <tr>
            <td class="px-4 py-3 text-gray-700 dark:text-gray-300">Peak Hour</td>
            <td class="px-4 py-3 text-right font-mono text-gray-900 dark:text-gray-100">
              {{ analytics.tokyoPeakHour }}:00
            </td>
            <td class="px-4 py-3 text-right font-mono text-gray-900 dark:text-gray-100">
              {{ analytics.kansaiPeakHour }}:00
            </td>
            <td class="px-4 py-3 text-right text-gray-600 dark:text-gray-400">
              {{ analytics.peakDiff }}h shift
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Recommendations -->
    <div class="p-4 rounded-lg bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800">
      <div class="flex items-start gap-2">
        <Lightbulb :size="18" class="text-amber-600 dark:text-amber-400 mt-0.5 flex-shrink-0" />
        <div class="flex-1">
          <h4 class="text-sm font-semibold text-amber-900 dark:text-amber-200 mb-2">Optimization Recommendations</h4>
          <ul class="space-y-1.5 text-xs text-amber-800 dark:text-amber-300">
            <li>• <strong>Cost savings:</strong> Consuming 1 MWh in {{ analytics.cheaperArea }} saves ¥{{ analytics.dailySavings.toFixed(0) }}/day</li>
            <li v-if="analytics.peakDiff >= 2">
              • <strong>Load shifting:</strong> Peak hours are {{ analytics.peakDiff }}h apart - opportunity for inter-regional arbitrage
            </li>
            <li v-if="analytics.maxArbitrage > 5">
              • <strong>Battery arbitrage:</strong> Up to ¥{{ analytics.maxArbitrage.toFixed(1) }}/kWh spread at {{ analytics.arbitrageHour }}:00
            </li>
            <li v-if="analytics.saferArea">
              • <strong>Reliability:</strong> {{ analytics.saferArea }} has {{ analytics.reserveDiff?.toFixed(1) }}% higher reserve margin
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>

  <!-- No data state -->
  <div v-else class="text-center py-12 text-gray-500 dark:text-gray-400">
    <div class="space-y-2">
      <p class="font-medium">Regional Analytics Unavailable</p>
      <p class="text-sm">Price data needed for comparison</p>
    </div>
  </div>
</template>
