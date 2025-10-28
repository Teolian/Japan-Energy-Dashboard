<script setup lang="ts">
import LoadingSkeleton from './LoadingSkeleton.vue'
</script>

<template>
  <div class="space-y-4">
    <!-- Header Skeleton -->
    <div class="flex items-center justify-between mb-4">
      <LoadingSkeleton type="text" height="1.5rem" width="200px" />
      <LoadingSkeleton type="badge" height="2rem" width="120px" rounded />
    </div>

    <!-- Chart Area Skeleton -->
    <div class="h-80 flex items-end gap-2 px-4">
      <!-- Simulated bar chart -->
      <LoadingSkeleton v-for="i in 24" :key="i" type="card" :height="`${Math.random() * 60 + 40}%`" width="3%" />
    </div>

    <!-- Metrics Skeleton -->
    <div class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700 grid grid-cols-3 gap-4">
      <div v-for="i in 3" :key="i" class="space-y-2">
        <LoadingSkeleton type="text" height="0.875rem" width="80px" />
        <LoadingSkeleton type="text" height="1.5rem" width="120px" />
      </div>
    </div>
  </div>
</template>
