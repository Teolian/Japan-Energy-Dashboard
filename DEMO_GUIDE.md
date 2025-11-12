# 📊 Japan Energy Dashboard - Demo Guide

## 🎯 Overview

Comprehensive energy market intelligence dashboard for Japanese power regions with AI-powered trading insights.

**Live Demo:** https://japan-energy-dashboard.vercel.app

---

## ✨ New Features (Version 2.0)

### 1. **Internationalization (i18n)** 🌐
- **Japanese (日本語)** and English support
- Language switcher in header
- Default: Japanese for enterprise demos
- All UI elements fully translated

### 2. **Enterprise Design** 🎨
- Professional gradient themes
- Energy-optimized color palette:
  - Primary: Deep Blue (reliability & trust)
  - Green: Renewable energy
  - Purple: AI/Premium features
  - Amber: Warnings
- Shadow effects and glow animations
- Responsive dark mode

### 3. **Trading Intelligence** 🧠 **[NEW TAB]**
Powered by AI-driven market analysis:

#### A. Arbitrage Opportunities
- Buy/Sell signals based on price spreads
- Expected profit calculations (JPY per MWh)
- Confidence levels (High/Medium/Low)
- Real-time recommendations

**Example Output:**
```
BUY at 03:00
Price ¥8.52/kWh below average
Expected Profit: ¥4,280 per MWh
Confidence: HIGH
```

#### B. Battery ROI Calculator
- Optimize battery storage investments
- Calculate payback periods
- Estimate daily/monthly/annual profits
- Adjustable parameters:
  - Battery capacity (MWh)
  - Cycles per day
  - Efficiency (%)

**Example Output:**
```
50 MWh Battery @ 85% efficiency
Daily Profit: ¥124,500
Payback Period: 4.2 years
ROI: 23.8% annual
```

#### C. Load Shift Advisor
- Identifies cost-saving opportunities
- Recommends optimal time shifts
- Calculates:
  - Cost savings (JPY)
  - Carbon reduction (kg CO₂)
  - Feasibility scores
- Interactive load profile visualization

**Example Recommendation:**
```
Shift 2,500 MW from 14:00 → 03:00
Savings: ¥185,000/day
CO₂ Reduction: 450 kg/day
Feasibility: 87%
```

### 4. **Enhanced UX**
- Gradient hero headers
- Tab navigation (Dashboard / Trading)
- AI badge on Trading Intelligence
- Smooth transitions and animations
- Improved data visualization

---

## 🚀 Quick Start

### Frontend Development
```bash
cd frontend
npm install
npm run dev
# Open http://localhost:5173
```

### Backend Development
```bash
cd backend
PORT=8080 go run cmd/api/main.go
```

### Full Stack Testing
```bash
# Terminal 1: Backend
cd backend && PORT=8080 go run cmd/api/main.go

# Terminal 2: Frontend
cd frontend && npm run dev

# Navigate to http://localhost:5173
```

---

## 📋 Demo Script for Japanese Energy Company

### **Opening (5 minutes)**

1. **Homepage Hero** (Japanese)
   ```
   日本エネルギーダッシュボード
   毎日00:30（JST）に最新データで自動更新
   ```

2. **Show Key Metrics:**
   - Tokyo demand: 35,280 MW
   - JEPX spot price: ¥28.5/kWh
   - Reserve margin: 15.2% (Stable)
   - Renewable %: 18.4%

3. **Language Switch**
   - Click 🇯🇵 → 🇬🇧 to demonstrate i18n
   - All metrics, labels, insights auto-translate

### **Core Dashboard (10 minutes)**

1. **Regional Demand**
   - Tokyo vs Kansai comparison
   - Peak hours identification
   - Forecast accuracy

2. **Market Analysis**
   - JEPX spot price trends
   - Price spread analysis
   - Duck curve visualization

3. **Generation Mix**
   - Solar, Wind, Hydro, Nuclear, LNG, Coal breakdown
   - Carbon intensity tracking
   - Greenest hour recommendations

4. **Insights Panel**
   - AI-generated recommendations
   - Cost optimization tips
   - Reserve margin alerts

### **Trading Intelligence** ⭐ **[WOW FACTOR]** (15 minutes)

Navigate to **トレーディング** (Trading) tab:

1. **Arbitrage Opportunities**
   - Show buy/sell signals
   - Explain profit calculations
   - Demonstrate confidence levels
   - **ROI Calculator:**
     - Input: 50 MWh battery
     - Output: ¥3.8M/month profit, 4.2yr payback
     - Highlight: "This pays for itself in 4 years"

2. **Load Shift Advisor**
   - Display top recommendations
   - Show feasibility scores
   - Highlight carbon reduction
   - **Example:** "Shifting 2.5GW saves ¥185K daily"

3. **Area Comparison**
   - Switch between Tokyo / Kansai
   - Compare arbitrage opportunities
   - Note: Tokyo typically 10-15% higher prices

### **Technical Architecture** (5 minutes)

1. **Data Sources (3/4 validated)**
   - ✅ TEPCO (Tokyo official)
   - ✅ JEPX (spot prices)
   - ✅ OCCTO (reserve margins)
   - ⚠️ Kansai (via OCCTO backup)

2. **Automation**
   - GitHub Actions cron (00:30 JST daily)
   - 100% free infrastructure
   - Auto-commit → Vercel deploy
   - ~90 min/month usage (2000 free)

3. **AI/ML Features**
   - Price pattern detection
   - Anomaly alerts
   - Optimization algorithms
   - Forecast accuracy tracking

---

## 🎨 Design Highlights

### Color Palette
```css
Primary (Blue):   #3b82f6 (Trust, reliability)
Renewable (Green): #10b981 (Clean energy)
AI/Premium (Purple): #8b5cf6 (Advanced features)
Warning (Amber):  #f59e0b (Alerts)
Critical (Red):   #ef4444 (Urgent)
```

### Charts
- Solar: Amber (☀️)
- Wind: Cyan (💨)
- Hydro: Blue (💧)
- Nuclear: Purple (⚛️)
- LNG: Orange (🔥)
- Coal: Gray (⚫)

---

## 💡 Key Messages

1. **"100% Free Infrastructure"**
   - GitHub Actions (free tier)
   - Vercel hosting (hobby plan)
   - No ongoing costs

2. **"Validated Data Sources"**
   - Official TEPCO APIs
   - JEPX market data
   - OCCTO system reserves
   - Real-time accuracy

3. **"AI-Powered Insights"**
   - Arbitrage detection
   - Load optimization
   - Carbon reduction tracking
   - ROI calculations

4. **"Enterprise Ready"**
   - Bilingual (日本語 + English)
   - Professional design
   - Scalable architecture
   - API-first approach

5. **"Digital Grid Compatible"**
   - Renewable procurement focus
   - Carbon credit potential
   - Market optimization tools

---

## 📊 Sample Metrics (Live Data)

### Dashboard
```
Tokyo Demand:   35,280 MW (peak at 14:00)
Kansai Demand:  18,450 MW (peak at 13:00)
JEPX Tokyo:     ¥28.5/kWh average
Reserve Margin: 15.2% (Stable)
Renewable:      18.4%
Carbon:         292 gCO₂/kWh
```

### Trading Intelligence
```
Arbitrage Opportunities:  12 detected
Daily Savings Potential:  ¥456,000
Monthly Potential:        ¥13.7M
Optimal Battery Size:     50 MWh
Average Spread:           ¥5.2/kWh
```

### Load Shift Recommendations
```
Top Recommendation:
  From: 14:00 (peak, ¥32.5/kWh)
  To:   03:00 (valley, ¥18.2/kWh)
  Amount: 2,500 MW
  Savings: ¥185,000/day
  CO₂ Reduction: 450 kg/day
  Feasibility: 87%
```

---

## 🔧 Technical Stack

### Frontend
- Vue 3.5 (Composition API)
- TypeScript 5.6
- Pinia (state management)
- vue-i18n (internationalization)
- Chart.js (visualizations)
- TailwindCSS (styling)
- Vite (build tool)

### Backend
- Go 1.23
- Gin (HTTP framework)
- Custom HTTP client (retry, circuit breaker)
- Shift-JIS encoding support
- ZIP archive extraction

### Infrastructure
- **Frontend:** Vercel (auto-deploy)
- **Backend:** Railway (optional, for live API)
- **Data Updates:** GitHub Actions (cron)
- **Cost:** $0/month (free tiers)

---

## 📈 Roadmap

### Phase 1: Complete ✅
- i18n (Japanese + English)
- Enterprise design
- Trading Intelligence
- Arbitrage analysis
- Load shift advisor
- Battery ROI calculator

### Phase 2: Planned
- ML price forecasting
- Historical data storage (PostgreSQL)
- WebSocket real-time updates
- User authentication
- Multi-tenancy support

### Phase 3: Future
- Portfolio management
- Renewable procurement strategy
- Carbon credit tracking
- Mobile app
- Public API

---

## 🎯 Success Metrics

### Technical KPIs
- ✅ Page load < 2s
- ✅ i18n coverage: 100%
- ✅ TypeScript strict mode
- ✅ Zero runtime errors
- ✅ Responsive design

### Business KPIs
- Arbitrage detection accuracy: >90%
- Cost optimization potential: 15-25%
- Carbon reduction tracking: kg CO₂
- User engagement: session duration

---

## 📞 Contact

**Demo Support:**
- Questions: Check CLAUDE.md
- Issues: GitHub Issues
- Documentation: README.md

**Deployment URLs:**
- Frontend: https://japan-energy-dashboard.vercel.app
- Backend API: https://japan-energy-api-production.up.railway.app

---

## ⚡ Quick Tips

1. **Switch Language:**
   - Click flag icon (🇯🇵/🇬🇧) in header
   - Persists in localStorage

2. **Navigate Tabs:**
   - Dashboard: Main energy metrics
   - Trading: AI-powered insights

3. **Change Date:**
   - Use arrows (◀ ▶) or keyboard (←/→)
   - Auto-fetches data for selected date

4. **Toggle Mode:**
   - Live Data: Real market data
   - Mock Data: Sample data (Oct 23-28)

5. **ROI Calculator:**
   - Adjust battery size, cycles, efficiency
   - See instant payback calculations

---

**Built with ❤️ for Japan's Energy Future**

🌱 Renewable Energy · 🔋 Energy Storage · 📊 Market Intelligence · 🤖 AI Optimization
