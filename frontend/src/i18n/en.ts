// English translations
export default {
  nav: {
    dashboard: 'Dashboard',
    trading: 'Trading Intelligence',
    analysis: 'Analysis',
    reports: 'Reports',
    settings: 'Settings'
  },
  header: {
    title: 'Japan Energy Dashboard',
    subtitle: 'Auto-updated daily at 00:30 JST with latest data',
    language: 'Language'
  },
  dashboard: {
    demand: 'Electricity Demand',
    prices: 'Spot Prices',
    reserves: 'Reserve Margin',
    generation: 'Generation Mix',
    weather: 'Weather & Solar',
    insights: 'Key Insights',
    comparison: 'Regional Comparison'
  },
  areas: {
    tokyo: 'Tokyo',
    kansai: 'Kansai',
    all: 'All Regions'
  },
  metrics: {
    peak: 'Peak',
    average: 'Average',
    current: 'Current',
    forecast: 'Forecast',
    actual: 'Actual',
    capacity: 'Capacity',
    utilization: 'Utilization',
    reserveMargin: 'Reserve Margin',
    renewable: 'Renewable',
    carbon: 'Carbon Intensity',
    price: 'Price',
    demand: 'Demand',
    generation: 'Generation'
  },
  units: {
    mw: 'MW',
    gwh: 'GWh',
    percent: '%',
    yenPerKwh: '¥/kWh',
    gco2PerKwh: 'gCO₂/kWh',
    hour: 'hour',
    day: 'day',
    week: 'week',
    month: 'month'
  },
  insights: {
    peakDemandTime: 'Peak Demand Time',
    peakDemandMessage: 'Peak demand occurs at {hour}:00 with {demand} MW in {area} area',
    costOptimization: 'Cost Optimization Opportunity',
    costOptimizationMessage: 'Prices are {diff}% higher during peak hours. Shifting consumption to night (0-6am) can reduce costs significantly',
    powerSupplyStable: 'Power Supply Stable',
    powerSupplyStableMessage: 'Reserve margin at {reserve}% - Power supply is comfortable with no constraints expected',
    reserveUnderWatch: 'Reserve Margin Under Watch',
    reserveUnderWatchMessage: 'Reserve at {reserve}% - Monitoring required. Consider load shifting if possible',
    tightPowerSupply: 'Tight Power Supply',
    tightPowerSupplyMessage: 'Reserve critically low at {reserve}% - Demand reduction recommended during peak hours'
  },
  trading: {
    title: 'Trading Intelligence',
    subtitle: 'AI-powered insights for energy market optimization',
    arbitrage: 'Arbitrage Opportunities',
    loadShift: 'Load Shift Advisor',
    renewable: 'Renewable Strategy',
    forecast: 'Market Forecast',
    opportunities: 'Opportunities',
    recommendations: 'Recommendations',
    strategy: 'Strategy',
    prediction: 'Prediction'
  },
  arbitrage: {
    title: 'Arbitrage Opportunities',
    buySignal: 'Buy Signal',
    sellSignal: 'Sell Signal',
    profit: 'Expected Profit',
    confidence: 'Confidence',
    high: 'High',
    medium: 'Medium',
    low: 'Low',
    recommendation: 'Recommendation',
    roi: 'ROI Calculator',
    batterySize: 'Battery Size (MWh)',
    cyclesPerDay: 'Cycles per Day',
    efficiency: 'Efficiency',
    estimatedProfit: 'Estimated Daily Profit'
  },
  loadShift: {
    title: 'Load Shift Advisor',
    currentLoad: 'Current Load Profile',
    optimizedLoad: 'Optimized Load Profile',
    savings: 'Cost Savings',
    carbonReduction: 'Carbon Reduction',
    feasibility: 'Feasibility Score',
    recommendations: 'Shift Recommendations',
    shiftFrom: 'Shift From',
    shiftTo: 'Shift To',
    amount: 'Amount',
    potentialSaving: 'Potential Saving',
    implement: 'Implement'
  },
  status: {
    loading: 'Loading data...',
    error: 'Error loading data',
    noData: 'No data available',
    updated: 'Last updated',
    stable: 'Stable',
    watch: 'Watch',
    tight: 'Tight',
    critical: 'Critical'
  },
  actions: {
    refresh: 'Refresh',
    export: 'Export',
    download: 'Download',
    settings: 'Settings',
    viewDetails: 'View Details',
    apply: 'Apply',
    cancel: 'Cancel',
    save: 'Save'
  },
  dataMode: {
    live: 'Live Data',
    mock: 'Mock Data',
    switchTo: 'Switch to {mode}'
  }
}
