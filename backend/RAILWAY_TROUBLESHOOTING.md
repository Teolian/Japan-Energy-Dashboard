# Railway Deployment Troubleshooting

## Issue: Deploy Failed

### Fix 1: Updated Dockerfile (DONE ✅)
- Added `curl` to alpine image for healthcheck
- Simplified railway.toml config

### Next Steps:

**1. Redeploy with fixed Dockerfile:**
```bash
railway up
```

**2. If still failing, add PORT environment variable:**

**Via Railway Dashboard:**
1. Open https://railway.app/dashboard
2. Select your project: "japan-energy-api"
3. Click on your service
4. Go to "Variables" tab
5. Click "Add Variable"
6. Add: `PORT` = `8080`
7. Click "Deploy" or run `railway up`

**Via CLI:**
```bash
railway variables set PORT=8080
railway up
```

**3. If healthcheck is causing issues, use simple Dockerfile:**
```bash
# Rename current Dockerfile
mv Dockerfile Dockerfile.backup

# Use simple version
cp Dockerfile.simple Dockerfile

# Redeploy
railway up

# Restore original if needed
mv Dockerfile.backup Dockerfile
```

## Common Issues

### Issue: Port not binding
**Solution:** Railway auto-injects PORT variable. Update backend code if needed:
```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
```

### Issue: Healthcheck failing
**Solution:**
1. Check logs: Open Railway dashboard → Deployments → View Logs
2. Use Dockerfile.simple (no healthcheck)
3. Verify /api/health endpoint works

### Issue: Build succeeds but deploy fails
**Possible causes:**
- Application crashes on startup
- Port binding issue
- Missing environment variables

**Debug:**
```bash
# Check logs in Railway dashboard
# Or use CLI (may need service selection):
railway logs --tail 100
```

## Verify Deployment

Once deployed successfully:

```bash
# Get your URL
railway status

# Test health endpoint
curl https://your-app.railway.app/api/health

# Expected: {"status":"ok","time":"..."}
```

## Alternative: Deploy without CLI

1. Push to GitHub
2. Connect Railway to your GitHub repo
3. Railway will auto-deploy on push

## Need Help?

- Railway Docs: https://docs.railway.app
- Railway Discord: https://discord.gg/railway
- Check backend/DEPLOY_CHECKLIST.md for full guide
