<script setup lang="ts">
import { ref } from 'vue'
import { RefreshCw } from 'lucide-vue-next'
import { useDemandStore } from '@/stores/demand'
import { useJEPXStore } from '@/stores/jepx'
import { useReserveStore } from '@/stores/reserve'
import RefreshSidebar from './RefreshSidebar.vue'

const demandStore = useDemandStore()
const jepxStore = useJEPXStore()
const reserveStore = useReserveStore()

const isRefreshing = ref(false)
const isSidebarOpen = ref(false)

interface DataSource {
  source: string
  status: 'pending' | 'loading' | 'success' | 'error'
  file_path?: string
  error?: string
  duration?: string
}

const sources = ref<DataSource[]>([])

interface RefreshResponse {
  success: boolean
  message: string
  results: Array<{
    source: string
    status: string
    file_path?: string
    error?: string
    duration: string
  }>
}

function initializeSources() {
  sources.value = [
    { source: 'tokyo-demand', status: 'pending' },
    { source: 'tokyo-jepx', status: 'pending' },
    { source: 'kansai-demand', status: 'pending' },
    { source: 'kansai-jepx', status: 'pending' },
    { source: 'reserve', status: 'pending' }
  ]
}

async function refreshData() {
  isRefreshing.value = true
  isSidebarOpen.value = true

  // Initialize sources
  initializeSources()

  // Set all to loading
  sources.value = sources.value.map(s => ({ ...s, status: 'loading' }))

  try {
    const response = await fetch('http://localhost:8080/api/data/refresh', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        date: demandStore.currentDate,
        areas: ['tokyo', 'kansai']
      })
    })

    const data: RefreshResponse = await response.json()

    // Update each source
    data.results.forEach(result => {
      const sourceIndex = sources.value.findIndex(s => s.source === result.source)
      if (sourceIndex !== -1) {
        sources.value[sourceIndex] = {
          source: result.source,
          status: result.status === 'success' ? 'success' : 'error',
          file_path: result.file_path,
          error: result.error,
          duration: result.duration
        }
      }
    })

    // Reload stores if any succeeded
    const hasSuccess = sources.value.some(s => s.status === 'success')
    if (hasSuccess) {
      try {
        await Promise.all([
          demandStore.fetchAllDemandData(demandStore.currentDate),
          jepxStore.fetchJEPXData('tokyo', demandStore.currentDate),
          jepxStore.fetchJEPXData('kansai', demandStore.currentDate),
          reserveStore.fetchReserveData(demandStore.currentDate)
        ])
      } catch (err) {
        console.warn('[RefreshButton] Some stores failed to reload:', err)
      }
    }
  } catch (error) {
    console.error('[RefreshButton] Failed to refresh data:', error)

    // Mark all as error
    sources.value = sources.value.map(s => ({
      ...s,
      status: 'error',
      error: error instanceof Error ? error.message : 'Network error'
    }))
  } finally {
    isRefreshing.value = false
  }
}

function closeSidebar() {
  isSidebarOpen.value = false
}
</script>

<template>
  <div>
    <!-- Compact Refresh Button -->
    <button
      @click="refreshData"
      :disabled="isRefreshing"
      class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-all
             bg-blue-600 hover:bg-blue-700 text-white
             disabled:opacity-50 disabled:cursor-not-allowed"
      :class="{'animate-pulse': isRefreshing}"
    >
      <RefreshCw
        :size="16"
        :class="{'animate-spin': isRefreshing}"
      />
      <span v-if="isRefreshing">Refreshing...</span>
      <span v-else>Refresh</span>
    </button>

    <!-- Slide-in Sidebar -->
    <RefreshSidebar
      :is-open="isSidebarOpen"
      :sources="sources"
      :is-refreshing="isRefreshing"
      @close="closeSidebar"
    />
  </div>
</template>
