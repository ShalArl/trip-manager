package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
)

// Config holds application configuration
type Config struct {
	CORSAllowedOrigins []string
	DatabaseURL        string
	ServerPort         string
	Environment        string
	JWTSecret          string
	TokenExpiration    time.Duration

	Storage        StorageConfig
	FirebaseConfig *firebase.Config
}

type StorageConfig struct {
	Type         string // "gcs", "s3", "local"
	SignedURLTTL time.Duration
	// S3 / MinIO
	S3 S3Settings
	// GCS
	GCS GCSSettings
}

type S3Settings struct {
	Endpoint  string
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
	PublicURL string
	UseSSL    bool
}

type GCSSettings struct {
	Bucket   string
	SignerSA string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Build database URL from individual components
	databaseURL := buildDatabaseURL()
	tokenExp := getEnvDuration("TOKEN_EXPIRATION", 7*24*time.Hour)

	firebaseCfg := loadFirebaseConfig()
	if firebaseCfg.ProjectID == "" {
		return nil, fmt.Errorf("FIREBASE_PROJECT_ID is required")
	}

	cfg := &Config{
		CORSAllowedOrigins: parseOrigins(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		DatabaseURL:        databaseURL,
		ServerPort:         getEnv("SERVER_PORT", "8000"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		TokenExpiration:    tokenExp,

		Storage:        loadStorageConfig(),
		FirebaseConfig: firebaseCfg,
	}

	return cfg, nil
}

func loadStorageConfig() StorageConfig {
	return StorageConfig{
		Type:         getEnv("STORAGE_TYPE", ""),
		SignedURLTTL: getEnvDuration("SIGNED_URL_TTL", 15*time.Minute),
		S3:           loadS3Config(),
		GCS:          loadGCSConfig(),
	}
}

func loadFirebaseConfig() *firebase.Config {
	return &firebase.Config{
		ProjectID: getEnv("FIREBASE_PROJECT_ID", ""),
	}
}

func loadGCSConfig() GCSSettings {
	return GCSSettings{
		Bucket:   getEnv("GCS_BUCKET", "trip-manager"),
		SignerSA: getEnv("GCS_SIGNER_SA", "minioadmin"),
	}
}

func loadS3Config() S3Settings {
	return S3Settings{
		Endpoint:  getEnv("S3_ENDPOINT", "http://localhost:9000"),
		Bucket:    getEnv("S3_BUCKET", "trip-manager"),
		Region:    getEnv("S3_REGION", "us-east-1"),
		AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
		SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
		PublicURL: getEnv("S3_PUBLIC_URL", "http://localhost:9000/trip-manager"),
		UseSSL:    getEnvBool("S3_USE_SSL", false),
	}
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

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err != nil {
			return defaultValue
		}
		return duration
	}
	return defaultValue
}

// parseOrigins parses a comma-separated list of origins into a slice
func parseOrigins(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// String returns a string representation of the config (without sensitive data)
func (c Config) String() string {
	return fmt.Sprintf(
		"Config{Port: %s, Environment: %s, StorageType: %s}",
		c.ServerPort,
		c.Environment,
		c.Storage.Type,
	)
}
