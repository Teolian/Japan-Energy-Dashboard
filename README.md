# Japan Energy Dashboard

Real-time electricity market monitoring for Japan. A pet project exploring energy data visualization and analytics.
LIVE-DEMO:  [
](https://japan-energy-dashboard.vercel.app/)## What it does

Dashboard for tracking electricity demand, spot prices, and system reserves across major Japanese power regions (Tokyo/TEPCO and Kansai).

**Features:**

**Demand Monitoring**
- Hourly electricity demand for Tokyo and Kansai regions
- Capacity utilization tracking
- Peak demand detection
- Historical comparison

**Spot Price Analysis**
- JEPX (Japan Electric Power Exchange) day-ahead prices
- 24-hour price curves
- Regional price comparison
- Price trend analysis

**System Reserves**
- OCCTO (Organization for Cross-regional Coordination) reserve data
- Multi-region coverage (10 Japanese power areas)
- Reserve margin calculations
- Critical threshold alerts

**Data Management**
- Mock/Live data mode toggle
- Manual data refresh
- Progress tracking for data fetching
- Date navigation

## Tech Stack

**Frontend**
- Vue 3.5 (Composition API)
- TypeScript 5.6
- Pinia (state management)
- Chart.js (data visualization)
- TailwindCSS
- Vite

**Backend**
- Go 1.23
- Gin framework
- HTTP adapters for TEPCO, Kansai Electric, JEPX, OCCTO APIs
- Custom HTTP client with retry logic
- JSON file caching

## Data Sources

**Live Mode:**
- TEPCO - Tokyo area demand
- Kansai Electric - Kansai area demand
- JEPX - day-ahead spot prices
- OCCTO - system reserve capacity

**Mock Mode:**
- Sample data in `frontend/public/data/jp/`
- Dates: 2025-10-23 to 2025-10-28
- 30 JSON files with realistic patterns

## Architecture

```
Frontend (Vue 3)
├── Pinia Stores (demand, jepx, reserve)
├── Data Client (mock/live switching)
└── Chart.js visualizations

Backend (Go)
├── REST API (port 8080)
├── HTTP Adapters (TEPCO, Kansai, JEPX, OCCTO)
└── CSV parsers with retry logic
```

## Project Structure

```
japan-energy-dashboard/
├── frontend/              # Vue 3 dashboard
│   ├── src/
│   │   ├── components/    # UI components
│   │   ├── stores/        # Pinia state
│   │   ├── services/      # Data fetching
│   │   └── views/         # Dashboard page
│   └── public/data/jp/    # Mock data (30 JSON files)
│
└── backend/               # Go API
    ├── cmd/api/           # REST server
    ├── internal/adapters/ # Data source parsers
    └── pkg/http/          # HTTP client
```

## Running Locally

### Frontend Only (uses mock data)

```bash
cd frontend
npm install
npm run dev
# Open http://localhost:5173
```

### Full Stack (with backend)

**Terminal 1 - Backend:**
```bash
cd backend
PORT=8080 go run cmd/api/main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```

Now you can switch to LIVE mode and refresh data from real sources.

## Building

```bash
cd frontend
npm run build
# Output: dist/ directory
```

## Data Formats

**Demand (demand-tokyo-2025-10-24.json):**
```json
{
  "date": "2025-10-24",
  "area": "tokyo",
  "demand_mw": [
    { "ts": "2025-10-24T00:00:00+09:00", "demand": 28500 }
  ],
  "capacity_mw": 45000
}
```

**JEPX Spot Prices (spot-tokyo-2025-10-24.json):**
```json
{
  "date": "2025-10-24",
  "area": "tokyo",
  "price_yen_per_kwh": [
    { "ts": "2025-10-24T00:00:00+09:00", "price": 24.32 }
  ]
}
```

**Reserve Capacity (reserve-2025-10-24.json):**
```json
{
  "date": "2025-10-24",
  "reserves": [
    {
      "region": "hokkaido",
      "demand_mw": 3500,
      "capacity_mw": 5000,
      "reserve_percent": 42.86
    }
  ]
}
```

## State Management

```typescript
// Demand Store
const demandStore = useDemandStore()
demandStore.tokyoData        // Tokyo demand
demandStore.kansaiData       // Kansai demand
demandStore.fetchDemandData(area, date)

// JEPX Store
const jepxStore = useJEPXStore()
jepxStore.tokyoData          // Spot prices
jepxStore.fetchJEPXData(area, date)

// Reserve Store
const reserveStore = useReserveStore()
reserveStore.data            // Reserve data
reserveStore.fetchReserveData(date)
```

## HTTP Client

Go HTTP client features:
- Exponential backoff retry
- Circuit breaker pattern
- Browser-like headers to avoid bot detection
- Connection pooling
- Configurable timeouts

```go
fetcher := http.NewFetcher(http.BrowserConfig())
body, err := fetcher.Fetch(url)
```

## License

MIT
