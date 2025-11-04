<script setup lang="ts">
import { computed } from 'vue'
import { useDemandStore } from '@/stores/demand'

const demandStore = useDemandStore()

const STORAGE_KEY = 'jp-energy-data-mode' // Same as dataClient.ts

const isMock = computed(() => demandStore.dataMode === 'mock')

function setMode(mode: 'mock' | 'live') {
  // Skip if already in this mode
  if ((mode === 'mock' && isMock.value) || (mode === 'live' && !isMock.value)) {
    return
  }

  // Update localStorage with correct key
  localStorage.setItem(STORAGE_KEY, mode)

  // Reload page to apply new mode
  window.location.reload()
}
</script>

<template>
  <div class="inline-flex rounded-lg border border-gray-300 dark:border-gray-600 overflow-hidden">
    <!-- MOCK Button -->
    <button
      @click="setMode('mock')"
      class="px-4 py-2 text-sm font-medium transition-colors"
      :class="isMock
        ? 'bg-gray-600 text-white'
        : 'bg-white text-gray-700 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'"
    >
      MOCK
    </button>

    <!-- LIVE Button -->
    <button
      @click="setMode('live')"
      class="px-4 py-2 text-sm font-medium border-l border-gray-300 dark:border-gray-600 transition-colors"
      :class="!isMock
        ? 'bg-blue-600 text-white'
        : 'bg-white text-gray-700 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'"
    >
      LIVE
    </button>
  </div>
</template>
