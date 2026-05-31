package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShalArl/trip-manager/backend/newsletter/internal/newsletter"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type config struct {
	Port           string
	Neo4jURI       string
	Neo4jUser      string
	Neo4jPassword  string
	AuthServiceURL string
}

func loadConfig() config {
	return config{
		Port:           getEnv("PORT", "8008"),
		Neo4jURI:       getEnv("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:      getEnv("NEO4J_USER", "neo4j"),
		Neo4jPassword:  getEnv("NEO4J_PASSWORD", "password"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8082"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	cfg := loadConfig()

	driver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jURI,
		neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""),
	)
	if err != nil {
		log.Fatalf("newsletter: failed to create neo4j driver: %v", err)
	}
	defer func(driver neo4j.DriverWithContext, ctx context.Context) {
		err := driver.Close(ctx)
		if err != nil {
			log.Printf("newsletter: failed to close neo4j driver: %v", err)
		}
	}(driver, context.Background())

	if err := driver.VerifyConnectivity(context.Background()); err != nil {
		log.Fatalf("newsletter: neo4j not reachable: %v", err)
	}
	log.Printf("newsletter: connected to neo4j at %s", cfg.Neo4jURI)

	repo := newsletter.NewRepository(driver)
	svc := newsletter.NewService(repo)

	authClient := authclient.NewClient(cfg.AuthServiceURL)
	requireAuth := authclient.RequireAuth(authClient)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			log.Printf("newsletter: error writing health response: %v", err)
			return
		}
	})

	mux.HandleFunc("GET /api/newsletter", requireAuth(newsletter.GetNewsletterHandler(svc)))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		<-sigch
		log.Println("newsletter: shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("newsletter: shutdown error: %v", err)
		}
	}()

	log.Printf("newsletter service starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("newsletter: server error: %v", err)
	}
	log.Println("newsletter: stopped")
}
