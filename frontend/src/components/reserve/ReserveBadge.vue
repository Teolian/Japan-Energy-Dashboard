<script setup lang="ts">
import { computed } from 'vue'
import type { Area } from '@/types/demand'
import type { AreaReserve } from '@/types/reserve'
import { STATUS_CONFIGS } from '@/types/reserve'
import Tooltip from '@/components/common/Tooltip.vue'

interface Props {
  area: Area
  data: AreaReserve | null
}

const props = defineProps<Props>()

const statusConfig = computed(() => {
  if (!props.data) return null
  return STATUS_CONFIGS[props.data.status]
})

const displayText = computed(() => {
  if (!props.data) return 'No data'
  return `${props.data.reserve_margin_pct.toFixed(1)}%`
})

const tooltipText = computed(() => {
  if (!props.data) return 'No reserve data available'

  const explanations = {
    stable: 'Stable: ≥8% reserve margin - Power supply is secure',
    watch: 'Watch: 5-8% reserve margin - Monitoring required',
    tight: 'Tight: <5% reserve margin - Power supply is constrained'
  }

  return explanations[props.data.status]
})
</script>

<template>
  <Tooltip v-if="data && statusConfig" :content="tooltipText" position="left">
    <div :class="['inline-flex items-center gap-2 px-3 py-1.5 rounded-lg border cursor-help', statusConfig.bgColor]">
      <span :class="['text-sm font-medium', statusConfig.color]">
        {{ statusConfig.icon }}
      </span>
      <span :class="['text-sm font-semibold', statusConfig.color]">
        {{ displayText }}
      </span>
      <span :class="['text-xs', statusConfig.color]">
        {{ statusConfig.label }}
      </span>
    </div>
  </Tooltip>
  <div v-else class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
    <span class="text-sm text-gray-500 dark:text-gray-400">
      Reserve: N/A
    </span>
  </div>
</template>
