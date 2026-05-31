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

	"github.com/ShalArl/trip-manager/backend/locations/config"
	"github.com/ShalArl/trip-manager/backend/locations/database"
	"github.com/ShalArl/trip-manager/backend/locations/internal/location"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/shared/middleware"
	"github.com/ShalArl/trip-manager/backend/shared/userclient"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	corsConfig := middleware.DefaultCORSConfig()
	corsConfig.AllowedOrigins = []string{
		"https://neatnode.xyz",
		"https://www.neatnode.xyz",
	}

	// DB
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
	}(db)

	// Migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Clients
	authClient := authclient.NewClient(cfg.AuthServiceURL)
	usersClient := userclient.NewUsersClient(cfg.UsersServiceURL)
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
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			return
		}
	})

	// Locations
	mux.HandleFunc("GET /{tripId}",
		optionalAuth(location.ListHandler(svc, cfg.S3Endpoint, cfg.S3Bucket)))
	mux.HandleFunc("POST /{tripId}",
		requireAuth(location.CreateHandler(svc, usersClient, cfg.S3Endpoint, cfg.S3Bucket)))
	mux.HandleFunc("PUT /{tripId}/{locationId}",
		requireAuth(location.UpdateHandler(svc, cfg.S3Endpoint, cfg.S3Bucket)))
	mux.HandleFunc("DELETE /{tripId}/{locationId}",
		requireAuth(location.DeleteHandler(svc)))

	// Location images
	mux.HandleFunc("POST /{tripId}/{locationId}/images",
		requireAuth(location.AddImageHandler(svc, cfg.S3Endpoint, cfg.S3Bucket)))
	mux.HandleFunc("DELETE /{tripId}/{locationId}/images/{imageId}",
		requireAuth(location.DeleteImageHandler(svc)))

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
			log.Fatalf("failed to shutdown server: %v", err)
		}
	}()

	log.Printf("locations service starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
	log.Println("Server stopped")
}
