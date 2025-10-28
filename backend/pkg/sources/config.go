// Package sources provides configuration for data source URLs.
// Supports environment variable overrides for real HTTP fetching.
package sources

import (
	"os"
)

// Config holds URLs for all external data sources.
type Config struct {
	TEPCO  SourceConfig
	Kansai SourceConfig
	OCCTO  SourceConfig
	JEPX   SourceConfig
}

// SourceConfig represents a single data source configuration.
type SourceConfig struct {
	URL         string // Base URL for the data source
	Name        string // Display name (e.g., "TEPCO")
	FallbackURL string // Optional fallback URL
}

// LoadConfig loads source URLs from environment variables with defaults.
// Environment variables:
//   - TEPCO_URL: Tokyo Electric Power Company demand data
//   - KANSAI_URL: Kansai Electric Power Company demand data
//   - OCCTO_URL: Organization for Cross-regional Coordination of Transmission Operators reserve margin
//   - JEPX_URL: Japan Electric Power Exchange spot prices
func LoadConfig() Config {
	return Config{
		TEPCO: SourceConfig{
			URL:  getEnv("TEPCO_URL", "https://www.tepco.co.jp/forecast/html/download-j.html"),
			Name: "TEPCO",
		},
		Kansai: SourceConfig{
			URL:  getEnv("KANSAI_URL", "https://www.kansai-td.co.jp/denkiyoho/download.html"),
			Name: "Kansai Electric",
		},
		OCCTO: SourceConfig{
			URL:  getEnv("OCCTO_URL", "https://www.occto.or.jp/"),
			Name: "OCCTO",
		},
		JEPX: SourceConfig{
			URL:  getEnv("JEPX_URL", "https://www.jepx.jp/"),
			Name: "JEPX",
		},
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
