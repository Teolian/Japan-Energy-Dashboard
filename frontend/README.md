# Japan Energy Dashboard - Frontend

Vue 3 + TypeScript dashboard for monitoring Japanese electricity market data.

## Quick Start

```bash
npm install
npm run dev        # Start dev server (port 5173)
npm run build      # Build for production
npm run preview    # Preview production build
```

## Features

**Demand Monitoring**
- Real-time electricity demand for Tokyo (TEPCO) and Kansai regions
- Hourly demand curves with 24-hour visualization
- Capacity utilization percentages
- Peak demand detection and trend analysis

**JEPX Spot Price Analysis**
- Day-ahead market prices from Japan Electric Power Exchange
- Regional price comparison (Tokyo vs Kansai)
- 24-hour price curves with peak identification
- Historical price data navigation

**System Reserve Capacity**
- OCCTO (Organization for Cross-regional Coordination) reserve data
- Multi-region coverage across all 10 Japanese power areas
- Reserve percentage calculations
- Capacity margin monitoring

**Data Management**
- Mock/Live mode toggle with localStorage persistence
- Manual data refresh with progress tracking sidebar
- Real-time fetch status indicators for all data sources
- Date navigation controls

## Technology Stack

- Vue 3.5 - Composition API with `<script setup>`
- TypeScript 5.6 - Full type safety
- Pinia 2.2 - State management
- TailwindCSS 3.4 - Utility-first CSS
- Chart.js 4.4 - Interactive visualizations
- Vite 5.4 - Build tool with HMR
- Lucide Icons - Modern SVG icons

## Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── common/            # Reusable UI components
│   │   ├── demand/            # Demand visualization
│   │   ├── reserve/           # Reserve capacity display
│   │   └── comparison/        # Regional comparison
│   ├── stores/
│   │   ├── demand.ts          # Demand data state
│   │   ├── jepx.ts            # JEPX spot price state
│   │   └── reserve.ts         # Reserve capacity state
│   ├── services/
│   │   ├── dataClient.ts      # Data fetching with mode switching
│   │   └── weatherClient.ts   # Weather data integration
│   ├── types/                 # TypeScript interfaces
│   └── views/
│       └── Dashboard.vue      # Main dashboard page
├── public/data/jp/            # Mock data (static JSON files)
│   ├── tokyo/                 # TEPCO demand files
│   ├── kansai/                # Kansai demand files
│   ├── jepx/                  # JEPX spot price files
│   └── system/                # OCCTO reserve files
└── vercel.json                # Deployment configuration
```

## Data Modes

### Mock Mode (Default)
- Uses pre-generated JSON files from `public/data/jp/`
- No backend server required
- Available dates: 2025-10-23 to 2025-10-28
- Ideal for development and demonstrations
- Data stored in browser localStorage

### Live Mode
- Fetches real-time data from backend API
- Requires Go backend server running on port 8080
- Click "Refresh" button to fetch latest data
- Data cached to `public/data/jp/` after successful fetch
- Supports manual data refresh with progress tracking
- **Note:** the hosted Vercel deployment is static-only; to use the refresh workflow you must run the Go backend locally alongside `npm run dev`.

## State Management

### Pinia Stores

**Demand Store (stores/demand.ts)**
```typescript
const demandStore = useDemandStore()

// State
demandStore.tokyoData        // Tokyo demand response
demandStore.kansaiData       // Kansai demand response
demandStore.currentDate      // Selected date
demandStore.loading          // Loading state
demandStore.dataMode         // 'mock' | 'live'

// Computed
demandStore.tokyoMetrics     // Peak, average, forecast accuracy
demandStore.kansaiMetrics    // Peak, average, forecast accuracy
demandStore.tokyoChartData   // Chart.js formatted data
demandStore.kansaiChartData  // Chart.js formatted data

// Actions
demandStore.fetchDemandData(area, date)
demandStore.fetchAllDemandData(date)
demandStore.setDate(date)
demandStore.nextDay()
demandStore.prevDay()
```

**JEPX Store (stores/jepx.ts)**
```typescript
const jepxStore = useJEPXStore()

// State
jepxStore.tokyoData          // Tokyo spot prices
jepxStore.kansaiData         // Kansai spot prices
jepxStore.loading            // Loading state

// Computed
jepxStore.pricesForArea(area)
jepxStore.priceValues(area)
jepxStore.priceAtHour(area, hour)
jepxStore.hasWarning(area)

// Actions
jepxStore.fetchJEPXData(area, date)
```

**Reserve Store (stores/reserve.ts)**
```typescript
const reserveStore = useReserveStore()

// State
reserveStore.data            // Reserve response
reserveStore.loading         // Loading state

// Computed
reserveStore.reserveForArea(area)
reserveStore.statusConfig(status)
reserveStore.hasWarning
reserveStore.warning

// Actions
reserveStore.fetchReserveData(date)
reserveStore.reset()
```

## Data Formats

### Demand Response
```typescript
interface DemandResponse {
  date: string
  area: Area
  timezone: string
  timescale: "hourly"
  series: Array<{
    ts: string             // ISO 8601 timestamp
    demand_mw: number      // Actual demand in megawatts
    forecast_mw?: number   // Forecast (if available)
  }>
  source: {
    name: string
    url: string
    accessed_at: string
  }
  meta?: {
    peak_mw: number
    avg_mw: number
    capacity_mw: number
  }
}
```

### JEPX Spot Price Response
```typescript
interface JEPXResponse {
  date: string
  area: string
  timescale: string
  price_yen_per_kwh: Array<{
    ts: string        // ISO 8601 timestamp
    price: number     // Spot price in JPY/kWh
  }>
  source: {
    name: string
    url: string
    accessed_at: string
  }
  meta?: {
    min_price: number
    max_price: number
    avg_price: number
  }
}
```

### Reserve Capacity Response
```typescript
interface ReserveResponse {
  date: string
  areas: Array<{
    area: Area
    demand_mw: number
    capacity_mw: number
    reserve_mw: number
    reserve_percent: number
    status: "critical" | "low" | "adequate" | "good"
  }>
  source: {
    name: string
    url: string
  }
  meta?: {
    warning?: string
  }
}
```

## Build Output

```
dist/
├── index.html                   0.46 KB
├── assets/
│   ├── index-[hash].css        31.38 KB (gzipped: 5.47 KB)
│   └── index-[hash].js        347.94 KB (gzipped: 118.33 KB)
└── data/jp/                    30 JSON files (~500 KB)

Total transferred: ~600 KB
```

## Development

### Adding New Components

```vue
<!-- src/components/example/NewComponent.vue -->
<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  title: string
  value: number
}

const props = defineProps<Props>()

const formattedValue = computed(() => props.value.toFixed(2))
</script>

<template>
  <div class="p-4 bg-white rounded-lg shadow">
    <h3 class="text-lg font-semibold">{{ title }}</h3>
    <p class="text-2xl">{{ formattedValue }}</p>
  </div>
</template>
```

### Adding New Stores

```typescript
// src/stores/example.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useExampleStore = defineStore('example', () => {
  const data = ref<string[]>([])
  const loading = ref(false)

  async function fetchData() {
    loading.value = true
    try {
      const response = await fetch('/api/example')
      data.value = await response.json()
    } finally {
      loading.value = false
    }
  }

  return { data, loading, fetchData }
})
```

## Configuration

### Vite Configuration

```typescript
// vite.config.ts
export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',  // Backend API
        changeOrigin: true,
      },
    },
  },
})
```

### Environment Variables

Create `.env.local` for custom configuration:

```bash
# Backend API URL (optional, defaults to localhost:8080)
VITE_API_BASE_URL=http://localhost:8080
```

## Troubleshooting

### Build Errors

```bash
# Clear dependencies and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear Vite cache
rm -rf .vite node_modules/.vite
```

### Data Not Loading

**Mock Mode:**
- Verify `public/data/jp/` directory exists
- Check browser DevTools Network tab for 404 errors
- Ensure date is within range (2025-10-23 to 2025-10-28)

**Live Mode:**
- Confirm backend is running on port 8080
- Check backend logs for errors
- Test API endpoint: `curl http://localhost:8080/api/data/refresh`

## License

MIT
