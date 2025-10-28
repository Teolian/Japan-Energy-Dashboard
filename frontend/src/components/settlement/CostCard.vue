<script setup lang="ts">
import { computed } from 'vue'
import { useSettlementStore } from '@/stores/settlement'
import { Zap, DollarSign, Info } from 'lucide-vue-next'

const settlementStore = useSettlementStore()

// Format PV offset percentage
const formattedPVOffset = computed(() => {
  if (!settlementStore.assumptions) return '—'
  return `${Math.round(settlementStore.assumptions.pv_offset_pct * 100)}%`
})

// Format area name
const areaName = computed(() => {
  if (!settlementStore.assumptions) return '—'
  const area = settlementStore.assumptions.area
  return area.charAt(0).toUpperCase() + area.slice(1)
})
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
        Settlement Cost
      </h2>
      <div class="text-xs text-gray-500 dark:text-gray-400">
        Settlement-lite
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="settlementStore.loading" class="text-center py-8">
      <div class="text-gray-600 dark:text-gray-400">Calculating...</div>
    </div>

    <!-- Error State -->
    <div v-else-if="settlementStore.error" class="text-center py-8">
      <div class="text-red-600 dark:text-red-400 text-sm">{{ settlementStore.error }}</div>
    </div>

    <!-- Content -->
    <div v-else-if="settlementStore.totals" class="space-y-6">
      <!-- Totals Grid -->
      <div class="grid grid-cols-2 gap-6">
        <!-- Total Consumption -->
        <div class="space-y-2">
          <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
            <Zap :size="16" class="text-yellow-500" />
            <span>Total Consumption</span>
          </div>
          <div class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ settlementStore.formattedTotalKWh }}
            <span class="text-sm font-normal text-gray-500 dark:text-gray-400 ml-1">kWh</span>
          </div>
        </div>

        <!-- Total Cost -->
        <div class="space-y-2">
          <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
            <DollarSign :size="16" class="text-green-500" />
            <span>Total Cost</span>
          </div>
          <div class="text-2xl font-bold text-gray-900 dark:text-white">
            ¥{{ settlementStore.formattedTotalCost }}
            <span class="text-sm font-normal text-gray-500 dark:text-gray-400 ml-1">JPY</span>
          </div>
        </div>
      </div>

      <!-- Divider -->
      <div class="border-t border-gray-200 dark:border-gray-700"></div>

      <!-- Assumptions -->
      <div class="space-y-3">
        <div class="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-gray-300">
          <Info :size="16" />
          <span>Assumptions</span>
        </div>

        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span class="text-gray-600 dark:text-gray-400">Area:</span>
            <span class="ml-2 font-medium text-gray-900 dark:text-white">{{ areaName }}</span>
          </div>
          <div>
            <span class="text-gray-600 dark:text-gray-400">PV Offset:</span>
            <span class="ml-2 font-medium text-gray-900 dark:text-white">{{ formattedPVOffset }}</span>
          </div>
        </div>
      </div>

      <!-- Source Attribution -->
      <div v-if="settlementStore.source" class="pt-3 border-t border-gray-200 dark:border-gray-700">
        <div class="text-xs text-gray-500 dark:text-gray-400">
          Price data: {{ settlementStore.source.name }}
        </div>
      </div>
    </div>

    <!-- No Data State -->
    <div v-else class="text-center py-8">
      <div class="text-gray-600 dark:text-gray-400">No settlement data available</div>
    </div>
  </div>
</template>
