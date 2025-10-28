<script setup lang="ts">
import { HelpCircle } from 'lucide-vue-next'
import Tooltip from './Tooltip.vue'

interface Props {
  content: string
  position?: 'top' | 'bottom' | 'left' | 'right'
}

withDefaults(defineProps<Props>(), {
  position: 'top'
})
</script>

<template>
  <Tooltip :content="content" :position="position">
    <HelpCircle
      :size="14"
      class="text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 cursor-help transition-colors"
    />
  </Tooltip>
</template>
