<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import SectionCard from '@/components/common/SectionCard.vue'
import PriceSpreadChart from '@/components/comparison/PriceSpreadChart.vue'
import DuckCurveAnalysis from '@/components/dashboard/DuckCurveAnalysis.vue'
import { TrendingUp } from 'lucide-vue-next'

type Tab = 'spread' | 'duck'

const { t } = useI18n()
const activeTab = ref<Tab>('spread')

const tabs = [
  { id: 'spread' as Tab, label: 'Price Spread', description: 'Tokyo vs Kansai price difference' },
  { id: 'duck' as Tab, label: 'Duck Curve', description: 'Solar impact on net demand' }
]
</script>

<template>
  <SectionCard padding="md">
    <template #header>
      <div class="flex items-center justify-between w-full">
        <div class="flex items-center gap-2">
          <TrendingUp :size="20" class="text-orange-600 dark:text-orange-400" />
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('dashboard.prices') }}
          </h2>
        </div>

        <!-- Tabs -->
        <div class="flex gap-1 bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'px-4 py-2 rounded-md text-sm font-medium transition-colors',
              activeTab === tab.id
                ? 'bg-white dark:bg-gray-800 text-gray-900 dark:text-white shadow-sm'
                : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
            ]"
          >
            {{ tab.label }}
          </button>
        </div>
      </div>
    </template>

    <!-- Tab Content -->
    <div class="mt-4">
      <!-- Price Spread Tab -->
      <div v-if="activeTab === 'spread'">
        <PriceSpreadChart />
      </div>

      <!-- Duck Curve Tab -->
      <div v-if="activeTab === 'duck'">
        <DuckCurveAnalysis />
      </div>
    </div>
  </SectionCard>
</template>
