package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/social/client"
	"github.com/ShalArl/trip-manager/backend/social/internal/comment"
	"github.com/ShalArl/trip-manager/backend/social/internal/config"
	"github.com/ShalArl/trip-manager/backend/social/internal/like"
	"github.com/ShalArl/trip-manager/backend/social/kafka"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()

	log.Printf("Starting Social Service on port %s\n", cfg.Port)

	firestoreClient, err := config.ConnectFirestore(ctx, cfg.FirestoreProject)
	if err != nil {
		log.Fatalf("Failed to connect to Firestore: %v", err)
	}
	defer func(firestoreClient *firestore.Client) {
		if err := firestoreClient.Close(); err != nil {
			log.Fatalf("Failed to close Firestore client: %v", err)
		}
	}(firestoreClient)

	// Kafka Producer
	var kafkaProducer *kafka.Producer
	if brokers := cfg.KafkaBrokers; brokers != "" {
		kafkaProducer = kafka.NewProducer(strings.Split(brokers, ","))
		defer kafkaProducer.Close()
		log.Printf("kafka producer connected to %s", brokers)
	} else {
		log.Println("warn: KAFKA_BROKERS not set, trip.liked/commented events will not be published")
	}

	authClient := authclient.NewClient(cfg.AuthClientConnectionString)
	usersClient := client.NewUsersClient(cfg.UsersServiceURL)

	likeRepo := like.NewLikeRepository(firestoreClient)
	likeService := like.NewServiceImpl(likeRepo)

	commentRepo := comment.NewCommentRepository(firestoreClient)
	commentService := comment.NewServiceImpl(commentRepo)

	mux := http.NewServeMux()

	// Like endpoints
	mux.HandleFunc("GET /api/trips/{tripId}/likes", authclient.OptionalAuth(authClient)(like.GetTripLikesHandler(likeService)))
	mux.HandleFunc("POST /api/trips/{tripId}/likes", authclient.RequireAuth(authClient)(like.LikeTripHandler(likeService, kafkaProducer)))
	mux.HandleFunc("DELETE /api/trips/{tripId}/likes", authclient.RequireAuth(authClient)(like.UnlikeTripHandler(likeService)))

	// Comment endpoints
	mux.HandleFunc("GET /api/trips/{tripId}/comments", comment.ListTripCommentsHandler(commentService))
	mux.HandleFunc("POST /api/trips/{tripId}/comments", authclient.RequireAuth(authClient)(comment.CreateTripCommentHandler(commentService, usersClient, kafkaProducer)))
	mux.HandleFunc("DELETE /api/trips/{tripId}/comments/{commentId}", authclient.RequireAuth(authClient)(comment.DeleteCommentHandler(commentService)))

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			return
		}
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

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

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
	log.Println("Server stopped")
}
