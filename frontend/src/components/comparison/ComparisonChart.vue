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
import { useDemandStore } from '@/stores/demand'
import { useJEPXStore } from '@/stores/jepx'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend)

const demandStore = useDemandStore()
const jepxStore = useJEPXStore()

const chartData = computed(() => {
  const labels = Array.from({ length: 24 }, (_, i) => `${i.toString().padStart(2, '0')}:00`)

  const tokyoDemand = demandStore.tokyoData?.series.map(s => s.demand_mw) || []
  const kansaiDemand = demandStore.kansaiData?.series.map(s => s.demand_mw) || []
  const tokyoPrices = jepxStore.priceValues('tokyo')
  const kansaiPrices = jepxStore.priceValues('kansai')

  return {
    labels,
    datasets: [
      {
        label: 'Tokyo Demand',
        data: tokyoDemand,
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        borderWidth: 2,
        tension: 0.4,
        yAxisID: 'y',
        pointRadius: 2,
        pointHoverRadius: 4
      },
      {
        label: 'Kansai Demand',
        data: kansaiDemand,
        borderColor: 'rgb(16, 185, 129)',
        backgroundColor: 'rgba(16, 185, 129, 0.1)',
        borderWidth: 2,
        tension: 0.4,
        yAxisID: 'y',
        pointRadius: 2,
        pointHoverRadius: 4
      },
      {
        label: 'Tokyo Price',
        data: tokyoPrices,
        borderColor: 'rgb(251, 146, 60)',
        backgroundColor: 'rgba(251, 146, 60, 0.1)',
        borderWidth: 2,
        borderDash: [5, 5],
        tension: 0.4,
        yAxisID: 'y1',
        pointRadius: 2,
        pointHoverRadius: 4
      },
      {
        label: 'Kansai Price',
        data: kansaiPrices,
        borderColor: 'rgb(234, 88, 12)',
        backgroundColor: 'rgba(234, 88, 12, 0.1)',
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
  animation: {
    duration: 750,
    easing: 'easeInOutQuart'
  },
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
        padding: 15
      }
    },
    tooltip: {
      enabled: true,
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      padding: 12,
      cornerRadius: 8,
      callbacks: {
        label: (context) => {
          const label = context.dataset.label || ''
          const value = context.parsed.y?.toFixed(1)
          // Check if this is price dataset
          if (label.includes('Price')) {
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
      title: {
        display: true,
        text: 'Demand (MW)'
      },
      ticks: {
        callback: (value) => `${value} MW`
      }
    },
    y1: {
      type: 'linear',
      display: true,
      position: 'right',
      title: {
        display: true,
        text: 'Price (JPY/kWh)'
      },
      ticks: {
        callback: (value: any) => `¥${value}`
      },
      grid: {
        drawOnChartArea: false
      }
    }
  }
}))
</script>

<template>
  <div class="h-80">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
