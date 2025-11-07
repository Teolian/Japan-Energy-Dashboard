# Japan Energy Dashboard - Business Analytics Roadmap

**Document Version:** 1.0
**Date:** 2025-11-06
**Status:** Planning Phase

---

## Executive Summary

Japan Energy Dashboard currently provides demand forecasting and spot price tracking for Tokyo and Kansai regions. This roadmap identifies **5 critical data sources** currently missing from our analytics that unlock advanced business insights: arbitrage opportunities, curtailment prediction, generation mix analysis, and cross-regional trading signals.

**Business Impact:** Adding these data sources enables transition from **descriptive analytics** (what happened) to **predictive analytics** (what will happen) and **prescriptive analytics** (what should we do).

---

## Current State (As of Nov 2025)

### ✅ Implemented Data Sources

| Data Source | Coverage | Update Frequency | Business Value |
|------------|----------|------------------|----------------|
| **TEPCO Demand** | Tokyo only | Daily (historical) | Demand forecasting, capacity planning |
| **JEPX Spot Prices** | Tokyo + Kansai | Daily (day-ahead) | Price forecasting, trading signals |
| **OCCTO Reserve Margins** | 10 regions | Daily | Grid stability monitoring |
| **Weather/Solar Forecast** | Tokyo + Kansai | Mock data | Renewable generation proxy |

### ❌ Missing Critical Data

**Problem Areas:**
1. **Kansai Demand** - No historical data (using mock)
2. **Renewable Curtailment** - No visibility into solar/wind shutdowns
3. **Intraday Prices** - Missing short-term arbitrage signals
4. **Generation Mix** - No breakdown of energy sources (solar/wind/nuclear/LNG)
5. **Inter-regional Flows** - No cross-border transmission data

---

## Business Opportunities - Detailed Analysis

### 🔴 #1: Renewable Energy Curtailment Tracking (出力抑制)

**What is it:**
Renewable curtailment = forced shutdown of solar/wind plants when grid cannot absorb excess supply. Happens primarily in Kyushu (high solar penetration) during midday.

**Business Impact:**
- **Price Prediction:** Curtailment events → spot prices crash (見える化の価値: predict negative prices)
- **Trading Strategy:** Short solar generation → long spot market when curtailment announced
- **Capacity Planning:** Identify oversupplied hours for industrial demand response

**Data Available:**
- **Kyushu Electric:** https://www.kyuden.co.jp/td_power_usages/out_ctrl_history.html (past curtailment records)
- **OCCTO System:** Cross-regional curtailment aggregation
- **Electrical Japan:** Annual curtailment statistics (1.76 TWh in FY2023)

**Key Metrics to Track:**
- Curtailment volume (MW/hour)
- Curtailment frequency (days/month)
- Curtailment correlation with JEPX spot prices
- Regional breakdown (Kyushu has 80%+ of national curtailment)

**Business Use Cases:**

1. **Arbitrage Alert System**
   - When Kyushu announces curtailment → Tokyo spot prices often drop 30-50% within 2 hours
   - Automated trading signals for intraday market

2. **Demand Response Optimization**
   - Alert industrial customers when free/cheap power available (curtailment periods)
   - Battery storage dispatch signals (charge during curtailment)

3. **Solar Investment Risk Assessment**
   - Calculate curtailment risk by region for new solar projects
   - ROI adjustments based on expected curtailment losses

**Implementation Priority:** 🟡 Medium (requires multi-source aggregation)

**Technical Requirements:**
- Scrape Kyushu Electric curtailment calendar (CSV)
- Parse OCCTO area supply breakdown (jhSybt=03)
- Build correlation model: curtailment volume → JEPX spot price delta

**Estimated Development:** 2-3 weeks

---

### 🔴 #2: JEPX Intraday Market Prices (当日市場)

**What is it:**
Hour-ahead market where power is traded 1 hour before delivery. Prices fluctuate based on real-time grid conditions (vs. day-ahead spot which locks prices 24h early).

**Business Impact:**
- **Arbitrage Signals:** Spot price = ¥10/kWh, Intraday = ¥15/kWh → 50% arbitrage opportunity
- **Forecast Accuracy:** Intraday prices reflect actual conditions (weather changes, outages)
- **Real-time Trading:** Enable algorithmic trading strategies

**Data Available:**
- **JEPX Official:** https://www.jepx.jp/electricpower/market-data/intraday/ (requires membership)
- **japanesepower.org:** https://japanesepower.org/jepxIntraday.csv (free CSV, historical)
- **API Option:** ICE Data Services (paid)

**Key Metrics:**
- Spot vs Intraday spread (¥/kWh)
- Spread volatility (standard deviation)
- Spread correlation with renewable generation forecast errors

**Business Use Cases:**

1. **Intraday Arbitrage Dashboard**
   - Display real-time spread: Spot - Intraday
   - Alert when spread > ¥5/kWh (profitable arbitrage threshold)
   - Historical spread statistics (avg, max, percentiles)

2. **Forecast Error Analysis**
   - Compare day-ahead spot vs intraday → measure forecast accuracy
   - Identify hours with highest forecast errors (typically 16-18h when solar drops)

3. **Trading Signal Generator**
   - If intraday > spot + transaction costs → buy spot, sell intraday
   - Backtesting: calculate profitability of arbitrage strategies

**Implementation Priority:** 🟡 Medium (high value, moderate complexity)

**Technical Requirements:**
- Fetch japanesepower.org/jepxIntraday.csv (similar to current spot price adapter)
- Store both spot and intraday in separate tables
- Calculate spread: `spread = intraday_price - spot_price`
- Add Spread Chart component in frontend

**Estimated Development:** 1-2 weeks

---

### 🔴 #3: Generation Mix Breakdown (電源構成)

**What is it:**
Real-time breakdown of electricity generation by source: Solar, Wind, Nuclear, LNG, Coal, Hydro, etc.

**Business Impact:**
- **Duck Curve Prediction:** High solar % → expect price drop at noon
- **Renewable Penetration Index:** Track Japan's clean energy transition
- **Carbon Intensity Tracking:** Calculate CO2 emissions per kWh for ESG reporting

**Data Available:**
- **OCCTO System:** jhSybt=03 (area supply capacity by source type)
- **Individual TSO APIs:**
  - TEPCO: https://www.tepco.co.jp/forecast/html/area_data-j.html
  - Kyushu: https://www.kyuden.co.jp/td_power_usages/pc.html (shows solar %)
- **Electricity Maps API:** https://api.electricitymaps.com/ (paid, aggregated)

**Key Metrics:**
- Solar generation % (real-time)
- Wind generation % (real-time)
- Renewable ratio (Solar + Wind + Hydro) / Total
- Carbon intensity (gCO2/kWh)

**Business Use Cases:**

1. **Duck Curve Severity Index**
   - Formula: `severity = (solar_peak_MW - solar_valley_MW) / demand_MW`
   - High severity → expect extreme spot price volatility
   - Visualize: "Today's duck curve is **78% more severe** than average"

2. **Clean Energy Score**
   - Display renewable % in real-time
   - Gamification: "Tokyo is 45% renewable right now - greenest hour of the day!"
   - ESG Reporting: Monthly average renewable penetration

3. **Price Prediction Model Input**
   - Feature: `solar_generation_pct` → strong negative correlation with spot price
   - ML Model: predict spot price using [demand, solar%, wind%, time_of_day]

4. **Carbon-Aware Load Shifting**
   - Alert: "Switch heavy workloads to 13:00-15:00 (70% solar, lowest carbon intensity)"
   - Data center operators, EV charging optimization

**Implementation Priority:** 🟡 Medium-High (foundational for ML models)

**Technical Requirements:**
- Parse OCCTO jhSybt=03 CSV (supply capacity by source)
- Aggregate by fuel type: renewable vs fossil vs nuclear
- Store time-series: `{timestamp, solar_mw, wind_mw, nuclear_mw, lng_mw, coal_mw}`
- Frontend: Stacked area chart showing generation mix over 24h

**Estimated Development:** 2-3 weeks

---

### 🔴 #4: Cross-Regional Power Flows (地域間連系線)

**What is it:**
Electricity transmitted between regions via interconnection lines (e.g., Tokyo → Kansai, Kyushu → Chugoku). Indicates regional supply/demand imbalances.

**Business Impact:**
- **Regional Arbitrage:** Tokyo price = ¥20, Kansai price = ¥10 → flow from Kansai to Tokyo
- **Grid Congestion Signals:** High flow → interconnection saturated → price spread widens
- **Supply Security:** Monitor dependency on imports from other regions

**Data Available:**
- **OCCTO System:** jhSybt=04 (inter-regional transmission flow data)
- **Real-time Display:** https://www.occto.or.jp/supply-demand/occto/supply-monitor.html

**Key Metrics:**
- Flow volume (MW) per interconnection line
- Flow direction (import/export)
- Capacity utilization (% of max transmission capacity)
- Correlation: flow volume → price spread

**Business Use Cases:**

1. **Regional Price Spread Dashboard**
   - Display: Tokyo spot = ¥15, Kansai spot = ¥8, Spread = ¥7
   - Flow: Tokyo ← 500 MW ← Kansai (importing)
   - Insight: "Tokyo importing because local prices high"

2. **Interconnection Congestion Alert**
   - When Tokyo-Kansai line at 95% capacity → spread likely to widen further
   - Trading signal: expect Tokyo prices to stay elevated (limited import ability)

3. **Supply Security Index**
   - Calculate: `import_dependency = imported_mw / total_demand_mw`
   - Tokyo importing 20% of demand → vulnerable to Kansai supply disruptions

4. **Arbitrage Opportunity Screener**
   - Scan all 10 regions for price spreads > ¥5/kWh
   - Check if interconnection capacity available (not congested)
   - Display: "Arbitrage opportunity: Buy Kyushu (¥5), Sell Tokyo (¥15), profit = ¥10"

**Implementation Priority:** 🟡 Medium (requires OCCTO system integration)

**Technical Requirements:**
- Fetch OCCTO jhSybt=04 CSV (interconnection flow data)
- Parse columns: `{from_area, to_area, flow_mw, capacity_mw, timestamp}`
- Calculate price spreads: join JEPX spot prices by area
- Frontend: Flow diagram (Sankey chart) showing power flows between regions

**Estimated Development:** 2-3 weeks

---

### 🟢 #5: Kansai Demand Data (OCCTO Source)

**What is it:**
Historical hourly electricity demand for Kansai region. Currently using mock data (unrealistic).

**Business Impact:**
- **Complete Coverage:** Enable Tokyo vs Kansai comparative analysis
- **Model Training:** Use real data for demand forecasting ML models
- **Customer Trust:** Show actual data instead of fake mock data

**Data Available:**
- **OCCTO API:** jhSybt=01 (hourly demand for all 10 regions, including Kansai)
- Format: Same as reserve data (CSV with area breakdown)
- Historical backfill: Available back to 2016

**Key Metrics:**
- Kansai hourly demand (MW)
- Peak demand time (typically 18:00 in winter)
- Demand correlation with Tokyo (usually 60-70% correlation)

**Business Use Cases:**

1. **Tokyo-Kansai Comparison**
   - Display: "Tokyo peak = 40,000 MW, Kansai peak = 16,000 MW"
   - Demand ratio: Kansai is typically 40% of Tokyo demand
   - Load shape comparison (Tokyo more volatile due to business district)

2. **Regional Demand Forecasting**
   - Train separate models for Tokyo and Kansai
   - Forecast Kansai demand based on weather, day-of-week, holidays

3. **Grid Stress Indicator**
   - Compare demand vs supply capacity for Kansai
   - Alert when demand approaches 90% of capacity (grid stress)

**Implementation Priority:** 🟢 High - Quick Win (30 minutes work)

**Technical Requirements:**
- Already implemented OCCTO adapter (used for reserve data)
- Change parameter: jhSybt=02 → jhSybt=01
- Parse demand column: `エリア需要(MW)`
- Update frontend to use real data instead of mock

**Estimated Development:** 30 minutes (TONIGHT!)

---

## Prioritization Matrix

### 🟢 Quick Wins (Week 1)

**#1: Kansai Demand via OCCTO** ✅ START HERE
- Effort: 30 minutes
- Impact: High (fixes data quality issue)
- Tech: Reuse existing OCCTO adapter

**#2: Tokyo-Kansai Price Spread Chart**
- Effort: 2 hours
- Impact: Medium (enables arbitrage analysis)
- Tech: Frontend component using existing JEPX data

---

### 🟡 Medium Term (Weeks 2-4)

**#3: JEPX Intraday Prices**
- Effort: 1-2 weeks
- Impact: High (real-time arbitrage signals)
- Tech: New adapter for japanesepower.org/jepxIntraday.csv

**#4: Generation Mix Breakdown**
- Effort: 2-3 weeks
- Impact: High (enables ML models, carbon tracking)
- Tech: Parse OCCTO jhSybt=03, aggregate by fuel type

---

### 🔴 Long Term (Months 2-3)

**#5: Renewable Curtailment Tracking**
- Effort: 2-3 weeks
- Impact: Very High (predict price crashes)
- Tech: Scrape Kyushu Electric, correlate with JEPX prices

**#6: Cross-Regional Power Flows**
- Effort: 2-3 weeks
- Impact: Medium (arbitrage identification)
- Tech: Parse OCCTO jhSybt=04, Sankey chart visualization

---

## Advanced Analytics Roadmap

### Phase 1: Data Foundation (Current → Week 4)
- ✅ Complete all 10-region data collection (demand, reserve, prices)
- ✅ Add intraday prices for real-time signals
- ✅ Implement generation mix tracking

### Phase 2: Correlation Analysis (Weeks 5-8)
- Solar generation % → Spot price correlation model
- Curtailment events → Price crash prediction
- Tokyo-Kansai demand correlation (identify lead/lag indicators)
- Weather → Demand/Solar forecasting

### Phase 3: Predictive Models (Weeks 9-16)
- **Demand Forecasting ML Model**
  - Features: `[hour, day_of_week, temperature, solar_forecast, historical_demand]`
  - Target: `demand_mw` (1-hour ahead, 24-hour ahead)
  - Model: LSTM (time-series) or XGBoost (gradient boosting)

- **Spot Price Forecasting ML Model**
  - Features: `[demand_mw, solar_generation_pct, reserve_margin, curtailment_flag, hour]`
  - Target: `spot_price_yen_kwh`
  - Model: Random Forest or Neural Network

- **Duck Curve Severity Predictor**
  - Input: Solar forecast (GFS weather model)
  - Output: Predicted midday price drop depth (%)
  - Use: Alert traders 24h before extreme duck curves

### Phase 4: Decision Support (Weeks 17-24)
- **Automated Trading Signals**
  - Buy/Sell recommendations based on arbitrage opportunities
  - Confidence scores (ML model predictions)
  - Backtesting framework to validate strategies

- **Demand Response Optimizer**
  - Recommend load shifting times for industrial customers
  - Calculate cost savings: shift 100 MW from 18:00 → 13:00 saves ¥X

- **Carbon-Aware Scheduling**
  - API endpoint: `GET /api/carbon-intensity/{hour}` → returns gCO2/kWh
  - Use case: Data centers schedule batch jobs during low-carbon hours

---

## Business Value Quantification

### Revenue Opportunities

**1. Arbitrage Trading (Spot vs Intraday)**
- Average spread: ¥3/kWh × 100 MW × 10 trades/day = ¥3,000/day = ¥1M/year
- Requires: Intraday price data + automated trading signals

**2. Demand Response Aggregation**
- Industrial customers pay for "cheap power alerts" (curtailment notifications)
- Revenue model: ¥100/month per customer × 1000 customers = ¥100k/month = ¥1.2M/year

**3. API-as-a-Service**
- Sell real-time generation mix data to ESG platforms
- Pricing: ¥50k/month per enterprise customer × 20 customers = ¥12M/year

**Total Potential Revenue:** ¥14M+/year

### Cost Savings for Customers

**1. Optimized Load Shifting**
- Industrial customer: 10 GWh/year consumption
- Shift 30% to low-price hours (average savings: ¥5/kWh → ¥2/kWh = ¥3/kWh)
- Savings: 3 GWh × ¥3 = ¥9M/year per customer

**2. Battery Storage Dispatch Optimization**
- Battery operator: 50 MW / 100 MWh storage
- Optimize charge/discharge based on intraday spreads
- Additional revenue: ¥1M/year per MW = ¥50M/year

---

## Technical Architecture Changes

### New OCCTO Endpoints to Integrate

```bash
# Demand data (all regions including Kansai)
jhSybt=01: エリア需要 (area demand)

# Supply capacity breakdown by source
jhSybt=03: 電源種別供給力 (generation mix)

# Inter-regional power flows
jhSybt=04: 地域間連系線潮流 (interconnection flows)

# Renewable curtailment
jhSybt=05: 出力抑制実績 (curtailment records)
```

### Database Schema Extensions

```sql
-- New table: Generation Mix
CREATE TABLE generation_mix (
  timestamp TIMESTAMP,
  area VARCHAR(20),
  solar_mw DECIMAL(10,2),
  wind_mw DECIMAL(10,2),
  nuclear_mw DECIMAL(10,2),
  lng_mw DECIMAL(10,2),
  coal_mw DECIMAL(10,2),
  hydro_mw DECIMAL(10,2),
  renewable_pct DECIMAL(5,2),
  carbon_intensity DECIMAL(10,2)  -- gCO2/kWh
);

-- New table: Intraday Prices
CREATE TABLE jepx_intraday (
  timestamp TIMESTAMP,
  area VARCHAR(20),
  price_yen_kwh DECIMAL(10,4),
  volume_mwh DECIMAL(10,2)
);

-- New table: Curtailment Events
CREATE TABLE renewable_curtailment (
  date DATE,
  area VARCHAR(20),
  curtailed_mw DECIMAL(10,2),
  curtailed_mwh DECIMAL(10,2),
  source_type VARCHAR(20),  -- 'solar', 'wind'
  reason VARCHAR(100)
);

-- New table: Inter-regional Flows
CREATE TABLE interconnection_flows (
  timestamp TIMESTAMP,
  from_area VARCHAR(20),
  to_area VARCHAR(20),
  flow_mw DECIMAL(10,2),
  capacity_mw DECIMAL(10,2),
  utilization_pct DECIMAL(5,2)
);
```

### API Endpoints to Add

```
GET /api/generation-mix/{area}/{date}
GET /api/jepx/intraday/{area}/{date}
GET /api/curtailment/{area}/{date}
GET /api/flows/{from_area}/{to_area}/{date}
GET /api/carbon-intensity/{area}/{hour}
GET /api/arbitrage/opportunities  # Returns price spreads > threshold
```

---

## Success Metrics (KPIs)

### Data Quality
- ✅ Kansai demand data accuracy: >95% vs OCCTO official
- ✅ Data completeness: 0 missing days for all regions
- ✅ Update latency: <30 min after official publication

### User Engagement
- Dashboard views: 10,000+/month (currently ~2,000/month)
- API calls: 100,000+/month
- Avg session time: 5+ minutes (currently ~2 min)

### Business Impact
- Arbitrage trades executed: 500+/month
- Customer cost savings: ¥50M+/year (aggregated)
- ML model forecast accuracy: MAE <5% for demand, <10% for prices

---

## Next Steps (Action Items)

### Week 1 - Foundation
- [ ] Implement Kansai demand via OCCTO jhSybt=01
- [ ] Add Tokyo-Kansai price spread chart
- [ ] Update CLAUDE.md with new data sources
- [ ] Test OCCTO jhSybt=03 (generation mix) endpoint

### Week 2-3 - Intraday Prices
- [ ] Build adapter for japanesepower.org/jepxIntraday.csv
- [ ] Add database schema for intraday prices
- [ ] Create Spread Chart component (Spot vs Intraday)
- [ ] Calculate historical spread statistics

### Week 4-6 - Generation Mix
- [ ] Parse OCCTO jhSybt=03 for fuel type breakdown
- [ ] Implement carbon intensity calculation
- [ ] Add Stacked Area Chart for generation mix
- [ ] Build correlation model: solar% → spot price

### Week 7-10 - Curtailment & Flows
- [ ] Scrape Kyushu Electric curtailment calendar
- [ ] Parse OCCTO jhSybt=04 for inter-regional flows
- [ ] Build curtailment alert system
- [ ] Add Sankey diagram for power flows

### Week 11+ - ML Models
- [ ] Collect 6+ months of complete data
- [ ] Train demand forecasting model (LSTM)
- [ ] Train spot price prediction model (XGBoost)
- [ ] Implement backtesting framework
- [ ] Deploy models to production API

---

## Risks & Mitigation

### Technical Risks

**Risk 1: Data Source Changes**
- OCCTO may change CSV format or URL structure
- **Mitigation:** Implement schema validation, alerting on parse failures

**Risk 2: API Rate Limits**
- japanesepower.org may rate-limit scraping
- **Mitigation:** Respect robots.txt, cache aggressively, consider paid API

**Risk 3: Storage Costs**
- 5-minute interval data = 288 records/day × 10 regions × 365 days = 1M+ records/year
- **Mitigation:** Aggregate to hourly for analytics, keep raw data for 90 days only

### Business Risks

**Risk 1: Data Licensing**
- JEPX data may require paid license for commercial use
- **Mitigation:** Consult JEPX terms of service, use free sources (japanesepower.org) for MVP

**Risk 2: Market Competition**
- Bloomberg, Reuters already provide similar analytics
- **Mitigation:** Focus on Japan-specific insights, free tier for individual traders

---

## Conclusion

Japan Energy Dashboard has strong foundation with demand and spot price tracking. Adding **5 critical data sources** (curtailment, intraday prices, generation mix, flows, Kansai demand) unlocks transition to predictive analytics.

**Immediate Action:** Start with Kansai demand (30 min work) → then prioritize intraday prices and generation mix (highest ROI).

**Timeline:** Phase 1 complete in 4 weeks, ML models in production by Week 16.

**Expected Outcome:** Platform becomes essential tool for energy traders, industrial demand response, and ESG reporting in Japan market.

---

**Document Owner:** Claude Code
**Review Cycle:** Monthly
**Next Review:** 2025-12-06
