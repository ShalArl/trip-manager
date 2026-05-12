package main

import (
	"log"
	"net/http"

	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/users/config"
	"github.com/ShalArl/trip-manager/backend/users/handler"
	"github.com/ShalArl/trip-manager/backend/users/repository"
	"github.com/ShalArl/trip-manager/backend/users/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	// DB
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Auth client
	authClient := authclient.NewClient(cfg.AuthServiceURL)

	// Wire up
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	// Middleware
	requireAuth := authclient.RequireAuth(authClient)

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users/provision", requireAuth(handler.ProvisionHandler(svc)))
	mux.HandleFunc("GET /api/users/me", requireAuth(handler.GetMeHandler(svc)))
	mux.HandleFunc("PUT /api/users/me", requireAuth(handler.UpdateMeHandler(svc)))

	log.Printf("users service starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
