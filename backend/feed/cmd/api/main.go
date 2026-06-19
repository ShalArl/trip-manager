package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	sharedotel "otel"
	"syscall"
	"time"

	"github.com/ShalArl/trip-manager/backend/feed/config"
	"github.com/ShalArl/trip-manager/backend/feed/internal/feed"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/shared/middleware"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	otelProvider, err := sharedotel.New(ctx, "feed", cfg.OTELCollectorEndpoint)
	if err != nil {
		log.Printf("warn: failed to initialize otel: %v", err)
	}
	var metrics *sharedotel.ServiceMetrics
	if otelProvider != nil {
		defer otelProvider.Shutdown(ctx)
		metrics, _ = sharedotel.NewServiceMetrics(otelProvider.Meter, "feed")
	}

	corsConfig := middleware.DefaultCORSConfig()
	allowedOrigins := cfg.CORSAllowedOrigins
	if len(allowedOrigins) == 0 {
		log.Fatalf("feed: no allowed origin configured")
	}
	corsConfig.AllowedOrigins = allowedOrigins

	// Neo4j
	driver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jURI,
		neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""),
	)
	if err != nil {
		log.Fatalf("feed: failed to create neo4j driver: %v", err)
	}
	defer func(driver neo4j.DriverWithContext, ctx context.Context) {
		err := driver.Close(ctx)
		if err != nil {
			log.Fatalf("feed: failed to close neo4j driver: %v", err)
		}
	}(driver, context.Background())

	if err := driver.VerifyConnectivity(context.Background()); err != nil {
		log.Fatalf("feed: neo4j not reachable: %v", err)
	}
	log.Printf("feed: connected to neo4j at %s", cfg.Neo4jURI)

	// Wire up
	feedRepo := feed.NewRepository(driver)
	feedSvc := feed.NewService(feedRepo)

	// Auth
	authClient := authclient.NewClient(cfg.AuthServiceURL)
	requireAuth := authclient.RequireAuth(authClient)

	// Router
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			log.Printf("feed: failed to write health response: %v", err)
			return
		}
	})

	// Globaler Feed – öffentlich, Gäste + eingeloggte User
	mux.HandleFunc("GET /", requireAuth(feed.GetGlobalFeedHandler(feedSvc)))

	// Personalisierter Feed – nur für eingeloggte User
	mux.HandleFunc("GET /personal", requireAuth(feed.GetPersonalFeedHandler(feedSvc)))

	// Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.CORS(corsConfig)(sharedotel.MetricsMiddleware(metrics, authClient)(mux)),
	}

	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		<-sigch
		log.Println("feed: shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("feed: shutdown error: %v", err)
		}
	}()

	log.Printf("feed service starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("feed: server error: %v", err)
	}
	log.Println("feed: stopped")
}
