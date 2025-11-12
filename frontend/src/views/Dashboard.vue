<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useDemandStore } from '@/stores/demand'
import { useReserveStore } from '@/stores/reserve'
import { useJEPXStore } from '@/stores/jepx'
import { useGenerationStore } from '@/stores/generation'
import { useWeatherStore } from '@/stores/weather'
import { useSettlementStore } from '@/stores/settlement'
import { useDarkMode } from '@/composables/useDarkMode'
import { useFeatureFlags } from '@/composables/useFeatureFlags'
import { useKeyboardNavigation } from '@/composables/useKeyboardNavigation'

// Components
import BaseButton from '@/components/common/BaseButton.vue'
import DataModeToggle from '@/components/common/DataModeToggle.vue'
import SummaryStatsBar from '@/components/dashboard/SummaryStatsBar.vue'
import InsightsPanel from '@/components/dashboard/InsightsPanel.vue'
import DataStatusIndicator from '@/components/dashboard/DataStatusIndicator.vue'
import DemandCard from '@/components/demand/DemandCard.vue'
import MarketAnalysisSection from '@/components/market/MarketAnalysisSection.vue'
import GenerationSection from '@/components/generation/GenerationSection.vue'
import RegionalComparisonSection from '@/components/comparison/RegionalComparisonSection.vue'
import CostCard from '@/components/settlement/CostCard.vue'
import ChartLoadingSkeleton from '@/components/common/ChartLoadingSkeleton.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import BaseCard from '@/components/common/BaseCard.vue'

import { Moon, Sun, ChevronLeft, ChevronRight } from 'lucide-vue-next'

// Stores
const demandStore = useDemandStore()
const reserveStore = useReserveStore()
const jepxStore = useJEPXStore()
const generationStore = useGenerationStore()
const weatherStore = useWeatherStore()
const settlementStore = useSettlementStore()

// Composables
const { isDark, toggleDarkMode } = useDarkMode()
const flags = useFeatureFlags()

// Keyboard navigation for date changes
useKeyboardNavigation(
  () => demandStore.prevDay(),
  () => demandStore.nextDay()
)

// Initial data fetch
onMounted(async () => {
  await demandStore.fetchAllDemandData()
  weatherStore.fetchBothAreas(demandStore.currentDate)

  if (flags.isReserveEnabled) {
    reserveStore.fetchReserveData(demandStore.currentDate)
  }

  if (flags.isJEPXEnabled) {
    await jepxStore.fetchBothAreas(demandStore.currentDate)
    generationStore.fetchTokyo(demandStore.currentDate)
  }

  if (flags.isSettlementEnabled && demandStore.tokyoData) {
    const profile = demandStore.tokyoData.series.map(s => ({
      ts: s.ts,
      kwh: s.demand_mw * 1000 // Convert MW to kW
    }))
    settlementStore.runSettlement(profile, 'tokyo', demandStore.currentDate, 0.15)
  }
})

// Watch for date changes
watch(() => demandStore.currentDate, async (newDate) => {
  weatherStore.fetchBothAreas(newDate)

  if (flags.isReserveEnabled) {
    reserveStore.fetchReserveData(newDate)
  }

  if (flags.isJEPXEnabled) {
    await jepxStore.fetchBothAreas(newDate)
    generationStore.fetchTokyo(newDate)
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
            Auto-updated daily at 00:30 JST with latest data
          </p>
        </div>
        <div class="flex items-center gap-3">
          <DataModeToggle />
          <button
            @click="toggleDarkMode"
            class="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
          >
            <Moon v-if="!isDark" :size="24" />
            <Sun v-else :size="24" />
          </button>
        </div>
      </div>

      <!-- Date Navigation -->
      <div class="flex items-center justify-center gap-4 mt-6">
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
    </header>

    <!-- Summary Stats Bar -->
    <SummaryStatsBar v-if="!demandStore.loading" class="mb-8 -mx-8" />

    <!-- Loading State -->
    <div v-if="demandStore.loading" class="space-y-8">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <BaseCard><ChartLoadingSkeleton /></BaseCard>
        <BaseCard><ChartLoadingSkeleton /></BaseCard>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="demandStore.error">
      <BaseCard>
        <EmptyState type="error" :message="demandStore.error" />
      </BaseCard>
    </div>

    <!-- Main Content -->
    <div v-else class="space-y-8">
      <!-- Level 1: Core Demand Data (Tokyo + Kansai) -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <DemandCard area="tokyo" />
        <DemandCard area="kansai" />
      </div>

      <!-- Level 2: Market Analysis (JEPX Prices) -->
      <div v-if="flags.isJEPXEnabled">
        <MarketAnalysisSection />
      </div>

      <!-- Level 3: Generation Mix & Carbon -->
      <div v-if="flags.isJEPXEnabled">
        <GenerationSection />
      </div>

      <!-- Level 4: Insights & Settlement -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <InsightsPanel />
        <CostCard v-if="flags.isSettlementEnabled" />
      </div>

      <!-- Level 5: Regional Comparison (Collapsible) -->
      <RegionalComparisonSection v-if="flags.isJEPXEnabled" />

      <!-- Data Status Footer -->
      <div class="flex justify-center pt-6 border-t border-gray-200 dark:border-gray-700">
        <DataStatusIndicator />
      </div>
    </div>
  </div>
</template>
