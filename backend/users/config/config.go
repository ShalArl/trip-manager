package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port                     string   `envconfig:"PORT" default:"8001"`
	AuthServiceURL           string   `envconfig:"AUTH_SERVICE_URL"`
	LogLevel                 string   `envconfig:"LOG_LEVEL"`
	CORSAllowedOrigins       []string `envconfig:"CORS_ALLOWED_ORIGINS"`
	FirebaseProjectID        string   `envconfig:"FIREBASE_PROJECT_ID"`
	FirebaseAuthEmulatorHost string   `envconfig:"FIREBASE_AUTH_EMULATOR_HOST" default:""`
	PrometheusURL            string   `envconfig:"PROMETHEUS_URL" default:""`
	OTELCollectorEndpoint    string   `envconfig:"OTEL_COLLECTOR_ENDPOINT" default:""`
	BaseUrl                  string   `envconfig:"BASE_URL" default:""`
	GCPProjectID             string   `envconfig:"GCP_PROJECT_ID" default:""`
	DatabaseURL              string   `envconfig:"DATABASE_URL"`
	MigrationDBURL           string   `envconfig:"MIGRATION_DB_URL"`
	AppDBPassword            string   `envconfig:"APP_DB_PASSWORD"`
	GitHubToken              string   `envconfig:"GITHUB_TOKEN" default:""`
	GitHubRepo               string   `envconfig:"GITHUB_REPO" default:""`
	GitHubBranch             string   `envconfig:"GITHUB_BRANCH" default:""`
}

func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
