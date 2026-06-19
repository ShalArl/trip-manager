package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port                  string   `envconfig:"PORT" default:"8080"`
	APIUrl                string   `envconfig:"API_URL"`
	RedisUrl              string   `envconfig:"REDIS_URL"`
	CORSAllowedOrigins    []string `envconfig:"CORS_ALLOWED_ORIGINS"`
	LocationsDBURL        string   `envconfig:"LOCATIONS_DB_URL"`
	ForecastDays          int      `envconfig:"FORECAST_DAYS" default:"3"`
	CacheTTLHours         int      `envconfig:"CACHE_TTL_HOURS" default:"24"`
	OTELCollectorEndpoint string   `envconfig:"OTEL_COLLECTOR_ENDPOINT" default:""`
	AuthServiceURL        string   `envconfig:"AUTH_SERVICE_URL"`
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
