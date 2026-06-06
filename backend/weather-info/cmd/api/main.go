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
	"github.com/ShalArl/trip-manager/backend/weather-info/config"
	"github.com/ShalArl/trip-manager/backend/weather-info/internal/cache"
	"github.com/ShalArl/trip-manager/backend/weather-info/internal/fetcher"
	"github.com/ShalArl/trip-manager/backend/weather-info/internal/handler"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Redis Cache
	weatherCache, err := cache.NewWeatherCache(cfg.RedisUrl, cfg.CacheTTLHours)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer func(weatherCache *cache.WeatherCache) {
		err := weatherCache.Close()
		if err != nil {
			log.Fatalf("failed to close cache: %v", err)
		}
	}(weatherCache)

	// Open-Meteo Client
	meteoClient := fetcher.NewClient(cfg.APIUrl, cfg.ForecastDays)

	// Handler
	h := handler.NewHandler(weatherCache, meteoClient)

	// Router
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("GET /", h.GetWeather)

	// CORS
	corsConfig := middleware.DefaultCORSConfig()
	allowedOrigins := cfg.CORSAllowedOrigins
	if len(allowedOrigins) == 0 {
		log.Fatalf("No allowed origin configured")
	}
	corsConfig.AllowedOrigins = allowedOrigins

	// Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.CORS(corsConfig)(mux),
	}

	// Graceful shutdown
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		<-sigch
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Printf("weather-info service starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
	log.Println("Server stopped")
}
