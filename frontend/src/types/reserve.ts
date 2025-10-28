// TypeScript types for OCCTO reserve margin data (AGENT_TECH_SPEC §3.2)

import type { Area } from './demand'

export type ReserveStatus = 'stable' | 'watch' | 'tight'

export interface Source {
  name: string
  url: string
}

export interface Meta {
  warning?: string
}

export interface AreaReserve {
  area: Area
  reserve_margin_pct: number
  status: ReserveStatus
}

// Backend JSON response structure
export interface ReserveResponse {
  date: string // YYYY-MM-DD
  areas: AreaReserve[]
  source: Source
  meta?: Meta
}

// UI helpers for status colors
export interface StatusConfig {
  color: string
  bgColor: string
  label: string
  icon: string
}

export const STATUS_CONFIGS: Record<ReserveStatus, StatusConfig> = {
  stable: {
    color: 'text-green-600 dark:text-green-400',
    bgColor: 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800',
    label: 'Stable',
    icon: '✓'
  },
  watch: {
    color: 'text-yellow-600 dark:text-yellow-400',
    bgColor: 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800',
    label: 'Watch',
    icon: '⚠'
  },
  tight: {
    color: 'text-red-600 dark:text-red-400',
    bgColor: 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800',
    label: 'Tight',
    icon: '⚡'
  }
}
