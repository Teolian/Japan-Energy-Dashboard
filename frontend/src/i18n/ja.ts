// Japanese translations (日本語)
export default {
  nav: {
    dashboard: 'ダッシュボード',
    trading: 'トレーディング',
    analysis: '分析',
    reports: 'レポート',
    settings: '設定'
  },
  header: {
    title: '日本エネルギーダッシュボード',
    subtitle: '毎日00:30（JST）に最新データで自動更新',
    language: '言語'
  },
  dashboard: {
    demand: '電力需要',
    prices: 'スポット価格',
    reserves: '供給予備率',
    generation: '電源構成',
    weather: '天気と太陽光',
    insights: '重要な洞察',
    comparison: '地域比較'
  },
  areas: {
    tokyo: '東京',
    kansai: '関西',
    all: '全地域'
  },
  metrics: {
    peak: 'ピーク',
    average: '平均',
    current: '現在',
    forecast: '予測',
    actual: '実績',
    capacity: '供給力',
    utilization: '使用率',
    reserveMargin: '予備率',
    renewable: '再生可能',
    carbon: '炭素強度',
    price: '価格',
    demand: '需要',
    generation: '発電'
  },
  units: {
    mw: 'MW',
    gwh: 'GWh',
    percent: '%',
    yenPerKwh: '円/kWh',
    gco2PerKwh: 'gCO₂/kWh',
    hour: '時間',
    day: '日',
    week: '週',
    month: '月'
  },
  insights: {
    peakDemandTime: 'ピーク需要時間',
    peakDemandMessage: '{area}エリアで{hour}:00にピーク需要{demand} MWが発生',
    costOptimization: 'コスト最適化の機会',
    costOptimizationMessage: 'ピーク時の価格は{diff}%高くなります。夜間（0-6時）へのシフトでコスト削減可能',
    powerSupplyStable: '電力供給安定',
    powerSupplyStableMessage: '予備率{reserve}% - 電力供給は余裕があり、制約は予想されません',
    reserveUnderWatch: '予備率監視中',
    reserveUnderWatchMessage: '予備率{reserve}% - 監視が必要です。可能であれば負荷シフトを検討してください',
    tightPowerSupply: '電力供給逼迫',
    tightPowerSupplyMessage: '予備率が{reserve}%と低水準 - ピーク時の需要削減を推奨'
  },
  trading: {
    title: 'トレーディングインテリジェンス',
    subtitle: 'AI駆動のエネルギー市場最適化インサイト',
    arbitrage: '裁定取引機会',
    loadShift: '負荷シフト推奨',
    renewable: '再エネ戦略',
    forecast: '市場予測',
    opportunities: '機会',
    recommendations: '推奨事項',
    strategy: '戦略',
    prediction: '予測'
  },
  arbitrage: {
    title: '裁定取引機会',
    buySignal: '買いシグナル',
    sellSignal: '売りシグナル',
    profit: '期待利益',
    confidence: '信頼度',
    high: '高',
    medium: '中',
    low: '低',
    recommendation: '推奨',
    roi: 'ROI計算機',
    batterySize: '蓄電池容量（MWh）',
    cyclesPerDay: '1日のサイクル数',
    efficiency: '効率',
    estimatedProfit: '推定日次利益'
  },
  loadShift: {
    title: '負荷シフトアドバイザー',
    currentLoad: '現在の負荷プロファイル',
    optimizedLoad: '最適化された負荷プロファイル',
    savings: 'コスト削減',
    carbonReduction: 'CO₂削減',
    feasibility: '実現可能性スコア',
    recommendations: 'シフト推奨',
    shiftFrom: 'シフト元',
    shiftTo: 'シフト先',
    amount: '量',
    potentialSaving: '潜在的削減額',
    implement: '実装'
  },
  status: {
    loading: 'データ読み込み中...',
    error: 'データ読み込みエラー',
    noData: 'データがありません',
    updated: '最終更新',
    stable: '安定',
    watch: '監視',
    tight: '逼迫',
    critical: '危機的'
  },
  actions: {
    refresh: '更新',
    export: 'エクスポート',
    download: 'ダウンロード',
    settings: '設定',
    viewDetails: '詳細を見る',
    apply: '適用',
    cancel: 'キャンセル',
    save: '保存'
  },
  dataMode: {
    live: 'ライブデータ',
    mock: 'モックデータ',
    switchTo: '{mode}に切り替え'
  }
}
