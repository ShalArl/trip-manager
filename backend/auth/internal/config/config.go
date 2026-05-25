package config

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

// Config holds application configuration
type Config struct {
	Port               string
	FirebaseProjectID  string
	LogLevel           string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:              getEnv("PORT", "8082"),
		FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", "trip-manager-local"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
	}
}

// InitializeFirebase initializes Firebase and returns auth client
func InitializeFirebase(ctx context.Context, cfg *Config) (*auth.Client, error) {
	// Check for emulator
	if emuHost := os.Getenv("FIREBASE_AUTH_EMULATOR_HOST"); emuHost != "" {
		log.Printf("Using Firebase Auth Emulator: %s\n", emuHost)
	}

	// Initialize Firebase app
	fbConfig := &firebase.Config{
		ProjectID: cfg.FirebaseProjectID,
	}

	app, err := firebase.NewApp(ctx, fbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Get Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Firebase Auth client: %w", err)
	}

	log.Println("Successfully initialized Firebase Auth")
	return authClient, nil
}

// getEnv returns environment variable value or default
func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultVal
}

