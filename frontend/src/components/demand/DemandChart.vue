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
  type ChartOptions
} from 'chart.js'
import type { ChartDataPoint } from '@/types/demand'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend)

interface Props {
  title: string
  data: ChartDataPoint[]
  prices?: number[] // Optional JEPX spot prices (JPY/kWh)
}

const props = defineProps<Props>()

const chartData = computed(() => {
  // Find peak hour index
  const actualValues = props.data.map(d => d.actual)
  const peakValue = Math.max(...actualValues)
  const peakIndex = actualValues.indexOf(peakValue)

  // Create point styles array - highlight peak hour
  const pointBackgroundColors = actualValues.map((_, i) =>
    i === peakIndex ? 'rgb(239, 68, 68)' : 'rgb(14, 165, 233)'
  )
  const pointBorderColors = actualValues.map((_, i) =>
    i === peakIndex ? 'rgb(220, 38, 38)' : 'rgb(14, 165, 233)'
  )
  const pointRadii = actualValues.map((_, i) => i === peakIndex ? 6 : 3)

  const datasets = [
    {
      label: 'Actual Demand',
      data: props.data.map(d => d.actual),
      borderColor: 'rgb(14, 165, 233)',
      backgroundColor: 'rgba(14, 165, 233, 0.1)',
      borderWidth: 2,
      tension: 0.4,
      yAxisID: 'y', // Left axis (MW)
      pointBackgroundColor: pointBackgroundColors,
      pointBorderColor: pointBorderColors,
      pointRadius: pointRadii,
      pointHoverRadius: pointRadii.map(r => r + 2)
    },
    {
      label: 'Forecast',
      data: props.data.map(d => d.forecast || null),
      borderColor: 'rgb(147, 51, 234)',
      borderWidth: 2,
      borderDash: [5, 5],
      tension: 0.4,
      yAxisID: 'y' // Left axis (MW)
    }
  ]

  // Add price overlay if prices provided
  if (props.prices && props.prices.length > 0) {
    datasets.push({
      label: 'Spot Price',
      data: props.prices,
      borderColor: 'rgb(251, 146, 60)', // Orange
      backgroundColor: 'rgba(251, 146, 60, 0.1)',
      borderWidth: 2,
      tension: 0.4,
      yAxisID: 'y1' // Right axis (JPY/kWh)
    } as any)
  }

  return {
    labels: props.data.map(d => d.time),
    datasets
  }
})

const chartOptions = computed<ChartOptions<'line'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  animation: {
    duration: 750,
    easing: 'easeInOutQuart'
  },
  transitions: {
    active: {
      animation: {
        duration: 300
      }
    }
  },
  interaction: {
    mode: 'index',
    intersect: false
  },
  hover: {
    mode: 'index',
    intersect: false,
    animationDuration: 200
  },
  plugins: {
    legend: {
      display: true,
      position: 'top',
      labels: {
        usePointStyle: true,
        padding: 15
      }
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
      displayColors: true,
      callbacks: {
        label: (context) => {
          const label = context.dataset.label || ''
          const value = context.parsed.y?.toFixed(1)
          // Check if this is the price dataset (yAxisID: 'y1')
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
      type: 'linear',
      display: true,
      position: 'left',
      beginAtZero: false,
      title: {
        display: true,
        text: 'Demand (MW)'
      },
      ticks: { callback: (value) => `${value} MW` }
    },
    ...(props.prices && props.prices.length > 0 ? {
      y1: {
        type: 'linear' as const,
        display: true,
        position: 'right' as const,
        beginAtZero: false,
        title: {
          display: true,
          text: 'Price (JPY/kWh)'
        },
        ticks: { callback: (value: any) => `¥${value}` },
        grid: {
          drawOnChartArea: false // Don't draw gridlines for right axis
        }
      }
    } : {})
  }
}))
</script>

<template>
  <div class="h-64">
    <h3 v-if="title" class="text-lg font-semibold mb-4">{{ title }}</h3>
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
