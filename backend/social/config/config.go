package config

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/api/option"
)

type Config struct {
	Port                         string   `envconfig:"PORT" default:"8080"`
	FirestoreProject             string   `envconfig:"FIRESTORE_PROJECT_ID" required:"true"`
	LogLevel                     string   `envconfig:"LOG_LEVEL" default:"info"`
	AuthClientConnectionString   string   `envconfig:"AUTH_CLIENT_CONNECTION_STRING" required:"true"`
	UsersServiceURL              string   `envconfig:"USERS_SERVICE_URL" required:"true"`
	GCPProjectID                 string   `envconfig:"GCP_PROJECT_ID" required:"true"`
	PubSubTopicID                string   `envconfig:"PUBSUB_TOPIC_ID" required:"true"`
	FirestoreEmulatorHost        string   `envconfig:"FIRESTORE_EMULATOR_HOST" default:""`
	GoogleApplicationCredentials string   `envconfig:"GOOGLE_APPLICATION_CREDENTIALS"`
	CORSAllowedOrigins           []string `envconfig:"CORS_ALLOWED_ORIGINS"`
	PubSubEmulatorHost           string   `envconfig:"PUBSUB_EMULATOR_HOST"`
	OTELCollectorEndpoint        string   `envconfig:"OTEL_COLLECTOR_ENDPOINT" default:""`
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func ConnectFirestore(ctx context.Context, cfg Config) (*firestore.Client, error) {
	if cfg.FirestoreProject == "" {
		return nil, fmt.Errorf("firestore project id required")
	}
	var opts []option.ClientOption
	if cfg.FirestoreEmulatorHost != "" {
		log.Printf("Using Firestore Emulator")
	}
	if cfg.GoogleApplicationCredentials != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.GoogleApplicationCredentials))
	}
	client, err := firestore.NewClient(ctx, cfg.FirestoreProject, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to firestore: %w", err)
	}
	log.Println("Successfully connected to firestore")
	return client, nil
}
