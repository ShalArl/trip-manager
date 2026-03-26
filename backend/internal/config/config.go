package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	DatabaseURL string
	ServerPort  string
	Environment string
	JWTSecret   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/trip_manager?sslmode=disable"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}

	return cfg, nil
}

// getEnv reads an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// String returns a string representation of the config (without sensitive data)
func (c Config) String() string {
	return fmt.Sprintf(
		"Config{Port: %s, Environment: %s}",
		c.ServerPort,
		c.Environment,
	)
}
