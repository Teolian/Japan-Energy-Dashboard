<script setup lang="ts">
import { computed } from 'vue'
import { useGenerationStore } from '@/stores/generation'
import { TrendingDown, TrendingUp, AlertCircle, CheckCircle } from 'lucide-vue-next'

const generationStore = useGenerationStore()

// Metrics
const metrics = computed(() => generationStore.tokyoMetrics)
const chartData = computed(() => generationStore.tokyoChartData)

// Carbon intensity classification
const carbonStatus = computed(() => {
  if (!metrics.value) return null

  const carbon = metrics.value.carbonIntensity
  return {
    value: Math.round(carbon.value),
    level: carbon.level,
    label: carbon.label,
    color: carbon.color
  }
})

// Find top 3 greenest hours
const greenHours = computed(() => {
  if (chartData.value.length === 0) return []

  return chartData.value
    .map((point, idx) => ({ ...point, idx }))
    .sort((a, b) => a.carbon_gco2_kwh - b.carbon_gco2_kwh)
    .slice(0, 3)
    .map(point => ({
      time: point.time,
      carbon: Math.round(point.carbon_gco2_kwh),
      renewable: point.renewable_pct.toFixed(1)
    }))
})

// Find top 3 dirtiest hours
const dirtyHours = computed(() => {
  if (chartData.value.length === 0) return []

  return chartData.value
    .map((point, idx) => ({ ...point, idx }))
    .sort((a, b) => b.carbon_gco2_kwh - a.carbon_gco2_kwh)
    .slice(0, 3)
    .map(point => ({
      time: point.time,
      carbon: Math.round(point.carbon_gco2_kwh),
      renewable: point.renewable_pct.toFixed(1)
    }))
})

// Gauge rotation (0-180 degrees)
const gaugeRotation = computed(() => {
  if (!carbonStatus.value) return 0

  // Map 0-500 gCO2 to 0-180 degrees
  const value = Math.min(carbonStatus.value.value, 500)
  return (value / 500) * 180
})
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div>
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
        Carbon Intensity Tracker
      </h3>
      <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
        For ESG reporting (RE100, CDP, Scope 2)
      </p>
    </div>

    <!-- Gauge -->
    <div v-if="carbonStatus" class="relative p-6 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-800 dark:to-gray-900 rounded-lg border border-gray-200 dark:border-gray-700">
      <div class="text-center">
        <!-- Gauge SVG -->
        <div class="relative inline-block">
          <svg width="200" height="120" viewBox="0 0 200 120">
            <!-- Background arc -->
            <path
              d="M 20 100 A 80 80 0 0 1 180 100"
              fill="none"
              stroke="#e5e7eb"
              stroke-width="12"
              class="dark:stroke-gray-700"
            />

            <!-- Colored segments -->
            <path
              d="M 20 100 A 80 80 0 0 1 65 30"
              fill="none"
              stroke="#10b981"
              stroke-width="12"
            />
            <path
              d="M 65 30 A 80 80 0 0 1 100 20"
              fill="none"
              stroke="#3b82f6"
              stroke-width="12"
            />
            <path
              d="M 100 20 A 80 80 0 0 1 135 30"
              fill="none"
              stroke="#f59e0b"
              stroke-width="12"
            />
            <path
              d="M 135 30 A 80 80 0 0 1 180 100"
              fill="none"
              stroke="#ef4444"
              stroke-width="12"
            />

            <!-- Needle -->
            <line
              x1="100"
              y1="100"
              x2="100"
              y2="30"
              stroke="#1f2937"
              stroke-width="3"
              stroke-linecap="round"
              :transform="`rotate(${gaugeRotation - 90}, 100, 100)`"
              class="dark:stroke-gray-100"
            />
            <circle cx="100" cy="100" r="6" fill="#1f2937" class="dark:fill-gray-100" />
          </svg>

          <!-- Center value -->
          <div class="absolute inset-x-0 bottom-2 text-center">
            <div class="text-3xl font-bold" :class="carbonStatus.color">
              {{ carbonStatus.value }}
            </div>
            <div class="text-xs text-gray-500 dark:text-gray-400">
              gCO₂/kWh
            </div>
          </div>
        </div>

        <!-- Status label -->
        <div class="mt-4">
          <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium" :class="[
            carbonStatus.level === 'low' ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400' :
            carbonStatus.level === 'medium' ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400' :
            carbonStatus.level === 'high' ? 'bg-orange-100 text-orange-800 dark:bg-orange-900/20 dark:text-orange-400' :
            'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400'
          ]">
            {{ carbonStatus.label }}
          </span>
        </div>
      </div>

      <!-- Scale labels -->
      <div class="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-2 px-4">
        <span>0</span>
        <span>200</span>
        <span>300</span>
        <span>500+</span>
      </div>
    </div>

    <!-- Green Hours (Load Shifting Recommendations) -->
    <div v-if="greenHours.length > 0" class="p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
      <div class="flex items-start gap-2 mb-3">
        <TrendingDown :size="18" class="text-green-600 dark:text-green-400 mt-0.5" />
        <div>
          <div class="font-medium text-green-900 dark:text-green-100 text-sm">
            Cleanest Hours (Load Shifting)
          </div>
          <div class="text-xs text-green-700 dark:text-green-300 mt-0.5">
            Schedule workloads during these hours to minimize carbon
          </div>
        </div>
      </div>

      <div class="space-y-2">
        <div
          v-for="(hour, idx) in greenHours"
          :key="idx"
          class="flex items-center justify-between p-2 bg-white dark:bg-gray-800 rounded border border-green-200 dark:border-green-700"
        >
          <div class="flex items-center gap-2">
            <CheckCircle :size="16" class="text-green-600 dark:text-green-400" />
            <span class="font-mono text-sm font-medium text-gray-900 dark:text-white">
              {{ hour.time }}
            </span>
          </div>
          <div class="text-right">
            <div class="text-sm font-bold text-green-600 dark:text-green-400">
              {{ hour.carbon }} gCO₂
            </div>
            <div class="text-xs text-gray-500">{{ hour.renewable }}% renewable</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Dirty Hours (Avoid) -->
    <div v-if="dirtyHours.length > 0" class="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
      <div class="flex items-start gap-2 mb-3">
        <TrendingUp :size="18" class="text-red-600 dark:text-red-400 mt-0.5" />
        <div>
          <div class="font-medium text-red-900 dark:text-red-100 text-sm">
            High Carbon Hours (Avoid)
          </div>
          <div class="text-xs text-red-700 dark:text-red-300 mt-0.5">
            Shift non-critical workloads away from these hours
          </div>
        </div>
      </div>

      <div class="space-y-2">
        <div
          v-for="(hour, idx) in dirtyHours"
          :key="idx"
          class="flex items-center justify-between p-2 bg-white dark:bg-gray-800 rounded border border-red-200 dark:border-red-700"
        >
          <div class="flex items-center gap-2">
            <AlertCircle :size="16" class="text-red-600 dark:text-red-400" />
            <span class="font-mono text-sm font-medium text-gray-900 dark:text-white">
              {{ hour.time }}
            </span>
          </div>
          <div class="text-right">
            <div class="text-sm font-bold text-red-600 dark:text-red-400">
              {{ hour.carbon }} gCO₂
            </div>
            <div class="text-xs text-gray-500">{{ hour.renewable }}% renewable</div>
          </div>
        </div>
      </div>
    </div>

    <!-- ESG Reporting Note -->
    <div class="p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg text-xs">
      <div class="font-medium text-blue-900 dark:text-blue-100 mb-1">
        ESG Reporting (RE100, CDP)
      </div>
      <div class="text-blue-700 dark:text-blue-300">
        Daily average carbon intensity can be used for Scope 2 emissions calculations.
        Shift loads to green hours to reduce your carbon footprint.
      </div>
    </div>
  </div>
</template>
