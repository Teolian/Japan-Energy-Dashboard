<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import SectionCard from '@/components/common/SectionCard.vue'
import MetricBadge from '@/components/common/MetricBadge.vue'
import DemandChart from './DemandChart.vue'
import ReserveBadge from '@/components/reserve/ReserveBadge.vue'
import { useDemandStore } from '@/stores/demand'
import { useReserveStore } from '@/stores/reserve'
import { useJEPXStore } from '@/stores/jepx'
import { useWeatherStore } from '@/stores/weather'
import { useFeatureFlags } from '@/composables/useFeatureFlags'
import { Sun, CloudRain } from 'lucide-vue-next'
import type { Area } from '@/types/weather'

interface Props {
  area: 'tokyo' | 'kansai'
}

const props = defineProps<Props>()

const { t } = useI18n()
const demandStore = useDemandStore()
const reserveStore = useReserveStore()
const jepxStore = useJEPXStore()
const weatherStore = useWeatherStore()
const flags = useFeatureFlags()

// Computed data for the area
const chartData = computed(() =>
  props.area === 'tokyo' ? demandStore.tokyoChartData : demandStore.kansaiChartData
)

const metrics = computed(() =>
  props.area === 'tokyo' ? demandStore.tokyoMetrics : demandStore.kansaiMetrics
)

const prevMetrics = computed(() =>
  props.area === 'tokyo' ? demandStore.prevTokyoMetrics : demandStore.prevKansaiMetrics
)

const priceValues = computed(() =>
  flags.isJEPXEnabled ? jepxStore.priceValues(props.area) : undefined
)

const reserve = computed(() => reserveStore.reserveForArea(props.area))

const forecast = computed(() => weatherStore.forecastForArea(props.area as Area))

// Calculate trend percentages
const peakTrend = computed(() => {
  if (!metrics.value || !prevMetrics.value) return undefined
  return ((metrics.value.peak - prevMetrics.value.peak) / prevMetrics.value.peak) * 100
})

const avgTrend = computed(() => {
  if (!metrics.value || !prevMetrics.value) return undefined
  return ((metrics.value.average - prevMetrics.value.average) / prevMetrics.value.average) * 100
})

// Weather metrics
const avgRadiation = computed(() => forecast.value?.avg_radiation || 0)
const peakRadiationHour = computed(() => forecast.value?.peak_radiation_hour || 12)

const areaName = computed(() => t(`areas.${props.area}`))
</script>

<template>
  <SectionCard :title="areaName">
    <template #actions>
      <ReserveBadge v-if="flags.isReserveEnabled && reserve" :area="area" :data="reserve" />
    </template>

    <!-- Demand Chart -->
    <DemandChart
      title=""
      :data="chartData"
      :prices="priceValues"
    />

    <!-- Compact Metrics Row -->
    <div v-if="metrics" class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
      <div class="grid grid-cols-4 gap-4">
        <MetricBadge
          :label="t('metrics.peak')"
          :value="metrics.peak"
          :unit="t('units.mw')"
          :trend="peakTrend"
          color="blue"
          compact
        />
        <MetricBadge
          :label="t('metrics.average')"
          :value="metrics.average"
          :unit="t('units.mw')"
          :trend="avgTrend"
          color="gray"
          compact
        />
        <MetricBadge
          v-if="metrics.forecastAccuracy"
          :label="t('metrics.forecast')"
          :value="metrics.forecastAccuracy"
          :unit="t('units.percent')"
          color="green"
          compact
        />
        <MetricBadge
          v-if="reserve"
          :label="t('metrics.reserveMargin')"
          :value="reserve.reserve_margin_pct.toFixed(1)"
          :unit="t('units.percent')"
          :color="reserve.reserve_margin_pct >= 8 ? 'green' : reserve.reserve_margin_pct >= 5 ? 'orange' : 'red'"
          compact
        />
      </div>
    </div>

    <!-- Inline Weather Summary -->
    <div v-if="forecast" class="mt-3 flex items-center gap-4 text-xs text-gray-600 dark:text-gray-400">
      <div class="flex items-center gap-1.5">
        <Sun :size="14" class="text-amber-500" />
        <span>Solar peak: {{ peakRadiationHour }}:00 ({{ avgRadiation.toFixed(0) }} W/m² avg)</span>
      </div>
      <div class="flex items-center gap-1.5">
        <CloudRain :size="14" class="text-blue-500" />
        <span>Total: {{ forecast.total_radiation_kwh_m2.toFixed(1) }} kWh/m²</span>
      </div>
    </div>
  </SectionCard>
</template>
