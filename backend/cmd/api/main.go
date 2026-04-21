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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   application.Config.CORSAllowedOrigins,
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
		r.Post("/auth/register", handler.CreateUserHandler(application))
		r.Post("/auth/login", handler.LoginHandler(application))
		r.Get("/trips/search", handler.SearchTripsHandler(application))
		r.Get("/trips/recent", handler.ListRecentTripsHandler(application))

		// ─── Optional Auth Routes (public but user context if available) ────────
		r.With(chimiddleware.OptionalAuthMiddleware(authManager)).Get("/trips/{tripId}", handler.GetTripHandler(application))

		// Protected routes - require JWT authentication
		r.Group(func(r chi.Router) {
			r.Use(chimiddleware.AuthMiddleware(authManager))

			// ─── Upload Routes ──────────────────────────────────────────────────────
			r.Post("/uploads/presigned", handler.GetPresignedURLHandler(application))

			// ─── User Routes ────────────────────────────────────────────────────────
			r.Get("/users/me", handler.GetMeHandler(application))
			r.Put("/users/me", handler.UpdateMeHandler(application))
			r.Put("/users/me/password", handler.ChangePasswordHandler(application))
			r.Get("/users/{userId}", handler.GetUserHandler(application))
			r.Put("/users/{userId}", handler.UpdateUserHandler(application))
			r.Delete("/users/{userId}", handler.DeleteUserHandler(application))

			// ─── Trip Routes ────────────────────────────────────────────────────────
			r.Get("/trips", handler.ListTripsHandler(application))
			r.Post("/trips", handler.CreateTripHandler(application))
			//r.Get("/trips/{tripId}", handler.GetTripHandler(application))
			r.Put("/trips/{tripId}", handler.UpdateTripHandler(application))
			r.Delete("/trips/{tripId}", handler.DeleteTripHandler(application))

			// ─── Location Routes ────────────────────────────────────────────────────
			r.Route("/trips/{tripId}/locations", func(r chi.Router) {
				r.Get("/", handler.ListLocationsHandler(application))
				r.Post("/", handler.CreateLocationHandler(application))
				r.Route("/{locationId}", func(r chi.Router) {
					r.Get("/", handler.GetLocationHandler(application))
					r.Put("/", handler.UpdateLocationHandler(application))
					r.Delete("/", handler.DeleteLocationHandler(application))
				})
			})

			// ─── Direct Location Routes (for individual location access) ────────────
			r.Get("/locations/{locationId}", handler.GetLocationHandler(application))

			// ─── Activity Routes ────────────────────────────────────────────────────
			r.Route("/trips/{tripId}/activities", func(r chi.Router) {
				r.Get("/", handler.ListActivitiesForTripHandler(application))
				r.Post("/", handler.CreateActivityHandler(application))
				r.Route("/{activityId}", func(r chi.Router) {
					r.Put("/", handler.UpdateActivityHandler(application))
					r.Delete("/", handler.DeleteActivityHandler(application))

					// ─── Activity Likes Routes ──────────────────────────────────────────
					r.Post("/likes", handler.LikeActivityHandler(application))
					r.Delete("/likes", handler.UnlikeActivityHandler(application))

					// ─── Comments Routes ────────────────────────────────────────────────
					r.Route("/comments", func(r chi.Router) {
						r.Get("/", handler.ListCommentsHandler(application))
						r.Post("/", handler.CreateCommentHandler(application))
						r.Route("/{commentId}", func(r chi.Router) {
							r.Put("/", handler.UpdateCommentHandler(application))
							r.Delete("/", handler.DeleteCommentHandler(application))

							// ─── Comment Likes Routes ──────────────────────────────────────
							r.Post("/likes", handler.LikeCommentHandler(application))
							r.Delete("/likes", handler.UnlikeCommentHandler(application))

							// ─── Nested Comments (Replies) Routes ───────────────────────────
							r.Route("/comments", func(r chi.Router) {
								r.Get("/", handler.ListRepliesHandler(application))
								r.Post("/", handler.CreateReplyHandler(application))
								r.Route("/{replyId}", func(r chi.Router) {
									r.Put("/", handler.UpdateReplyHandler(application))
									r.Delete("/", handler.DeleteReplyHandler(application))

									// ─── Reply Likes Routes ─────────────────────────────────────
									r.Post("/likes", handler.LikeReplyHandler(application))
									r.Delete("/likes", handler.UnlikeReplyHandler(application))
								})
							})
						})
					})
				})
			})

			// ─── Activity by Location Route ──────────────────────────────────────────
			r.Get("/locations/{locationId}/activities", handler.ListActivitiesForLocationHandler(application))

			// ─── Direct Activity Routes (for individual activity access) ──────────────
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

	// Start server
	addr := fmt.Sprintf(":%s", application.Config.ServerPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
		// Timeouts Production
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Run in Goroutine
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

	// Graceful Shutdown with Timeout (10 seconds)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		application.Logger.Printf("graceful shutdown failed: %v", err)
	} else {
		application.Logger.Println("server shut down cleanly")
	}
}
