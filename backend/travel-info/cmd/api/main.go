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

	"github.com/ShalArl/trip-manager/backend/travel-info/config"
	"github.com/ShalArl/trip-manager/backend/travel-info/internal/cache"
	"github.com/ShalArl/trip-manager/backend/travel-info/internal/fetcher"
	"github.com/ShalArl/trip-manager/backend/travel-info/internal/handler"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 1. Gemeinsames OpenTelemetry Setup
	otelProvider, err := sharedotel.New(ctx, "travel-and-weather-info", cfg.OTELCollectorEndpoint)
	if err != nil {
		log.Printf("warn: failed to initialize otel: %v", err)
	}
	var metrics *sharedotel.ServiceMetrics
	if otelProvider != nil {
		defer otelProvider.Shutdown(ctx)
		metrics, _ = sharedotel.NewServiceMetrics(otelProvider.Meter, "travel-and-weather-info")
	}

	// 2. Gemeinsamer Redis-Client-Pool für beide Domänen
	sharedCache, err := cache.NewCache(cfg.RedisUrl, cfg.WeatherCacheTTLHours)
	if err != nil {
		log.Fatal(err)
	}
	defer sharedCache.Close()

	meteoClient := fetcher.NewOpenMeteoClient(cfg.WeatherAPIUrl, cfg.WeatherForecastDays)

	warningHandler := handler.NewWarningHandler(sharedCache)
	weatherHandler := handler.NewWeatherHandler(sharedCache, meteoClient)

	// 5. Router aufsetzen und Routen mappen
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	// Explizite Endpunkte zur Vermeidung von Namenskonflikten
	mux.HandleFunc("GET /info/warning/{countryCode}", warningHandler.GetWarning)
	mux.HandleFunc("GET /info/weather", weatherHandler.GetWeather)

	// 6. CORS & Auth-Client konfigurieren (identisch für beide)
	corsConfig := middleware.DefaultCORSConfig()
	if len(cfg.CORSAllowedOrigins) == 0 {
		log.Fatalf("No allowed origin configured")
	}
	corsConfig.AllowedOrigins = cfg.CORSAllowedOrigins

	authClient := authclient.NewClient(cfg.AuthServiceURL)

	// 7. Der einzige HTTP-Server, der ab jetzt im Pod läuft
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.CORS(corsConfig)(sharedotel.MetricsMiddleware(metrics, authClient)(mux)),
	}

	// 8. Graceful Shutdown
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		<-sigch
		log.Println("Shutting down integrated server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
	}()

	log.Printf("travel-and-weather-info: listening on port %s\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: could not listen on port %s: %v\n", cfg.Port, err)
	}

	log.Println("Server stopped")
}
