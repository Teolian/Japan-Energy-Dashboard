<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { Sun, Cloud, Zap, TrendingUp } from 'lucide-vue-next'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
  type ChartOptions
} from 'chart.js'
import { useWeatherStore } from '@/stores/weather'
import { useDemandStore } from '@/stores/demand'
import type { Area } from '@/types/weather'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler)

const props = defineProps<{
  area: Area
}>()

const weatherStore = useWeatherStore()
const demandStore = useDemandStore()

const forecast = computed(() => weatherStore.forecastForArea(props.area))

// Chart data combining solar radiation and demand
const chartData = computed(() => {
  const labels = Array.from({ length: 24 }, (_, i) => `${i.toString().padStart(2, '0')}:00`)

  const radiationData = weatherStore.radiationValues(props.area)
  const demandData = props.area === 'tokyo'
    ? (demandStore.tokyoData?.series.map(s => s.demand_mw) || [])
    : (demandStore.kansaiData?.series.map(s => s.demand_mw) || [])

  return {
    labels,
    datasets: [
      {
        label: 'Solar Radiation (W/m²)',
        data: radiationData,
        borderColor: 'rgb(251, 191, 36)', // Amber for sun
        backgroundColor: 'rgba(251, 191, 36, 0.1)',
        borderWidth: 2,
        tension: 0.4,
        yAxisID: 'y',
        fill: true,
        pointRadius: 2,
        pointHoverRadius: 4
      },
      {
        label: 'Demand (MW)',
        data: demandData,
        borderColor: 'rgb(59, 130, 246)', // Blue for demand
        backgroundColor: 'rgba(59, 130, 246, 0.05)',
        borderWidth: 2,
        borderDash: [5, 5],
        tension: 0.4,
        yAxisID: 'y1',
        pointRadius: 2,
        pointHoverRadius: 4
      }
    ]
  }
})

const chartOptions = computed<ChartOptions<'line'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    mode: 'index',
    intersect: false
  },
  plugins: {
    legend: {
      display: true,
      position: 'top',
      labels: {
        usePointStyle: true,
        padding: 10,
        font: { size: 11 }
      }
    },
    tooltip: {
      enabled: true,
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      padding: 10,
      cornerRadius: 6
    }
  },
  scales: {
    y: {
      type: 'linear',
      display: true,
      position: 'left',
      title: {
        display: true,
        text: 'Solar Radiation (W/m²)',
        font: { size: 11 }
      },
      ticks: {
        font: { size: 10 }
      }
    },
    y1: {
      type: 'linear',
      display: true,
      position: 'right',
      title: {
        display: true,
        text: 'Demand (MW)',
        font: { size: 11 }
      },
      ticks: {
        font: { size: 10 }
      },
      grid: {
        drawOnChartArea: false
      }
    }
  }
}))

// Fetch forecast on mount
onMounted(() => {
  weatherStore.fetchForecast(props.area, demandStore.currentDate)
})
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
        <Sun :size="20" class="text-amber-500" />
        Solar Forecast - {{ forecast?.location }}
      </h3>
      <button
        @click="weatherStore.toggleMockData"
        class="text-xs px-2 py-1 rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600"
      >
        {{ weatherStore.useMockData ? 'Mock' : 'Live' }} Data
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="weatherStore.loading" class="text-center py-8 text-gray-500 dark:text-gray-400">
      <div class="animate-spin inline-block w-6 h-6 border-2 border-current border-t-transparent rounded-full"></div>
      <p class="mt-2 text-sm">Loading solar forecast...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="weatherStore.error" class="text-center py-8 text-red-600 dark:text-red-400">
      <p class="text-sm">{{ weatherStore.error }}</p>
    </div>

    <!-- Forecast Data -->
    <div v-else-if="forecast">
      <!-- Key Metrics -->
      <div class="grid grid-cols-4 gap-3">
        <!-- Peak Hour -->
        <div class="p-3 rounded-lg bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800">
          <div class="flex items-center gap-2 mb-1">
            <TrendingUp :size="14" class="text-amber-600 dark:text-amber-400" />
            <span class="text-xs font-medium text-amber-900 dark:text-amber-200">Peak Hour</span>
          </div>
          <div class="text-lg font-bold text-amber-900 dark:text-amber-100">
            {{ forecast.peak_radiation_hour }}:00
          </div>
          <div class="text-xs text-amber-700 dark:text-amber-300">
            {{ forecast.data[forecast.peak_radiation_hour]?.ghi.toFixed(0) }} W/m²
          </div>
        </div>

        <!-- Daily Total -->
        <div class="p-3 rounded-lg bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800">
          <div class="flex items-center gap-2 mb-1">
            <Sun :size="14" class="text-yellow-600 dark:text-yellow-400" />
            <span class="text-xs font-medium text-yellow-900 dark:text-yellow-200">Daily Total</span>
          </div>
          <div class="text-lg font-bold text-yellow-900 dark:text-yellow-100">
            {{ forecast.total_radiation_kwh_m2.toFixed(1) }}
          </div>
          <div class="text-xs text-yellow-700 dark:text-yellow-300">
            kWh/m²
          </div>
        </div>

        <!-- Avg Radiation -->
        <div class="p-3 rounded-lg bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800">
          <div class="flex items-center gap-2 mb-1">
            <Zap :size="14" class="text-orange-600 dark:text-orange-400" />
            <span class="text-xs font-medium text-orange-900 dark:text-orange-200">Avg Radiation</span>
          </div>
          <div class="text-lg font-bold text-orange-900 dark:text-orange-100">
            {{ forecast.avg_radiation.toFixed(0) }}
          </div>
          <div class="text-xs text-orange-700 dark:text-orange-300">
            W/m²
          </div>
        </div>

        <!-- Cloud Cover -->
        <div class="p-3 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
          <div class="flex items-center gap-2 mb-1">
            <Cloud :size="14" class="text-blue-600 dark:text-blue-400" />
            <span class="text-xs font-medium text-blue-900 dark:text-blue-200">Avg Clouds</span>
          </div>
          <div class="text-lg font-bold text-blue-900 dark:text-blue-100">
            {{ (forecast.data.reduce((sum, d) => sum + d.cloud_cover, 0) / 24).toFixed(0) }}%
          </div>
          <div class="text-xs text-blue-700 dark:text-blue-300">
            Coverage
          </div>
        </div>
      </div>

      <!-- Solar vs Demand Chart -->
      <div class="h-64 mt-4">
        <Line :data="chartData" :options="chartOptions" />
      </div>

      <!-- Insight -->
      <div class="p-3 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 text-sm">
        <p class="text-green-900 dark:text-green-100 font-medium mb-1">💡 Solar Insight</p>
        <p class="text-green-800 dark:text-green-200 text-xs">
          Peak solar generation at {{ forecast.peak_radiation_hour }}:00 can offset demand.
          Best time for PV-powered operations: {{ Math.max(8, forecast.peak_radiation_hour - 2) }}:00 - {{ Math.min(18, forecast.peak_radiation_hour + 2) }}:00.
        </p>
      </div>
    </div>

    <!-- No Data State -->
    <div v-else class="text-center py-8 text-gray-500 dark:text-gray-400">
      <p class="text-sm">No forecast data available</p>
    </div>
  </div>
</template>
