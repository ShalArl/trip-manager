package config

import "os"

type Config struct {
	Port               string
	GCPProjectID       string
	PubSubSubscription string
	Neo4jURI           string
	Neo4jUser          string
	Neo4jPassword      string
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		GCPProjectID:       getEnv("GCP_PROJECT_ID", ""),
		PubSubSubscription: getEnv("PUBSUB_SUBSCRIPTION_ID", "trip-events-sub"),
		Neo4jURI:           getEnv("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:          getEnv("NEO4J_USERNAME", "neo4j"),
		Neo4jPassword:      getEnv("NEO4J_PASSWORD", "neo4jpassword"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
