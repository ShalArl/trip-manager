package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/auth"
	"github.com/ShalArl/trip-manager/internal/config"
	"github.com/ShalArl/trip-manager/internal/handler"
	chimiddleware "github.com/ShalArl/trip-manager/internal/middleware"
)

func startUp() (*app.App, error) {
	cfg, err := config.LoadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	return application, nil
}

func main() {
	application, err := startUp()
	if err != nil {
		log.Fatalf("Application startup failed: %v", err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			application.Logger.Printf("Error closing app: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   application.Config.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	firebaseAuth, err := auth.NewFirebaseAuth(ctx, application.Config.FirebaseConfig.ProjectID)
	if err != nil {
		application.Logger.Printf("Error initializing firebase auth: %v", err)
		os.Exit(1)
	}

	userResolver := application.Services.User

	r.Route("/api", func(r chi.Router) {
		// ─── Public Routes ───────────────────────────────────────────────────────
		r.Get("/trips/search", handler.SearchTripsHandler(application))
		r.Get("/trips/recent", handler.ListRecentTripsHandler(application))

		// ─── Provision ───────────────────────────────────────────────────────────
		r.With(chimiddleware.ProvisionMiddleware(firebaseAuth)).
			Post("/users/provision", handler.ProvisionMeHandler(application))

		// ─── Optional Auth ────────────────────────────────────────────────────────
		r.With(chimiddleware.OptionalFirebaseAuthMiddleware(firebaseAuth, userResolver)).
			Get("/trips/{tripId}", handler.GetTripHandler(application))

		// ─── Location Routes ──────────────────────────────────────────────────────
		r.Route("/trips/{tripId}/locations", func(r chi.Router) {
			r.With(chimiddleware.OptionalFirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Get("/", handler.ListLocationsHandler(application))
			r.With(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Post("/", handler.CreateLocationHandler(application))
			r.Route("/{locationId}", func(r chi.Router) {
				r.Use(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver))
				r.Get("/", handler.GetLocationHandler(application))
				r.Put("/", handler.UpdateLocationHandler(application))
				r.Delete("/", handler.DeleteLocationHandler(application))
				// ─── Location Image Routes ────────────────────────────────────────
				r.Post("/images", handler.AddLocationImageHandler(application))
				r.Delete("/images/{imageId}", handler.DeleteLocationImageHandler(application))
			})
		})

		// ─── Transport Routes ─────────────────────────────────────────────────────
		r.Route("/trips/{tripId}/transports", func(r chi.Router) {
			r.With(chimiddleware.OptionalFirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Get("/", handler.ListTransportsHandler(application))
			r.With(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Post("/", handler.CreateTransportHandler(application))
			r.Route("/{transportId}", func(r chi.Router) {
				r.Use(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver))
				r.Get("/", handler.GetTransportHandler(application))
				r.Put("/", handler.UpdateTransportHandler(application))
				r.Delete("/", handler.DeleteTransportHandler(application))
			})
		})

		// ─── Accommodation Routes ─────────────────────────────────────────────────
		r.Route("/trips/{tripId}/accommodations", func(r chi.Router) {
			r.With(chimiddleware.OptionalFirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Get("/", handler.ListAccommodationsHandler(application))
			r.With(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Post("/", handler.CreateAccommodationHandler(application))
			r.Route("/{accommodationId}", func(r chi.Router) {
				r.Use(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver))
				r.Put("/", handler.UpdateAccommodationHandler(application))
				r.Delete("/", handler.DeleteAccommodationHandler(application))
			})
		})

		// ─── Social Routes ────────────────────────────────────────────────────────
		r.Route("/trips/{tripId}/likes", func(r chi.Router) {
			r.With(chimiddleware.OptionalFirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Get("/", handler.GetTripLikesHandler(application))
			r.With(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Post("/", handler.LikeTripHandler(application))
			r.With(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Delete("/", handler.UnlikeTripHandler(application))
		})

		r.Route("/trips/{tripId}/comments", func(r chi.Router) {
			r.With(chimiddleware.OptionalFirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Get("/", handler.ListTripCommentsHandler(application))
			r.With(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Post("/", handler.CreateTripCommentHandler(application))
			r.With(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, userResolver)).
				Delete("/{commentId}", handler.DeleteTripCommentHandler(application))
		})

		// ─── Protected Routes ─────────────────────────────────────────────────────
		r.Group(func(r chi.Router) {
			r.Use(chimiddleware.FirebaseAuthMiddleware(firebaseAuth, application.Services.User))

			// Upload
			r.Post("/uploads/presigned", handler.GetPresignedURLHandler(application))

			// Users
			r.Get("/users/me", handler.GetMeHandler(application))
			r.Put("/users/me", handler.UpdateMeHandler(application))
			r.Get("/users/{userId}", handler.GetUserHandler(application))
			r.Put("/users/{userId}", handler.UpdateUserHandler(application))
			r.Delete("/users/{userId}", handler.DeleteUserHandler(application))

			// Trips
			r.Get("/trips", handler.ListTripsHandler(application))
			r.Post("/trips", handler.CreateTripHandler(application))
			r.Put("/trips/{tripId}", handler.UpdateTripHandler(application))
			r.Delete("/trips/{tripId}", handler.DeleteTripHandler(application))

			// Locations
			r.Get("/locations/{locationId}", handler.GetLocationHandler(application))

			// Activities
			r.Route("/trips/{tripId}/activities", func(r chi.Router) {
				r.Get("/", handler.ListActivitiesForTripHandler(application))
				r.Post("/", handler.CreateActivityHandler(application))
				r.Route("/{activityId}", func(r chi.Router) {
					r.Put("/", handler.UpdateActivityHandler(application))
					r.Delete("/", handler.DeleteActivityHandler(application))
				})
			})

			r.Get("/locations/{locationId}/activities", handler.ListActivitiesForLocationHandler(application))
			r.Get("/activities/{activityId}", handler.GetActivityHandler(application))
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			log.Printf("Error writing health check response: %v", err)
		}
	})

	addr := fmt.Sprintf(":%s", application.Config.ServerPort)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		application.Logger.Printf("🚀 Server starting on http://localhost%s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		application.Logger.Fatalf("Server error: %v", err)
	case <-ctx.Done():
		application.Logger.Println("shutting down gracefully")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		application.Logger.Printf("graceful shutdown failed: %v", err)
	} else {
		application.Logger.Println("server shut down cleanly")
	}
}
