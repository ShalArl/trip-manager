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

	"github.com/ShalArl/trip-manager/backend/newsletter/internal/newsletter"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type config struct {
	Port            string
	NewsletterDBURL string
	AuthServiceURL  string
}

func loadConfig() config {
	return config{
		Port:            getEnv("PORT", "8008"),
		NewsletterDBURL: getEnv("NEWSLETTER_DB_URL", "postgres://postgres:postgres@localhost:5432/newsletter_db?sslmode=disable"),
		AuthServiceURL:  getEnv("AUTH_SERVICE_URL", "http://localhost:8082"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	cfg := loadConfig()

	db, err := sqlx.Connect("postgres", cfg.NewsletterDBURL)
	if err != nil {
		log.Fatalf("newsletter: failed to connect to newsletter-db: %v", err)
	}
	defer db.Close()
	log.Println("newsletter: connected to newsletter-db")

	repo := newsletter.NewRepository(db)
	svc := newsletter.NewService(repo)

	authClient := authclient.NewClient(cfg.AuthServiceURL)
	requireAuth := authclient.RequireAuth(authClient)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("GET /api/newsletter", requireAuth(newsletter.GetNewsletterHandler(svc)))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		<-sigch
		log.Println("newsletter: shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("newsletter: shutdown error: %v", err)
		}
	}()

	log.Printf("newsletter service starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("newsletter: server error: %v", err)
	}
	log.Println("newsletter: stopped")
}
