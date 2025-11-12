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
  Legend
} from 'chart.js'
import { useChartConfig } from '@/composables/useChartConfig'
import type { ChartDataPoint } from '@/types/demand'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend)

interface Props {
  title: string
  data: ChartDataPoint[]
  prices?: number[] // Optional JEPX spot prices (JPY/kWh)
}

const props = defineProps<Props>()
const { colors, dualAxisConfig, lineDatasetDefaults } = useChartConfig()

const chartData = computed(() => {
  // Find peak hour index
  const actualValues = props.data.map(d => d.actual)
  const peakValue = Math.max(...actualValues)
  const peakIndex = actualValues.indexOf(peakValue)

  // Create point styles array - highlight peak hour only
  const pointBackgroundColors = actualValues.map((_, i) =>
    i === peakIndex ? colors.peak : 'transparent'
  )
  const pointBorderColors = actualValues.map((_, i) =>
    i === peakIndex ? colors.peak : 'transparent'
  )
  const pointRadii = actualValues.map((_, i) => i === peakIndex ? 5 : 0)

  const datasets = [
    {
      label: 'Actual Demand',
      data: props.data.map(d => d.actual),
      borderColor: colors.tokyo,
      backgroundColor: 'transparent',
      ...lineDatasetDefaults,
      yAxisID: 'y',
      pointBackgroundColor: pointBackgroundColors,
      pointBorderColor: pointBorderColors,
      pointRadius: pointRadii,
      pointHoverRadius: pointRadii.map(r => r > 0 ? r + 3 : 5)
    },
    {
      label: 'Forecast',
      data: props.data.map(d => d.forecast || null),
      borderColor: colors.forecast,
      backgroundColor: 'transparent',
      borderWidth: 2,
      borderDash: [8, 4],
      tension: 0.3,
      yAxisID: 'y',
      pointRadius: 0,
      pointHoverRadius: 5,
      pointHoverBackgroundColor: colors.forecast,
      pointHoverBorderColor: '#fff',
      pointHoverBorderWidth: 2
    }
  ]

  // Add price overlay if prices provided
  if (props.prices && props.prices.length > 0) {
    datasets.push({
      label: 'Spot Price',
      data: props.prices,
      borderColor: colors.price,
      backgroundColor: 'transparent',
      ...lineDatasetDefaults,
      yAxisID: 'y1'
    } as any)
  }

  return {
    labels: props.data.map(d => d.time),
    datasets
  }
})

const chartOptions = computed(() => {
  return dualAxisConfig('Demand (MW)', 'Price (JPY/kWh)', {
    plugins: {
      tooltip: {
        callbacks: {
          label: (context) => {
            const label = context.dataset.label || ''
            const value = context.parsed.y?.toFixed(1)
            if ((context.dataset as any).yAxisID === 'y1') {
              return `${label}: ¥${value}/kWh`
            }
            return `${label}: ${value} MW`
          }
        }
      }
    },
    scales: {
      y: {
        ticks: {
          callback: (value) => `${value.toLocaleString()}`
        }
      },
      y1: props.prices && props.prices.length > 0 ? {
        ticks: {
          callback: (value: any) => `¥${value.toLocaleString()}`
        }
      } : undefined
    }
  })
})
</script>

<template>
  <div class="h-64">
    <h3 v-if="title" class="text-lg font-semibold mb-4">{{ title }}</h3>
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
