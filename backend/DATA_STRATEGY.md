# Data Strategy - Japan Energy Dashboard

## Current State (October 2025)

### Data Modes

The dashboard supports 3 data modes:

1. **MOCK** - Frontend-generated synthetic data
   - Always available
   - Realistic patterns
   - No external dependencies

2. **LIVE** - Static JSON files in `/public/data/jp/`
   - Pre-fetched data from real sources
   - Currently available dates: 2025-10-23, 2025-10-24, 2025-10-25
   - Graceful fallback to MOCK if files missing

3. **API Refresh** - On-demand fetching via backend
   - Triggered by "Refresh Data" button
   - Runs CLI tools with `--use-http` flag
   - Attempts real HTTP, falls back to testdata

---

## Current Limitations

### Problem 1: HTTP Source Blocking

**TEPCO (Tokyo Demand):**
```
URL: https://www.tepco.co.jp/forecast/html/download-j.html
Error: HTTP 403 Forbidden
```
**Reason:** Website blocks automated requests without proper headers

**Kansai Electric (Kansai Demand):**
```
URL: https://www.kansai-td.co.jp/denkiyoho/download.html
Error: HTTP 404 Not Found
```
**Reason:** URL may have changed or requires authentication

**JEPX, OCCTO:**
- Currently use testdata only
- Real API integration not yet implemented

### Problem 2: Limited Testdata

**Available dates in testdata:**
- 2025-10-23 (original sample)
- 2025-10-24 (added for demo)
- 2025-10-25 (copied from 24)

**Requesting other dates:**
- CLI tools fail with "no data found for date YYYY-MM-DD"
- Frontend falls back to MOCK mode

---

## Solutions

### Short-term (Demo/Portfolio)

✅ **Option 1: Expand Testdata** (Current approach)
- Copy existing JSON files for more dates
- Update timestamps with `sed` or manual editing
- Good for: demos, screenshots, portfolio

✅ **Option 2: Use MOCK Mode**
- Set `VITE_DATA_MODE=mock` in `.env`
- Always generates data for any date
- Good for: development, testing UI

### Mid-term (Production MVP)

**Option 3: GitHub Actions Automation** (P2.1 in roadmap)
```yaml
name: Fetch Daily Data
on:
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM JST
jobs:
  fetch:
    runs-on: ubuntu-latest
    steps:
      - run: go run cmd/fetch-demand-http/main.go --use-http
      - run: git commit -am "data: update $(date +%Y-%m-%d)"
      - run: git push
```

**Benefits:**
- Automatic daily updates
- No manual intervention
- Historical data accumulates
- Works even if HTTP sources are blocked (uses testdata)

**Option 4: Fix HTTP Headers**
```go
// In pkg/http/fetcher.go
req.Header.Set("User-Agent", "Mozilla/5.0 ...")
req.Header.Set("Referer", "https://www.tepco.co.jp/")
req.Header.Set("Accept-Language", "ja-JP")
```

Try different User-Agent strings:
- Desktop browser: `Mozilla/5.0 (Windows NT 10.0; Win64; x64)...`
- Mobile browser: `Mozilla/5.0 (iPhone; CPU iPhone OS 15_0...)...`
- Generic crawler: `curl/7.68.0`

### Long-term (Scale)

**Option 5: Proxy Service**
- Use proxy rotation service (e.g., ScraperAPI, BrightData)
- Bypass IP-based blocking
- Cost: ~$50-150/month for 100k requests

**Option 6: Official APIs**
- Contact TEPCO, Kansai Electric, OCCTO, JEPX
- Request official API access
- May require business registration

**Option 7: Third-party Aggregator**
- Use existing services like:
  - EnergyCharts.info
  - ENTSO-E (Europe, but good pattern)
  - Energy Data Exchange (US DOE)
- Evaluate if Japanese equivalents exist

---

## Recommended Approach for Your Use Case

### For Portfolio/Demo (Now - 1 week)
```bash
# 1. Create data for 7-day window
for date in 2025-10-{23..29}; do
  cp public/data/jp/tokyo/demand-2025-10-24.json \
     public/data/jp/tokyo/demand-$date.json
  sed -i '' "s/2025-10-24/$date/g" \
     public/data/jp/tokyo/demand-$date.json
done

# 2. Same for kansai, jepx, reserve

# 3. Demo ready!
```

**Benefits:**
- Works immediately
- No HTTP issues
- Clean for screenshots
- Sufficient for portfolio

### For Production (1-2 months)

1. **Week 1-2:** Implement GitHub Actions daily fetch
2. **Week 3:** Fix HTTP headers, test with real sources
3. **Week 4:** Add error notifications (email/Slack)
4. **Week 5-8:** Explore official API partnerships

---

## Implementation Status

| Component | Status | Notes |
|-----------|--------|-------|
| MOCK mode | ✅ Complete | Always works |
| LIVE mode | ✅ Complete | 3 dates available |
| API Refresh Button | ✅ Complete | Shows detailed status |
| RefreshStatusModal | ✅ Complete | Visual progress per source |
| HTTP fetching | ⚠️ Partial | Blocked by 403/404 |
| GitHub Actions | ❌ Planned | P2.1 in roadmap |
| Official APIs | ❌ Research | Future consideration |

---

## How to Test

### Test MOCK Mode
```bash
# In frontend-japan/.env
VITE_DATA_MODE=mock

# Any date will work
# Navigate to 2025-11-01, 2026-01-15, etc.
```

### Test LIVE Mode
```bash
# In frontend-japan/.env
VITE_DATA_MODE=live

# Only works for: 2025-10-23, 24, 25
# Other dates fall back to MOCK
```

### Test API Refresh
```bash
# Start servers
cd backend && go run cmd/api/main.go
cd frontend-japan && npm run dev

# Click "Refresh Data" button
# Modal shows status for all 5 sources
```

### Simulate Date Range
```bash
# Generate data for Oct 20-30 (11 days)
cd /Users/teo/projeck/aversome
./scripts/generate-date-range.sh 2025-10-20 2025-10-30
```

---

## FAQ

**Q: Why do HTTP requests fail with 403?**
A: Websites detect automated requests. Need proper User-Agent and headers.

**Q: Can we use web scraping?**
A: Legal gray area in Japan. Better to request official access.

**Q: How often should data update?**
A: Daily is sufficient. Real-time not needed for historical data.

**Q: What about missing dates?**
A: Frontend gracefully falls back to MOCK mode with warning indicator.

**Q: Should we cache API responses?**
A: Yes, but currently static JSON serves this purpose.

**Q: What's the cost of official APIs?**
A: Usually free for research/non-commercial. Commercial may require licensing.

---

## Next Steps

**Priority 1 (This Week):**
- ✅ Implement RefreshStatusModal with visual feedback
- ✅ Create data for 2025-10-25
- ⏳ Document data strategy (this file)

**Priority 2 (Next Week):**
- Test different User-Agent strings
- Research TEPCO/Kansai API documentation
- Set up GitHub Actions workflow

**Priority 3 (Next Month):**
- Contact data providers for official access
- Implement caching layer
- Add data quality monitoring

---

Last updated: 2025-10-28
Maintainer: Teo
