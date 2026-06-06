package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port               string `envconfig:"PORT" default:"8080"`
	GCPProjectID       string `envconfig:"GCP_PROJECT_ID"`
	PubSubSubscription string `envconfig:"PUBSUB_SUBSCRIPTION_ID"`
	Neo4jURI           string `envconfig:"NEO4J_URI"`
	Neo4jUser          string `envconfig:"NEO4J_USERNAME"`
	Neo4jPassword      string `envconfig:"NEO4J_PASSWORD"`
	LogLevel           string `envconfig:"LOG_LEVEL"`
	PubSubEmulatorHost string `envconfig:"PUBSUB_EMULATOR_HOST"`
}

func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
