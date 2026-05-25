package config

import "os"

type Config struct {
	Port            string
	DatabaseURL     string
	AuthServiceURL  string
	UsersServiceURL string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8002"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		AuthServiceURL:  getEnv("AUTH_SERVICE_URL", "http://localhost:8082"),
		UsersServiceURL: getEnv("USERS_SERVICE_URL", "http://localhost:8001"),
	}
}
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
