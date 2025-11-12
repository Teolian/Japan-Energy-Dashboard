// Trading Intelligence Types

export type SignalType = 'buy' | 'sell' | 'hold'
export type ConfidenceLevel = 'high' | 'medium' | 'low'
export type TimeHorizon = 'short' | 'medium' | 'long'

// Arbitrage Opportunity
export interface ArbitrageOpportunity {
  time: string // ISO timestamp
  hour: number // 0-23
  type: SignalType
  currentPrice: number // JPY/kWh
  targetPrice: number // JPY/kWh
  spread: number // JPY/kWh
  expectedProfit: number // JPY per MWh
  confidence: ConfidenceLevel
  recommendation: string
  reasoning: string[]
}

// Battery ROI Calculation
export interface BatteryROI {
  batteryCapacityMWh: number
  cyclesPerDay: number
  efficiency: number // 0.0-1.0
  capitalCostJPY: number
  dailyProfitJPY: number
  monthlyProfitJPY: number
  yearlyProfitJPY: number
  paybackYears: number
  roi: number // %
}

// Load Shift Recommendation
export interface LoadShiftRecommendation {
  id: string
  fromHour: number
  toHour: number
  amountMW: number
  currentCost: number // JPY
  optimizedCost: number // JPY
  savings: number // JPY
  carbonReduction: number // kg CO2
  feasibility: number // 0-100
  priority: 'high' | 'medium' | 'low'
  reason: string
}

// Load Profile
export interface LoadProfile {
  hour: number
  currentLoad: number // MW
  optimizedLoad: number // MW
  price: number // JPY/kWh
  carbonIntensity: number // gCO2/kWh
}

// Renewable Procurement Strategy
export interface RenewableStrategy {
  currentRenewablePct: number
  targetRenewablePct: number
  timeHorizon: TimeHorizon
  actionPlan: RenewableAction[]
  estimatedCost: number // JPY
  carbonCreditsValue: number // JPY
  netCost: number // JPY
  timeline: string
}

export interface RenewableAction {
  id: string
  type: 'solar_ppa' | 'wind_ppa' | 'battery' | 'green_certificate' | 'demand_response'
  description: string
  capacityMW: number
  costJPY: number
  timelineMonths: number
  carbonImpactTonnes: number
  priority: number // 1-10
  status: 'planned' | 'in_progress' | 'completed'
}

// Market Forecast
export interface MarketForecast {
  date: string
  area: 'tokyo' | 'kansai'
  horizon: '1h' | '6h' | '24h'
  forecasts: ForecastPoint[]
  accuracy: number // 0-100
  confidenceInterval: {
    lower: number[]
    upper: number[]
  }
  alerts: ForecastAlert[]
}

export interface ForecastPoint {
  timestamp: string
  forecastPrice: number // JPY/kWh
  lowerBound: number
  upperBound: number
  confidence: number // 0-100
}

export interface ForecastAlert {
  type: 'spike' | 'drop' | 'unusual_pattern' | 'opportunity'
  severity: 'info' | 'warning' | 'critical'
  time: string
  message: string
  recommendation: string
}

// Trading Metrics (for dashboard summary)
export interface TradingMetrics {
  totalOpportunities: number
  estimatedDailySavings: number // JPY
  estimatedMonthlySavings: number // JPY
  carbonReductionPotential: number // kg CO2
  optimalBatterySize: number // MWh
  averageArbitrageSpread: number // JPY/kWh
}
