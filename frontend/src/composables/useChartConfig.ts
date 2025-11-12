import { computed } from 'vue'
import type { ChartOptions } from 'chart.js'
import { useDarkMode } from './useDarkMode'

/**
 * Unified chart configuration for consistent styling across all charts
 * Provides base config + helpers for line, bar, and area charts
 */
export function useChartConfig() {
  const { isDark } = useDarkMode()

  // Unified color palette
  const colors = {
    tokyo: 'rgb(59, 130, 246)',
    kansai: 'rgb(16, 185, 129)',
    price: 'rgb(234, 88, 12)',
    forecast: 'rgb(156, 163, 175)',
    peak: 'rgb(220, 38, 38)',

    // Generation mix
    solar: 'rgb(251, 191, 36)',
    wind: 'rgb(59, 130, 246)',
    hydro: 'rgb(14, 165, 233)',
    nuclear: 'rgb(168, 85, 247)',
    lng: 'rgb(156, 163, 175)',
    coal: 'rgb(75, 85, 99)',
    other: 'rgb(107, 114, 128)'
  }

  // Grid color based on theme
  const gridColor = computed(() =>
    isDark.value ? 'rgba(255, 255, 255, 0.06)' : 'rgba(0, 0, 0, 0.06)'
  )

  // Text color based on theme
  const textColor = computed(() =>
    isDark.value ? 'rgb(156, 163, 175)' : 'rgb(75, 85, 99)'
  )

  /**
   * Base chart options (common for all chart types)
   */
  const baseOptions: Partial<ChartOptions<any>> = {
    responsive: true,
    maintainAspectRatio: false,
    animation: {
      duration: 600,
      easing: 'easeOutQuart'
    },
    interaction: {
      mode: 'index',
      intersect: false
    },
    hover: {
      mode: 'index',
      intersect: false,
      animationDuration: 150
    },
    plugins: {
      legend: {
        display: true,
        position: 'top',
        labels: {
          usePointStyle: true,
          padding: 12,
          font: { size: 11 },
          color: textColor.value
        }
      },
      tooltip: {
        enabled: true,
        backgroundColor: 'rgba(0, 0, 0, 0.85)',
        padding: 12,
        cornerRadius: 8,
        titleFont: {
          size: 13,
          weight: 'bold'
        },
        bodyFont: {
          size: 12
        },
        displayColors: true
      }
    }
  }

  /**
   * Line chart specific config
   */
  function lineChartConfig(overrides?: Partial<ChartOptions<'line'>>): ChartOptions<'line'> {
    return {
      ...baseOptions,
      scales: {
        x: {
          grid: {
            display: true,
            color: gridColor.value,
            lineWidth: 1
          },
          ticks: {
            font: {
              size: 11,
              family: 'ui-monospace, monospace'
            },
            color: textColor.value
          }
        },
        y: {
          type: 'linear',
          display: true,
          position: 'left',
          beginAtZero: false,
          grid: {
            display: true,
            color: gridColor.value,
            lineWidth: 1
          },
          ticks: {
            font: {
              size: 11,
              family: 'ui-monospace, monospace'
            },
            color: textColor.value
          }
        }
      },
      ...overrides
    } as ChartOptions<'line'>
  }

  /**
   * Dual-axis line chart config (demand + price)
   */
  function dualAxisConfig(
    yAxisLabel: string,
    y1AxisLabel: string,
    overrides?: Partial<ChartOptions<'line'>>
  ): ChartOptions<'line'> {
    return {
      ...baseOptions,
      scales: {
        x: {
          grid: {
            display: true,
            color: gridColor.value,
            lineWidth: 1
          },
          ticks: {
            font: {
              size: 11,
              family: 'ui-monospace, monospace'
            },
            color: textColor.value
          }
        },
        y: {
          type: 'linear',
          display: true,
          position: 'left',
          beginAtZero: false,
          title: {
            display: true,
            text: yAxisLabel,
            font: {
              size: 12,
              weight: 'bold'
            },
            color: textColor.value
          },
          grid: {
            display: true,
            color: gridColor.value,
            lineWidth: 1
          },
          ticks: {
            font: {
              size: 11,
              family: 'ui-monospace, monospace'
            },
            color: textColor.value
          }
        },
        y1: {
          type: 'linear',
          display: true,
          position: 'right',
          beginAtZero: false,
          title: {
            display: true,
            text: y1AxisLabel,
            font: {
              size: 12,
              weight: 'bold'
            },
            color: textColor.value
          },
          grid: {
            drawOnChartArea: false
          },
          ticks: {
            font: {
              size: 11,
              family: 'ui-monospace, monospace'
            },
            color: textColor.value
          }
        }
      },
      ...overrides
    } as ChartOptions<'line'>
  }

  /**
   * Stacked area chart config (for generation mix)
   */
  function stackedAreaConfig(overrides?: Partial<ChartOptions<'line'>>): ChartOptions<'line'> {
    return {
      ...baseOptions,
      scales: {
        x: {
          stacked: true,
          grid: {
            display: false
          },
          ticks: {
            font: { size: 10 },
            color: textColor.value
          }
        },
        y: {
          stacked: true,
          type: 'linear',
          display: true,
          position: 'left',
          grid: {
            color: gridColor.value
          },
          ticks: {
            font: { size: 11 },
            color: textColor.value
          }
        }
      },
      ...overrides
    } as ChartOptions<'line'>
  }

  /**
   * Default dataset style for line charts (unified)
   */
  const lineDatasetDefaults = {
    borderWidth: 2.5,
    tension: 0.3,
    pointRadius: 0,
    pointHoverRadius: 5,
    pointHoverBorderColor: '#fff',
    pointHoverBorderWidth: 2
  }

  /**
   * Helper to create opacity from rgb
   */
  function withOpacity(rgb: string, opacity: number): string {
    return rgb.replace('rgb', 'rgba').replace(')', `, ${opacity})`)
  }

  return {
    colors,
    gridColor,
    textColor,
    baseOptions,
    lineChartConfig,
    dualAxisConfig,
    stackedAreaConfig,
    lineDatasetDefaults,
    withOpacity
  }
}
