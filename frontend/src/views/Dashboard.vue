<script setup lang="ts">
import { onMounted, watch, computed } from 'vue'
import { useDemandStore } from '@/stores/demand'
import { useReserveStore } from '@/stores/reserve'
import { useJEPXStore } from '@/stores/jepx'
import { useSettlementStore } from '@/stores/settlement'
import { useDarkMode } from '@/composables/useDarkMode'
import { useFeatureFlags } from '@/composables/useFeatureFlags'
import { useExport } from '@/composables/useExport'
import { useKeyboardNavigation } from '@/composables/useKeyboardNavigation'
import BaseCard from '@/components/common/BaseCard.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import DemandChart from '@/components/demand/DemandChart.vue'
import ReserveBadge from '@/components/reserve/ReserveBadge.vue'
import CostCard from '@/components/settlement/CostCard.vue'
import ChartLoadingSkeleton from '@/components/common/ChartLoadingSkeleton.vue'
import InfoIcon from '@/components/common/InfoIcon.vue'
import Sparkline from '@/components/common/Sparkline.vue'
import TrendIndicator from '@/components/common/TrendIndicator.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import SummaryStatsBar from '@/components/dashboard/SummaryStatsBar.vue'
import InsightsPanel from '@/components/dashboard/InsightsPanel.vue'
import DataStatusIndicator from '@/components/dashboard/DataStatusIndicator.vue'
import PriceHeatmap from '@/components/comparison/PriceHeatmap.vue'
import ComparisonChart from '@/components/comparison/ComparisonChart.vue'
import ComparisonAnalytics from '@/components/comparison/ComparisonAnalytics.vue'
import PriceSpreadChart from '@/components/comparison/PriceSpreadChart.vue'
import DataModeToggle from '@/components/common/DataModeToggle.vue'
import WeatherPanel from '@/components/weather/WeatherPanel.vue'
import { Moon, Sun, ChevronLeft, ChevronRight, Download } from 'lucide-vue-next'
import { demandToProfile } from '@/types/settlement'

const demandStore = useDemandStore()
const reserveStore = useReserveStore()
const jepxStore = useJEPXStore()
const settlementStore = useSettlementStore()
const { isDark, toggleDarkMode } = useDarkMode()
const flags = useFeatureFlags()
const { exportCombinedCSV } = useExport()

// Enable keyboard navigation for date changes
useKeyboardNavigation(
  () => demandStore.prevDay(),
  () => demandStore.nextDay()
)

// Export combined report for Tokyo
function handleExport() {
  if (!demandStore.tokyoData || !jepxStore.tokyoData) {
    alert('Data not loaded yet. Please wait for the dashboard to load.')
    return
  }

  exportCombinedCSV(
    demandStore.tokyoData,
    jepxStore.tokyoData,
    settlementStore.data,
    'tokyo',
    demandStore.currentDate
  )
}

// Generate consumption profile from Tokyo demand data
const tokyoProfile = computed(() => {
  if (!demandStore.tokyoData) return []

  // Extract demand_mw and timestamps
  const demandMw = demandStore.tokyoData.series.map(s => s.demand_mw)
  const timestamps = demandStore.tokyoData.series.map(s => s.ts)

  return demandToProfile(demandMw, timestamps)
})

// Sparkline data for metrics cards
const tokyoSparklineData = computed(() => {
  return demandStore.tokyoData?.series.map(s => s.demand_mw) || []
})

const kansaiSparklineData = computed(() => {
  return demandStore.kansaiData?.series.map(s => s.demand_mw) || []
})

// Run settlement calculation
function runTokyoSettlement() {
  if (tokyoProfile.value.length === 0) return

  settlementStore.runSettlement(
    tokyoProfile.value,
    'tokyo',
    demandStore.currentDate,
    0.15 // 15% PV offset default
  )
}

onMounted(async () => {
  await demandStore.fetchAllDemandData()

  // Fetch reserve data only if feature enabled
  if (flags.isReserveEnabled) {
    reserveStore.fetchReserveData(demandStore.currentDate)
  }

  // Fetch JEPX prices only if feature enabled
  if (flags.isJEPXEnabled) {
    jepxStore.fetchBothAreas(demandStore.currentDate)
  }

  // Run settlement only if feature enabled
  if (flags.isSettlementEnabled) {
    runTokyoSettlement()
  }
})

// Refetch data when date changes
watch(() => demandStore.currentDate, (newDate) => {
  if (flags.isReserveEnabled) {
    reserveStore.fetchReserveData(newDate)
  }

  if (flags.isJEPXEnabled) {
    jepxStore.fetchBothAreas(newDate)
  }
})

// Refetch settlement when Tokyo data changes
watch(() => demandStore.tokyoData, () => {
  if (flags.isSettlementEnabled) {
    runTokyoSettlement()
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 p-8">
    <!-- Header -->
    <header class="mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-4xl font-bold text-gray-900 dark:text-white">
            Japan Energy Dashboard
          </h1>
          <p class="text-gray-600 dark:text-gray-400 mt-2">
            📊 Auto-updated daily at 00:30 JST with latest data
          </p>
        </div>
        <div class="flex items-center gap-3">
          <DataModeToggle />
          <button @click="toggleDarkMode" class="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700">
            <Moon v-if="!isDark" :size="24" />
            <Sun v-else :size="24" />
          </button>
        </div>
      </div>

      <!-- Date Navigation & Actions -->
      <div class="flex items-center justify-between mt-6 gap-4">
        <!-- Center: Date Navigation -->
        <div class="flex-1 flex items-center justify-center gap-4">
          <BaseButton size="sm" @click="demandStore.prevDay">
            <ChevronLeft :size="16" />
          </BaseButton>
          <div class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ demandStore.currentDate }}
          </div>
          <BaseButton size="sm" @click="demandStore.nextDay">
            <ChevronRight :size="16" />
          </BaseButton>
        </div>

        <!-- Right: Export Button -->
        <div class="flex items-center gap-3">
          <BaseButton size="sm" @click="handleExport" v-if="!demandStore.loading">
            <Download :size="16" class="mr-2" />
            <span>Export CSV</span>
          </BaseButton>
        </div>
      </div>
    </header>

    <!-- Summary Stats Bar -->
    <SummaryStatsBar v-if="!demandStore.loading" class="mb-8 -mx-8" />

    <!-- Loading State -->
    <div v-if="demandStore.loading" class="space-y-8">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <BaseCard>
          <ChartLoadingSkeleton />
        </BaseCard>
        <BaseCard>
          <ChartLoadingSkeleton />
        </BaseCard>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="demandStore.error">
      <BaseCard>
        <EmptyState
          type="error"
          :message="demandStore.error"
        />
      </BaseCard>
    </div>

    <!-- Content -->
    <div v-else class="space-y-8">
      <!-- Demand Row: Kansai + Tokyo -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Kansai Demand -->
        <BaseCard>
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Kansai Demand</h2>
          <ReserveBadge v-if="flags.isReserveEnabled" area="kansai" :data="reserveStore.reserveForArea('kansai')" />
        </div>
        <DemandChart
          title=""
          :data="demandStore.kansaiChartData"
          :prices="flags.isJEPXEnabled ? jepxStore.priceValues('kansai') : undefined"
        />
        <div v-if="demandStore.kansaiMetrics" class="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700 grid grid-cols-3 gap-3">
          <div>
            <div class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 mb-1">
              <span>Peak</span>
              <InfoIcon content="Highest demand value in the 24-hour period" />
            </div>
            <div class="flex items-center gap-1.5">
              <div class="text-lg font-bold text-gray-900 dark:text-white">{{ demandStore.kansaiMetrics.peak.toFixed(0) }}</div>
              <span class="text-xs text-gray-500">MW</span>
              <TrendIndicator
                v-if="demandStore.prevKansaiMetrics"
                :current="demandStore.kansaiMetrics.peak"
                :previous="demandStore.prevKansaiMetrics.peak"
                format="percentage"
                :size="12"
              />
            </div>
            <Sparkline
              v-if="kansaiSparklineData.length > 0"
              :data="kansaiSparklineData"
              :width="60"
              :height="16"
              color="rgb(16, 185, 129)"
              fill-color="rgb(16, 185, 129)"
              class="mt-1"
            />
          </div>
          <div>
            <div class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 mb-1">
              <span>Average</span>
              <InfoIcon content="Mean demand across all 24 hours" />
            </div>
            <div class="flex items-center gap-1.5">
              <div class="text-lg font-bold text-gray-900 dark:text-white">{{ demandStore.kansaiMetrics.average.toFixed(0) }}</div>
              <span class="text-xs text-gray-500">MW</span>
              <TrendIndicator
                v-if="demandStore.prevKansaiMetrics"
                :current="demandStore.kansaiMetrics.average"
                :previous="demandStore.prevKansaiMetrics.average"
                format="percentage"
                :size="12"
              />
            </div>
          </div>
          <div v-if="demandStore.kansaiMetrics.forecastAccuracy">
            <div class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 mb-1">
              <span>Accuracy</span>
              <InfoIcon content="How closely the forecast matches actual demand" />
            </div>
            <div class="flex items-center gap-1.5">
              <div class="text-lg font-bold text-gray-900 dark:text-white">{{ demandStore.kansaiMetrics.forecastAccuracy }}</div>
              <span class="text-xs text-gray-500">%</span>
            </div>
          </div>
        </div>
      </BaseCard>

        <!-- Tokyo Demand -->
        <BaseCard>
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Tokyo Demand (TEPCO)</h2>
          <ReserveBadge v-if="flags.isReserveEnabled" area="tokyo" :data="reserveStore.reserveForArea('tokyo')" />
        </div>
        <DemandChart
          title=""
          :data="demandStore.tokyoChartData"
          :prices="flags.isJEPXEnabled ? jepxStore.priceValues('tokyo') : undefined"
        />
        <div v-if="demandStore.tokyoMetrics" class="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700 grid grid-cols-3 gap-3">
          <div>
            <div class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 mb-1">
              <span>Peak</span>
              <InfoIcon content="Highest demand value in the 24-hour period" />
            </div>
            <div class="flex items-center gap-1.5">
              <div class="text-lg font-bold text-gray-900 dark:text-white">{{ demandStore.tokyoMetrics.peak.toFixed(0) }}</div>
              <span class="text-xs text-gray-500">MW</span>
              <TrendIndicator
                v-if="demandStore.prevTokyoMetrics"
                :current="demandStore.tokyoMetrics.peak"
                :previous="demandStore.prevTokyoMetrics.peak"
                format="percentage"
                :size="12"
              />
            </div>
            <Sparkline
              v-if="tokyoSparklineData.length > 0"
              :data="tokyoSparklineData"
              :width="60"
              :height="16"
              color="rgb(59, 130, 246)"
              fill-color="rgb(59, 130, 246)"
              class="mt-1"
            />
          </div>
          <div>
            <div class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 mb-1">
              <span>Average</span>
              <InfoIcon content="Mean demand across all 24 hours" />
            </div>
            <div class="flex items-center gap-1.5">
              <div class="text-lg font-bold text-gray-900 dark:text-white">{{ demandStore.tokyoMetrics.average.toFixed(0) }}</div>
              <span class="text-xs text-gray-500">MW</span>
              <TrendIndicator
                v-if="demandStore.prevTokyoMetrics"
                :current="demandStore.tokyoMetrics.average"
                :previous="demandStore.prevTokyoMetrics.average"
                format="percentage"
                :size="12"
              />
            </div>
          </div>
          <div v-if="demandStore.tokyoMetrics.forecastAccuracy">
            <div class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 mb-1">
              <span>Accuracy</span>
              <InfoIcon content="How closely the forecast matches actual demand" />
            </div>
            <div class="flex items-center gap-1.5">
              <div class="text-lg font-bold text-gray-900 dark:text-white">{{ demandStore.tokyoMetrics.forecastAccuracy }}</div>
              <span class="text-xs text-gray-500">%</span>
            </div>
          </div>
        </div>
      </BaseCard>
      </div>

      <!-- Solar Forecast Row: Kansai + Tokyo -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Kansai Weather -->
        <BaseCard>
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Solar Forecast - Kansai</h2>
          <WeatherPanel area="kansai" />
        </BaseCard>

        <!-- Tokyo Weather -->
        <BaseCard>
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Solar Forecast - Tokyo</h2>
          <WeatherPanel area="tokyo" />
        </BaseCard>
      </div>

      <!-- Price Spread Analysis -->
      <div v-if="flags.isJEPXEnabled" class="w-full">
        <BaseCard>
          <PriceSpreadChart />
        </BaseCard>
      </div>

      <!-- Insights & Settlement Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Key Insights -->
        <InsightsPanel />

        <!-- Settlement Cost Card -->
        <div v-if="flags.isSettlementEnabled">
          <CostCard />
        </div>
      </div>

      <!-- Regional Comparison & Analytics -->
      <div v-if="flags.isJEPXEnabled" class="space-y-6">
        <div class="border-t border-gray-200 dark:border-gray-700 pt-8">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-6">Regional Comparison & Analytics</h2>

          <!-- Price Heatmap -->
          <BaseCard class="mb-6">
            <PriceHeatmap />
          </BaseCard>

          <!-- Comparison Chart & Analytics Grid -->
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <!-- Comparison Chart (2 columns) -->
            <BaseCard class="lg:col-span-2">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Multi-Area Overlay</h3>
              <ComparisonChart />
            </BaseCard>

            <!-- Analytics (1 column) -->
            <BaseCard>
              <ComparisonAnalytics />
            </BaseCard>
          </div>
        </div>
      </div>

      <!-- Data Status Footer -->
      <div class="flex justify-center pt-6 border-t border-gray-200 dark:border-gray-700">
        <DataStatusIndicator />
      </div>
    </div>
  </div>
</template>
