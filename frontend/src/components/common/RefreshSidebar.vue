<script setup lang="ts">
import { computed } from 'vue'
import { Check, X, Loader2, Clock, ChevronRight } from 'lucide-vue-next'

interface DataSource {
  source: string
  status: 'pending' | 'loading' | 'success' | 'error'
  file_path?: string
  error?: string
  duration?: string
}

const props = defineProps<{
  isOpen: boolean
  sources: DataSource[]
  isRefreshing: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const progressPercent = computed(() => {
  const completed = props.sources.filter(s => s.status === 'success' || s.status === 'error').length
  const total = props.sources.length
  return total > 0 ? (completed / total) * 100 : 0
})

const successCount = computed(() => props.sources.filter(s => s.status === 'success').length)
const errorCount = computed(() => props.sources.filter(s => s.status === 'error').length)

function getSourceName(source: string): string {
  const names: Record<string, string> = {
    'tokyo-demand': 'Tokyo Demand',
    'kansai-demand': 'Kansai Demand',
    'tokyo-jepx': 'Tokyo JEPX',
    'kansai-jepx': 'Kansai JEPX',
    'reserve': 'Reserve Margin'
  }
  return names[source] || source
}

function getSourceIcon(source: string): string {
  const icons: Record<string, string> = {
    'tokyo-demand': '🗼',
    'kansai-demand': '🏯',
    'tokyo-jepx': '⚡',
    'kansai-jepx': '⚡',
    'reserve': '🔋'
  }
  return icons[source] || '📊'
}
</script>

<template>
  <!-- Backdrop -->
  <Transition
    enter-active-class="transition-opacity duration-300"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition-opacity duration-300"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isOpen"
      class="fixed inset-0 bg-black/20 z-40"
      @click="emit('close')"
    />
  </Transition>

  <!-- Slide-in Panel -->
  <Transition
    enter-active-class="transition-transform duration-300 ease-out"
    enter-from-class="translate-x-full"
    enter-to-class="translate-x-0"
    leave-active-class="transition-transform duration-300 ease-in"
    leave-from-class="translate-x-0"
    leave-to-class="translate-x-full"
  >
    <div
      v-if="isOpen"
      class="fixed right-0 top-0 h-full w-96 bg-white dark:bg-gray-800 shadow-2xl z-50 flex flex-col"
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ isRefreshing ? 'Refreshing Data...' : 'Refresh Complete' }}
          </h3>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {{ successCount }}/{{ sources.length }} sources successful
          </p>
        </div>
        <button
          @click="emit('close')"
          class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        >
          <ChevronRight :size="20" class="text-gray-500" />
        </button>
      </div>

      <!-- Progress Bar -->
      <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
        <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
          <div
            class="h-full transition-all duration-500 ease-out"
            :class="{
              'bg-blue-600': isRefreshing,
              'bg-green-600': !isRefreshing && errorCount === 0,
              'bg-yellow-600': !isRefreshing && errorCount > 0 && successCount > 0,
              'bg-red-600': !isRefreshing && errorCount > 0 && successCount === 0
            }"
            :style="{ width: `${progressPercent}%` }"
          />
        </div>
        <div class="flex justify-between items-center mt-2">
          <span class="text-xs font-medium text-gray-600 dark:text-gray-400">
            Progress
          </span>
          <span class="text-xs font-semibold text-gray-900 dark:text-white">
            {{ progressPercent.toFixed(0) }}%
          </span>
        </div>
      </div>

      <!-- Sources List -->
      <div class="flex-1 overflow-y-auto p-6 space-y-3">
        <div
          v-for="source in sources"
          :key="source.source"
          class="p-4 rounded-lg border transition-all"
          :class="{
            'border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50': source.status === 'pending',
            'border-blue-300 dark:border-blue-700 bg-blue-50 dark:bg-blue-900/20': source.status === 'loading',
            'border-green-300 dark:border-green-700 bg-green-50 dark:bg-green-900/20': source.status === 'success',
            'border-red-300 dark:border-red-700 bg-red-50 dark:bg-red-900/20': source.status === 'error'
          }"
        >
          <!-- Header -->
          <div class="flex items-center gap-3 mb-2">
            <!-- Checkbox -->
            <div
              class="w-5 h-5 rounded flex items-center justify-center flex-shrink-0"
              :class="{
                'border-2 border-gray-300 dark:border-gray-600': source.status === 'pending',
                'bg-blue-500': source.status === 'loading',
                'bg-green-500': source.status === 'success',
                'bg-red-500': source.status === 'error'
              }"
            >
              <Loader2
                v-if="source.status === 'loading'"
                :size="14"
                class="text-white animate-spin"
              />
              <Check
                v-else-if="source.status === 'success'"
                :size="14"
                class="text-white"
              />
              <X
                v-else-if="source.status === 'error'"
                :size="14"
                class="text-white"
              />
            </div>

            <!-- Name -->
            <div class="flex-1">
              <div class="flex items-center gap-2">
                <span>{{ getSourceIcon(source.source) }}</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">
                  {{ getSourceName(source.source) }}
                </span>
              </div>
            </div>

            <!-- Duration -->
            <div
              v-if="source.duration"
              class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400"
            >
              <Clock :size="12" />
              <span>{{ source.duration }}</span>
            </div>
          </div>

          <!-- File Path -->
          <div
            v-if="source.status === 'success' && source.file_path"
            class="text-xs text-gray-600 dark:text-gray-400 font-mono bg-gray-100 dark:bg-gray-800 p-2 rounded break-all"
          >
            {{ source.file_path }}
          </div>

          <!-- Error -->
          <div v-if="source.status === 'error' && source.error">
            <details class="text-xs">
              <summary class="cursor-pointer text-red-700 dark:text-red-300 hover:underline mb-1">
                Show error details
              </summary>
              <pre class="mt-2 p-2 bg-red-100 dark:bg-red-900/30 rounded text-xs overflow-x-auto whitespace-pre-wrap break-all">{{ source.error }}</pre>
            </details>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="p-6 border-t border-gray-200 dark:border-gray-700">
        <button
          v-if="!isRefreshing"
          @click="emit('close')"
          class="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
        >
          Close
        </button>
        <div v-else class="text-center text-sm text-gray-500 dark:text-gray-400">
          Please wait...
        </div>
      </div>
    </div>
  </Transition>
</template>
