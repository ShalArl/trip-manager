package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port               string   `envconfig:"PORT" default:"8001"`
	DatabaseURL        string   `envconfig:"DATABASE_URL"`
	AuthServiceURL     string   `envconfig:"AUTH_SERVICE_URL"`
	LogLevel           string   `envconfig:"LOG_LEVEL"`
	CORSAllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS"`
	FirebaseProjectID  string   `envconfig:"FIREBASE_PROJECT_ID"`
}

func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
