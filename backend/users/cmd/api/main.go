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

	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/shared/middleware"
	"github.com/ShalArl/trip-manager/backend/users/config"
	"github.com/ShalArl/trip-manager/backend/users/database"
	"github.com/ShalArl/trip-manager/backend/users/handler"
	"github.com/ShalArl/trip-manager/backend/users/repository"
	"github.com/ShalArl/trip-manager/backend/users/service"
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
	defer db.Close()

	// Migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Auth client
	authClient := authclient.NewClient(cfg.AuthServiceURL)

	// Wire up
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	// Middleware
	requireAuth := authclient.RequireAuth(authClient)

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("POST /provision", requireAuth(handler.ProvisionHandler(svc)))
	mux.HandleFunc("GET /me", requireAuth(handler.GetMeHandler(svc)))
	mux.HandleFunc("PUT /me", requireAuth(handler.UpdateMeHandler(svc)))
	mux.HandleFunc("GET /{id}", handler.GetByIDHandler(svc))

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
			log.Fatalf("Failed to shutdown server: %v", err)
		}
	}()

	log.Printf("users service starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
	log.Println("Server stopped")
}
