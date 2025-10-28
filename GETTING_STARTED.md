# Getting Started with Japan Energy Dashboard

## Project Location

**New dedicated repository:** https://github.com/Teolian/Japan-Energy-Dashboard

**Local development:**
```
/Users/teo/projeck/japan-energy-dashboard/
├── frontend/     # Vue 3 dashboard
└── backend/      # Go API server
```

## Quick Start (5 minutes)

### 1. Run Frontend (Mock Mode - No Backend Needed)

```bash
cd /Users/teo/projeck/japan-energy-dashboard/frontend
npm install
npm run dev
```

Open http://localhost:5173

**What you can do:**
- View Tokyo and Kansai demand data
- Analyze JEPX spot prices
- Monitor reserve capacity
- Navigate dates (2025-10-23 to 2025-10-28)
- Toggle between Mock/Live modes

### 2. Run Full Stack (Optional - For Live Mode)

**Terminal 1: Backend**
```bash
cd /Users/teo/projeck/japan-energy-dashboard/backend
PORT=8080 go run cmd/api/main.go
```

**Terminal 2: Frontend**
```bash
cd /Users/teo/projeck/japan-energy-dashboard/frontend
npm run dev
```

Now you can:
- Switch to LIVE mode (toggle in top-right)
- Click "Refresh" button to fetch real-time data
- See progress sidebar with fetch status

## Deploy to Vercel

### Option 1: Vercel CLI (Fastest)

```bash
# Install Vercel CLI (one-time)
npm i -g vercel

# Deploy
cd /Users/teo/projeck/japan-energy-dashboard/frontend
vercel login
vercel --prod
```

**Important:** When Vercel asks:
- Root Directory: `.` (current directory, already in frontend/)
- Build Command: `npm run build`
- Output Directory: `dist`

### Option 2: GitHub + Vercel UI

1. Go to https://vercel.com/new
2. Import repository: `Teolian/Japan-Energy-Dashboard`
3. Configure:
   - **Root Directory**: `frontend`
   - **Build Command**: `npm run build`
   - **Output Directory**: `dist`
   - **Framework Preset**: Vite
4. Click "Deploy"

Vercel will automatically:
- Build your Vue app
- Deploy to global CDN
- Give you a URL like: `https://japan-energy-dashboard.vercel.app`

### Vercel Configuration

Already configured in `frontend/vercel.json`:
```json
{
  "buildCommand": "npm run build",
  "outputDirectory": "dist",
  "framework": "vite"
}
```

## Project Structure

```
japan-energy-dashboard/
├── README.md                    # Main documentation
├── LICENSE                      # MIT license
├── .gitignore                   # Git ignore rules
│
├── frontend/                    # Vue 3 Dashboard (deploy this)
│   ├── src/
│   │   ├── components/          # Vue components
│   │   ├── stores/              # Pinia state management
│   │   ├── services/            # Data fetching
│   │   └── views/               # Dashboard page
│   ├── public/data/jp/          # Mock data (30 JSON files)
│   ├── vercel.json              # Vercel config
│   └── package.json             # Dependencies
│
└── backend/                     # Go API (optional for Live mode)
    ├── cmd/api/                 # REST API server
    ├── internal/adapters/       # TEPCO, Kansai, JEPX, OCCTO
    ├── pkg/http/                # Custom HTTP client
    └── docs/                    # Technical specs
```

## Common Tasks

### Build for Production
```bash
cd frontend
npm run build
# Output: dist/ directory (600 KB total)
```

### Preview Production Build
```bash
npm run preview
# Open http://localhost:4173
```

### Type Check
```bash
npm run type-check
```

### Backend Tests
```bash
cd backend
go test ./internal/adapters/...
go test ./pkg/http/...
```

## Switching Between Projects

**Old project (Corporate Energy Benchmark):**
```bash
cd /Users/teo/projeck/aversome
git checkout main
```

**New project (Japan Energy Dashboard):**
```bash
cd /Users/teo/projeck/japan-energy-dashboard
git checkout main
```

## Git Commands

```bash
# Check status
git status

# Add changes
git add .

# Commit
git commit -m "feat: your feature description"

# Push to GitHub
git push origin main

# Pull latest changes
git pull origin main
```

## Updating GitHub README

After deploying to Vercel, update the live URL in README.md:

```markdown
## Live Demo

**View Demo:** https://your-actual-vercel-url.vercel.app
```

Then commit and push:
```bash
git add README.md
git commit -m "docs: add Vercel deployment URL"
git push origin main
```

## Troubleshooting

### Frontend won't start
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run dev
```

### Build fails
```bash
npm run type-check  # Check for TypeScript errors
```

### Vercel deployment fails
- Check Root Directory is set to `frontend`
- Verify Build Command is `npm run build`
- Check build logs in Vercel dashboard

### Data not loading (Mock mode)
- Ensure you're on dates 2025-10-23 to 2025-10-28
- Check browser console for errors
- Verify `frontend/public/data/jp/` exists with 30 JSON files

## Next Steps

1. **Test locally:**
   ```bash
   cd frontend && npm run dev
   ```

2. **Deploy to Vercel:**
   ```bash
   cd frontend && vercel --prod
   ```

3. **Share your dashboard:**
   - Get your Vercel URL
   - Update README.md with live link
   - Share on LinkedIn/Twitter

## Documentation

- [Main README](README.md) - Full project documentation
- [Frontend README](frontend/README.md) - Detailed frontend docs
- [Backend Docs](backend/) - API and adapter documentation

## Support

**Repository:** https://github.com/Teolian/Japan-Energy-Dashboard
**Issues:** https://github.com/Teolian/Japan-Energy-Dashboard/issues

---

**Current Directory:** `/Users/teo/projeck/japan-energy-dashboard/`
**Repository:** `https://github.com/Teolian/Japan-Energy-Dashboard.git`
**Status:** Ready to deploy
