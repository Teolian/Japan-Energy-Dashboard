// i18n configuration
import { createI18n } from 'vue-i18n'
import en from './en'
import ja from './ja'

// Get locale from localStorage or default to Japanese (for demo)
const getStoredLocale = (): string => {
  const stored = localStorage.getItem('jp-energy-locale')
  return stored || 'ja' // Default to Japanese for demo
}

export type MessageSchema = typeof en

const i18n = createI18n<[MessageSchema], 'en' | 'ja'>({
  legacy: false, // Use Composition API
  locale: getStoredLocale(),
  fallbackLocale: 'en',
  messages: {
    en,
    ja
  }
})

export default i18n

// Helper to switch locale
export function setLocale(locale: 'en' | 'ja') {
  (i18n.global.locale as any).value = locale
  localStorage.setItem('jp-energy-locale', locale)
}

export function getCurrentLocale(): 'en' | 'ja' {
  return (i18n.global.locale as any).value as 'en' | 'ja'
}
