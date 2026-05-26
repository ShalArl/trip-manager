package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port               string   `envconfig:"PORT" default:"8080"`
	APIUrl             string   `envconfig:"API_URL"`
	RedisUrl           string   `envconfig:"REDIS_URL"`
	CORSAllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS"`
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
