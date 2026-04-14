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

	// Storage configuration
	StorageType string
	UploadDir   string // For local storage

	// S3 configuration (for MinIO or AWS S3)
	S3Endpoint   string
	S3Bucket     string
	S3Region     string
	S3AccessKey  string
	S3SecretKey  string
	S3PublicURL  string
	S3UseSSL     bool
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

		// Storage configuration
		StorageType: getEnv("STORAGE_TYPE", "local"),
		UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),

		// S3 configuration (defaults for local MinIO development)
		S3Endpoint:  getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3Bucket:    getEnv("S3_BUCKET", "trip-manager"),
		S3Region:    getEnv("S3_REGION", "us-east-1"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
		S3PublicURL: getEnv("S3_PUBLIC_URL", "http://localhost:9000/trip-manager"),
		S3UseSSL:    getEnvBool("S3_USE_SSL", false),
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

// getEnvBool reads a boolean environment variable with a fallback default
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

// String returns a string representation of the config (without sensitive data)
func (c Config) String() string {
	return fmt.Sprintf(
		"Config{Port: %s, Environment: %s, StorageType: %s}",
		c.ServerPort,
		c.Environment,
		c.StorageType,
	)
}
