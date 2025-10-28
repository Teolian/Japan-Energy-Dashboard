# HTTP Integration Guide

**Status:** P1.3 Complete ✅
**Branch:** feat/japan-dashboard

---

## Overview

Real HTTP fetching integration with exponential backoff retry, circuit breaker, and graceful fallback to testdata.

## Architecture

### Components

1. **HTTP Fetcher** (`pkg/http/fetcher.go`)
   - Exponential backoff retry (max 3 attempts)
   - Configurable timeout (default: 30s)
   - User-Agent headers

2. **Circuit Breaker** (`pkg/http/circuit_breaker.go`)
   - Tracks consecutive failures
   - Opens circuit after N failures (default: 5)
   - Half-open state for recovery testing

3. **Sources Configuration** (`pkg/sources/config.go`)
   - Environment variable support
   - Default URLs for JP sources
   - Fallback URL support

4. **Structured Logger** (`pkg/logger/structured.go`)
   - JSON format logging
   - Per-fetch metrics: source, duration_ms, status, artifact
   - Human-readable fallback

---

## Usage

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Data source URLs (optional - falls back to testdata)
TEPCO_URL=https://www.tepco.co.jp/forecast/html/download-j.html
KANSAI_URL=https://www.kansai-td.co.jp/denkiyoho/download.html
OCCTO_URL=https://www.occto.or.jp/
JEPX_URL=https://www.jepx.jp/

# HTTP configuration
HTTP_MAX_RETRIES=3
HTTP_TIMEOUT=30

# Circuit breaker
CIRCUIT_BREAKER_MAX_FAILURES=5
CIRCUIT_BREAKER_TIMEOUT=60

# Logging
LOG_JSON=false
USE_HTTP_FETCH=false
```

### CLI Tools with HTTP Support

#### fetch-demand-http

Fetch demand data with HTTP or testdata fallback:

```bash
# Testdata mode (default)
go run cmd/fetch-demand-http/main.go --area tokyo --date 2025-10-24

# HTTP mode
go run cmd/fetch-demand-http/main.go --area tokyo --date 2025-10-24 --use-http

# With JSON logging
go run cmd/fetch-demand-http/main.go --area kansai --date 2025-10-24 --json-log
```

**Output:**
```
2025/10/25 13:43:39 [info] Fetching tokyo demand data for 2025-10-24 (HTTP: false)
2025/10/25 13:43:39 [info] Using testdata: internal/adapters/testdata/tepco-sample.csv
2025/10/25 13:43:39 [info] Parsed 24 data points
2025/10/25 13:43:39 [info] TEPCO (testdata): Successfully wrote demand data (testdata mode) (0ms)
2025/10/25 13:43:39 ✓ Successfully wrote ../public/data/jp/tokyo/demand-2025-10-24.json (2791 bytes)
```

**JSON Logging Output:**
```json
{"timestamp":"2025-10-25T13:43:54+09:00","level":"info","source":"Kansai (testdata)","duration_ms":0,"status":"success","artifact":"../public/data/jp/kansai/demand-2025-10-24.json","message":"Successfully wrote demand data (testdata mode)"}
```

---

## Features

### 1. Exponential Backoff Retry

Automatic retry with increasing delays:
- Attempt 1: Immediate
- Attempt 2: 500ms delay
- Attempt 3: 1s delay
- Attempt 4: 2s delay

**Max backoff:** 30s
**Total attempts:** 4 (1 initial + 3 retries)

### 2. Circuit Breaker

Prevents cascading failures:

**States:**
- **Closed**: Normal operation, requests flow
- **Open**: Circuit tripped, requests blocked
- **Half-Open**: Testing recovery

**Configuration:**
- Max failures: 5 consecutive errors
- Timeout: 60s before retry
- Auto-reset on success

### 3. Graceful Fallback

If HTTP fetch fails:
1. Log failure with structured output
2. Fallback to testdata
3. Continue processing normally
4. No service disruption

### 4. Structured Logging

**JSON Format:**
```json
{
  "timestamp": "2025-10-25T13:43:54+09:00",
  "level": "info",
  "source": "TEPCO",
  "duration_ms": 1234,
  "status": "success",
  "artifact": "public/data/jp/tokyo/demand-2025-10-24.json",
  "message": "HTTP fetch successful"
}
```

**Human Format:**
```
[info] TEPCO: HTTP fetch successful (1234ms)
```

---

## Integration Checklist

- ✅ HTTP fetcher with retry logic
- ✅ Circuit breaker implementation
- ✅ Sources configuration with env vars
- ✅ Structured JSON logging
- ✅ Testdata fallback mode
- ✅ CLI tool integration (fetch-demand-http)
- ⏳ Real source URL validation
- ⏳ Production deployment config

---

## Testing

### Test Testdata Mode

```bash
cd backend
go run cmd/fetch-demand-http/main.go --area tokyo --date 2025-10-24
```

### Test JSON Logging

```bash
go run cmd/fetch-demand-http/main.go --area kansai --date 2025-10-24 --json-log
```

### Test Circuit Breaker

Set invalid URL and watch circuit open:

```bash
export TEPCO_URL=https://invalid-url.example.com
go run cmd/fetch-demand-http/main.go --area tokyo --date 2025-10-24 --use-http
```

---

## Production Considerations

### 1. Rate Limiting

Current implementation is polite:
- User-Agent: `Corporate-Energy-Platform/1.0 (Educational Project)`
- Exponential backoff prevents hot-loops
- Circuit breaker stops repeated failures

### 2. Source Availability

JP sources may:
- Change URL structure
- Update HTML/CSV format
- Add authentication requirements
- Have rate limits

**Mitigation:**
- Environment variable URLs (easy updates)
- Testdata fallback (always works)
- Adapter versioning

### 3. Monitoring

Log metrics to track:
- `duration_ms`: Fetch performance
- `status`: Success/failure rates
- `source`: Which source is failing
- `artifact`: Output verification

### 4. Caching

Not implemented yet (P1.5):
- Add Redis/disk cache
- TTL per source type
- Cache invalidation strategy

---

## Next Steps (Future)

1. **Real Source Testing**
   - Validate actual TEPCO/Kansai URLs
   - Handle HTML parsing (if needed)
   - Add authentication support

2. **Cache Layer**
   - Redis integration
   - TTL configuration
   - Cache warming

3. **Monitoring Dashboard**
   - Parse JSON logs
   - Grafana/Prometheus
   - Alert on circuit breaker opens

4. **GitHub Actions**
   - Scheduled fetching
   - Commit JSON artifacts
   - Deploy to Vercel

---

## Files

**Core:**
- `pkg/http/fetcher.go` — HTTP client with retry
- `pkg/http/circuit_breaker.go` — Circuit breaker
- `pkg/sources/config.go` — Source URLs configuration
- `pkg/logger/structured.go` — JSON logging

**CLI:**
- `cmd/fetch-demand-http/main.go` — Demand data fetcher with HTTP

**Config:**
- `.env.example` — Environment variables template

---

**Last Updated:** 2025-10-25
**Status:** Production-ready (testdata mode), Beta (HTTP mode)
