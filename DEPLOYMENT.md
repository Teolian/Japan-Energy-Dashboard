# Deployment Guide

## Backend Deployment (Railway)

### Prerequisites
- Railway CLI installed: `npm i -g @railway/cli`
- Railway account (free tier available)

### Steps

1. **Login to Railway**
```bash
railway login
```

2. **Navigate to backend directory**
```bash
cd backend
```

3. **Initialize Railway project**
```bash
railway init
# Select: "Create a new project"
# Name it: "japan-energy-api" (or your preferred name)
```

4. **Deploy**
```bash
railway up
```

5. **Get your deployment URL**
```bash
railway status
# Look for: "Service URL: https://japan-energy-api.railway.app"
```

6. **Test health endpoint**
```bash
curl https://your-app.railway.app/api/health
# Expected: {"status":"ok","time":"..."}
```

### Environment Variables (Optional)

Set in Railway dashboard or via CLI:
```bash
railway variables set PORT=8080
railway variables set GIN_MODE=release
```

---

## Frontend Deployment (Vercel)

### Prerequisites
- Vercel CLI installed: `npm i -g vercel`
- Vercel account

### Steps

1. **Update environment variables**

Create `frontend/.env.production`:
```bash
VITE_DATA_MODE=live
VITE_API_BASE_URL=https://your-app.railway.app
VITE_FEAT_RESERVE=true
VITE_FEAT_JEPX=true
VITE_FEAT_SETTLEMENT=true
```

2. **Build and deploy**
```bash
cd frontend
npm run build
vercel --prod
```

3. **Set environment variables in Vercel dashboard**
- Go to: Project Settings → Environment Variables
- Add all variables from `.env.production`

4. **Redeploy**
```bash
vercel --prod
```

---

## Local Testing

### Test Backend

```bash
cd backend

# Build
go build -o api cmd/api/main.go

# Run
PORT=8080 ./api

# In another terminal, test endpoints
curl http://localhost:8080/api/health
curl -X POST http://localhost:8080/api/data/refresh \
  -H "Content-Type: application/json" \
  -d '{"date":"2025-10-24","areas":["tokyo","kansai"]}'
```

### Test Frontend with Backend

```bash
# Terminal 1: Backend
cd backend
PORT=8080 go run cmd/api/main.go

# Terminal 2: Frontend
cd frontend
npm run dev
```

Open http://localhost:5173 and:
1. Toggle to "LIVE" mode
2. Click "Refresh" button
3. Check browser console for API calls

---

## Docker Testing (Optional)

### Build and run backend container

```bash
cd backend

# Build
docker build -t japan-energy-api .

# Run
docker run -p 8080:8080 japan-energy-api

# Test
curl http://localhost:8080/api/health
```

---

## Troubleshooting

### CORS Errors
- Check that Vercel URL is in `backend/cmd/api/main.go` CORS config
- Ensure `VITE_API_BASE_URL` points to Railway URL (not localhost)

### 404 on /api/data/refresh
- Verify backend is deployed and running
- Check Railway logs: `railway logs`

### Build Errors
- Go version: Must be 1.23+
- Node version: Must be 20.19+
- Run `go mod tidy` and `npm install`

### Health Check Failing
- Check Railway logs for startup errors
- Verify port 8080 is exposed
- Check firewall rules

---

## Costs

### Railway (Backend)
- Free tier: 500 hours/month
- Estimated usage: ~720 hours/month (always-on)
- Cost: ~$5/month after free tier

### Vercel (Frontend)
- Free tier: Unlimited for personal projects
- Commercial use: $20/month

---

## Next Steps

After deployment:
1. ✅ Test all features in production
2. ✅ Monitor Railway logs for errors
3. ✅ Set up custom domain (optional)
4. ✅ Configure production secrets
5. ✅ Set up monitoring/alerts
