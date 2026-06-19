package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                  string   `envconfig:"PORT" default:"8007"`
	AuthServiceURL        string   `envconfig:"AUTH_CLIENT_CONNECTION_STRING"`
	Neo4jURI              string   `envconfig:"NEO4J_URI"`
	Neo4jUser             string   `envconfig:"NEO4J_USERNAME"`
	Neo4jPassword         string   `envconfig:"NEO4J_PASSWORD"`
	LogLevel              string   `envconfig:"LOG_LEVEL"`
	CORSAllowedOrigins    []string `envconfig:"CORS_ALLOWED_ORIGINS"`
	OTELCollectorEndpoint string   `envconfig:"OTEL_COLLECTOR_ENDPOINT" default:""`
}

func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
