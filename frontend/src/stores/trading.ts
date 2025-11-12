// Trading Intelligence Store
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  ArbitrageOpportunity,
  LoadShiftRecommendation,
  LoadProfile,
  BatteryROI,
  TradingMetrics,
  SignalType,
  ConfidenceLevel
} from '@/types/trading'
import { useJEPXStore } from './jepx'
import { useDemandStore } from './demand'

export const useTradingStore = defineStore('trading', () => {
  // State
  const opportunities = ref<ArbitrageOpportunity[]>([])
  const loadShiftRecommendations = ref<LoadShiftRecommendation[]>([])
  const loadProfile = ref<LoadProfile[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Stores
  const jepxStore = useJEPXStore()
  const demandStore = useDemandStore()

  // Computed: Trading Metrics Summary
  const metrics = computed((): TradingMetrics | null => {
    if (opportunities.value.length === 0) return null

    const totalOpportunities = opportunities.value.filter(o => o.type !== 'hold').length
    const estimatedDailySavings = loadShiftRecommendations.value.reduce((sum, rec) => sum + rec.savings, 0)
    const carbonReductionPotential = loadShiftRecommendations.value.reduce((sum, rec) => sum + rec.carbonReduction, 0)
    const averageArbitrageSpread = opportunities.value.reduce((sum, o) => sum + o.spread, 0) / opportunities.value.length

    return {
      totalOpportunities,
      estimatedDailySavings,
      estimatedMonthlySavings: estimatedDailySavings * 30,
      carbonReductionPotential,
      optimalBatterySize: calculateOptimalBatterySize(),
      averageArbitrageSpread
    }
  })

  // Computed: Best Opportunities (top 5)
  const bestOpportunities = computed(() => {
    return [...opportunities.value]
      .filter(o => o.type !== 'hold')
      .sort((a, b) => b.expectedProfit - a.expectedProfit)
      .slice(0, 5)
  })

  // Computed: High Priority Load Shifts
  const priorityLoadShifts = computed(() => {
    return loadShiftRecommendations.value
      .filter(rec => rec.feasibility > 70)
      .sort((a, b) => b.savings - a.savings)
      .slice(0, 5)
  })

  // Actions: Analyze Arbitrage Opportunities
  function analyzeArbitrageOpportunities(area: 'tokyo' | 'kansai') {
    loading.value = true
    error.value = null

    try {
      const prices = jepxStore.priceData(area)
      if (!prices || prices.length === 0) {
        throw new Error('No price data available')
      }

      const opps: ArbitrageOpportunity[] = []

      for (let i = 0; i < prices.length; i++) {
        const current = prices[i]!
        const hour = new Date(current.ts).getHours()

        // Calculate if this is a buy or sell opportunity
        const avgPrice = prices.reduce((sum, p) => sum + p.price, 0) / prices.length
        const signal: SignalType = current.price < avgPrice * 0.85 ? 'buy'
                                   : current.price > avgPrice * 1.15 ? 'sell'
                                   : 'hold'

        if (signal !== 'hold') {
          // Find opposite signal in next 12 hours
          const futureWindow = prices.slice(i + 1, i + 13)
          const bestTarget = signal === 'buy'
            ? futureWindow.reduce((max, p) => p.price > max.price ? p : max, futureWindow[0] || current)
            : futureWindow.reduce((min, p) => p.price < min.price ? p : min, futureWindow[0] || current)

          if (bestTarget) {
            const spread = signal === 'buy'
              ? bestTarget.price - current.price
              : current.price - bestTarget.price

            const expectedProfit = spread * 1000 // JPY per MWh

            // Confidence based on spread size and price volatility
            const confidence: ConfidenceLevel = Math.abs(spread) > avgPrice * 0.2 ? 'high'
                                               : Math.abs(spread) > avgPrice * 0.1 ? 'medium'
                                               : 'low'

            opps.push({
              time: current.ts,
              hour,
              type: signal,
              currentPrice: current.price,
              targetPrice: bestTarget.price,
              spread: Math.abs(spread),
              expectedProfit: Math.round(expectedProfit),
              confidence,
              recommendation: generateRecommendation(signal, spread, hour),
              reasoning: generateReasoning(signal, current.price, avgPrice, spread)
            })
          }
        }
      }

      opportunities.value = opps
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to analyze arbitrage'
      console.error('Arbitrage analysis error:', err)
    } finally {
      loading.value = false
    }
  }

  // Actions: Generate Load Shift Recommendations
  function generateLoadShiftRecommendations(area: 'tokyo' | 'kansai') {
    loading.value = true
    error.value = null

    try {
      const prices = jepxStore.priceData(area)
      const demand = area === 'tokyo' ? demandStore.tokyoData : demandStore.kansaiData

      if (!prices || !demand) {
        throw new Error('Missing price or demand data')
      }

      const recommendations: LoadShiftRecommendation[] = []
      const profile: LoadProfile[] = []

      // Build load profile
      for (let i = 0; i < 24; i++) {
        const price = prices.find(p => new Date(p.ts).getHours() === i)
        const demandPoint = demand.series.find(d => new Date(d.ts).getHours() === i)

        if (price && demandPoint) {
          profile.push({
            hour: i,
            currentLoad: demandPoint.demand_mw,
            optimizedLoad: demandPoint.demand_mw, // Will be calculated
            price: price.price,
            carbonIntensity: 300 // Simplified, should come from generation mix
          })
        }
      }

      // Find high price hours and low price hours
      const sortedByPrice = [...profile].sort((a, b) => b.price - a.price)
      const highPriceHours = sortedByPrice.slice(0, 6)
      const lowPriceHours = sortedByPrice.slice(-6)

      // Generate shift recommendations
      highPriceHours.forEach(high => {
        lowPriceHours.forEach(low => {
          // Don't shift to adjacent hours (not practical)
          if (Math.abs(high.hour - low.hour) < 2) return

          const shiftAmount = high.currentLoad * 0.1 // Shift 10% of load
          const currentCost = high.currentLoad * high.price
          const optimizedCost = (high.currentLoad - shiftAmount) * high.price + shiftAmount * low.price
          const savings = currentCost - optimizedCost

          if (savings > 0) {
            recommendations.push({
              id: `shift-${high.hour}-to-${low.hour}`,
              fromHour: high.hour,
              toHour: low.hour,
              amountMW: Math.round(shiftAmount),
              currentCost: Math.round(currentCost),
              optimizedCost: Math.round(optimizedCost),
              savings: Math.round(savings),
              carbonReduction: Math.round(shiftAmount * (high.carbonIntensity - low.carbonIntensity) * 0.001),
              feasibility: calculateFeasibility(high.hour, low.hour, shiftAmount),
              priority: savings > 100000 ? 'high' : savings > 50000 ? 'medium' : 'low',
              reason: `Shift ${Math.round(shiftAmount)}MW from ${high.hour}:00 (¥${high.price.toFixed(2)}/kWh) to ${low.hour}:00 (¥${low.price.toFixed(2)}/kWh)`
            })
          }
        })
      })

      loadShiftRecommendations.value = recommendations
      loadProfile.value = profile
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to generate load shift recommendations'
      console.error('Load shift error:', err)
    } finally {
      loading.value = false
    }
  }

  // Actions: Calculate Battery ROI
  function calculateBatteryROI(
    batteryCapacityMWh: number,
    cyclesPerDay: number,
    efficiency: number,
    capitalCostJPY: number
  ): BatteryROI {
    // Use average arbitrage spread
    const avgSpread = opportunities.value.length > 0
      ? opportunities.value.reduce((sum, o) => sum + o.spread, 0) / opportunities.value.length
      : 5 // Default 5 JPY/kWh spread

    // Daily profit = capacity * cycles * spread * efficiency
    const dailyProfitJPY = batteryCapacityMWh * 1000 * cyclesPerDay * avgSpread * efficiency

    const monthlyProfitJPY = dailyProfitJPY * 30
    const yearlyProfitJPY = dailyProfitJPY * 365

    const paybackYears = capitalCostJPY / yearlyProfitJPY
    const roi = (yearlyProfitJPY / capitalCostJPY) * 100

    return {
      batteryCapacityMWh,
      cyclesPerDay,
      efficiency,
      capitalCostJPY,
      dailyProfitJPY: Math.round(dailyProfitJPY),
      monthlyProfitJPY: Math.round(monthlyProfitJPY),
      yearlyProfitJPY: Math.round(yearlyProfitJPY),
      paybackYears: Math.round(paybackYears * 10) / 10,
      roi: Math.round(roi * 10) / 10
    }
  }

  // Helper: Calculate optimal battery size
  function calculateOptimalBatterySize(): number {
    if (opportunities.value.length === 0) return 0

    // Find largest spread opportunity
    const maxSpread = Math.max(...opportunities.value.map(o => o.spread))

    // Optimal size is proportional to max spread (simplified)
    return Math.round(maxSpread * 10) // MWh
  }

  // Helper: Calculate feasibility of load shift
  function calculateFeasibility(fromHour: number, toHour: number, amountMW: number): number {
    // Factors:
    // 1. Time distance (closer = easier)
    const timeDistance = Math.abs(fromHour - toHour)
    const timeFactor = Math.max(0, 100 - timeDistance * 5)

    // 2. Shift amount (smaller = easier)
    const amountFactor = amountMW < 1000 ? 100 : amountMW < 5000 ? 80 : 60

    // 3. Time of day (night shifts easier)
    const nightShift = (fromHour >= 22 || fromHour <= 6) || (toHour >= 22 || toHour <= 6)
    const todFactor = nightShift ? 20 : 0

    return Math.min(100, Math.round((timeFactor + amountFactor + todFactor) / 2.2))
  }

  // Helper: Generate recommendation text
  function generateRecommendation(signal: SignalType, spread: number, hour: number): string {
    if (signal === 'buy') {
      return `Buy at ${hour}:00. Price ${Math.abs(spread).toFixed(2)} JPY/kWh below average. Sell when prices recover.`
    } else {
      return `Sell at ${hour}:00. Price ${Math.abs(spread).toFixed(2)} JPY/kWh above average. Buy back when prices drop.`
    }
  }

  // Helper: Generate reasoning
  function generateReasoning(signal: SignalType, currentPrice: number, avgPrice: number, spread: number): string[] {
    const reasons: string[] = []

    if (signal === 'buy') {
      reasons.push(`Current price (¥${currentPrice.toFixed(2)}/kWh) is ${Math.round((1 - currentPrice / avgPrice) * 100)}% below daily average`)
      reasons.push(`Expected price recovery creates arbitrage opportunity`)
      if (Math.abs(spread) > avgPrice * 0.2) {
        reasons.push(`Large price deviation indicates strong opportunity`)
      }
    } else {
      reasons.push(`Current price (¥${currentPrice.toFixed(2)}/kWh) is ${Math.round((currentPrice / avgPrice - 1) * 100)}% above daily average`)
      reasons.push(`Expected price normalization creates selling opportunity`)
      if (Math.abs(spread) > avgPrice * 0.2) {
        reasons.push(`Price spike provides favorable selling conditions`)
      }
    }

    return reasons
  }

  return {
    // State
    opportunities,
    loadShiftRecommendations,
    loadProfile,
    loading,
    error,

    // Computed
    metrics,
    bestOpportunities,
    priorityLoadShifts,

    // Actions
    analyzeArbitrageOpportunities,
    generateLoadShiftRecommendations,
    calculateBatteryROI
  }
})
