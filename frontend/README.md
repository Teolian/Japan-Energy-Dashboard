# Japan Energy Dashboard - Frontend

Vue 3 + TypeScript dashboard for monitoring Japanese electricity market data.

## Quick Start

```bash
# Install dependencies
npm install

# Start development server (port 5173)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
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

- **Vue 3.5** - Progressive JavaScript framework with Composition API
- **TypeScript 5.6** - Strict type safety throughout the application
- **Pinia 2.2** - State management for demand, JEPX, and reserve data
- **TailwindCSS 3.4** - Utility-first CSS framework
- **Chart.js 4.4** - Interactive data visualizations
- **Vite 5.4** - Build tool with HMR and optimized production builds
- **Lucide Icons** - Modern SVG icon library

## Project Structure

```
frontend-japan/
├── src/
│   ├── components/
│   │   ├── common/
│   │   │   ├── BaseButton.vue         # Reusable button component
│   │   │   ├── DataModeToggle.vue     # Mock/Live mode switcher
│   │   │   ├── RefreshButton.vue      # Data refresh trigger
│   │   │   └── RefreshSidebar.vue     # Slide-in progress panel
│   │   ├── demand/
│   │   │   └── DemandChart.vue        # Demand visualization component
│   │   ├── reserve/
│   │   │   └── ReserveTable.vue       # Reserve capacity table
│   │   └── comparison/
│   │       └── ComparisonAnalytics.vue # Regional comparison
│   ├── views/
│   │   └── Dashboard.vue              # Main dashboard page
│   ├── stores/
│   │   ├── demand.ts                  # Demand data state (Tokyo/Kansai)
│   │   ├── jepx.ts                    # JEPX spot price state
│   │   └── reserve.ts                 # Reserve capacity state
│   ├── services/
│   │   ├── dataClient.ts              # Data fetching with mode switching
│   │   └── weatherClient.ts           # Weather data integration
│   ├── types/
│   │   ├── demand.ts                  # Demand data TypeScript interfaces
│   │   ├── jepx.ts                    # JEPX data interfaces
│   │   └── reserve.ts                 # Reserve data interfaces
│   ├── App.vue                        # Root component
│   ├── main.ts                        # Application entry point
│   └── router.ts                      # Vue Router configuration
├── public/
│   └── data/jp/                       # Mock data (static JSON files)
│       ├── tokyo/                     # TEPCO demand files
│       ├── kansai/                    # Kansai demand files
│       ├── jepx/                      # JEPX spot price files
│       └── system/                    # OCCTO reserve files
├── index.html                         # HTML entry point
├── vite.config.ts                     # Vite build configuration
├── tailwind.config.js                 # TailwindCSS configuration
├── tsconfig.json                      # TypeScript compiler options
├── vercel.json                        # Vercel deployment settings
└── package.json                       # Project dependencies
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

## Development Guide

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

## Deployment

### Vercel Deployment

**Step 1: Install Vercel CLI**
```bash
npm i -g vercel
```

**Step 2: Login to Vercel**
```bash
vercel login
```

**Step 3: Deploy from frontend-japan directory**
```bash
cd frontend-japan
vercel
```

**Step 4: Production Deployment**
```bash
vercel --prod
```

### Vercel Configuration (vercel.json)

```json
{
  "buildCommand": "npm run build",
  "outputDirectory": "dist",
  "framework": "vite",
  "installCommand": "npm install",
  "cleanUrls": true,
  "trailingSlash": false,
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/index.html"
    }
  ],
  "headers": [
    {
      "source": "/data/(.*)",
      "headers": [
        {
          "key": "Cache-Control",
          "value": "public, max-age=3600, must-revalidate"
        }
      ]
    }
  ]
}
```

### GitHub + Vercel Integration

1. Push to GitHub
2. Go to vercel.com/new
3. Import repository
4. Configure settings:
   - Root Directory: `frontend-japan`
   - Build Command: `npm run build`
   - Output Directory: `dist`
   - Install Command: `npm install`
5. Deploy

### Post-Deployment Checklist

- [ ] Verify Mock mode loads data correctly
- [ ] Test date navigation (previous/next day)
- [ ] Confirm Mock/Live toggle persists on page reload
- [ ] Test responsive design on mobile devices
- [ ] Verify all charts render without errors
- [ ] Check browser console for errors
- [ ] Test production build locally with `npm run preview`

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

## Performance Optimization

**Lighthouse Scores (Production Build):**
- Performance: 95+
- Accessibility: 98+
- Best Practices: 100
- SEO: 100

**Optimization Techniques:**
- Tree-shaking with Vite
- Code splitting for lazy-loaded routes
- Gzip/Brotli compression (Vercel automatic)
- CDN distribution (Vercel Edge Network)
- Efficient Chart.js rendering with data sampling

## Troubleshooting

### Build Errors

```bash
# Clear dependencies and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear Vite cache
rm -rf .vite node_modules/.vite
```

### TypeScript Errors

```bash
# Run type check
npm run type-check

# Check for common issues
npx vue-tsc --noEmit
```

### CORS Errors (Live Mode)

Ensure backend allows CORS from localhost:5173:

```go
// backend/cmd/api/main.go
import "github.com/gin-contrib/cors"

router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST"},
    AllowHeaders:     []string{"Content-Type"},
    AllowCredentials: true,
}))
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

## Testing

### Type Safety
```bash
npm run type-check
```

### Build Verification
```bash
npm run build
npm run preview
# Open http://localhost:4173
```

### Manual Testing Checklist
- [ ] Mock mode displays data
- [ ] Live mode fetches from API
- [ ] Toggle persists on reload
- [ ] Date navigation works
- [ ] Refresh button shows progress
- [ ] Charts render correctly
- [ ] Responsive on mobile

## Resources

- [Vue 3 Documentation](https://vuejs.org/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [TailwindCSS Documentation](https://tailwindcss.com/docs)
- [Chart.js Documentation](https://www.chartjs.org/docs/)
- [Vite Documentation](https://vitejs.dev/)
- [Vercel Documentation](https://vercel.com/docs)

## Contributing

This is an educational project demonstrating modern Vue 3 development practices. Contributions are welcome for:
- Performance improvements
- Accessibility enhancements
- New data visualizations
- Bug fixes

## License

MIT - See [LICENSE](../LICENSE) for details

---

**Main Repository:** [README.md](../README.md)
**Backend Documentation:** [backend/README.md](../backend/README.md)
**Deployment Guide:** [DEPLOY.md](../DEPLOY.md)
