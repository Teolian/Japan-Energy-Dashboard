<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { LayoutDashboard, Brain } from 'lucide-vue-next'

const { t } = useI18n()
const route = useRoute()

const navItems = [
  { path: '/', name: 'dashboard', label: t('nav.dashboard'), icon: LayoutDashboard },
  { path: '/trading', name: 'trading', label: t('nav.trading'), icon: Brain, badge: 'AI' }
]

const isActive = (path: string) => {
  return route.path === path
}
</script>

<template>
  <nav class="mb-8 border-b border-gray-200 dark:border-gray-700">
    <div class="flex gap-1">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        :class="[
          'flex items-center gap-2 px-6 py-3 border-b-2 transition-all font-medium',
          isActive(item.path)
            ? 'border-blue-600 text-blue-600 dark:text-blue-400'
            : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:border-gray-300 dark:hover:border-gray-600'
        ]"
      >
        <component :is="item.icon" :size="20" />
        <span>{{ item.label }}</span>
        <span
          v-if="item.badge"
          class="ml-1 px-2 py-0.5 text-xs font-bold bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-full"
        >
          {{ item.badge }}
        </span>
      </router-link>
    </div>
  </nav>
</template>
