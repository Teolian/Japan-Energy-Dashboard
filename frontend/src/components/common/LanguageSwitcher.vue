<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { setLocale, getCurrentLocale } from '@/i18n'
import { Languages } from 'lucide-vue-next'

const { locale } = useI18n()

const languages = [
  { code: 'en', label: 'English', flag: '🇬🇧' },
  { code: 'ja', label: '日本語', flag: '🇯🇵' }
]

const currentLanguage = () => {
  return languages.find(lang => lang.code === getCurrentLocale()) || languages[0]
}

const switchLanguage = (langCode: 'en' | 'ja') => {
  setLocale(langCode)
}
</script>

<template>
  <div class="relative group">
    <button
      class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
      :title="$t('header.language')"
    >
      <Languages :size="20" class="text-gray-700 dark:text-gray-300" />
      <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ currentLanguage().flag }}
      </span>
    </button>

    <!-- Dropdown -->
    <div
      class="absolute right-0 mt-2 w-40 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50"
    >
      <div class="py-1">
        <button
          v-for="lang in languages"
          :key="lang.code"
          @click="switchLanguage(lang.code as 'en' | 'ja')"
          :class="[
            'w-full px-4 py-2 text-left flex items-center gap-2 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors',
            locale === lang.code ? 'bg-blue-50 dark:bg-blue-900/20' : ''
          ]"
        >
          <span class="text-lg">{{ lang.flag }}</span>
          <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ lang.label }}
          </span>
          <span
            v-if="locale === lang.code"
            class="ml-auto text-blue-600 dark:text-blue-400 text-xs"
          >
            ✓
          </span>
        </button>
      </div>
    </div>
  </div>
</template>
