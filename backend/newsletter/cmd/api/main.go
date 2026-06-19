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

	"github.com/ShalArl/trip-manager/backend/newsletter/internal/newsletter"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/shared/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

type config struct {
	Port                  string   `envconfig:"PORT" default:"8008"`
	NewsletterDBURL       string   `envconfig:"NEWSLETTER_DB_URL"`
	AuthServiceURL        string   `envconfig:"AUTH_SERVICE_URL"`
	LogLevel              string   `envconfig:"LOG_LEVEL"`
	CORSAllowedOrigins    []string `envconfig:"CORS_ALLOWED_ORIGINS"`
	OTELCollectorEndpoint string   `envconfig:"OTEL_COLLECTOR_ENDPOINT" default:""`
}

func load() (*config, error) {
	var config config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	ctx := context.Background()
	cfg, err := load()
	if err != nil {
		log.Fatal("Failed to load config")
	}

	otelProvider, err := sharedotel.New(ctx, "newsletter", cfg.OTELCollectorEndpoint)
	if err != nil {
		log.Printf("warn: failed to initialize otel: %v", err)
	}
	var metrics *sharedotel.ServiceMetrics
	if otelProvider != nil {
		defer otelProvider.Shutdown(ctx)
		metrics, _ = sharedotel.NewServiceMetrics(otelProvider.Meter, "newsletter")
	}

	corsConfig := middleware.DefaultCORSConfig()
	allowedOrigins := cfg.CORSAllowedOrigins
	if len(allowedOrigins) == 0 {
		log.Fatal("No allowed origins configured")
	}
	corsConfig.AllowedOrigins = allowedOrigins

	db, err := sqlx.Connect("postgres", cfg.NewsletterDBURL)
	if err != nil {
		log.Fatalf("newsletter: failed to connect to newsletter-db: %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("newsletter: failed to close db: %v", err)
		}
	}(db)
	log.Println("newsletter: connected to newsletter-db")

	repo := newsletter.NewRepository(db)
	svc := newsletter.NewService(repo)

	authClient := authclient.NewClient(cfg.AuthServiceURL)
	requireAuth := authclient.RequireAuth(authClient)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			log.Fatalf("Failed to write health check response: %v", err)
			return
		}
	})

	mux.HandleFunc("GET /", requireAuth(newsletter.GetNewsletterHandler(svc)))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.CORS(corsConfig)(sharedotel.MetricsMiddleware(metrics, authClient)(mux)),
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
