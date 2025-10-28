<script setup lang="ts">
import { Database, AlertCircle, FileQuestion, Inbox } from 'lucide-vue-next'
import { computed } from 'vue'

interface Props {
  type?: 'no-data' | 'error' | 'not-found' | 'empty'
  title?: string
  message?: string
  icon?: any
}

const props = withDefaults(defineProps<Props>(), {
  type: 'no-data'
})

const defaultConfig = computed(() => {
  switch (props.type) {
    case 'error':
      return {
        icon: AlertCircle,
        title: 'Error Loading Data',
        message: 'Something went wrong while loading the data. Please try again later.',
        color: 'text-red-600 dark:text-red-400'
      }
    case 'not-found':
      return {
        icon: FileQuestion,
        title: 'Data Not Found',
        message: 'The requested data could not be found. Try selecting a different date.',
        color: 'text-yellow-600 dark:text-yellow-400'
      }
    case 'empty':
      return {
        icon: Inbox,
        title: 'No Data Available',
        message: 'There is no data to display at the moment.',
        color: 'text-gray-400 dark:text-gray-500'
      }
    default: // no-data
      return {
        icon: Database,
        title: 'No Data Available',
        message: 'Data for this date is not available yet. Please select another date.',
        color: 'text-blue-600 dark:text-blue-400'
      }
  }
})

const displayIcon = computed(() => props.icon || defaultConfig.value.icon)
const displayTitle = computed(() => props.title || defaultConfig.value.title)
const displayMessage = computed(() => props.message || defaultConfig.value.message)
const displayColor = computed(() => defaultConfig.value.color)
</script>

<template>
  <div class="flex flex-col items-center justify-center py-12 px-4 text-center">
    <component
      :is="displayIcon"
      :size="64"
      :class="['mb-4', displayColor]"
      stroke-width="1.5"
    />
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
      {{ displayTitle }}
    </h3>
    <p class="text-sm text-gray-600 dark:text-gray-400 max-w-md">
      {{ displayMessage }}
    </p>
    <slot />
  </div>
</template>
