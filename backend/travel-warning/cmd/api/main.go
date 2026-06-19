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

	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/shared/middleware"
	"github.com/ShalArl/trip-manager/backend/travel-warning/config"
	"github.com/ShalArl/trip-manager/backend/travel-warning/internal/cache"
	"github.com/ShalArl/trip-manager/backend/travel-warning/internal/handler"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	otelProvider, err := sharedotel.New(ctx, "travel-warning", cfg.OTELCollectorEndpoint)
	if err != nil {
		log.Printf("warn: failed to initialize otel: %v", err)
	}
	var metrics *sharedotel.ServiceMetrics
	if otelProvider != nil {
		defer otelProvider.Shutdown(ctx)
		metrics, _ = sharedotel.NewServiceMetrics(otelProvider.Meter, "travel-warning")
	}

	warningCache, err := cache.NewWarningCache(cfg.RedisUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer warningCache.Close()

	h := handler.NewHandler(warningCache)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status": "ok"}`))
		if err != nil {
			return
		}
	})

	mux.HandleFunc("GET /{countryCode}", h.GetWarning)

	// CORS
	corsConfig := middleware.DefaultCORSConfig()
	allowedOrigins := cfg.CORSAllowedOrigins
	if len(allowedOrigins) == 0 {
		log.Fatalf("No allowed origin configured")
	}
	corsConfig.AllowedOrigins = allowedOrigins

	// Auth client
	authClient := authclient.NewClient(cfg.AuthServiceURL)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.CORS(corsConfig)(sharedotel.MetricsMiddleware(metrics, authClient)(mux)),
	}

	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		<-sigch
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
	}()

	log.Printf("travel-warning: listening on port %s\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: could not listen on port %s: %v\n", cfg.Port, err)
	}

	log.Println("Server stopped")
}
