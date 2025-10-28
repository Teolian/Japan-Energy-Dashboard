<script setup lang="ts">
import { computed } from 'vue'
import { Database, CheckCircle, AlertTriangle, Clock } from 'lucide-vue-next'
import { useDemandStore } from '@/stores/demand'
import { useJEPXStore } from '@/stores/jepx'
import { useReserveStore } from '@/stores/reserve'
import Tooltip from '@/components/common/Tooltip.vue'

const demandStore = useDemandStore()
const jepxStore = useJEPXStore()
const reserveStore = useReserveStore()

// Check data freshness (how old is the data)
const dataAge = computed(() => {
  const currentDate = new Date(demandStore.currentDate)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  currentDate.setHours(0, 0, 0, 0)

  const diffDays = Math.floor((today.getTime() - currentDate.getTime()) / (1000 * 60 * 60 * 24))

  if (diffDays === 0) return { label: 'Today', color: 'text-green-600 dark:text-green-400', status: 'fresh' }
  if (diffDays === 1) return { label: 'Yesterday', color: 'text-blue-600 dark:text-blue-400', status: 'recent' }
  if (diffDays <= 7) return { label: `${diffDays} days ago`, color: 'text-yellow-600 dark:text-yellow-400', status: 'aging' }
  return { label: `${diffDays} days ago`, color: 'text-red-600 dark:text-red-400', status: 'stale' }
})

// Check data completeness (all sources loaded)
const dataCompleteness = computed(() => {
  const sources = []

  if (demandStore.tokyoData) sources.push('Tokyo Demand')
  if (demandStore.kansaiData) sources.push('Kansai Demand')
  if (reserveStore.data) sources.push('Reserve Margin')
  if (jepxStore.tokyoData) sources.push('Tokyo Prices')
  if (jepxStore.kansaiData) sources.push('Kansai Prices')

  const total = 5
  const loaded = sources.length
  const percentage = (loaded / total) * 100

  return {
    loaded,
    total,
    percentage,
    sources,
    isComplete: loaded === total,
    isMostly: loaded >= 4,
    label: loaded === total ? 'Complete' : `${loaded}/${total} sources`
  }
})

const dataMode = computed(() => {
  return demandStore.dataMode.toUpperCase()
})
</script>

<template>
  <div class="flex items-center gap-6 text-sm">
    <!-- Data Mode -->
    <div class="flex items-center gap-2">
      <Database :size="16" class="text-gray-500 dark:text-gray-400" />
      <Tooltip :content="`Data source mode: ${dataMode}`">
        <span class="text-gray-600 dark:text-gray-400 cursor-help">
          {{ dataMode }} mode
        </span>
      </Tooltip>
    </div>

    <!-- Data Freshness -->
    <div class="flex items-center gap-2">
      <Clock :size="16" :class="dataAge.color" />
      <Tooltip
        :content="`Data is ${dataAge.label.toLowerCase()}${dataAge.status === 'stale' ? ' - Consider updating' : ''}`"
      >
        <span :class="['cursor-help', dataAge.color]">
          {{ dataAge.label }}
        </span>
      </Tooltip>
    </div>

    <!-- Data Completeness -->
    <div class="flex items-center gap-2">
      <component
        :is="dataCompleteness.isComplete ? CheckCircle : AlertTriangle"
        :size="16"
        :class="dataCompleteness.isComplete ? 'text-green-600 dark:text-green-400' : 'text-yellow-600 dark:text-yellow-400'"
      />
      <Tooltip
        :content="`Loaded: ${dataCompleteness.sources.join(', ')}${!dataCompleteness.isComplete ? ' - Some data missing' : ''}`"
      >
        <span
          :class="[
            'cursor-help',
            dataCompleteness.isComplete ? 'text-green-600 dark:text-green-400' : 'text-yellow-600 dark:text-yellow-400'
          ]"
        >
          {{ dataCompleteness.label }}
        </span>
      </Tooltip>
    </div>
  </div>
</template>
