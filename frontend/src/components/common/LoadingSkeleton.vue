<script setup lang="ts">
interface Props {
  type?: 'text' | 'chart' | 'card' | 'badge'
  height?: string
  width?: string
  rounded?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  height: '1rem',
  width: '100%',
  rounded: false
})

const heightClass = props.height
const widthClass = props.width
</script>

<template>
  <div
    :class="[
      'animate-pulse bg-gray-200 dark:bg-gray-700',
      rounded && 'rounded-full',
      !rounded && type === 'card' && 'rounded-lg',
      !rounded && type === 'badge' && 'rounded-md',
      !rounded && type === 'text' && 'rounded'
    ]"
    :style="{ height: heightClass, width: widthClass }"
  />
</template>
