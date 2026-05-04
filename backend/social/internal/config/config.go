package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// Config holds application configuration
type Config struct {
	Port                       string
	FirestoreProject           string
	LogLevel                   string
	AuthClientConnectionString string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:                       getEnv("PORT", "8080"),
		FirestoreProject:           getEnv("FIRESTORE_PROJECT_ID", "trip-manager-local"),
		LogLevel:                   getEnv("LOG_LEVEL", "info"),
		AuthClientConnectionString: getEnv("AUTH_CLIENT_CONNECTION_STRING", ""),
	}
}

// ConnectFirestore establishes a connection to Firestore
func ConnectFirestore(ctx context.Context, projectID string) (*firestore.Client, error) {
	if projectID == "" {
		return nil, fmt.Errorf("firestore project id required")
	}

	var opts []option.ClientOption

	// Check for local emulator
	if emuHost := os.Getenv("FIRESTORE_EMULATOR_HOST"); emuHost != "" {
		log.Printf("Using Firestore Emulator: %s\n", emuHost)
	}

	// Check for credentials file
	if credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credFile != "" {
		opts = append(opts, option.WithCredentialsFile(credFile))
	}

	client, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to firestore: %w", err)
	}

	log.Println("Successfully connected to firestore")
	return client, nil
}

// getEnv returns environment variable value or default
func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultVal
}
