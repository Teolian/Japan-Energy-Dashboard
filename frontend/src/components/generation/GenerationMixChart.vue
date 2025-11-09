<script setup lang="ts">
import { computed } from 'vue'
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
import { useGenerationStore } from '@/stores/generation'
import { Leaf, Zap, Factory } from 'lucide-vue-next'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler)

const generationStore = useGenerationStore()

// Chart data - stacked area for generation mix
const chartData = computed(() => {
  const data = generationStore.tokyoChartData

  if (data.length === 0) {
    return { labels: [], datasets: [] }
  }

  const labels = data.map(d => d.time)

  return {
    labels,
    datasets: [
      {
        label: 'Solar',
        data: data.map(d => d.solar),
        borderColor: 'rgb(251, 191, 36)',
        backgroundColor: 'rgba(251, 191, 36, 0.6)',
        borderWidth: 0,
        fill: true,
        tension: 0.4,
        pointRadius: 0
      },
      {
        label: 'Wind',
        data: data.map(d => d.wind),
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.6)',
        borderWidth: 0,
        fill: true,
        tension: 0.4,
        pointRadius: 0
      },
      {
        label: 'Hydro',
        data: data.map(d => d.hydro),
        borderColor: 'rgb(14, 165, 233)',
        backgroundColor: 'rgba(14, 165, 233, 0.6)',
        borderWidth: 0,
        fill: true,
        tension: 0.4,
        pointRadius: 0
      },
      {
        label: 'Nuclear',
        data: data.map(d => d.nuclear),
        borderColor: 'rgb(168, 85, 247)',
        backgroundColor: 'rgba(168, 85, 247, 0.6)',
        borderWidth: 0,
        fill: true,
        tension: 0.4,
        pointRadius: 0
      },
      {
        label: 'LNG',
        data: data.map(d => d.lng),
        borderColor: 'rgb(156, 163, 175)',
        backgroundColor: 'rgba(156, 163, 175, 0.6)',
        borderWidth: 0,
        fill: true,
        tension: 0.4,
        pointRadius: 0
      },
      {
        label: 'Coal',
        data: data.map(d => d.coal),
        borderColor: 'rgb(75, 85, 99)',
        backgroundColor: 'rgba(75, 85, 99, 0.6)',
        borderWidth: 0,
        fill: true,
        tension: 0.4,
        pointRadius: 0
      },
      {
        label: 'Other',
        data: data.map(d => d.other),
        borderColor: 'rgb(107, 114, 128)',
        backgroundColor: 'rgba(107, 114, 128, 0.4)',
        borderWidth: 0,
        fill: true,
        tension: 0.4,
        pointRadius: 0
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
      position: 'bottom',
      labels: {
        usePointStyle: true,
        padding: 12,
        font: { size: 11 }
      }
    },
    tooltip: {
      enabled: true,
      backgroundColor: 'rgba(0, 0, 0, 0.85)',
      padding: 12,
      cornerRadius: 8,
      callbacks: {
        title: (context) => {
          return `Time: ${context[0].label}`
        },
        label: (context) => {
          const value = context.parsed.y
          const percentage = generationStore.tokyoChartData[context.dataIndex]
          let pct = 0

          if (percentage && percentage.total > 0) {
            pct = (value / percentage.total) * 100
          }

          return `${context.dataset.label}: ${Math.round(value).toLocaleString()} MW (${pct.toFixed(1)}%)`
        },
        footer: (context) => {
          const idx = context[0].dataIndex
          const point = generationStore.tokyoChartData[idx]
          if (!point) return ''

          return [
            `Total: ${Math.round(point.total).toLocaleString()} MW`,
            `Renewable: ${point.renewable_pct.toFixed(1)}%`,
            `Carbon: ${Math.round(point.carbon_gco2_kwh)} gCO₂/kWh`
          ]
        }
      }
    }
  },
  scales: {
    x: {
      stacked: true,
      grid: {
        display: false
      },
      ticks: {
        font: { size: 10 }
      }
    },
    y: {
      stacked: true,
      title: {
        display: true,
        text: 'Generation (MW)',
        font: { size: 12, weight: 'bold' as const }
      },
      grid: {
        color: 'rgba(0, 0, 0, 0.06)'
      },
      ticks: {
        callback: (value) => {
          return `${Math.round(Number(value) / 1000)}k`
        }
      }
    }
  }
}))

// Metrics
const metrics = computed(() => generationStore.tokyoMetrics)
const greenest = computed(() => generationStore.greenestHour)
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex items-start justify-between">
      <div>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
          <Leaf :size="20" class="text-green-500" />
          Generation Mix
        </h2>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
          Estimated from demand + price correlation
        </p>
      </div>

      <!-- Quick Metrics -->
      <div v-if="metrics" class="flex gap-4 text-sm">
        <div class="text-center">
          <div class="text-xs text-gray-500 dark:text-gray-400">Renewable</div>
          <div class="text-lg font-bold text-green-600 dark:text-green-400">
            {{ metrics.renewablePct.toFixed(1) }}%
          </div>
        </div>
        <div class="text-center">
          <div class="text-xs text-gray-500 dark:text-gray-400">Carbon</div>
          <div class="text-lg font-bold" :class="metrics.carbonIntensity.color">
            {{ Math.round(metrics.carbonIntensity.value) }}
          </div>
          <div class="text-xs text-gray-500">gCO₂/kWh</div>
        </div>
      </div>
    </div>

    <!-- Chart -->
    <div class="h-80 bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
      <Line :data="chartData" :options="chartOptions" />
    </div>

    <!-- Greenest Hour Highlight -->
    <div v-if="greenest" class="p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
      <div class="flex items-start gap-2">
        <Zap :size="18" class="text-green-600 dark:text-green-400 mt-0.5" />
        <div class="flex-1">
          <div class="text-sm font-medium text-green-900 dark:text-green-100">
            Greenest Hour: {{ greenest.time }}
          </div>
          <div class="text-xs text-green-700 dark:text-green-300 mt-1">
            {{ greenest.renewablePct.toFixed(1) }}% renewable
            • {{ Math.round(greenest.carbonGCO2) }} gCO₂/kWh
            • Best time for green certificate procurement
          </div>
        </div>
      </div>
    </div>

    <!-- Legend Info -->
    <div class="grid grid-cols-3 gap-3 text-xs">
      <div class="p-2 bg-amber-50 dark:bg-amber-900/10 rounded border border-amber-200 dark:border-amber-800">
        <div class="flex items-center gap-1 text-amber-700 dark:text-amber-400 font-medium mb-1">
          <Leaf :size="14" />
          Renewables
        </div>
        <div class="text-gray-600 dark:text-gray-400">
          Solar, Wind, Hydro<br/>
          Zero carbon emissions
        </div>
      </div>

      <div class="p-2 bg-purple-50 dark:bg-purple-900/10 rounded border border-purple-200 dark:border-purple-800">
        <div class="flex items-center gap-1 text-purple-700 dark:text-purple-400 font-medium mb-1">
          <Zap :size="14" />
          Nuclear
        </div>
        <div class="text-gray-600 dark:text-gray-400">
          Base load<br/>
          Zero carbon emissions
        </div>
      </div>

      <div class="p-2 bg-gray-50 dark:bg-gray-900/10 rounded border border-gray-200 dark:border-gray-700">
        <div class="flex items-center gap-1 text-gray-700 dark:text-gray-400 font-medium mb-1">
          <Factory :size="14" />
          Fossil Fuels
        </div>
        <div class="text-gray-600 dark:text-gray-400">
          LNG, Coal, Other<br/>
          350-850 gCO₂/kWh
        </div>
      </div>
    </div>
  </div>
</template>
