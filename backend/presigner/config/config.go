package config

import (
	"os"
	"time"
)

type Config struct {
	Port     string
	Bucket   string
	TTL      time.Duration
	LogLevel string
}

type StorageConfig struct {
	Type         string
	SignedURLTTL time.Duration
	S3           S3Settings
	GCS          GCSSettings
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
func LoadConfig() *Config {
	ttl := 15 * time.Minute
	if ttlStr := os.Getenv("PRESIGNER_URL_EXPIRATION"); ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr); err == nil {
			ttl = d
		}
	}

	return &Config{
		Port:     getEnv("PORT", "8081"),
		Bucket:   getEnv("STORAGE_BUCKET", "trip-manager"),
		TTL:      ttl,
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// LoadStorageConfig loads storage provider configuration from environment variables
func LoadStorageConfig() StorageConfig {
	ttl := 15 * time.Minute
	if ttlStr := os.Getenv("PRESIGNER_URL_EXPIRATION"); ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr); err == nil {
			ttl = d
		}
	}

	return StorageConfig{
		Type:         getEnv("STORAGE_PROVIDER", "s3"),
		SignedURLTTL: ttl,
		S3: S3Settings{
			Endpoint:  getEnv("S3_ENDPOINT", "http://localhost:9000"),
			Bucket:    getEnv("STORAGE_BUCKET", "trip-manager"),
			Region:    getEnv("S3_REGION", "us-east-1"),
			AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
			PublicURL: getEnv("S3_PUBLIC_URL", "http://localhost:9000"),
			UseSSL:    getEnv("S3_USE_SSL", "false") == "true",
		},
		GCS: GCSSettings{
			Bucket:   getEnv("GCS_BUCKET", "trip-manager"),
			SignerSA: getEnv("GCS_SIGNER_SA", ""),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
