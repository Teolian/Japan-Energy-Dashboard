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
import { useJEPXStore } from '@/stores/jepx'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler)

const jepxStore = useJEPXStore()

// Calculate price spread (Tokyo - Kansai)
const spreadData = computed(() => {
  if (!jepxStore.tokyoData || !jepxStore.kansaiData) return []

  const tokyo = jepxStore.tokyoData.price_yen_per_kwh
  const kansai = jepxStore.kansaiData.price_yen_per_kwh

  return tokyo.map((t, idx) => {
    const k = kansai[idx]
    if (!k) return null

    const spread = t.price - k.price
    return {
      ts: t.ts,
      time: new Date(t.ts).toLocaleTimeString('ja-JP', { hour: '2-digit', minute: '2-digit', timeZone: 'Asia/Tokyo' }),
      spread
    }
  }).filter(s => s !== null)
})

// Metrics
const avgSpread = computed(() => {
  if (spreadData.value.length === 0) return 0
  const sum = spreadData.value.reduce((acc, d) => acc + d.spread, 0)
  return Math.round((sum / spreadData.value.length) * 10) / 10
})

const maxSpread = computed(() => {
  if (spreadData.value.length === 0) return { value: 0, time: '' }
  let max = spreadData.value[0]
  if (!max) return { value: 0, time: '' }
  spreadData.value.forEach(d => {
    if (d && d.spread > max!.spread) max = d
  })
  return { value: Math.round(max.spread * 10) / 10, time: max.time }
})

const minSpread = computed(() => {
  if (spreadData.value.length === 0) return { value: 0, time: '' }
  let min = spreadData.value[0]
  if (!min) return { value: 0, time: '' }
  spreadData.value.forEach(d => {
    if (d && d.spread < min!.spread) min = d
  })
  return { value: Math.round(min.spread * 10) / 10, time: min.time }
})

// Arbitrage detection
const hasArbitrage = computed(() => Math.abs(maxSpread.value.value) > 5 || Math.abs(minSpread.value.value) > 5)

const chartData = computed(() => {
  const spreads = spreadData.value.map(d => d.spread)
  const labels = spreadData.value.map(d => d.time)

  return {
    labels,
    datasets: [
      {
        label: 'Price Spread (Tokyo - Kansai)',
        data: spreads,
        borderColor: 'rgb(0, 102, 204)',
        backgroundColor: (context: any) => {
          // Gradient fill based on positive/negative spread
          if (!context.chart.chartArea) return 'transparent'

          const { ctx, chartArea: { top, bottom } } = context.chart
          const gradient = ctx.createLinearGradient(0, top, 0, bottom)

          // Positive spread = Tokyo more expensive (red gradient)
          gradient.addColorStop(0, 'rgba(220, 38, 38, 0.15)')
          gradient.addColorStop(0.5, 'rgba(0, 0, 0, 0.02)')
          gradient.addColorStop(1, 'rgba(34, 197, 94, 0.15)')

          return gradient
        },
        borderWidth: 2.5,
        tension: 0,
        fill: true,
        pointRadius: 0,
        pointHoverRadius: 5,
        pointHoverBackgroundColor: 'rgb(0, 102, 204)',
        pointHoverBorderColor: '#fff',
        pointHoverBorderWidth: 2
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
      display: false
    },
    tooltip: {
      enabled: true,
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      padding: 12,
      cornerRadius: 8,
      titleFont: {
        size: 14,
        weight: 'bold'
      },
      bodyFont: {
        size: 13
      },
      callbacks: {
        label: (context) => {
          const value = context.parsed.y
          if (value === null || value === undefined) return ''
          const sign = value >= 0 ? '+' : ''
          return `Spread: ${sign}¥${value.toFixed(2)}/kWh`
        },
        afterLabel: (context) => {
          const value = context.parsed.y
          if (value === null || value === undefined) return ''
          if (value > 0) return '(Tokyo more expensive)'
          if (value < 0) return '(Kansai more expensive)'
          return '(Equal prices)'
        }
      }
    }
  },
  scales: {
    x: {
      grid: {
        display: true,
        drawBorder: true,
        color: 'rgba(0, 0, 0, 0.06)',
        lineWidth: 1
      },
      ticks: {
        font: {
          size: 11,
          family: 'ui-monospace, monospace'
        }
      }
    },
    y: {
      type: 'linear',
      display: true,
      position: 'left',
      title: {
        display: true,
        text: 'Spread (JPY/kWh)',
        font: {
          size: 12,
          weight: 'bold' as const
        }
      },
      grid: {
        display: true,
        drawBorder: true,
        color: (context) => {
          // Highlight zero line
          if (context.tick.value === 0) return 'rgba(0, 0, 0, 0.3)'
          return 'rgba(0, 0, 0, 0.06)'
        },
        lineWidth: (context) => {
          if (context.tick.value === 0) return 2
          return 1
        }
      },
      ticks: {
        callback: (value: string | number) => {
          const numValue = Number(value)
          const sign = numValue >= 0 ? '+' : ''
          return `${sign}¥${numValue}`
        },
        font: {
          size: 11,
          family: 'ui-monospace, monospace'
        }
      }
    }
  }
}))
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div>
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Regional Price Spread</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
        Tokyo Spot Price - Kansai Spot Price (positive = Tokyo more expensive)
      </p>
    </div>

    <!-- Chart -->
    <div class="h-64">
      <Line :data="chartData" :options="chartOptions" />
    </div>

    <!-- Metrics -->
    <div class="grid grid-cols-3 gap-4 pt-4 border-t border-gray-200 dark:border-gray-700">
      <div>
        <div class="text-xs text-gray-500 dark:text-gray-400 mb-1">Average Spread</div>
        <div class="text-2xl font-bold" :class="avgSpread >= 0 ? 'text-red-600' : 'text-green-600'">
          {{ avgSpread >= 0 ? '+' : '' }}¥{{ avgSpread }}
        </div>
        <div class="text-xs text-gray-500">per kWh</div>
      </div>

      <div>
        <div class="text-xs text-gray-500 dark:text-gray-400 mb-1">Max Spread ({{ maxSpread.time }})</div>
        <div class="text-2xl font-bold text-red-600">
          +¥{{ maxSpread.value }}
        </div>
        <div v-if="maxSpread.value > 5" class="text-xs text-red-600 font-medium">
          ⚡ Arbitrage opportunity
        </div>
        <div v-else class="text-xs text-gray-500">
          Tokyo peak
        </div>
      </div>

      <div>
        <div class="text-xs text-gray-500 dark:text-gray-400 mb-1">Min Spread ({{ minSpread.time }})</div>
        <div class="text-2xl font-bold text-green-600">
          {{ minSpread.value }}
        </div>
        <div v-if="Math.abs(minSpread.value) > 5" class="text-xs text-green-600 font-medium">
          ⚡ Reverse arbitrage
        </div>
        <div v-else class="text-xs text-gray-500">
          Kansai peak
        </div>
      </div>
    </div>

    <!-- Insights -->
    <div v-if="hasArbitrage" class="p-4 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg">
      <div class="flex items-start gap-2">
        <span class="text-amber-600 dark:text-amber-400">💡</span>
        <div>
          <div class="font-medium text-amber-900 dark:text-amber-100 text-sm">Arbitrage Opportunity Detected</div>
          <div class="text-sm text-amber-700 dark:text-amber-300 mt-1">
            <span v-if="maxSpread.value > 5">
              Tokyo is +¥{{ maxSpread.value }}/kWh more expensive at {{ maxSpread.time }}.
              Consider buying from Kansai and selling to Tokyo region.
            </span>
            <span v-else-if="Math.abs(minSpread.value) > 5">
              Kansai is ¥{{ Math.abs(minSpread.value) }}/kWh more expensive at {{ minSpread.time }}.
              Consider reverse arbitrage strategy.
            </span>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
      <div class="text-sm text-gray-600 dark:text-gray-400">
        💡 Price spread is within normal range (±¥5). No significant arbitrage opportunities detected today.
      </div>
    </div>
  </div>
</template>
