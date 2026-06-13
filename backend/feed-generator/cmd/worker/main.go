package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShalArl/trip-manager/backend/feed-generator/config"
	"github.com/ShalArl/trip-manager/backend/feed-generator/internal/consumer"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("feed-generator: failed to load config: %v", err)
	}
	log.Printf("feed-generator: starting with config: %+v", cfg)

	// Neo4j
	driver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jURI,
		neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""),
	)

	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		err := http.ListenAndServe(":"+cfg.Port, nil)
		if err != nil {
			return
		}
	}()

	if err != nil {
		log.Fatalf("feed-generator: failed to create neo4j driver: %v", err)
	}
	defer func(driver neo4j.DriverWithContext, ctx context.Context) {
		err := driver.Close(ctx)
		if err != nil {
			log.Printf("feed-generator: failed to close neo4j driver: %v", err)
		}
	}(driver, context.Background())

	// Verbindung testen
	if err := driver.VerifyConnectivity(context.Background()); err != nil {
		log.Fatalf("feed-generator: neo4j not reachable: %v", err)
	}
	log.Printf("feed-generator: connected to neo4j at %s", cfg.Neo4jURI)

	// Schema / Indizes anlegen
	if err := consumer.SetupSchema(context.Background(), driver); err != nil {
		log.Fatalf("feed-generator: failed to setup neo4j schema: %v", err)
	}

	c, err := consumer.New(driver, cfg.GCPProjectID, cfg.PubSubSubscription)
	if err != nil {
		log.Fatalf("failed to initialize pubsub consumer: %v", err)
	}
	defer func(c *consumer.Consumer) {
		err := c.Close()
		if err != nil {
			log.Printf("feed-generator: failed to close pubsub consumer: %v", err)
		}
	}(c)

	// Context mit Graceful Shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		<-sigch
		log.Println("feed-generator: received shutdown signal")
		cancel()
	}()

	log.Println("feed-generator: starting consumers")
	c.Start(ctx) // blockiert bis ctx cancelled
	log.Println("feed-generator: stopped")
}
