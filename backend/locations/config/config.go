package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port            string `envconfig:"PORT" default:"8005"`
	DatabaseURL     string `envconfig:"DATABASE_URL"`
	AuthServiceURL  string `envconfig:"AUTH_SERVICE_URL"`
	UsersServiceURL string `envconfig:"NEO4J_URI"`
	S3Endpoint      string `envconfig:"S3_ENDPOINT"`
	S3Bucket        string `envconfig:"S3_BUCKET"`
	LogLevel        string `envconfig:"LOG_LEVEL"`
}

func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
