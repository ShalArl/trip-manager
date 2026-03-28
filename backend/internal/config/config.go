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
	// Build database URL from individual components
	databaseURL := buildDatabaseURL()

	cfg := &Config{
		DatabaseURL: databaseURL,
		ServerPort:  getEnv("SERVER_PORT", "8000"),
		Environment: getEnv("ENVIRONMENT", "development"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}

	return cfg, nil
}

// buildDatabaseURL constructs the database URL from environment variables or defaults
func buildDatabaseURL() string {
	// Allow complete DATABASE_URL override
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}

	// Otherwise build from individual components
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	database := getEnv("POSTGRES_DB", "trip_manager")
	sslMode := getEnv("POSTGRES_SSLMODE", "disable")

	// Build connection string
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, database, sslMode,
	)
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
