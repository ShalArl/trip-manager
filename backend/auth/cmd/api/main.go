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

	"github.com/ShalArl/trip-manager/backend/auth/internal/config"
	"github.com/ShalArl/trip-manager/backend/auth/internal/handler"
	"github.com/ShalArl/trip-manager/backend/auth/internal/service"
)

func main() {
	ctx := context.Background()

	// Load config
	cfg := config.LoadConfig()
	log.Printf("Starting Auth Service on port %s\n", cfg.Port)

	// Initialize Firebase
	authClient, err := config.InitializeFirebase(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Setup service
	authService := service.NewService(authClient)
	authHandler := handler.NewHandler(authService)

	// Setup routes
	mux := http.NewServeMux()

	// Auth endpoints
	mux.HandleFunc("POST /validate-token", authHandler.ValidateToken)
	mux.HandleFunc("GET /validate-token", authHandler.ValidateTokenFromHeader)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Start server
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
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
