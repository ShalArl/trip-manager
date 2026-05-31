package config

import "os"

type Config struct {
	KafkaBrokers  string
	KafkaGroupID  string
	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string
}

func Load() *Config {
	return &Config{
		KafkaBrokers:  getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaGroupID:  getEnv("KAFKA_GROUP_ID", "feed-generator"),
		Neo4jURI:      getEnv("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:     getEnv("NEO4J_USERNAME", "neo4j"),
		Neo4jPassword: getEnv("NEO4J_PASSWORD", "neo4jpassword"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
