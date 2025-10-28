# Japan Energy Dashboard

Real-time monitoring and analysis of Japanese electricity market data.

## Overview

Professional dashboard for monitoring Japan's electricity market, providing hourly demand data, spot pricing, and system reserve capacity analysis across major power regions (Tokyo/TEPCO and Kansai).

## Features

**Demand Monitoring**
- Hourly electricity demand for Tokyo (TEPCO) and Kansai regions
- Real-time capacity utilization percentages
- Peak demand detection and trend analysis
- Historical comparison with previous periods

**Spot Price Analysis**
- JEPX (Japan Electric Power Exchange) day-ahead market prices in JPY/kWh
- 24-hour price curves with peak detection
- Regional price differentials (Tokyo vs Kansai)
- Price trend indicators

**Reserve Capacity**
- OCCTO (Organization for Cross-regional Coordination) system-wide reserve monitoring
- Coverage of all 10 Japanese power regions
- Reserve percentage calculations
- Critical threshold warnings

**Data Management**
- Mock/Live mode toggle with localStorage persistence
- Manual data refresh with progress tracking sidebar
- Real-time fetch status indicators for all data sources
- Date navigation controls

## Technology Stack

**Frontend** (frontend/)
- Vue 3.5 with TypeScript 5.6 and Composition API
- Pinia state management for demand, pricing, and reserve data
- Chart.js for interactive data visualizations
- TailwindCSS for responsive UI
- Vite build system with hot module replacement

**Backend** (backend/)
- Go 1.23 with Gin framework
- HTTP adapters for TEPCO, Kansai Electric, JEPX, and OCCTO APIs
- Custom HTTP client with retry logic and circuit breaker pattern
- Browser-like headers to bypass bot detection
- JSON file caching for offline mode

## Data Sources

**Live Mode:**
- TEPCO (Tokyo Electric Power Company) - Kanto region demand data
- Kansai Electric Power Company - Kansai region demand data
- JEPX (Japan Electric Power Exchange) - Day-ahead spot market prices
- OCCTO (Organization for Cross-regional Coordination) - System reserve capacity

**Mock Mode:**
- Pre-generated sample data in frontend/public/data/jp/
- Covers dates 2025-10-23 to 2025-10-28
- 30+ JSON files with realistic hourly patterns

## Project Structure

```
japan-energy-dashboard/
├── frontend/                    # Vue 3 dashboard application
│   ├── src/
│   │   ├── components/
│   │   │   ├── common/          # Shared UI components
│   │   │   ├── demand/          # Demand visualization
│   │   │   └── reserve/         # Reserve capacity display
│   │   ├── stores/              # Pinia state management
│   │   ├── services/            # Data fetching logic
│   │   ├── types/               # TypeScript interfaces
│   │   └── views/               # Dashboard page
│   ├── public/data/jp/          # Mock data (JSON files)
│   └── vercel.json              # Vercel deployment config
│
├── backend/                     # Go API server
│   ├── cmd/
│   │   ├── api/                 # REST API server (port 8080)
│   │   ├── fetch-demand-http/   # TEPCO/Kansai fetcher CLI
│   │   ├── fetch-jepx/          # JEPX fetcher CLI
│   │   └── run-settlement/      # Settlement calculator
│   ├── internal/
│   │   ├── adapters/            # Data source adapters
│   │   │   ├── tepco.go         # TEPCO demand parser
│   │   │   ├── kansai.go        # Kansai demand parser
│   │   │   ├── jepx.go          # JEPX CSV parser
│   │   │   └── occto.go         # OCCTO reserve parser
│   │   ├── jepx/                # JEPX domain models
│   │   └── settlement/          # Market settlement logic
│   ├── pkg/
│   │   ├── http/                # Custom HTTP client
│   │   └── sources/             # Data source configs
│   └── docs/                    # Technical documentation
│
└── docs/
    ├── PROJECT_ROADMAP.md       # Development roadmap
    ├── AGENT_TECH_SPEC.md       # Technical specifications
    └── AGENT_BIZ_SPEC.md        # Business requirements
```

## Quick Start

### Frontend Only (Mock Mode)

```bash
cd frontend
npm install
npm run dev
# Open http://localhost:5173
```

### Full Stack (Live Mode)

```bash
# Terminal 1: Start backend API
cd backend
PORT=8080 go run cmd/api/main.go

# Terminal 2: Start frontend
cd frontend
npm run dev
```

### Build for Production

```bash
cd frontend
npm run build
# Output: dist/ (348 KB JS + 31 KB CSS gzipped)
```

## Data Formats

### Demand Response
```json
{
  "date": "2025-10-24",
  "area": "tokyo",
  "timescale": "hourly",
  "demand_mw": [
    { "ts": "2025-10-24T00:00:00+09:00", "demand": 28500 }
  ],
  "capacity_mw": 45000,
  "source": { "name": "TEPCO", "url": "https://www.tepco.co.jp" }
}
```

### JEPX Spot Prices
```json
{
  "date": "2025-10-24",
  "area": "tokyo",
  "price_yen_per_kwh": [
    { "ts": "2025-10-24T00:00:00+09:00", "price": 24.32 }
  ],
  "source": { "name": "JEPX", "url": "https://www.jepx.jp" }
}
```

### Reserve Capacity
```json
{
  "date": "2025-10-24",
  "reserves": [
    {
      "region": "hokkaido",
      "demand_mw": 3500,
      "capacity_mw": 5000,
      "reserve_mw": 1500,
      "reserve_percent": 42.86
    }
  ]
}
```

## Deployment

### Vercel (Frontend Only)

```bash
cd frontend
npm i -g vercel
vercel login
vercel --prod
```

**Configuration:**
- Root Directory: `frontend`
- Build Command: `npm run build`
- Output Directory: `dist`
- Framework: Vite

See [frontend/README.md](frontend/README.md) for detailed deployment instructions.

## API Endpoints

**Backend API (http://localhost:8080):**

```
POST /api/data/refresh
Body: { "date": "2025-10-24", "areas": ["tokyo", "kansai"] }
Response: { "success": true, "results": [...] }
```

Triggers data fetch from all sources and saves to frontend/public/data/jp/.

## State Management

### Pinia Stores

```typescript
// Demand Store (stores/demand.ts)
const demandStore = useDemandStore()
demandStore.tokyoData        // Tokyo demand response
demandStore.kansaiData       // Kansai demand response
demandStore.currentDate      // Selected date
demandStore.fetchDemandData(area, date)

// JEPX Store (stores/jepx.ts)
const jepxStore = useJEPXStore()
jepxStore.tokyoData          // Tokyo spot prices
jepxStore.kansaiData         // Kansai spot prices
jepxStore.fetchJEPXData(area, date)

// Reserve Store (stores/reserve.ts)
const reserveStore = useReserveStore()
reserveStore.data            // Reserve response
reserveStore.fetchReserveData(date)
```

## HTTP Adapters

Custom Go HTTP client with:
- Retry logic with exponential backoff
- Circuit breaker pattern for repeated failures
- Browser-like headers (User-Agent, Accept-Language, etc.)
- Connection pooling and timeout handling

```go
// pkg/http/fetcher.go
fetcher := http.NewFetcher(http.BrowserConfig())
body, err := fetcher.Fetch(url)
```

## Performance

**Build Output:**
- HTML: 0.46 KB
- CSS: 31.38 KB (gzipped: 5.47 KB)
- JS: 347.94 KB (gzipped: 118.33 KB)
- Data: 30 JSON files (~500 KB total)
- Total: ~600 KB transferred

**Lighthouse Score (production build):**
- Performance: 95+
- Accessibility: 98+
- Best Practices: 100
- SEO: 100

## Testing

```bash
# Frontend type check
cd frontend
npm run type-check

# Build verification
npm run build
npm run preview

# Backend tests
cd backend
go test ./internal/adapters/...
go test ./pkg/http/...
```

## Documentation

- [Frontend README](frontend/README.md) - Detailed frontend documentation
- [Backend Docs](backend/) - API and adapter documentation
- [Project Roadmap](backend/PROJECT_ROADMAP.md) - Development plan
- [Tech Spec](backend/AGENT_TECH_SPEC.md) - Technical specifications

## Contributing

This is an educational project demonstrating:
- Modern Vue 3 patterns (Composition API, TypeScript)
- Go backend with HTTP client best practices
- Real-world data fetching with retry logic
- Clean architecture and separation of concerns

## License

MIT - See [LICENSE](LICENSE) for details

---

**Repository:** https://github.com/Teolian/Japan-Energy-Dashboard
