<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTradingStore } from '@/stores/trading'
import { TrendingUp, TrendingDown, Zap, Calculator } from 'lucide-vue-next'
import type { BatteryROI } from '@/types/trading'

const { t } = useI18n()
const tradingStore = useTradingStore()

// Battery ROI Calculator inputs
const batterySize = ref(50) // MWh
const cyclesPerDay = ref(1.5)
const efficiency = ref(0.85)
const capitalCost = ref(5000000000) // 50億円 for 50MWh

// Computed ROI
const roi = computed((): BatteryROI | null => {
  if (!tradingStore.opportunities || tradingStore.opportunities.length === 0) return null

  return tradingStore.calculateBatteryROI(
    batterySize.value,
    cyclesPerDay.value,
    efficiency.value,
    capitalCost.value
  )
})

const getSignalIcon = (type: string) => {
  return type === 'buy' ? TrendingUp : TrendingDown
}

const getSignalColor = (type: string) => {
  return type === 'buy'
    ? 'text-green-600 dark:text-green-400 bg-green-50 dark:bg-green-900/20'
    : 'text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20'
}

const getConfidenceBadge = (confidence: string) => {
  const colors = {
    high: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300',
    medium: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300',
    low: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
  }
  return colors[confidence as keyof typeof colors] || colors.medium
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <Zap :size="24" class="text-energy-purple" />
          {{ t('arbitrage.title') }}
        </h3>
        <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
          {{ tradingStore.bestOpportunities.length }} {{ t('trading.opportunities') }}
        </p>
      </div>
    </div>

    <!-- Opportunities List -->
    <div class="space-y-3">
      <div
        v-for="(opp, index) in tradingStore.bestOpportunities"
        :key="index"
        :class="[
          'p-4 rounded-lg border-2 transition-all hover:shadow-lg',
          getSignalColor(opp.type),
          opp.confidence === 'high' ? 'border-opacity-50' : 'border-opacity-30'
        ]"
      >
        <!-- Signal Header -->
        <div class="flex items-start justify-between mb-3">
          <div class="flex items-center gap-3">
            <component
              :is="getSignalIcon(opp.type)"
              :size="28"
              :class="[
                'flex-shrink-0',
                opp.type === 'buy' ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
              ]"
            />
            <div>
              <div class="flex items-center gap-2">
                <span class="font-bold text-lg uppercase">
                  {{ opp.type === 'buy' ? t('arbitrage.buySignal') : t('arbitrage.sellSignal') }}
                </span>
                <span :class="['text-xs px-2 py-1 rounded-full font-medium', getConfidenceBadge(opp.confidence)]">
                  {{ t(`arbitrage.${opp.confidence}`) }}
                </span>
              </div>
              <div class="text-sm opacity-80 mt-0.5">
                {{ opp.hour }}:00 JST
              </div>
            </div>
          </div>

          <!-- Expected Profit -->
          <div class="text-right">
            <div class="text-xs opacity-70">{{ t('arbitrage.profit') }}</div>
            <div class="text-xl font-bold">
              ¥{{ opp.expectedProfit.toLocaleString() }}
            </div>
            <div class="text-xs opacity-70">per MWh</div>
          </div>
        </div>

        <!-- Price Info -->
        <div class="grid grid-cols-2 gap-3 mb-3 text-sm">
          <div class="bg-white dark:bg-gray-800 bg-opacity-50 p-2 rounded">
            <div class="opacity-70 text-xs">Current Price</div>
            <div class="font-semibold">¥{{ opp.currentPrice.toFixed(2) }}/kWh</div>
          </div>
          <div class="bg-white dark:bg-gray-800 bg-opacity-50 p-2 rounded">
            <div class="opacity-70 text-xs">Target Price</div>
            <div class="font-semibold">¥{{ opp.targetPrice.toFixed(2) }}/kWh</div>
          </div>
        </div>

        <!-- Recommendation -->
        <div class="text-sm font-medium mb-2">
          {{ opp.recommendation }}
        </div>

        <!-- Reasoning -->
        <div class="text-xs opacity-80 space-y-1">
          <div v-for="(reason, idx) in opp.reasoning" :key="idx" class="flex items-start gap-1">
            <span class="opacity-50">•</span>
            <span>{{ reason }}</span>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div
        v-if="tradingStore.bestOpportunities.length === 0"
        class="text-center py-8 text-gray-500 dark:text-gray-400"
      >
        <Zap :size="48" class="mx-auto mb-2 opacity-30" />
        <p>No arbitrage opportunities detected</p>
        <p class="text-sm mt-1">Price spreads are minimal in current market conditions</p>
      </div>
    </div>

    <!-- Battery ROI Calculator -->
    <div v-if="roi" class="mt-6 p-5 bg-gradient-to-br from-purple-50 to-blue-50 dark:from-purple-900/20 dark:to-blue-900/20 rounded-xl border border-purple-200 dark:border-purple-800">
      <div class="flex items-center gap-2 mb-4">
        <Calculator :size="20" class="text-energy-purple" />
        <h4 class="font-bold text-gray-900 dark:text-white">{{ t('arbitrage.roi') }}</h4>
      </div>

      <!-- Inputs -->
      <div class="grid grid-cols-2 gap-4 mb-4">
        <div>
          <label class="text-xs text-gray-600 dark:text-gray-400">{{ t('arbitrage.batterySize') }}</label>
          <input
            v-model.number="batterySize"
            type="number"
            class="w-full mt-1 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-sm"
            min="1"
            max="500"
          />
        </div>
        <div>
          <label class="text-xs text-gray-600 dark:text-gray-400">{{ t('arbitrage.cyclesPerDay') }}</label>
          <input
            v-model.number="cyclesPerDay"
            type="number"
            step="0.1"
            class="w-full mt-1 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-sm"
            min="0.5"
            max="3"
          />
        </div>
        <div>
          <label class="text-xs text-gray-600 dark:text-gray-400">{{ t('arbitrage.efficiency') }} (%)</label>
          <input
            v-model.number="efficiency"
            type="number"
            step="0.01"
            class="w-full mt-1 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-sm"
            min="0.5"
            max="1"
          />
        </div>
        <div>
          <label class="text-xs text-gray-600 dark:text-gray-400">Capital Cost (億円)</label>
          <input
            v-model.number="capitalCost"
            type="number"
            step="100000000"
            class="w-full mt-1 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-sm"
          />
        </div>
      </div>

      <!-- Results -->
      <div class="grid grid-cols-3 gap-3">
        <div class="bg-white dark:bg-gray-800 p-3 rounded-lg text-center">
          <div class="text-xs text-gray-600 dark:text-gray-400 mb-1">{{ t('arbitrage.estimatedProfit') }}</div>
          <div class="text-lg font-bold text-green-600 dark:text-green-400">
            ¥{{ roi.dailyProfitJPY.toLocaleString() }}
          </div>
          <div class="text-xs opacity-70">/day</div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-3 rounded-lg text-center">
          <div class="text-xs text-gray-600 dark:text-gray-400 mb-1">Payback Period</div>
          <div class="text-lg font-bold text-blue-600 dark:text-blue-400">
            {{ roi.paybackYears }}
          </div>
          <div class="text-xs opacity-70">years</div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-3 rounded-lg text-center">
          <div class="text-xs text-gray-600 dark:text-gray-400 mb-1">ROI</div>
          <div class="text-lg font-bold text-purple-600 dark:text-purple-400">
            {{ roi.roi }}%
          </div>
          <div class="text-xs opacity-70">annual</div>
        </div>
      </div>
    </div>
  </div>
</template>
