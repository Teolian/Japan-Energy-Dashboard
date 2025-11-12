<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTradingStore } from '@/stores/trading'
import { Clock, ArrowRight, TrendingDown, Leaf, CheckCircle2 } from 'lucide-vue-next'

const { t } = useI18n()
const tradingStore = useTradingStore()

// Computed: Total potential savings
const totalSavings = computed(() => {
  return tradingStore.priorityLoadShifts.reduce((sum, rec) => sum + rec.savings, 0)
})

const totalCarbonReduction = computed(() => {
  return tradingStore.priorityLoadShifts.reduce((sum, rec) => sum + rec.carbonReduction, 0)
})

const getPriorityColor = (priority: string) => {
  const colors = {
    high: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300',
    medium: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300',
    low: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300'
  }
  return colors[priority as keyof typeof colors] || colors.medium
}

const getFeasibilityColor = (score: number) => {
  if (score >= 80) return 'text-green-600 dark:text-green-400'
  if (score >= 60) return 'text-yellow-600 dark:text-yellow-400'
  return 'text-gray-600 dark:text-gray-400'
}

const getFeasibilityBg = (score: number) => {
  if (score >= 80) return 'bg-green-500'
  if (score >= 60) return 'bg-yellow-500'
  return 'bg-gray-400'
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header with Summary Stats -->
    <div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2 mb-4">
        <Clock :size="24" class="text-energy-cyan" />
        {{ t('loadShift.title') }}
      </h3>

      <!-- Summary Cards -->
      <div class="grid grid-cols-3 gap-4 mb-6">
        <div class="bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 p-4 rounded-lg border border-green-200 dark:border-green-800">
          <div class="flex items-center gap-2 mb-2">
            <TrendingDown :size="18" class="text-green-600 dark:text-green-400" />
            <span class="text-xs text-gray-600 dark:text-gray-400">{{ t('loadShift.savings') }}</span>
          </div>
          <div class="text-2xl font-bold text-green-600 dark:text-green-400">
            ¥{{ totalSavings.toLocaleString() }}
          </div>
          <div class="text-xs text-gray-600 dark:text-gray-400 mt-1">/day potential</div>
        </div>

        <div class="bg-gradient-to-br from-blue-50 to-cyan-50 dark:from-blue-900/20 dark:to-cyan-900/20 p-4 rounded-lg border border-blue-200 dark:border-blue-800">
          <div class="flex items-center gap-2 mb-2">
            <Leaf :size="18" class="text-blue-600 dark:text-blue-400" />
            <span class="text-xs text-gray-600 dark:text-gray-400">{{ t('loadShift.carbonReduction') }}</span>
          </div>
          <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">
            {{ totalCarbonReduction.toLocaleString() }}
          </div>
          <div class="text-xs text-gray-600 dark:text-gray-400 mt-1">kg CO₂/day</div>
        </div>

        <div class="bg-gradient-to-br from-purple-50 to-pink-50 dark:from-purple-900/20 dark:to-pink-900/20 p-4 rounded-lg border border-purple-200 dark:border-purple-800">
          <div class="flex items-center gap-2 mb-2">
            <CheckCircle2 :size="18" class="text-purple-600 dark:text-purple-400" />
            <span class="text-xs text-gray-600 dark:text-gray-400">{{ t('loadShift.recommendations') }}</span>
          </div>
          <div class="text-2xl font-bold text-purple-600 dark:text-purple-400">
            {{ tradingStore.priorityLoadShifts.length }}
          </div>
          <div class="text-xs text-gray-600 dark:text-gray-400 mt-1">high feasibility</div>
        </div>
      </div>
    </div>

    <!-- Recommendations List -->
    <div class="space-y-3">
      <div
        v-for="rec in tradingStore.priorityLoadShifts"
        :key="rec.id"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 hover:shadow-lg transition-shadow"
      >
        <!-- Header -->
        <div class="flex items-start justify-between mb-3">
          <div class="flex items-center gap-3 flex-1">
            <!-- Time Shift Visualization -->
            <div class="flex items-center gap-2 text-sm font-mono">
              <div class="px-3 py-2 bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300 rounded font-bold">
                {{ rec.fromHour.toString().padStart(2, '0') }}:00
              </div>
              <ArrowRight :size="20" class="text-gray-400" />
              <div class="px-3 py-2 bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 rounded font-bold">
                {{ rec.toHour.toString().padStart(2, '0') }}:00
              </div>
            </div>

            <!-- Priority Badge -->
            <span :class="['text-xs px-2 py-1 rounded-full font-medium uppercase', getPriorityColor(rec.priority)]">
              {{ rec.priority }}
            </span>
          </div>

          <!-- Savings Highlight -->
          <div class="text-right">
            <div class="text-xs text-gray-600 dark:text-gray-400">Savings</div>
            <div class="text-xl font-bold text-green-600 dark:text-green-400">
              ¥{{ rec.savings.toLocaleString() }}
            </div>
          </div>
        </div>

        <!-- Details Grid -->
        <div class="grid grid-cols-3 gap-3 mb-3 text-sm">
          <div class="bg-gray-50 dark:bg-gray-700/50 p-2 rounded">
            <div class="text-xs text-gray-600 dark:text-gray-400">{{ t('loadShift.amount') }}</div>
            <div class="font-semibold">{{ rec.amountMW.toLocaleString() }} MW</div>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700/50 p-2 rounded">
            <div class="text-xs text-gray-600 dark:text-gray-400">CO₂ Reduction</div>
            <div class="font-semibold">{{ rec.carbonReduction }} kg</div>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700/50 p-2 rounded">
            <div class="text-xs text-gray-600 dark:text-gray-400">{{ t('loadShift.feasibility') }}</div>
            <div :class="['font-semibold', getFeasibilityColor(rec.feasibility)]">
              {{ rec.feasibility }}%
            </div>
          </div>
        </div>

        <!-- Feasibility Bar -->
        <div class="mb-3">
          <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <div
              :class="['h-full transition-all', getFeasibilityBg(rec.feasibility)]"
              :style="{ width: `${rec.feasibility}%` }"
            ></div>
          </div>
        </div>

        <!-- Reason -->
        <div class="text-sm text-gray-600 dark:text-gray-400">
          {{ rec.reason }}
        </div>
      </div>

      <!-- Empty State -->
      <div
        v-if="tradingStore.priorityLoadShifts.length === 0"
        class="text-center py-8 text-gray-500 dark:text-gray-400"
      >
        <Clock :size="48" class="mx-auto mb-2 opacity-30" />
        <p>No load shift opportunities available</p>
        <p class="text-sm mt-1">Generate recommendations by analyzing demand and price data</p>
      </div>
    </div>

    <!-- Load Profile Visualization (simplified) -->
    <div v-if="tradingStore.loadProfile.length > 0" class="bg-gradient-to-br from-gray-50 to-blue-50 dark:from-gray-800 dark:to-blue-900/20 p-5 rounded-xl border border-gray-200 dark:border-gray-700">
      <h4 class="font-semibold text-gray-900 dark:text-white mb-3">{{ t('loadShift.currentLoad') }}</h4>

      <!-- Simple bar chart -->
      <div class="grid grid-cols-24 gap-0.5 h-32">
        <div
          v-for="profile in tradingStore.loadProfile"
          :key="profile.hour"
          class="relative group"
        >
          <div
            class="absolute bottom-0 w-full bg-blue-400 dark:bg-blue-600 rounded-t transition-all hover:bg-blue-500"
            :style="{
              height: `${(profile.currentLoad / Math.max(...tradingStore.loadProfile.map(p => p.currentLoad))) * 100}%`
            }"
          ></div>

          <!-- Tooltip on hover -->
          <div class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-2 py-1 bg-gray-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none z-10">
            {{ profile.hour }}:00<br/>
            {{ Math.round(profile.currentLoad).toLocaleString() }} MW<br/>
            ¥{{ profile.price.toFixed(2) }}/kWh
          </div>
        </div>
      </div>

      <!-- Hour labels -->
      <div class="grid grid-cols-24 gap-0.5 mt-1 text-xs text-gray-600 dark:text-gray-400 text-center">
        <div v-for="h in 24" :key="h">{{ (h - 1).toString().padStart(2, '0') }}</div>
      </div>
    </div>
  </div>
</template>
