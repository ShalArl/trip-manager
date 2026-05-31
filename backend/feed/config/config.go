package config

import "os"

type Config struct {
	Port           string
	AuthServiceURL string
	Neo4jURI       string
	Neo4jUser      string
	Neo4jPassword  string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8007"),
		AuthServiceURL: getEnv("AUTH_CLIENT_CONNECTION_STRING", "http://localhost:8082"),
		Neo4jURI:       getEnv("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:      getEnv("NEO4J_USERNAME", "neo4j"),
		Neo4jPassword:  getEnv("NEO4J_PASSWORD", "neo4jpassword"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
