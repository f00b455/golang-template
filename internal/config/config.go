package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	Port          string
	Environment   string
	SpiegelRSSURL string
}

// Load creates a new Config instance with values from environment variables.
func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "3002"),
		Environment:   getEnv("ENV", "development"),
		SpiegelRSSURL: getEnv("SPIEGEL_RSS_URL", "https://www.spiegel.de/schlagzeilen/index.rss"),
	}
}

// getEnv returns the value of the environment variable or the default value if not set.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
