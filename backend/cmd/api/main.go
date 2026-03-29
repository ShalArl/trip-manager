package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/ShalArl/trip-manager/internal/api"
	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/auth"
	"github.com/ShalArl/trip-manager/internal/config"
	chimiddleware "github.com/ShalArl/trip-manager/internal/middleware"
)

func startUp() (*app.App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	return application, nil
}

func main() {
	// Load configuration and initialize application
	application, err := startUp()
	if err != nil {
		log.Fatalf("Application startup failed: %v", err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			application.Logger.Printf("Error closing app: %v", err)
		}
	}()

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api", func(r chi.Router) {
		// Initialize auth manager for middleware (7 day token expiration)
		authManager := auth.NewAuthManager(application.Config.JWTSecret, 7*24*time.Hour)

		// ─── Auth Routes (no auth required) ────────────────────────────────────
		r.Post("/auth/register", api.CreateUserHandler(application))
		r.Post("/auth/login", api.LoginHandler(application))

		// Protected routes - require JWT authentication
		r.Group(func(r chi.Router) {
			r.Use(chimiddleware.AuthMiddleware(authManager))

			// ─── User Routes ────────────────────────────────────────────────────────
			r.Get("/users/me", api.GetMeHandler(application))
			r.Put("/users/me", api.UpdateMeHandler(application))
			r.Put("/users/me/password", api.ChangePasswordHandler(application))
			r.Get("/users/{userId}", api.GetUserHandler(application))
			r.Put("/users/{userId}", api.UpdateUserHandler(application))
			r.Delete("/users/{userId}", api.DeleteUserHandler(application))

			// ─── Trip Routes ────────────────────────────────────────────────────────
			r.Get("/trips", api.ListTripsHandler(application))
			r.Post("/trips", api.CreateTripHandler(application))
			r.Get("/trips/{tripId}", api.GetTripHandler(application))
			r.Put("/trips/{tripId}", api.UpdateTripHandler(application))
			r.Delete("/trips/{tripId}", api.DeleteTripHandler(application))

			// ─── Location Routes ────────────────────────────────────────────────────
			r.Route("/trips/{tripId}/locations", func(r chi.Router) {
				r.Get("/", api.ListLocationsHandler(application))
				r.Post("/", api.CreateLocationHandler(application))
				r.Route("/{locationId}", func(r chi.Router) {
					r.Get("/", api.GetLocationHandler(application))
					r.Put("/", api.UpdateLocationHandler(application))
					r.Delete("/", api.DeleteLocationHandler(application))
				})
			})

			// ─── Direct Location Routes (for individual location access) ────────────
			r.Get("/locations/{locationId}", api.GetLocationHandler(application))

			// ─── Activity Routes ────────────────────────────────────────────────────
			r.Route("/trips/{tripId}/activities", func(r chi.Router) {
				r.Get("/", api.ListActivitiesForTripHandler(application))
				r.Post("/", api.CreateActivityHandler(application))
				r.Route("/{activityId}", func(r chi.Router) {
					r.Put("/", api.UpdateActivityHandler(application))
					r.Delete("/", api.DeleteActivityHandler(application))
				})
			})

			// ─── Activity by Location Route ──────────────────────────────────────────
			r.Get("/locations/{locationId}/activities", api.ListActivitiesForLocationHandler(application))

			// ─── Direct Activity Routes (for individual activity access) ──────────────
			r.Get("/activities/{activityId}", api.GetActivityHandler(application))
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
	// Start server
	addr := fmt.Sprintf(":%s", application.Config.ServerPort)
	application.Logger.Printf("🚀 Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		application.Logger.Fatalf("Server error: %v", err)
	}
}
