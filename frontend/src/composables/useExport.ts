// Export utilities for CSV and PNG downloads

import type { DemandResponse } from '@/types/demand'
import type { JEPXResponse } from '@/types/jepx'
import type { SettlementResponse } from '@/types/settlement'

export function useExport() {
  // Helper to trigger file download
  const downloadFile = (content: string | Blob, filename: string, mimeType: string) => {
    const blob = content instanceof Blob ? content : new Blob([content], { type: mimeType })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  // Export demand data to CSV
  const exportDemandCSV = (data: DemandResponse, area: string, date: string) => {
    const headers = ['Timestamp', 'Hour', 'Demand_MW', 'Forecast_MW']
    const rows = data.series.map((point) => {
      const hour = new Date(point.ts).getHours()
      return [
        point.ts,
        hour.toString(),
        point.demand_mw.toString(),
        point.forecast_mw?.toString() || ''
      ]
    })

    const csv = [headers, ...rows]
      .map(row => row.join(','))
      .join('\n')

    const filename = `demand-${area}-${date}.csv`
    downloadFile(csv, filename, 'text/csv')
  }

  // Export JEPX prices to CSV
  const exportPricesCSV = (data: JEPXResponse, area: string, date: string) => {
    const headers = ['Timestamp', 'Hour', 'Price_JPY_per_kWh']
    const rows = data.price_yen_per_kwh.map((point) => {
      const hour = new Date(point.ts).getHours()
      return [
        point.ts,
        hour.toString(),
        point.price.toString()
      ]
    })

    const csv = [headers, ...rows]
      .map(row => row.join(','))
      .join('\n')

    const filename = `prices-${area}-${date}.csv`
    downloadFile(csv, filename, 'text/csv')
  }

  // Export settlement results to CSV
  const exportSettlementCSV = (data: SettlementResponse, area: string, date: string) => {
    // Summary section
    const summary = [
      ['Settlement Cost Report'],
      ['Date', date],
      ['Area', area],
      ['PV Offset', `${(data.assumptions.pv_offset_pct * 100).toFixed(1)}%`],
      [''],
      ['Total Consumption (kWh)', data.totals.kwh.toString()],
      ['Total Cost (JPY)', data.totals.cost_yen.toString()],
      [''],
      ['Hourly Breakdown'],
      ['Timestamp', 'Hour', 'Consumption_kWh', 'Price_JPY_per_kWh', 'Cost_JPY']
    ]

    const hourly = data.by_hour.map((point) => {
      const hour = new Date(point.ts).getHours()
      return [
        point.ts,
        hour.toString(),
        point.kwh.toString(),
        point.price.toString(),
        point.cost.toString()
      ]
    })

    const csv = [...summary, ...hourly]
      .map(row => row.join(','))
      .join('\n')

    const filename = `settlement-${area}-${date}.csv`
    downloadFile(csv, filename, 'text/csv')
  }

  // Export combined report (demand + prices + settlement)
  const exportCombinedCSV = (
    demand: DemandResponse,
    prices: JEPXResponse,
    settlement: SettlementResponse | null,
    area: string,
    date: string
  ) => {
    const headers = [
      'Hour',
      'Timestamp',
      'Demand_MW',
      'Forecast_MW',
      'Spot_Price_JPY_per_kWh',
      'Settlement_Cost_JPY'
    ]

    const rows = demand.series.map((point, index) => {
      const hour = new Date(point.ts).getHours()
      const pricePoint = prices.price_yen_per_kwh[index]
      const settlementPoint = settlement?.by_hour[index]

      return [
        hour.toString(),
        point.ts,
        point.demand_mw.toString(),
        point.forecast_mw?.toString() || '',
        pricePoint?.price.toString() || '',
        settlementPoint?.cost.toString() || ''
      ]
    })

    const csv = [headers, ...rows]
      .map(row => row.join(','))
      .join('\n')

    const filename = `energy-report-${area}-${date}.csv`
    downloadFile(csv, filename, 'text/csv')
  }

  return {
    exportDemandCSV,
    exportPricesCSV,
    exportSettlementCSV,
    exportCombinedCSV
  }
}
