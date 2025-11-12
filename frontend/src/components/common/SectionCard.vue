<script setup lang="ts">
interface Props {
  title?: string
  subtitle?: string
  padding?: 'sm' | 'md' | 'lg' | 'none'
  noBorder?: boolean
  transparent?: boolean
}

withDefaults(defineProps<Props>(), {
  padding: 'md',
  noBorder: false,
  transparent: false
})

const paddingClasses = {
  none: '',
  sm: 'p-4',
  md: 'p-6',
  lg: 'p-8'
}
</script>

<template>
  <div
    :class="[
      'rounded-lg',
      transparent ? 'bg-transparent' : 'bg-white dark:bg-gray-800',
      noBorder ? '' : 'border border-gray-200 dark:border-gray-700 shadow-sm',
      paddingClasses[padding]
    ]"
  >
    <!-- Header -->
    <div v-if="title || $slots.header" class="mb-4">
      <slot name="header">
        <div v-if="title" class="flex items-center justify-between">
          <div>
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ title }}
            </h2>
            <p v-if="subtitle" class="text-xs text-gray-500 dark:text-gray-400 mt-1">
              {{ subtitle }}
            </p>
          </div>
          <slot name="actions" />
        </div>
      </slot>
    </div>

    <!-- Content -->
    <slot />

    <!-- Footer -->
    <div v-if="$slots.footer" class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
      <slot name="footer" />
    </div>
  </div>
</template>
