package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                  string   `envconfig:"PORT" default:"8080"`
	LogLevel              string   `envconfig:"LOG_LEVEL" default:"info"`
	CORSAllowedOrigins    []string `envconfig:"CORS_ALLOWED_ORIGINS"`
	OTELCollectorEndpoint string   `envconfig:"OTEL_COLLECTOR_ENDPOINT" default:""`
	AuthServiceURL        string   `envconfig:"AUTH_SERVICE_URL"`

	Neo4jURI      string `envconfig:"NEO4J_URI"`
	Neo4jUser     string `envconfig:"NEO4J_USERNAME"`
	Neo4jPassword string `envconfig:"NEO4J_PASSWORD"`

	GCPProjectID       string `envconfig:"GCP_PROJECT_ID"`
	PubSubSubscription string `envconfig:"PUBSUB_SUBSCRIPTION_ID"`
	PubSubEmulatorHost string `envconfig:"PUBSUB_EMULATOR_HOST"`
}

func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
