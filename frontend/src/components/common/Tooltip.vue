<script setup lang="ts">
interface Props {
  content: string
  position?: 'top' | 'bottom' | 'left' | 'right'
}

withDefaults(defineProps<Props>(), {
  position: 'bottom'
})
</script>

<template>
  <div class="relative inline-flex group">
    <!-- Trigger (slot) -->
    <slot />

    <!-- Tooltip -->
    <div
      :class="[
        'absolute z-50 invisible group-hover:visible opacity-0 group-hover:opacity-100 transition-opacity duration-200',
        'px-3 py-2 text-sm text-white bg-gray-900 dark:bg-gray-700 rounded-lg shadow-lg',
        'whitespace-nowrap pointer-events-none',
        position === 'top' && 'bottom-full left-1/2 -translate-x-1/2 mb-2',
        position === 'bottom' && 'top-full left-1/2 -translate-x-1/2 mt-2',
        position === 'left' && 'right-full top-1/2 -translate-y-1/2 mr-2',
        position === 'right' && 'left-full top-1/2 -translate-y-1/2 ml-2'
      ]"
    >
      {{ content }}

      <!-- Arrow -->
      <div
        :class="[
          'absolute w-2 h-2 bg-gray-900 dark:bg-gray-700 rotate-45',
          position === 'top' && 'bottom-[-4px] left-1/2 -translate-x-1/2',
          position === 'bottom' && 'top-[-4px] left-1/2 -translate-x-1/2',
          position === 'left' && 'right-[-4px] top-1/2 -translate-y-1/2',
          position === 'right' && 'left-[-4px] top-1/2 -translate-y-1/2'
        ]"
      />
    </div>
  </div>
</template>
