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
	"github.com/ShalArl/trip-manager/backend/shared/userclient"
	"github.com/ShalArl/trip-manager/backend/trips/config"
	"github.com/ShalArl/trip-manager/backend/trips/database"
	"github.com/ShalArl/trip-manager/backend/trips/internal/accommodation"
	"github.com/ShalArl/trip-manager/backend/trips/internal/transport"
	"github.com/ShalArl/trip-manager/backend/trips/internal/trip"
	"github.com/ShalArl/trip-manager/backend/trips/pubsub"
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

	// PubSub Producer
	var pubsubProducer *pubsub.Producer
	if cfg.GCPProjectID != "" && cfg.PubSubTopicID != "" {
		var err error
		pubsubProducer, err = pubsub.NewProducer(cfg.GCPProjectID, cfg.PubSubTopicID)
		if err != nil {
			log.Fatalf("failed to initialize pubsub producer: %v", err)
		}
		defer func(pubsubProducer *pubsub.Producer) {
			err := pubsubProducer.Close()
			if err != nil {
				log.Fatalf("failed to close pubsub producer: %v", err)
			}
		}(pubsubProducer)
		log.Printf("Pub/Sub producer initialized for project %s on topic %s", cfg.GCPProjectID, cfg.PubSubTopicID)
	} else {
		log.Println("warn: GCP_PROJECT_ID or PUBSUB_TOPIC_ID not set, trip.created events will not be published")
	}

	// Clients
	authClient := authclient.NewClient(cfg.AuthServiceURL)
	usersClient := userclient.NewUsersClient(cfg.UsersServiceURL)
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
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			return
		}
	})

	// Trips
	mux.HandleFunc("GET /", requireAuth(trip.ListTripsHandler(tripSvc, usersClient)))
	mux.HandleFunc("POST /", requireAuth(trip.CreateTripHandler(tripSvc, usersClient, pubsubProducer)))
	mux.HandleFunc("GET /recent", optionalAuth(trip.ListRecentTripsHandler(tripSvc, usersClient)))
	mux.HandleFunc("GET /search", optionalAuth(trip.SearchTripsHandler(tripSvc, usersClient)))
	mux.HandleFunc("GET /{tripId}", optionalAuth(trip.GetTripHandler(tripSvc, usersClient)))
	mux.HandleFunc("PUT /{tripId}", requireAuth(trip.UpdateTripHandler(tripSvc, usersClient)))
	mux.HandleFunc("DELETE /{tripId}", requireAuth(trip.DeleteTripHandler(tripSvc, usersClient)))

	// Transports
	mux.HandleFunc("GET /{tripId}/transports", optionalAuth(transport.ListHandler(transportSvc)))
	mux.HandleFunc("POST /{tripId}/transports", requireAuth(transport.CreateHandler(transportSvc, usersClient)))
	mux.HandleFunc("PUT /{tripId}/transports/{transportId}", requireAuth(transport.UpdateHandler(transportSvc, usersClient)))
	mux.HandleFunc("DELETE /{tripId}/transports/{transportId}", requireAuth(transport.DeleteHandler(transportSvc, usersClient)))

	// Accommodations
	mux.HandleFunc("GET /{tripId}/accommodations", optionalAuth(accommodation.ListHandler(accommodationSvc)))
	mux.HandleFunc("POST /{tripId}/accommodations", requireAuth(accommodation.CreateHandler(accommodationSvc, usersClient)))
	mux.HandleFunc("PUT /{tripId}/accommodations/{accommodationId}", requireAuth(accommodation.UpdateHandler(accommodationSvc, usersClient)))
	mux.HandleFunc("DELETE /{tripId}/accommodations/{accommodationId}", requireAuth(accommodation.DeleteHandler(accommodationSvc, usersClient)))

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

	log.Printf("trips service starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
	log.Println("Server stopped")
}
