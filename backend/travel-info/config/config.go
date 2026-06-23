package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port                  string   `envconfig:"PORT" default:"8080"`
	RedisUrl              string   `envconfig:"REDIS_URL"`
	CORSAllowedOrigins    []string `envconfig:"CORS_ALLOWED_ORIGINS"`
	OTELCollectorEndpoint string   `envconfig:"OTEL_COLLECTOR_ENDPOINT" default:""`
	AuthServiceURL        string   `envconfig:"AUTH_SERVICE_URL"`
	// Travel warning
	TravelWarningUrl string `envconfig:"TRAVEL_WARNING_URL"`
	// Weather
	WeatherAPIUrl        string `envconfig:"WEATHER_API_URL"`
	WeatherCacheTTLHours int    `envconfig:"FORECAST_DAYS" default:"3"`
	WeatherForecastDays  int    `envconfig:"CACHE_TTL_HOURS" default:"24"`
	LocationsDBURL       string `envconfig:"LOCATIONS_DB_URL"`
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
