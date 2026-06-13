package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port               string   `envconfig:"PORT" default:"8002"`
	DatabaseURL        string   `envconfig:"DATABASE_URL"`
	AuthServiceURL     string   `envconfig:"AUTH_SERVICE_URL"`
	UsersServiceURL    string   `envconfig:"USERS_SERVICE_URL"`
	GCPProjectID       string   `envconfig:"GCP_PROJECT_ID"`
	PubSubTopicID      string   `envconfig:"PUBSUB_TOPIC_ID" default:""`
	CORSAllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS"`
	PubSubEmulatorHost string   `envconfig:"PUBSUB_EMULATOR_HOST"`
}

func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
