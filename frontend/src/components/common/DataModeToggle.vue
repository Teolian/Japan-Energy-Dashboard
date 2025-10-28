<script setup lang="ts">
import { computed } from 'vue'
import { useDemandStore } from '@/stores/demand'

const demandStore = useDemandStore()

const STORAGE_KEY = 'jp-energy-data-mode' // Same as dataClient.ts

const isMock = computed(() => demandStore.dataMode === 'mock')

function toggleMode() {
  const newMode = isMock.value ? 'live' : 'mock'

  // Update localStorage with correct key
  localStorage.setItem(STORAGE_KEY, newMode)

  // Reload page to apply new mode
  window.location.reload()
}
</script>

<template>
  <div class="flex items-center gap-3">
    <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
      MOCK
    </span>

    <!-- Toggle Switch -->
    <button
      @click="toggleMode"
      class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
      :class="isMock ? 'bg-gray-400' : 'bg-blue-600'"
      role="switch"
      :aria-checked="!isMock"
    >
      <span
        class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform"
        :class="isMock ? 'translate-x-1' : 'translate-x-6'"
      />
    </button>

    <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
      LIVE
    </span>
  </div>
</template>
