# Deployment Checklist

## ✅ Pre-Deployment (Completed)

- [x] README.md updated with Weather & Solar Forecast
- [x] Backend code compiled successfully (`go build`)
- [x] All tests passing (`go test ./...`)
- [x] Dockerfile created for Railway
- [x] .dockerignore added
- [x] railway.toml configuration created
- [x] CORS updated with Vercel production URL
- [x] .env.example created for frontend
- [x] DEPLOYMENT.md guide created
- [x] Backend tested locally (health endpoint works)

## 📋 Deployment Steps

### 1. Deploy Backend to Railway

```bash
# Install Railway CLI (if not installed)
npm i -g @railway/cli

# Login
railway login

# Navigate to backend
cd backend

# Initialize project
railway init
# Select: Create new project
# Name: japan-energy-api

# Deploy
railway up

# Get deployment URL
railway status
# Copy the Service URL (e.g., https://japan-energy-api.railway.app)
```

### 2. Verify Backend Deployment

```bash
# Test health endpoint
curl https://your-app.railway.app/api/health
# Expected: {"status":"ok","time":"..."}

# Test refresh endpoint
curl -X POST 'https://your-app.railway.app/api/data/refresh' \
  -H 'Content-Type: application/json' \
  -d '{"date":"2025-10-24","areas":["tokyo","kansai"]}'
```

### 3. Update Frontend Environment

Create `frontend/.env.production`:
```bash
VITE_DATA_MODE=live
VITE_API_BASE_URL=https://your-app.railway.app
VITE_FEAT_RESERVE=true
VITE_FEAT_JEPX=true
VITE_FEAT_SETTLEMENT=true
```

### 4. Deploy Frontend to Vercel

```bash
cd frontend

# Build locally to check for errors
npm run build

# Deploy
vercel --prod

# Set environment variables in Vercel dashboard
# Go to: Project Settings → Environment Variables
# Add all variables from .env.production

# Redeploy
vercel --prod
```

### 5. Final Verification

- [ ] Visit https://japan-energy-dashboard.vercel.app
- [ ] Toggle to "LIVE" mode
- [ ] Click "Refresh" button
- [ ] Check browser console for API calls
- [ ] Verify all charts load correctly
- [ ] Check Weather panel shows forecast data
- [ ] Test dark mode toggle
- [ ] Test date navigation

## 🔍 Troubleshooting

### Backend Not Starting
```bash
# Check Railway logs
railway logs

# Common issues:
# - Port binding: Ensure PORT env var is set to 8080
# - Build errors: Check Dockerfile syntax
# - Dependencies: Run `go mod tidy`
```

### CORS Errors
```bash
# Verify CORS in backend/cmd/api/main.go includes:
# "https://japan-energy-dashboard.vercel.app"

# Redeploy backend:
railway up
```

### Frontend Build Errors
```bash
# Check TypeScript errors
npm run build

# Check environment variables
cat .env.production

# Verify API_BASE_URL is correct
```

## 📊 Monitoring

### Railway Dashboard
- View logs: `railway logs --follow`
- View metrics: https://railway.app/dashboard
- Check uptime and response times

### Vercel Dashboard
- View deployments: https://vercel.com/dashboard
- Check build logs
- Monitor bandwidth usage

## 💰 Cost Estimates

### Railway
- Free tier: 500 hours/month
- After free tier: ~$5/month for always-on service
- Recommended: Set up sleep on inactivity

### Vercel
- Free tier: Unlimited for personal projects
- No cost expected for this project

## 🚀 Next Steps

After successful deployment:

1. **Custom Domain** (Optional)
   - Railway: Add custom domain in dashboard
   - Vercel: Add custom domain in project settings
   - Update CORS in backend

2. **Monitoring** (Optional)
   - Set up UptimeRobot for backend health checks
   - Configure alerts for downtime

3. **Caching** (Optional)
   - Add Redis for API response caching
   - Reduce API calls to external sources

4. **Analytics** (Optional)
   - Add Vercel Analytics
   - Track user engagement

## 📝 Notes

- Mock data dates: 2025-10-23 to 2025-10-28
- Live mode requires backend running
- Weather data: Open-Meteo API (free, no auth)
- JEPX data: May fail if CSV format changes
- TEPCO/Kansai: May require browser headers
