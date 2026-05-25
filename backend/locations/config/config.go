package config

import "os"

type Config struct {
	Port            string
	DatabaseURL     string
	AuthServiceURL  string
	UsersServiceURL string
	S3Endpoint      string
	S3Bucket        string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8005"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/locations_db?sslmode=disable"),
		AuthServiceURL:  getEnv("AUTH_SERVICE_URL", "http://localhost:8082"),
		UsersServiceURL: getEnv("USERS_SERVICE_URL", "http://localhost:8001"),
		S3Endpoint:      getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3Bucket:        getEnv("S3_BUCKET", "trip-manager"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
