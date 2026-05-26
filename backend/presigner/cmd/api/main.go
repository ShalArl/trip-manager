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

	"github.com/ShalArl/trip-manager/backend/presigner/config"
	"github.com/ShalArl/trip-manager/backend/presigner/internal/handler"
	"github.com/ShalArl/trip-manager/backend/presigner/internal/provider"
	"github.com/ShalArl/trip-manager/backend/presigner/internal/service"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/shared/middleware"
)

func main() {
	ctx := context.Background()

	// Load cache
	cfg := config.LoadConfig()
	log.Printf("Starting Presigner Service on port %s\n", cfg.Port)

	corsConfig := middleware.DefaultCORSConfig()
	corsConfig.AllowedOrigins = []string{
		"https://neatnode.xyz",
		"https://www.neatnode.xyz",
	}

	// Load storage provider configuration
	storageCfg := config.LoadStorageConfig()

	// Create storage provider (S3 or GCS)
	storage, err := provider.NewFromEnv(ctx, storageCfg)
	if err != nil {
		log.Fatalf("Failed to create storage provider: %v", err)
	}

	// Create presigner service
	presignerService := service.NewService(storage, cfg.TTL)
	authClient := authclient.NewClient(os.Getenv("AUTH_SERVICE_URL"))

	// Setup routes
	mux := http.NewServeMux()

	// Presigner endpoints
	mux.HandleFunc("POST /uploads/presigned",
		authclient.RequireAuth(authClient)(handler.GetPresignedUploadURLHandler(presignerService, authClient)))

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Start server
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
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
