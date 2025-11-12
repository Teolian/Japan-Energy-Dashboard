import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '@/views/Dashboard.vue'
import TradingIntelligence from '@/views/TradingIntelligence.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: Dashboard
    },
    {
      path: '/trading',
      name: 'trading',
      component: TradingIntelligence
    }
  ]
})

export default router
