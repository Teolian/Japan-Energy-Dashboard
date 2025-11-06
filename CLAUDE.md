# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Japan Energy Dashboard - Real-time electricity market monitoring for Japan. This is a full-stack application displaying demand, spot prices, and reserve capacity for Japanese power regions (Tokyo/Kansai).

**Live Demo:** https://japan-energy-dashboard.vercel.app
**Backend API:** https://japan-energy-api-production.up.railway.app

## Architecture

### Dual-Mode System
The application operates in two modes:
- **MOCK mode**: Uses static JSON files from `frontend/public/data/jp/` (dates: 2025-10-23 to 2025-10-28)
- **LIVE mode**: Fetches real-time data from Railway-hosted Go API, which scrapes TEPCO, Kansai Electric, JEPX, and OCCTO sources

Mode selection hierarchy:
1. URL query parameter `?mode=live` or `?mode=mock`
2. localStorage key `jp-energy-data-mode`
3. Build-time env `VITE_DATA_MODE`
4. Default: mock

### Data Flow (LIVE mode - PLANNED)

**Current Status:** ✅ 3 out of 4 sources validated and production-ready (2025-11-05)

**Planned Architecture (100% Free):**
```
GitHub Actions Cron (daily at 00:30 JST)
  → Go fetch binaries (compile & run)
    → HTTP adapters (TEPCO/Kansai/JEPX/OCCTO)
      → CSV parsers (Shift-JIS → UTF-8)
        → JSON files → frontend/public/data/jp/
          → Git commit + push
            → Vercel auto-deploy
              → Users see fresh data
```

**Cost:** $0 (GitHub Actions 2000 min/month free, uses ~90 min/month)

**Old approach (Railway):**
```
Frontend (Vercel)
  → Railway API (POST /api/data/refresh)
    → fetch-* binaries (Go)
      → HTTP adapters
        → JSON files (ephemeral storage)
```
**Why changed:** Railway costs money, ephemeral storage. GitHub Actions is free forever.

### Backend Architecture
- **API Server**: Gin framework REST API on port 8080
- **Fetch Binaries**: Separate executables invoked by API
  - `fetch-demand`: TEPCO (ZIP archives, Shift-JIS) / Kansai (direct CSV)
  - `fetch-jepx`: JEPX day-ahead spot prices
  - `fetch-reserve`: OCCTO system reserves (10 regions)
- **HTTP Client**: Custom fetcher with exponential backoff, circuit breaker, browser-like headers
- **Adapters**: Pure functions converting CSV → typed structs

### File Storage Pattern
API writes fetched data to filesystem:
```
backend/public/data/jp/
├── tokyo/demand-2025-10-24.json
├── kansai/demand-2025-10-24.json
├── jepx/spot-tokyo-2025-10-24.json
└── system/reserve-2025-10-24.json
```

**Important**: Railway uses ephemeral storage. Data persists during deployment but resets on redeploy. For production, consider:
- Adding persistent volume mount
- Using cloud storage (S3/GCS)
- Implementing automatic daily fetches with cron

## Development Commands

### Frontend (Vue 3 + TypeScript)
```bash
cd frontend
npm install
npm run dev          # Start dev server (http://localhost:5173)
npm run build        # Build for production
npm run preview      # Preview production build
```

### Backend (Go 1.23)
```bash
cd backend

# Run API server
PORT=8080 go run cmd/api/main.go

# Build all binaries
go build -o fetch-demand cmd/fetch-demand-http/main.go
go build -o fetch-jepx cmd/fetch-jepx-http/main.go
go build -o fetch-reserve cmd/fetch-reserve-http/main.go

# Run tests
go test ./...

# Tidy dependencies
go mod tidy
```

### Full-Stack Local Testing
```bash
# Terminal 1: Backend
cd backend
PORT=8080 go run cmd/api/main.go

# Terminal 2: Frontend
cd frontend
npm run dev

# Open http://localhost:5173
# Toggle to LIVE mode and click Refresh
```

### Manual Data Fetch (Testing)
```bash
cd backend

# Fetch Tokyo demand for specific date
./fetch-demand -area tokyo -date 2025-11-05 --use-http

# Fetch JEPX spot prices
./fetch-jepx -area tokyo -date 2025-11-05 --use-http

# Fetch system reserves
./fetch-reserve -date 2025-11-05
```

## Deployment

### Backend to Railway
```bash
cd backend
railway login
railway up
railway status  # Get deployment URL
```

Update `frontend/.env.production` with Railway URL.

### Frontend to Vercel
```bash
cd frontend
npm run build
vercel --prod
```

**Critical**: Ensure CORS origins in `backend/cmd/api/main.go` include your Vercel domain.

## Key Technical Details

### Data Encoding
- **TEPCO CSV**: Shift-JIS encoding, requires `golang.org/x/text/encoding/japanese`
- **ZIP extraction**: TEPCO publishes monthly data as ZIP, extract with pattern matching
- **CSV parsing**: Auto-detect headers by name (column order varies)

### Date Handling
- All dates in JST (Asia/Tokyo, UTC+9)
- Format: YYYY-MM-DD
- Timestamps: ISO 8601 with timezone (`2025-10-24T15:00:00+09:00`)
- Frontend persists selected date in localStorage (`jp-energy-last-date`)

### State Management (Pinia)
Stores: `demand`, `jepx`, `reserve`, `weather`, `settlement`

Each store:
- Fetches data via dataClient
- Computes metrics (peak, average, forecast accuracy)
- Transforms to chart-ready format
- Handles loading/error states

### API Endpoints
```
GET  /api/health                         # Health check
POST /api/data/refresh                   # Trigger data fetch
     Body: {"date": "YYYY-MM-DD", "areas": ["tokyo", "kansai"]}
GET  /api/demand/:area/:date             # Get demand data (auto-fetch if missing)
GET  /api/jepx/:area/:date               # Get JEPX prices
GET  /api/reserve/:date                  # Get reserve margins
```

### HTTP Client Features
- Exponential backoff (500ms → 30s)
- Circuit breaker pattern
- Browser-like User-Agent to avoid bot detection
- Automatic gzip decompression
- ZIP archive extraction
- Connection pooling
- Configurable timeouts (default: 45s)

### Data Validation Checks
When verifying data freshness:
1. Check file modification time vs. expected date
2. For historical data: expect data up to Nov 4, 2025 (based on project context)
3. Mock data: Oct 23-28, 2025 only
4. Use `ls -la backend/public/data/jp/*/` to inspect fetched files

## Current State & Limitations

### What Works
- Frontend displays mock data (Oct 23-28)
- Backend API server runs and serves health endpoint
- Data fetching binaries execute successfully
- CORS configured for localhost and Vercel

### Known Issues
1. **Railway Storage**: Ephemeral filesystem, data lost on redeploy
2. **Data Freshness**: No automated daily fetch (requires cron or scheduler)
3. **Error Handling**: Fetch failures silently fall back to mock data
4. **Rate Limiting**: No rate limiting on external API calls (TEPCO/JEPX/OCCTO)

### User Request Context (from /init)
User wants to:
1. Make live data work in production (backend updates automatically)
2. See live data reflected in UI
3. Understand where to store data after updates
4. Verify data is current (validation: data up to Nov 4)

## Common Workflows

### Adding a New Data Source
1. Create adapter in `backend/internal/adapters/` (e.g., `new_source.go`)
2. Implement `ParseCSV(io.Reader, date string)` method
3. Add fetch command in `backend/cmd/fetch-new-source/`
4. Build binary and add to Dockerfile
5. Update API server to invoke binary
6. Add Pinia store in `frontend/src/stores/`
7. Create TypeScript types in `frontend/src/types/`
8. Update dataClient with fetch functions

### Debugging Data Fetch Issues
1. Check Railway logs: `railway logs`
2. Test fetch binary locally: `./fetch-demand -area tokyo -date 2025-11-05 --use-http`
3. Inspect output JSON: `cat backend/public/data/jp/tokyo/demand-2025-11-05.json`
4. Verify HTTP adapter: Check browser Network tab for external API responses
5. Test with curl: `curl http://localhost:8080/api/demand/tokyo/2025-11-05`

### Switching Between Mock/Live
```javascript
// In browser console
localStorage.setItem('jp-energy-data-mode', 'live')
location.reload()

// Or use URL
window.location.href = '/?mode=live'
```

## File Locations

### Backend
- API server: `backend/cmd/api/main.go`
- Fetch commands: `backend/cmd/fetch-*/`
- Adapters: `backend/internal/adapters/` (tepco.go, kansai.go, jepx.go, occto.go)
- HTTP client: `backend/pkg/http/fetcher.go`
- Data types: `backend/internal/demand/`, `backend/internal/jepx/`

### Frontend
- Main view: `frontend/src/views/DashboardView.vue`
- Stores: `frontend/src/stores/` (demand.ts, jepx.ts, reserve.ts, weather.ts)
- Data client: `frontend/src/services/dataClient.ts`
- Types: `frontend/src/types/`
- Components: `frontend/src/components/`

### Configuration
- Backend env: `backend/.env.example`
- Frontend env: `frontend/.env.production` (Vercel), `frontend/.env` (local)
- Railway: `backend/railway.toml`, `backend/Dockerfile`
- Vercel: `frontend/vercel.json`

## Testing & Validation

⚠️ **IMPORTANT**: Data validation is in progress. Do NOT enable automatic data fetching until all sources are validated.

### Data Sources Status (2025-11-05)

**Production Ready Sources (3/4):**
- ✅ **TEPCO (Tokyo demand)**: Official ZIP/CSV, Shift-JIS encoding, 24 hourly points
- ✅ **JEPX (Spot prices)**: Third-party CSV (japanesepower.org), 24 hourly points, Tokyo + Kansai
- ✅ **OCCTO (Reserve capacity)**: Official HTTP API, 10 regions, 30-min intervals aggregated to daily

**Not Production Ready:**
- ⚠️ **Kansai Electric**: No simple CSV API available, using testdata only

### Data Validation Workflow

1. **Run validation test:**
```bash
cd backend
./test-data-fetch.sh 2025-11-05
```

2. **Check validation guide:**
   See `DATA_VALIDATION.md` for detailed checklist per source:
   - TEPCO (Tokyo demand) - Shift-JIS encoding, ZIP extraction, unit conversion
   - Kansai Electric - CSV format, unit verification
   - JEPX Spot Prices - Day-ahead prices, reasonable ranges
   - OCCTO Reserve - 10 regions, formula validation

3. **Manual verification:**
```bash
# Test individual source
cd backend
./fetch-demand -area tokyo -date 2025-11-05 --use-http
cat public/data/jp/tokyo/demand-2025-11-05.json | jq

# Compare with official website
open https://www.tepco.co.jp/forecast/html/index-j.html
```

4. **Test Results (2025-11-03):**
   ```bash
   ./test-data-fetch.sh 2025-11-03

   ✅ Tokyo demand: 21,270-30,730 MW (24 points)
   ❌ Kansai demand: HTTP 404 (testdata fallback)
   ✅ JEPX Tokyo: 7.76-15.53 JPY/kWh (24 points)
   ✅ JEPX Kansai: 2.01-12.16 JPY/kWh (24 points)
   ✅ Reserve: 10 regions, all stable (15.71% Tokyo, 18.39% Kansai)
   ```

5. **Ready for Automation:**
   - GitHub Actions can be enabled for TEPCO, JEPX, OCCTO
   - Kansai Electric will remain on testdata until CSV API becomes available
   - See `DATA_SOURCES.md` for detailed documentation

### Backend Tests
```bash
cd backend
go test ./internal/adapters/...  # Test CSV parsers
go test ./pkg/http/...           # Test HTTP fetcher
```

### Frontend Type Checking
```bash
cd frontend
vue-tsc -b                       # TypeScript type check
```

### End-to-End Validation
1. Start backend: `PORT=8080 go run cmd/api/main.go`
2. Test health: `curl http://localhost:8080/api/health`
3. Trigger refresh: `curl -X POST http://localhost:8080/api/data/refresh -H "Content-Type: application/json" -d '{"date":"2025-11-05","areas":["tokyo"]}'`
4. Verify files created: `ls -la backend/public/data/jp/tokyo/`
5. Test GET endpoint: `curl http://localhost:8080/api/demand/tokyo/2025-11-05`

## Dependencies

### Backend
- `gin-gonic/gin` - HTTP framework
- `gin-contrib/cors` - CORS middleware
- `golang.org/x/text` - Character encoding (Shift-JIS)

### Frontend
- `vue` 3.5 - Framework
- `pinia` - State management
- `chart.js` + `vue-chartjs` - Charts
- `lucide-vue-next` - Icons
- `tailwindcss` - Styling

## Environment Variables

### Backend (.env)
```bash
PORT=8080              # API server port
GIN_MODE=release       # Gin mode (debug/release)
```

### Frontend (.env.production)
```bash
VITE_DATA_MODE=live                           # Data mode (mock/live)
VITE_API_BASE_URL=https://...railway.app      # Backend API URL
VITE_FEAT_RESERVE=true                        # Feature flags
VITE_FEAT_JEPX=true
VITE_FEAT_SETTLEMENT=false
```
