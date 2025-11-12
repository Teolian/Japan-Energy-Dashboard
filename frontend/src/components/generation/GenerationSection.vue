<script setup lang="ts">
import SectionCard from '@/components/common/SectionCard.vue'
import GenerationMixChart from './GenerationMixChart.vue'
import CarbonIntensityPanel from './CarbonIntensityPanel.vue'
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Generation Mix Chart (2 columns) -->
    <div class="lg:col-span-2">
      <SectionCard padding="md">
        <GenerationMixChart />
      </SectionCard>
    </div>

    <!-- Carbon Intensity Panel (1 column) -->
    <div>
      <SectionCard padding="md">
        <CarbonIntensityPanel />
      </SectionCard>
    </div>
  </div>
</template>
