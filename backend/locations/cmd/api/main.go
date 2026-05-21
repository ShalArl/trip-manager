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

	"github.com/ShalArl/trip-manager/backend/locations/client"
	"github.com/ShalArl/trip-manager/backend/locations/config"
	"github.com/ShalArl/trip-manager/backend/locations/database"
	"github.com/ShalArl/trip-manager/backend/locations/internal/location"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	// DB
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Clients
	authClient := authclient.NewClient(cfg.AuthServiceURL)
	usersClient := client.NewUsersClient(cfg.UsersServiceURL)
	requireAuth := authclient.RequireAuth(authClient)
	optionalAuth := authclient.OptionalAuth(authClient)

	// Wire up
	repo := location.NewRepository(db)
	svc := location.NewService(repo)

	// Router
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Locations
	mux.HandleFunc("GET /api/trips/{tripId}/locations",
		optionalAuth(location.ListHandler(svc, cfg.S3Endpoint, cfg.S3Bucket)))
	mux.HandleFunc("POST /api/trips/{tripId}/locations",
		requireAuth(location.CreateHandler(svc, usersClient, cfg.S3Endpoint, cfg.S3Bucket)))
	mux.HandleFunc("PUT /api/trips/{tripId}/locations/{locationId}",
		requireAuth(location.UpdateHandler(svc, cfg.S3Endpoint, cfg.S3Bucket)))
	mux.HandleFunc("DELETE /api/trips/{tripId}/locations/{locationId}",
		requireAuth(location.DeleteHandler(svc)))

	// Location images
	mux.HandleFunc("POST /api/trips/{tripId}/locations/{locationId}/images",
		requireAuth(location.AddImageHandler(svc, cfg.S3Endpoint, cfg.S3Bucket)))
	mux.HandleFunc("DELETE /api/trips/{tripId}/locations/{locationId}/images/{imageId}",
		requireAuth(location.DeleteImageHandler(svc)))

	// Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
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
			log.Fatalf("failed to shutdown server: %v", err)
		}
	}()

	log.Printf("locations service starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
	log.Println("Server stopped")
}
