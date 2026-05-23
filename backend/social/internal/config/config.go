package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

type Config struct {
	Port                       string
	FirestoreProject           string
	LogLevel                   string
	AuthClientConnectionString string
	UsersServiceURL            string
	KafkaBrokers               string
}

func LoadConfig() *Config {
	return &Config{
		Port:                       getEnv("PORT", "8080"),
		FirestoreProject:           getEnv("FIRESTORE_PROJECT_ID", "trip-manager-local"),
		LogLevel:                   getEnv("LOG_LEVEL", "info"),
		AuthClientConnectionString: getEnv("AUTH_CLIENT_CONNECTION_STRING", ""),
		UsersServiceURL:            getEnv("USERS_SERVICE_URL", "http://localhost:8001"),
		KafkaBrokers:               getEnv("KAFKA_BROKERS", ""),
	}
}

func ConnectFirestore(ctx context.Context, projectID string) (*firestore.Client, error) {
	if projectID == "" {
		return nil, fmt.Errorf("firestore project id required")
	}
	if emuHost := os.Getenv("FIRESTORE_EMULATOR_HOST"); emuHost != "" {
		log.Printf("Using Firestore Emulator: %s\n", emuHost)
	}
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to firestore: %w", err)
	}
	log.Println("Successfully connected to firestore")
	return client, nil
}

func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultVal
}
