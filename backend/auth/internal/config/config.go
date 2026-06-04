package config

import (
	"context"
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

// Config holds application configuration
type Config struct {
	Port                     string `envconfig:"PORT" default:"8080"`
	FirebaseProjectID        string `envconfig:"FIREBASE_PROJECT_ID"`
	LogLevel                 string `envconfig:"LOG_LEVEL"`
	FirebaseAuthEmulatorHost string `envconfig:"FIREBASE_AUTH_EMULATOR_HOST" default:""`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// InitializeFirebase initializes Firebase and returns auth client
func InitializeFirebase(ctx context.Context, cfg *Config) (*auth.Client, error) {
	// Check for emulator
	if emuHost := cfg.FirebaseAuthEmulatorHost; emuHost != "" {
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
