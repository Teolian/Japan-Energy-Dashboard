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
import { useDemandStore } from '@/stores/demand'
import { useJEPXStore } from '@/stores/jepx'
import { Battery, TrendingDown, TrendingUp, Zap } from 'lucide-vue-next'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler)

const demandStore = useDemandStore()
const jepxStore = useJEPXStore()

// Duck Curve Analysis
const duckCurveData = computed(() => {
  if (!demandStore.tokyoData || !jepxStore.tokyoData) return null

  const demand = demandStore.tokyoData.series.map(s => s.demand_mw)
  const prices = jepxStore.tokyoData.price_yen_per_kwh.map(p => p.price)
  const times = demandStore.tokyoData.series.map(s =>
    new Date(s.ts).getHours()
  )

  // Find morning peak (6-10h), midday valley (11-15h), evening peak (17-21h)
  const morningDemand = demand.slice(6, 11)
  const morningPeak = Math.max(...morningDemand)
  const morningPeakHour = 6 + morningDemand.indexOf(morningPeak)

  const middayDemand = demand.slice(11, 16)
  const middayValley = Math.min(...middayDemand)
  const middayValleyHour = 11 + middayDemand.indexOf(middayValley)

  const eveningDemand = demand.slice(17, 22)
  const eveningPeak = Math.max(...eveningDemand)
  const eveningPeakHour = 17 + eveningDemand.indexOf(eveningPeak)

  // Duck Curve Severity: (morning_peak - midday_valley) / morning_peak
  const severity = ((morningPeak - middayValley) / morningPeak) * 100

  // Find global price min/max for battery arbitrage
  const minPrice = Math.min(...prices)
  const maxPrice = Math.max(...prices)
  const chargeHour = prices.indexOf(minPrice)
  const dischargeHour = prices.indexOf(maxPrice)

  // Calculate spread (only profitable if discharge after charge)
  let priceSpread = maxPrice - minPrice
  let isProfitable = dischargeHour > chargeHour && priceSpread > 0

  // If not profitable, set spread to 0
  if (!isProfitable) {
    priceSpread = 0
  }

  return {
    severity,
    morningPeak,
    morningPeakHour,
    middayValley,
    middayValleyHour,
    eveningPeak,
    eveningPeakHour,
    chargeHour,
    dischargeHour,
    chargePrice: minPrice,
    dischargePrice: maxPrice,
    priceSpread,
    isProfitable,
    demand,
    prices,
    times
  }
})

// Severity level classification
const severityLevel = computed(() => {
  if (!duckCurveData.value) return { level: 'low', color: 'text-gray-500', label: 'Normal' }

  const severity = duckCurveData.value.severity

  if (severity >= 30) return {
    level: 'high',
    color: 'text-red-600 dark:text-red-400',
    label: 'High',
    bgColor: 'bg-red-50 dark:bg-red-900/20',
    borderColor: 'border-red-200 dark:border-red-800'
  }
  if (severity >= 20) return {
    level: 'moderate',
    color: 'text-orange-600 dark:text-orange-400',
    label: 'Moderate',
    bgColor: 'bg-orange-50 dark:bg-orange-900/20',
    borderColor: 'border-orange-200 dark:border-orange-800'
  }
  return {
    level: 'low',
    color: 'text-green-600 dark:text-green-400',
    label: 'Low',
    bgColor: 'bg-green-50 dark:bg-green-900/20',
    borderColor: 'border-green-200 dark:border-green-800'
  }
})

// Battery optimization ROI calculation
const batteryROI = computed(() => {
  if (!duckCurveData.value) return 0

  // Assume 100 MWh battery capacity
  const capacityMWh = 100
  const profitPerCycle = capacityMWh * duckCurveData.value.priceSpread

  return Math.round(profitPerCycle)
})

// Chart data
const chartData = computed(() => {
  if (!duckCurveData.value) return { labels: [], datasets: [] }

  const data = duckCurveData.value
  const labels = data.times.map(h => `${h.toString().padStart(2, '0')}:00`)

  return {
    labels,
    datasets: [
      {
        label: 'Demand (MW)',
        data: data.demand,
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        borderWidth: 2.5,
        tension: 0.3,
        fill: true,
        yAxisID: 'y',
        pointRadius: 0,
        pointHoverRadius: 6
      },
      {
        label: 'Price (¥/kWh)',
        data: data.prices,
        borderColor: 'rgb(249, 115, 22)',
        backgroundColor: 'transparent',
        borderWidth: 2,
        borderDash: [5, 5],
        tension: 0.3,
        fill: false,
        yAxisID: 'y1',
        pointRadius: 0,
        pointHoverRadius: 6
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
        padding: 15,
        font: { size: 12 }
      }
    },
    tooltip: {
      enabled: true,
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      padding: 12,
      cornerRadius: 8
    }
  },
  scales: {
    x: {
      grid: {
        display: true,
        color: 'rgba(0, 0, 0, 0.06)'
      },
      ticks: {
        font: { size: 11 }
      }
    },
    y: {
      type: 'linear',
      display: true,
      position: 'left',
      title: {
        display: true,
        text: 'Demand (MW)',
        font: { size: 12, weight: 'bold' as const }
      },
      grid: {
        display: true,
        color: 'rgba(0, 0, 0, 0.06)'
      }
    },
    y1: {
      type: 'linear',
      display: true,
      position: 'right',
      title: {
        display: true,
        text: 'Price (¥/kWh)',
        font: { size: 12, weight: 'bold' as const }
      },
      grid: {
        display: false
      }
    }
  }
}))
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div>
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
        <Zap :size="20" class="text-amber-500" />
        Duck Curve Analysis
      </h2>
      <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
        Demand pattern & battery optimization signals
      </p>
    </div>

    <!-- Severity Card -->
    <div
      v-if="duckCurveData"
      class="p-4 rounded-lg border"
      :class="[severityLevel.bgColor, severityLevel.borderColor]"
    >
      <div class="flex items-start justify-between">
        <div>
          <div class="text-sm font-medium text-gray-600 dark:text-gray-300">
            Duck Curve Severity
          </div>
          <div class="mt-1 flex items-baseline gap-2">
            <span class="text-3xl font-bold" :class="severityLevel.color">
              {{ Math.round(duckCurveData.severity) }}%
            </span>
            <span class="text-sm font-medium" :class="severityLevel.color">
              {{ severityLevel.label }}
            </span>
          </div>
          <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">
            ({{ duckCurveData.morningPeak.toLocaleString() }} MW → {{ duckCurveData.middayValley.toLocaleString() }} MW drop)
          </div>
        </div>
        <div class="text-right">
          <div class="text-xs text-gray-500 dark:text-gray-400 mb-1">Price Volatility</div>
          <div class="text-lg font-bold text-gray-900 dark:text-white">
            ¥{{ duckCurveData.priceSpread.toFixed(2) }}
          </div>
          <div class="text-xs text-gray-500">spread</div>
        </div>
      </div>
    </div>

    <!-- Chart -->
    <div class="h-80 bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
      <Line :data="chartData" :options="chartOptions" />
    </div>

    <!-- Battery Optimization Signals -->
    <div
      v-if="duckCurveData"
      class="grid grid-cols-2 gap-4"
    >
      <!-- Charge Signal -->
      <div class="p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
        <div class="flex items-start gap-3">
          <div class="p-2 bg-green-100 dark:bg-green-800 rounded-lg">
            <TrendingDown :size="20" class="text-green-600 dark:text-green-400" />
          </div>
          <div class="flex-1">
            <div class="text-sm font-medium text-green-900 dark:text-green-100">
              Charge Battery
            </div>
            <div class="mt-1 text-2xl font-bold text-green-600 dark:text-green-400">
              {{ duckCurveData.chargeHour.toString().padStart(2, '0') }}:00
            </div>
            <div class="mt-1 text-xs text-green-700 dark:text-green-300">
              Low price: ¥{{ duckCurveData.chargePrice.toFixed(2) }}/kWh
            </div>
          </div>
        </div>
      </div>

      <!-- Discharge Signal -->
      <div class="p-4 bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800 rounded-lg">
        <div class="flex items-start gap-3">
          <div class="p-2 bg-orange-100 dark:bg-orange-800 rounded-lg">
            <TrendingUp :size="20" class="text-orange-600 dark:text-orange-400" />
          </div>
          <div class="flex-1">
            <div class="text-sm font-medium text-orange-900 dark:text-orange-100">
              Discharge Battery
            </div>
            <div class="mt-1 text-2xl font-bold text-orange-600 dark:text-orange-400">
              {{ duckCurveData.dischargeHour.toString().padStart(2, '0') }}:00
            </div>
            <div class="mt-1 text-xs text-orange-700 dark:text-orange-300">
              Peak price: ¥{{ duckCurveData.dischargePrice.toFixed(2) }}/kWh
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ROI Insight -->
    <div
      v-if="duckCurveData && duckCurveData.isProfitable"
      class="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg"
    >
      <div class="flex items-start gap-2">
        <Battery :size="20" class="text-blue-600 dark:text-blue-400 mt-0.5" />
        <div class="flex-1">
          <div class="font-medium text-blue-900 dark:text-blue-100 text-sm">
            DGP Battery Optimization
          </div>
          <div class="mt-1 text-sm text-blue-700 dark:text-blue-300">
            Based on today's duck curve pattern: charge at
            <span class="font-semibold">{{ duckCurveData.chargeHour.toString().padStart(2, '0') }}:00</span>
            (¥{{ duckCurveData.chargePrice.toFixed(2) }}), discharge at
            <span class="font-semibold">{{ duckCurveData.dischargeHour.toString().padStart(2, '0') }}:00</span>
            (¥{{ duckCurveData.dischargePrice.toFixed(2) }}).
          </div>
          <div class="mt-2 inline-flex items-baseline gap-2 px-3 py-1 bg-blue-100 dark:bg-blue-800 rounded-full">
            <span class="text-xs text-blue-600 dark:text-blue-300">Estimated profit (100 MWh):</span>
            <span class="text-sm font-bold text-blue-700 dark:text-blue-200">¥{{ batteryROI.toLocaleString() }}/cycle</span>
          </div>
        </div>
      </div>
    </div>

    <!-- No Arbitrage Opportunity Warning -->
    <div
      v-if="duckCurveData && !duckCurveData.isProfitable"
      class="p-4 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg"
    >
      <div class="flex items-start gap-2">
        <Battery :size="20" class="text-gray-400 mt-0.5" />
        <div class="flex-1">
          <div class="font-medium text-gray-900 dark:text-gray-100 text-sm">
            No Battery Arbitrage Opportunity
          </div>
          <div class="mt-1 text-sm text-gray-600 dark:text-gray-400">
            Price pattern does not support profitable charge/discharge cycle today.
            Price range: ¥{{ duckCurveData.chargePrice.toFixed(2) }} - ¥{{ duckCurveData.dischargePrice.toFixed(2) }}/kWh.
          </div>
        </div>
      </div>
    </div>

    <!-- Explanation -->
    <div class="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg text-sm text-gray-600 dark:text-gray-400">
      <p class="font-medium text-gray-900 dark:text-white mb-2">What is Duck Curve?</p>
      <p>
        The "duck curve" shows how solar generation creates a demand valley at midday,
        followed by a steep evening ramp-up. Higher severity = greater battery arbitrage opportunity.
      </p>
    </div>
  </div>
</template>
