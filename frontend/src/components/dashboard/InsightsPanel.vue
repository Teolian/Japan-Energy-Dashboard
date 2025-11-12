<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Lightbulb, TrendingUp, TrendingDown, AlertCircle, CheckCircle } from 'lucide-vue-next'
import { useDemandStore } from '@/stores/demand'
import { useJEPXStore } from '@/stores/jepx'
import { useReserveStore } from '@/stores/reserve'
import { useSettlementStore } from '@/stores/settlement'

const { t } = useI18n()
const demandStore = useDemandStore()
const jepxStore = useJEPXStore()
const reserveStore = useReserveStore()
const settlementStore = useSettlementStore()

interface Insight {
  type: 'success' | 'warning' | 'info' | 'tip'
  icon: any
  title: string
  message: string
}

const insights = computed((): Insight[] => {
  const results: Insight[] = []

  // Find peak hour for Tokyo
  if (demandStore.tokyoData) {
    const series = demandStore.tokyoData.series
    const maxDemand = Math.max(...series.map(s => s.demand_mw))
    const peakPoint = series.find(s => s.demand_mw === maxDemand)

    if (peakPoint) {
      const hour = new Date(peakPoint.ts).getHours()
      results.push({
        type: 'info',
        icon: TrendingUp,
        title: t('insights.peakDemandTime'),
        message: t('insights.peakDemandMessage', { hour, demand: maxDemand.toLocaleString(), area: t('areas.tokyo') })
      })
    }
  }

  // Price analysis
  const tokyoPrices = jepxStore.priceValues('tokyo')
  if (tokyoPrices.length > 0) {
    const nightPrices = tokyoPrices.slice(0, 6) // 0-6am
    const dayPrices = tokyoPrices.slice(9, 20) // 9am-8pm

    if (nightPrices.length > 0 && dayPrices.length > 0) {
      const avgNight = nightPrices.reduce((a, b) => a + b, 0) / nightPrices.length
      const avgDay = dayPrices.reduce((a, b) => a + b, 0) / dayPrices.length
      const diff = ((avgDay - avgNight) / avgNight * 100).toFixed(0)

      if (parseFloat(diff) > 20) {
        results.push({
          type: 'tip',
          icon: Lightbulb,
          title: t('insights.costOptimization'),
          message: t('insights.costOptimizationMessage', { diff })
        })
      }
    }
  }

  // Reserve margin analysis
  const tokyoReserve = reserveStore.reserveForArea('tokyo')
  const kansaiReserve = reserveStore.reserveForArea('kansai')

  if (tokyoReserve && kansaiReserve) {
    const minReserve = Math.min(tokyoReserve.reserve_margin_pct, kansaiReserve.reserve_margin_pct)

    if (minReserve >= 8) {
      results.push({
        type: 'success',
        icon: CheckCircle,
        title: t('insights.powerSupplyStable'),
        message: t('insights.powerSupplyStableMessage', { reserve: minReserve.toFixed(1) })
      })
    } else if (minReserve >= 5) {
      results.push({
        type: 'warning',
        icon: AlertCircle,
        title: t('insights.reserveUnderWatch'),
        message: t('insights.reserveUnderWatchMessage', { reserve: minReserve.toFixed(1) })
      })
    } else {
      results.push({
        type: 'warning',
        icon: AlertCircle,
        title: t('insights.tightPowerSupply'),
        message: t('insights.tightPowerSupplyMessage', { reserve: minReserve.toFixed(1) })
      })
    }
  }

  // Settlement cost insight
  if (settlementStore.totals) {
    const cost = settlementStore.totals.cost_yen
    const kwh = settlementStore.totals.kwh

    if (cost && kwh) {
      const avgRate = cost / kwh
      results.push({
        type: 'info',
        icon: TrendingDown,
        title: 'Settlement Cost',
        message: `Effective rate: ¥${avgRate.toFixed(2)}/kWh for ${kwh.toLocaleString()} kWh consumption (¥${cost.toLocaleString()} total)`
      })
    }
  }

  return results
})

const getIconColor = (type: string) => {
  switch (type) {
    case 'success': return 'text-green-600 dark:text-green-400'
    case 'warning': return 'text-yellow-600 dark:text-yellow-400'
    case 'tip': return 'text-purple-600 dark:text-purple-400'
    default: return 'text-blue-600 dark:text-blue-400'
  }
}

const getBgColor = (type: string) => {
  switch (type) {
    case 'success': return 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800'
    case 'warning': return 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800'
    case 'tip': return 'bg-purple-50 dark:bg-purple-900/20 border-purple-200 dark:border-purple-800'
    default: return 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800'
  }
}
</script>

<template>
  <div v-if="insights.length > 0" class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
    <div class="flex items-center gap-2 mb-4">
      <Lightbulb :size="20" class="text-gray-700 dark:text-gray-300" />
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
        {{ t('dashboard.insights') }}
      </h2>
    </div>

    <div class="space-y-3">
      <div
        v-for="(insight, index) in insights"
        :key="index"
        :class="['flex gap-3 p-3 rounded-lg border', getBgColor(insight.type)]"
      >
        <component :is="insight.icon" :size="18" :class="['flex-shrink-0 mt-0.5', getIconColor(insight.type)]" />
        <div>
          <div :class="['text-sm font-semibold mb-1', getIconColor(insight.type)]">
            {{ insight.title }}
          </div>
          <div class="text-sm text-gray-700 dark:text-gray-300">
            {{ insight.message }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
