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
	"github.com/ShalArl/trip-manager/backend/trips/internal/accommodation"
	"github.com/ShalArl/trip-manager/backend/trips/internal/transport"
	"github.com/ShalArl/trip-manager/backend/trips/internal/trip"
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

	// Wire up – trips
	tripRepo := trip.NewRepository(db)
	tripSvc := trip.NewService(tripRepo)

	// Wire up – transport
	transportRepo := transport.NewRepository(db)
	transportSvc := transport.NewService(transportRepo)

	// Wire up – accommodation
	accommodationRepo := accommodation.NewRepository(db)
	accommodationSvc := accommodation.NewService(accommodationRepo)

	// Router
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Trips
	mux.HandleFunc("GET /", requireAuth(trip.ListTripsHandler(tripSvc, usersClient)))
	mux.HandleFunc("POST /", requireAuth(trip.CreateTripHandler(tripSvc, usersClient)))
	mux.HandleFunc("GET /recent", optionalAuth(trip.ListRecentTripsHandler(tripSvc)))
	mux.HandleFunc("GET /search", optionalAuth(trip.SearchTripsHandler(tripSvc)))
	mux.HandleFunc("GET /{tripId}", optionalAuth(trip.GetTripHandler(tripSvc)))
	mux.HandleFunc("PUT /{tripId}", requireAuth(trip.UpdateTripHandler(tripSvc, usersClient)))
	mux.HandleFunc("DELETE /{tripId}", requireAuth(trip.DeleteTripHandler(tripSvc, usersClient)))

	// Transports
	mux.HandleFunc("GET /api/trips/{tripId}/transports", optionalAuth(transport.ListHandler(transportSvc)))
	mux.HandleFunc("POST /api/trips/{tripId}/transports", requireAuth(transport.CreateHandler(transportSvc, usersClient)))
	mux.HandleFunc("PUT /api/trips/{tripId}/transports/{transportId}", requireAuth(transport.UpdateHandler(transportSvc, usersClient)))
	mux.HandleFunc("DELETE /api/trips/{tripId}/transports/{transportId}", requireAuth(transport.DeleteHandler(transportSvc, usersClient)))

	// Accommodations
	mux.HandleFunc("GET /api/trips/{tripId}/accommodations", optionalAuth(accommodation.ListHandler(accommodationSvc)))
	mux.HandleFunc("POST /api/trips/{tripId}/accommodations", requireAuth(accommodation.CreateHandler(accommodationSvc, usersClient)))
	mux.HandleFunc("PUT /api/trips/{tripId}/accommodations/{accommodationId}", requireAuth(accommodation.UpdateHandler(accommodationSvc, usersClient)))
	mux.HandleFunc("DELETE /api/trips/{tripId}/accommodations/{accommodationId}", requireAuth(accommodation.DeleteHandler(accommodationSvc, usersClient)))

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
