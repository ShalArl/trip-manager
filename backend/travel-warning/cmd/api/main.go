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

	"github.com/ShalArl/trip-manager/backend/shared/middleware"
	"github.com/ShalArl/trip-manager/backend/travel-warning/config"
	"github.com/ShalArl/trip-manager/backend/travel-warning/internal/cache"
	"github.com/ShalArl/trip-manager/backend/travel-warning/internal/handler"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
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
	corsConfig.AllowedOrigins = []string{
		os.Getenv("CORS_ALLOWED_ORIGINS"),
	}

	if corsConfig.AllowedOrigins[0] == "" {
		corsConfig.AllowedOrigins = []string{"*"}
	}

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.CORS(corsConfig)(mux),
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
