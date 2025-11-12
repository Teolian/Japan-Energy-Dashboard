<script setup lang="ts">
import { computed } from 'vue'
import ChartLoadingSkeleton from './ChartLoadingSkeleton.vue'
import EmptyState from './EmptyState.vue'

interface Props {
  loading?: boolean
  error?: string | null
  empty?: boolean
  emptyMessage?: string
  height?: string
  title?: string
  subtitle?: string
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  error: null,
  empty: false,
  emptyMessage: 'No data available',
  height: '20rem' // h-80 equivalent
})

const heightStyle = computed(() => ({
  height: props.height
}))
</script>

<template>
  <div class="space-y-3">
    <!-- Header -->
    <div v-if="title || $slots.header" class="flex items-center justify-between">
      <slot name="header">
        <div v-if="title">
          <h3 class="text-base font-semibold text-gray-900 dark:text-white">
            {{ title }}
          </h3>
          <p v-if="subtitle" class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
            {{ subtitle }}
          </p>
        </div>
      </slot>
      <slot name="actions" />
    </div>

    <!-- Chart Container -->
    <div
      :style="heightStyle"
      class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4"
    >
      <!-- Loading State -->
      <div v-if="loading" class="h-full">
        <ChartLoadingSkeleton />
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="h-full">
        <EmptyState type="error" :message="error" />
      </div>

      <!-- Empty State -->
      <div v-else-if="empty" class="h-full flex items-center justify-center">
        <p class="text-sm text-gray-500 dark:text-gray-400">
          {{ emptyMessage }}
        </p>
      </div>

      <!-- Chart Content -->
      <div v-else class="h-full">
        <slot />
      </div>
    </div>

    <!-- Footer -->
    <div v-if="$slots.footer">
      <slot name="footer" />
    </div>
  </div>
</template>
