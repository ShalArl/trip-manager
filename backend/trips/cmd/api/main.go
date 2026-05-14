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
	"github.com/ShalArl/trip-manager/backend/trips/client"
	"github.com/ShalArl/trip-manager/backend/trips/config"
	"github.com/ShalArl/trip-manager/backend/trips/database"
	"github.com/ShalArl/trip-manager/backend/trips/handler"
	"github.com/ShalArl/trip-manager/backend/trips/repository"
	"github.com/ShalArl/trip-manager/backend/trips/service"
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
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	// Router
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("GET /api/trips", requireAuth(handler.ListTripsHandler(svc, usersClient)))
	mux.HandleFunc("POST /api/trips", requireAuth(handler.CreateTripHandler(svc, usersClient)))
	mux.HandleFunc("GET /api/trips/recent", optionalAuth(handler.ListRecentTripsHandler(svc)))
	mux.HandleFunc("GET /api/trips/search", optionalAuth(handler.SearchTripsHandler(svc)))
	mux.HandleFunc("GET /api/trips/{tripId}", optionalAuth(handler.GetTripHandler(svc)))
	mux.HandleFunc("PUT /api/trips/{tripId}", requireAuth(handler.UpdateTripHandler(svc, usersClient)))
	mux.HandleFunc("DELETE /api/trips/{tripId}", requireAuth(handler.DeleteTripHandler(svc, usersClient)))

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

	log.Printf("trips service starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
	log.Println("Server stopped")
}
