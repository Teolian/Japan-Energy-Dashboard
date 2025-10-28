<script setup lang="ts">
import { computed } from 'vue'
import { useJEPXStore } from '@/stores/jepx'
import { TrendingDown, TrendingUp } from 'lucide-vue-next'

const jepxStore = useJEPXStore()

interface HeatmapCell {
  hour: number
  price: number
  color: string
  intensity: 'low' | 'medium' | 'high' | 'peak'
}

// Generate heatmap data for area
function generateHeatmap(area: 'tokyo' | 'kansai'): HeatmapCell[] {
  const prices = jepxStore.priceValues(area)
  if (prices.length === 0) return []

  const min = Math.min(...prices)
  const max = Math.max(...prices)
  const range = max - min

  return prices.map((price, hour) => {
    // Normalize price to 0-1
    const normalized = range > 0 ? (price - min) / range : 0.5

    // Determine intensity and color
    let intensity: HeatmapCell['intensity']
    let color: string

    if (normalized < 0.25) {
      intensity = 'low'
      color = 'bg-green-500 dark:bg-green-600'
    } else if (normalized < 0.5) {
      intensity = 'medium'
      color = 'bg-yellow-500 dark:bg-yellow-600'
    } else if (normalized < 0.75) {
      intensity = 'high'
      color = 'bg-orange-500 dark:bg-orange-600'
    } else {
      intensity = 'peak'
      color = 'bg-red-500 dark:bg-red-600'
    }

    return { hour, price, color, intensity }
  })
}

const tokyoHeatmap = computed(() => generateHeatmap('tokyo'))
const kansaiHeatmap = computed(() => generateHeatmap('kansai'))

// Find cheapest and most expensive hours
const tokyoCheapest = computed(() => {
  if (tokyoHeatmap.value.length === 0) return null
  return tokyoHeatmap.value.reduce((min, cell) => cell.price < min.price ? cell : min)
})

const tokyoMostExpensive = computed(() => {
  if (tokyoHeatmap.value.length === 0) return null
  return tokyoHeatmap.value.reduce((max, cell) => cell.price > max.price ? cell : max)
})

const kansaiCheapest = computed(() => {
  if (kansaiHeatmap.value.length === 0) return null
  return kansaiHeatmap.value.reduce((min, cell) => cell.price < min.price ? cell : min)
})

const kansaiMostExpensive = computed(() => {
  if (kansaiHeatmap.value.length === 0) return null
  return kansaiHeatmap.value.reduce((max, cell) => cell.price > max.price ? cell : max)
})
</script>

<template>
  <div class="space-y-3">
    <div class="flex items-center justify-between">
      <h3 class="text-base font-semibold text-gray-900 dark:text-white">Spot Price Patterns</h3>
      <div class="flex items-center gap-3 text-xs">
        <div class="flex items-center gap-1">
          <div class="w-2 h-2 rounded bg-green-500"></div>
          <span class="text-gray-600 dark:text-gray-400">Low</span>
        </div>
        <div class="flex items-center gap-1">
          <div class="w-2 h-2 rounded bg-yellow-500"></div>
          <span class="text-gray-600 dark:text-gray-400">Mid</span>
        </div>
        <div class="flex items-center gap-1">
          <div class="w-2 h-2 rounded bg-orange-500"></div>
          <span class="text-gray-600 dark:text-gray-400">High</span>
        </div>
        <div class="flex items-center gap-1">
          <div class="w-2 h-2 rounded bg-red-500"></div>
          <span class="text-gray-600 dark:text-gray-400">Peak</span>
        </div>
      </div>
    </div>

    <!-- Tokyo Heatmap -->
    <div>
      <div class="flex items-center justify-between mb-2">
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Tokyo</span>
        <div class="flex items-center gap-3 text-xs">
          <div v-if="tokyoCheapest" class="flex items-center gap-1 text-green-600 dark:text-green-400">
            <TrendingDown :size="14" />
            <span>Best: {{ tokyoCheapest.hour }}:00 (¥{{ tokyoCheapest.price.toFixed(1) }})</span>
          </div>
          <div v-if="tokyoMostExpensive" class="flex items-center gap-1 text-red-600 dark:text-red-400">
            <TrendingUp :size="14" />
            <span>Peak: {{ tokyoMostExpensive.hour }}:00 (¥{{ tokyoMostExpensive.price.toFixed(1) }})</span>
          </div>
        </div>
      </div>
      <div v-if="tokyoHeatmap.length > 0" class="flex gap-1">
        <div
          v-for="cell in tokyoHeatmap"
          :key="`tokyo-${cell.hour}`"
          :class="[cell.color, 'h-8 flex-1 min-w-0 rounded cursor-pointer hover:opacity-80 transition-opacity']"
          :title="`${cell.hour}:00 - ¥${cell.price.toFixed(1)}/kWh`"
        />
      </div>
      <div v-else class="h-8 flex items-center justify-center text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-800 rounded">
        No data
      </div>
    </div>

    <!-- Kansai Heatmap -->
    <div>
      <div class="flex items-center justify-between mb-2">
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Kansai</span>
        <div class="flex items-center gap-3 text-xs">
          <div v-if="kansaiCheapest" class="flex items-center gap-1 text-green-600 dark:text-green-400">
            <TrendingDown :size="14" />
            <span>Best: {{ kansaiCheapest.hour }}:00 (¥{{ kansaiCheapest.price.toFixed(1) }})</span>
          </div>
          <div v-if="kansaiMostExpensive" class="flex items-center gap-1 text-red-600 dark:text-red-400">
            <TrendingUp :size="14" />
            <span>Peak: {{ kansaiMostExpensive.hour }}:00 (¥{{ kansaiMostExpensive.price.toFixed(1) }})</span>
          </div>
        </div>
      </div>
      <div v-if="kansaiHeatmap.length > 0" class="flex gap-1">
        <div
          v-for="cell in kansaiHeatmap"
          :key="`kansai-${cell.hour}`"
          :class="[cell.color, 'h-8 flex-1 min-w-0 rounded cursor-pointer hover:opacity-80 transition-opacity']"
          :title="`${cell.hour}:00 - ¥${cell.price.toFixed(1)}/kWh`"
        />
      </div>
      <div v-else class="h-8 flex items-center justify-center text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-800 rounded">
        No data
      </div>
    </div>

    <!-- Hour labels -->
    <div v-if="tokyoHeatmap.length > 0 || kansaiHeatmap.length > 0" class="flex gap-1 text-[10px] text-gray-500 dark:text-gray-400 text-center mt-1">
      <span v-for="hour in 24" :key="hour" class="flex-1 min-w-0">{{ hour - 1 }}</span>
    </div>
  </div>
</template>
