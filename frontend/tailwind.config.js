/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Enterprise primary colors (Deep Blue - reliability & trust)
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
          950: '#172554'
        },
        // Accent colors for energy states
        energy: {
          // Renewable/Good (Emerald)
          green: {
            light: '#34d399',
            DEFAULT: '#10b981',
            dark: '#059669'
          },
          // Warning/Watch (Amber)
          yellow: {
            light: '#fbbf24',
            DEFAULT: '#f59e0b',
            dark: '#d97706'
          },
          // Critical/Alert (Red)
          red: {
            light: '#f87171',
            DEFAULT: '#ef4444',
            dark: '#dc2626'
          },
          // Premium/AI features (Purple)
          purple: {
            light: '#a78bfa',
            DEFAULT: '#8b5cf6',
            dark: '#7c3aed'
          },
          // Information (Cyan)
          cyan: {
            light: '#67e8f9',
            DEFAULT: '#06b6d4',
            dark: '#0891b2'
          }
        },
        // Chart colors (optimized for data visualization)
        chart: {
          solar: '#f59e0b', // Amber (sun)
          wind: '#06b6d4', // Cyan (sky)
          hydro: '#3b82f6', // Blue (water)
          nuclear: '#8b5cf6', // Purple (advanced)
          lng: '#f97316', // Orange (thermal)
          coal: '#6b7280', // Gray (traditional)
          renewable: '#10b981', // Green (clean)
          fossil: '#64748b' // Slate (conventional)
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        // Japanese support
        ja: ['Noto Sans JP', 'Inter', 'sans-serif']
      },
      boxShadow: {
        'energy': '0 4px 20px -2px rgba(59, 130, 246, 0.3)',
        'energy-lg': '0 10px 40px -5px rgba(59, 130, 246, 0.4)',
        'glow-green': '0 0 20px rgba(16, 185, 129, 0.5)',
        'glow-amber': '0 0 20px rgba(245, 158, 11, 0.5)',
        'glow-red': '0 0 20px rgba(239, 68, 68, 0.5)'
      },
      backgroundImage: {
        'gradient-energy': 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        'gradient-renewable': 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
        'gradient-primary': 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)'
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'glow': 'glow 2s ease-in-out infinite alternate'
      },
      keyframes: {
        glow: {
          '0%': { boxShadow: '0 0 5px rgba(59, 130, 246, 0.5)' },
          '100%': { boxShadow: '0 0 20px rgba(59, 130, 246, 0.8)' }
        }
      },
      gridTemplateColumns: {
        '24': 'repeat(24, minmax(0, 1fr))'
      }
    }
  },
  plugins: []
}
