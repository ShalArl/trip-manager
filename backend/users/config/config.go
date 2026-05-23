package config

import (
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	AuthServiceURL string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8001"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
