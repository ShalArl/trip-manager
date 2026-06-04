package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port               string        `envconfig:"LOG_LEVEL"`
	Bucket             string        `envconfig:"BUCKET"`
	TTL                time.Duration `envconfig:"TTL" default:"15m"`
	LogLevel           string        `envconfig:"LOG_LEVEL"`
	Type               string        `envconfig:"STORAGE_PROVIDER"`
	Endpoint           string        `envconfig:"S3_ENDPOINT"`
	Region             string        `envconfig:"S3_REGION"`
	AccessKey          string        `envconfig:"S3_ACCESS_KEY"`
	SecretKey          string        `envconfig:"S3_SECRET_KEY"`
	PublicURL          string        `envconfig:"S3_PUBLIC_URL"`
	UseSSL             bool          `envconfig:"S3_USE_SSL" default:"false"`
	SignerSA           string        `envconfig:"GCS_SIGNER_SA"`
	CORSAllowedOrigins []string      `envconfig:"CORS_ALLOWED_ORIGINS"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
