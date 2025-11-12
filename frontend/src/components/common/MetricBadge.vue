<script setup lang="ts">
import { computed } from 'vue'
import { TrendingUp, TrendingDown } from 'lucide-vue-next'

interface Props {
  label: string
  value: string | number
  unit?: string
  trend?: number // Percentage change (e.g., +2.3 or -1.5)
  color?: 'blue' | 'green' | 'orange' | 'red' | 'gray' | 'purple'
  size?: 'sm' | 'md' | 'lg'
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  color: 'gray',
  size: 'md',
  compact: false
})

const colorClasses = computed(() => {
  const colors = {
    blue: 'text-blue-600 dark:text-blue-400',
    green: 'text-green-600 dark:text-green-400',
    orange: 'text-orange-600 dark:text-orange-400',
    red: 'text-red-600 dark:text-red-400',
    gray: 'text-gray-900 dark:text-white',
    purple: 'text-purple-600 dark:text-purple-400'
  }
  return colors[props.color]
})

const sizeClasses = computed(() => {
  if (props.compact) {
    return {
      label: 'text-xs',
      value: 'text-sm',
      unit: 'text-xs'
    }
  }

  const sizes = {
    sm: {
      label: 'text-xs',
      value: 'text-base',
      unit: 'text-xs'
    },
    md: {
      label: 'text-xs',
      value: 'text-lg',
      unit: 'text-xs'
    },
    lg: {
      label: 'text-sm',
      value: 'text-2xl',
      unit: 'text-sm'
    }
  }
  return sizes[props.size]
})

const trendColor = computed(() => {
  if (!props.trend) return ''
  return props.trend > 0
    ? 'text-green-600 dark:text-green-400'
    : 'text-red-600 dark:text-red-400'
})

const formattedValue = computed(() => {
  if (typeof props.value === 'number') {
    // Format large numbers with k/M suffix
    if (props.value >= 1000000) {
      return (props.value / 1000000).toFixed(1) + 'M'
    } else if (props.value >= 1000) {
      return (props.value / 1000).toFixed(1) + 'k'
    }
    return props.value.toLocaleString()
  }
  return props.value
})
</script>

<template>
  <div class="flex flex-col" :class="compact ? 'gap-0.5' : 'gap-1'">
    <!-- Label -->
    <div :class="[sizeClasses.label, 'text-gray-500 dark:text-gray-400 font-medium']">
      {{ label }}
    </div>

    <!-- Value + Unit -->
    <div class="flex items-baseline gap-1.5">
      <span :class="[sizeClasses.value, colorClasses, 'font-bold']">
        {{ formattedValue }}
      </span>
      <span v-if="unit" :class="[sizeClasses.unit, 'text-gray-500 dark:text-gray-400']">
        {{ unit }}
      </span>
    </div>

    <!-- Trend Indicator -->
    <div v-if="trend !== undefined" class="flex items-center gap-1">
      <TrendingUp v-if="trend > 0" :size="12" :class="trendColor" />
      <TrendingDown v-else-if="trend < 0" :size="12" :class="trendColor" />
      <span :class="[sizeClasses.unit, trendColor, 'font-medium']">
        {{ trend > 0 ? '+' : '' }}{{ trend.toFixed(1) }}%
      </span>
    </div>
  </div>
</template>
