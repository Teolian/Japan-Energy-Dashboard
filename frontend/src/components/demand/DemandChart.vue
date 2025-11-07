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

  // Create point styles array - highlight peak hour only (professional style)
  const pointBackgroundColors = actualValues.map((_, i) =>
    i === peakIndex ? 'rgb(220, 38, 38)' : 'transparent'
  )
  const pointBorderColors = actualValues.map((_, i) =>
    i === peakIndex ? 'rgb(220, 38, 38)' : 'transparent'
  )
  const pointRadii = actualValues.map((_, i) => i === peakIndex ? 5 : 0) // Only peak visible

  const datasets = [
    {
      label: 'Actual Demand',
      data: props.data.map(d => d.actual),
      borderColor: 'rgb(0, 102, 204)',      // Professional blue
      backgroundColor: 'transparent',        // No fill
      borderWidth: 2.5,                      // Thicker line
      tension: 0,                            // Straight lines (professional)
      yAxisID: 'y',
      pointBackgroundColor: pointBackgroundColors,
      pointBorderColor: pointBorderColors,
      pointRadius: pointRadii,               // Show peak point only
      pointHoverRadius: pointRadii.map(r => r + 3),
      pointHoverBorderColor: '#fff',
      pointHoverBorderWidth: 2
    },
    {
      label: 'Forecast',
      data: props.data.map(d => d.forecast || null),
      borderColor: 'rgb(107, 114, 128)',    // Gray (professional)
      backgroundColor: 'transparent',
      borderWidth: 2,
      borderDash: [8, 4],                    // Longer dashes
      tension: 0,
      yAxisID: 'y',
      pointRadius: 0,
      pointHoverRadius: 5,
      pointHoverBackgroundColor: 'rgb(107, 114, 128)',
      pointHoverBorderColor: '#fff',
      pointHoverBorderWidth: 2
    }
  ]

  // Add price overlay if prices provided
  if (props.prices && props.prices.length > 0) {
    datasets.push({
      label: 'Spot Price',
      data: props.prices,
      borderColor: 'rgb(220, 38, 38)',      // Red (financial alert color)
      backgroundColor: 'transparent',
      borderWidth: 2.5,
      tension: 0,
      yAxisID: 'y1',
      pointRadius: 0,
      pointHoverRadius: 5,
      pointHoverBackgroundColor: 'rgb(220, 38, 38)',
      pointHoverBorderColor: '#fff',
      pointHoverBorderWidth: 2
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
          family: 'ui-monospace, monospace'  // Professional monospace
        }
      }
    },
    y: {
      type: 'linear',
      display: true,
      position: 'left',
      beginAtZero: false,
      title: {
        display: true,
        text: 'Demand (MW)',
        font: {
          size: 12,
          weight: 'bold' as const
        }
      },
      grid: {
        display: true,
        drawBorder: true,
        color: 'rgba(0, 0, 0, 0.06)',
        lineWidth: 1
      },
      ticks: {
        callback: (value) => `${value.toLocaleString()}`,
        font: {
          size: 11,
          family: 'ui-monospace, monospace'
        }
      }
    },
    ...(props.prices && props.prices.length > 0 ? {
      y1: {
        type: 'linear' as const,
        display: true,
        position: 'right' as const,
        beginAtZero: false,
        title: {
          display: true,
          text: 'Price (JPY/kWh)',
          font: {
            size: 12,
            weight: 'bold' as const
          }
        },
        ticks: {
          callback: (value: any) => `¥${value.toLocaleString()}`,
          font: {
            size: 11,
            family: 'ui-monospace, monospace'
          }
        },
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
