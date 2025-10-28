<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  data: number[]
  width?: number
  height?: number
  color?: string
  fillColor?: string
  strokeWidth?: number
}

const props = withDefaults(defineProps<Props>(), {
  width: 80,
  height: 24,
  color: 'currentColor',
  fillColor: 'none',
  strokeWidth: 1.5
})

// Calculate SVG path from data points
const path = computed(() => {
  if (props.data.length === 0) return ''

  const min = Math.min(...props.data)
  const max = Math.max(...props.data)
  const range = max - min || 1 // Avoid division by zero

  const points = props.data.map((value, index) => {
    const x = (index / (props.data.length - 1)) * props.width
    const y = props.height - ((value - min) / range) * props.height
    return `${x},${y}`
  })

  return `M ${points.join(' L ')}`
})

// Calculate area fill path (optional)
const areaPath = computed(() => {
  if (props.data.length === 0 || props.fillColor === 'none') return ''

  const min = Math.min(...props.data)
  const max = Math.max(...props.data)
  const range = max - min || 1

  const points = props.data.map((value, index) => {
    const x = (index / (props.data.length - 1)) * props.width
    const y = props.height - ((value - min) / range) * props.height
    return `${x},${y}`
  })

  const firstX = 0
  const lastX = props.width
  const bottomY = props.height

  return `M ${firstX},${bottomY} L ${points.join(' L ')} L ${lastX},${bottomY} Z`
})
</script>

<template>
  <svg
    :width="width"
    :height="height"
    class="inline-block"
    preserveAspectRatio="none"
    viewBox="`0 0 ${width} ${height}`"
  >
    <!-- Area fill (optional) -->
    <path
      v-if="fillColor !== 'none'"
      :d="areaPath"
      :fill="fillColor"
      opacity="0.2"
    />

    <!-- Line path -->
    <path
      :d="path"
      :stroke="color"
      :stroke-width="strokeWidth"
      fill="none"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</template>
