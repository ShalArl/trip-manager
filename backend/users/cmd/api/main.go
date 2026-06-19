package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	sharedotel "otel"
	"syscall"
	"time"

	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/shared/firebaseclient"
	"github.com/ShalArl/trip-manager/backend/shared/middleware"
	"github.com/ShalArl/trip-manager/backend/users/config"
	"github.com/ShalArl/trip-manager/backend/users/database"
	"github.com/ShalArl/trip-manager/backend/users/handler"
	"github.com/ShalArl/trip-manager/backend/users/internal/tenant"
	"github.com/ShalArl/trip-manager/backend/users/repository"
	"github.com/ShalArl/trip-manager/backend/users/service"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	corsConfig := middleware.DefaultCORSConfig()
	allowedOrigins := cfg.CORSAllowedOrigins
	if len(allowedOrigins) == 0 {
		log.Fatalf("No allowed origin configured")
	}
	corsConfig.AllowedOrigins = allowedOrigins

	otelProvider, err := sharedotel.New(ctx, "users", cfg.OTELCollectorEndpoint)
	if err != nil {
		log.Printf("warn: failed to initialize otel: %v", err)
	}
	var metrics *sharedotel.ServiceMetrics
	if otelProvider != nil {
		defer otelProvider.Shutdown(ctx)
		metrics, _ = sharedotel.NewServiceMetrics(otelProvider.Meter, "users")
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

	// Auth client
	authClient := authclient.NewClient(cfg.AuthServiceURL)

	// Firebase client
	fbClient, err := firebaseclient.New(ctx, cfg.FirebaseProjectID)
	if err != nil {
		log.Fatalf("failed to init firebase client: %v", err)
	}

	// Wire up
	repo := repository.NewRepository(db)
	tenantRepo := tenant.NewRepository(db)

	svc := service.NewService(repo, fbClient)

	// Middleware
	requireAuth := authclient.RequireAuth(authClient)

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			log.Printf("failed to write health response: %v", err)
			return
		}
	})

	prometheusURL := cfg.PrometheusURL

	mux.HandleFunc("POST /provision", requireAuth(handler.ProvisionHandler(svc)))
	mux.HandleFunc("GET /me", requireAuth(handler.GetMeHandler(svc)))
	mux.HandleFunc("PUT /me", requireAuth(handler.UpdateMeHandler(svc)))
	mux.HandleFunc("GET /{id}", handler.GetByIDHandler(svc))
	mux.HandleFunc("POST /tenants/register", requireAuth(tenant.RegisterHandler(tenantRepo, svc)))
	mux.HandleFunc("GET /tenants/me", requireAuth(tenant.GetTenantHandler(tenantRepo)))
	mux.HandleFunc("GET /tenants/by-slug/{slug}", tenant.GetTenantBySlugHandler(tenantRepo))
	mux.HandleFunc("GET /tenants/me/branding", requireAuth(tenant.GetBrandingHandler(tenantRepo)))
	mux.HandleFunc("PUT /tenants/me/branding", requireAuth(tenant.UpdateBrandingHandler(tenantRepo)))
	mux.HandleFunc("PUT /tenants/me/tier", requireAuth(tenant.UpgradeTierHandler(tenantRepo)))
	mux.HandleFunc("GET /tenants/me/usage", requireAuth(tenant.GetUsageHandler(tenantRepo, prometheusURL)))
	mux.HandleFunc("GET /tenants/me/settings", requireAuth(tenant.GetSettingsHandler(tenantRepo)))
	mux.HandleFunc("PUT /tenants/me/settings", requireAuth(tenant.UpdateSettingsHandler(tenantRepo)))

	// Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.CORS(corsConfig)(sharedotel.MetricsMiddleware(metrics, authClient)(mux)),
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
