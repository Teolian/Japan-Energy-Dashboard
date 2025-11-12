<script setup lang="ts">
import { ref } from 'vue'
import SectionCard from '@/components/common/SectionCard.vue'
import ComparisonChart from './ComparisonChart.vue'
import ComparisonAnalytics from './ComparisonAnalytics.vue'
import PriceHeatmap from './PriceHeatmap.vue'
import { ChevronDown, ChevronUp, BarChart3 } from 'lucide-vue-next'

const isExpanded = ref(false)

function toggleExpanded() {
  isExpanded.value = !isExpanded.value
}
</script>

<template>
  <div class="border-t border-gray-200 dark:border-gray-700 pt-8">
    <!-- Collapsible Header -->
    <button
      @click="toggleExpanded"
      class="w-full flex items-center justify-between mb-6 group hover:bg-gray-50 dark:hover:bg-gray-800 p-4 rounded-lg transition-colors"
    >
      <div class="flex items-center gap-3">
        <BarChart3 :size="24" class="text-gray-700 dark:text-gray-300" />
        <div class="text-left">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">
            Regional Comparison & Analytics
          </h2>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
            Deep dive into Tokyo vs Kansai patterns
          </p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-sm text-gray-500 dark:text-gray-400">
          {{ isExpanded ? 'Hide' : 'Show' }} Details
        </span>
        <ChevronDown v-if="!isExpanded" :size="20" class="text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-200 transition-colors" />
        <ChevronUp v-else :size="20" class="text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-200 transition-colors" />
      </div>
    </button>

    <!-- Collapsible Content -->
    <div
      v-if="isExpanded"
      class="space-y-6 animate-in fade-in duration-300"
    >
      <!-- Price Heatmap -->
      <SectionCard padding="md">
        <PriceHeatmap />
      </SectionCard>

      <!-- Comparison Chart & Analytics Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Comparison Chart (2 columns) -->
        <div class="lg:col-span-2">
          <SectionCard title="Multi-Area Overlay" padding="md">
            <ComparisonChart />
          </SectionCard>
        </div>

        <!-- Analytics (1 column) -->
        <div>
          <SectionCard padding="md">
            <ComparisonAnalytics />
          </SectionCard>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-in {
  animation: fade-in 0.3s ease-out;
}
</style>
