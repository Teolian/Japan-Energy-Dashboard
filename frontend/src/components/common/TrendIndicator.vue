<script setup lang="ts">
import { computed } from 'vue'
import { TrendingUp, TrendingDown, Minus } from 'lucide-vue-next'

interface Props {
  current: number
  previous?: number
  format?: 'number' | 'percentage'
  size?: number
}

const props = withDefaults(defineProps<Props>(), {
  format: 'number',
  size: 16
})

// Calculate change and trend
const trend = computed(() => {
  if (props.previous === undefined || props.previous === 0) {
    return { type: 'neutral', change: 0, formatted: '—' }
  }

  const diff = props.current - props.previous
  const percentChange = (diff / props.previous) * 100

  const type = diff > 0 ? 'up' : diff < 0 ? 'down' : 'neutral'

  const formatted = props.format === 'percentage'
    ? `${Math.abs(percentChange).toFixed(1)}%`
    : Math.abs(diff).toFixed(0)

  return { type, change: percentChange, formatted }
})

const icon = computed(() => {
  switch (trend.value.type) {
    case 'up': return TrendingUp
    case 'down': return TrendingDown
    default: return Minus
  }
})

const colorClass = computed(() => {
  switch (trend.value.type) {
    case 'up': return 'text-green-600 dark:text-green-400'
    case 'down': return 'text-red-600 dark:text-red-400'
    default: return 'text-gray-400 dark:text-gray-500'
  }
})

const sign = computed(() => {
  if (trend.value.type === 'up') return '+'
  if (trend.value.type === 'down') return '-'
  return ''
})
</script>

<template>
  <div :class="['flex items-center gap-1 text-sm', colorClass]">
    <component :is="icon" :size="size" />
    <span>{{ sign }}{{ trend.formatted }}</span>
  </div>
</template>
