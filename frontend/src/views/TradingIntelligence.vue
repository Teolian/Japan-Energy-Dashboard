<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTradingStore } from '@/stores/trading'
import { useDemandStore } from '@/stores/demand'
import { Brain, TrendingUp, Zap } from 'lucide-vue-next'

// Components
import ArbitragePanel from '@/components/trading/ArbitragePanel.vue'
import LoadShiftAdvisor from '@/components/trading/LoadShiftAdvisor.vue'
import BaseCard from '@/components/common/BaseCard.vue'
import ChartLoadingSkeleton from '@/components/common/ChartLoadingSkeleton.vue'

const { t } = useI18n()
const tradingStore = useTradingStore()
const demandStore = useDemandStore()

const selectedArea = ref<'tokyo' | 'kansai'>('tokyo')
const activeTab = ref<'arbitrage' | 'loadshift'>('arbitrage')

const tabs = [
  { id: 'arbitrage' as const, label: t('trading.arbitrage'), icon: Zap },
  { id: 'loadshift' as const, label: t('trading.loadShift'), icon: TrendingUp }
]

// Initial analysis
onMounted(() => {
  runAnalysis()
})

// Re-run analysis when area changes
watch(selectedArea, () => {
  runAnalysis()
})

// Re-run when date changes
watch(() => demandStore.currentDate, () => {
  runAnalysis()
})

function runAnalysis() {
  tradingStore.analyzeArbitrageOpportunities(selectedArea.value)
  tradingStore.generateLoadShiftRecommendations(selectedArea.value)
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 p-8">
    <!-- Hero Header -->
    <div class="mb-8 relative overflow-hidden rounded-2xl bg-gradient-to-br from-purple-600 via-blue-600 to-cyan-500 p-8 shadow-energy-lg">
      <div class="relative z-10">
        <div class="flex items-center gap-3 mb-2">
          <Brain :size="40" class="text-white" />
          <h1 class="text-4xl font-bold text-white">
            {{ t('trading.title') }}
          </h1>
        </div>
        <p class="text-white/90 text-lg">
          {{ t('trading.subtitle') }}
        </p>

        <!-- Quick Stats -->
        <div v-if="tradingStore.metrics" class="grid grid-cols-4 gap-4 mt-6">
          <div class="bg-white/10 backdrop-blur-sm rounded-lg p-4">
            <div class="text-white/70 text-sm mb-1">{{ t('trading.opportunities') }}</div>
            <div class="text-3xl font-bold text-white">
              {{ tradingStore.metrics.totalOpportunities }}
            </div>
          </div>

          <div class="bg-white/10 backdrop-blur-sm rounded-lg p-4">
            <div class="text-white/70 text-sm mb-1">Daily Savings</div>
            <div class="text-3xl font-bold text-white">
              ¥{{ (tradingStore.metrics.estimatedDailySavings / 1000).toFixed(0) }}K
            </div>
          </div>

          <div class="bg-white/10 backdrop-blur-sm rounded-lg p-4">
            <div class="text-white/70 text-sm mb-1">Monthly Potential</div>
            <div class="text-3xl font-bold text-white">
              ¥{{ (tradingStore.metrics.estimatedMonthlySavings / 1000000).toFixed(1) }}M
            </div>
          </div>

          <div class="bg-white/10 backdrop-blur-sm rounded-lg p-4">
            <div class="text-white/70 text-sm mb-1">Optimal Battery</div>
            <div class="text-3xl font-bold text-white">
              {{ tradingStore.metrics.optimalBatterySize }} MWh
            </div>
          </div>
        </div>
      </div>

      <!-- Background decoration -->
      <div class="absolute top-0 right-0 w-96 h-96 bg-white/5 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 w-80 h-80 bg-cyan-400/10 rounded-full blur-3xl"></div>
    </div>

    <!-- Area Selector -->
    <div class="mb-6 flex items-center gap-4">
      <label class="text-sm font-medium text-gray-700 dark:text-gray-300">Area:</label>
      <div class="flex gap-2">
        <button
          @click="selectedArea = 'tokyo'"
          :class="[
            'px-4 py-2 rounded-lg font-medium transition-all',
            selectedArea === 'tokyo'
              ? 'bg-blue-600 text-white shadow-lg'
              : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
          ]"
        >
          {{ t('areas.tokyo') }}
        </button>
        <button
          @click="selectedArea = 'kansai'"
          :class="[
            'px-4 py-2 rounded-lg font-medium transition-all',
            selectedArea === 'kansai'
              ? 'bg-blue-600 text-white shadow-lg'
              : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
          ]"
        >
          {{ t('areas.kansai') }}
        </button>
      </div>
    </div>

    <!-- Tabs -->
    <div class="mb-6">
      <div class="flex gap-2 bg-gray-100 dark:bg-gray-800 rounded-lg p-1 inline-flex">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            'flex items-center gap-2 px-6 py-3 rounded-md font-medium transition-all',
            activeTab === tab.id
              ? 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white shadow-md'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
          ]"
        >
          <component :is="tab.icon" :size="20" />
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Content -->
    <div v-if="tradingStore.loading">
      <BaseCard>
        <ChartLoadingSkeleton />
      </BaseCard>
    </div>

    <div v-else>
      <!-- Arbitrage Panel -->
      <BaseCard v-if="activeTab === 'arbitrage'" class="shadow-energy">
        <ArbitragePanel />
      </BaseCard>

      <!-- Load Shift Advisor -->
      <BaseCard v-if="activeTab === 'loadshift'" class="shadow-energy">
        <LoadShiftAdvisor />
      </BaseCard>
    </div>

    <!-- Footer Info -->
    <div class="mt-8 text-center text-sm text-gray-500 dark:text-gray-400">
      <p>
        💡 Analysis based on JEPX spot prices and demand patterns for {{ demandStore.currentDate }}
      </p>
      <p class="mt-1">
        AI-powered recommendations · {{ t('dataMode.' + demandStore.dataMode) }}
      </p>
    </div>
  </div>
</template>
