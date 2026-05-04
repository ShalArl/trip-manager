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

	"cloud.google.com/go/firestore"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/social/internal/comment"
	"github.com/ShalArl/trip-manager/backend/social/internal/config"
	"github.com/ShalArl/trip-manager/backend/social/internal/like"
)

func main() {
	ctx := context.Background()

	// Load config
	cfg := config.LoadConfig()
	log.Printf("Starting Social Service on port %s\n", cfg.Port)

	// Connect to Firestore
	firestoreClient, err := config.ConnectFirestore(ctx, cfg.FirestoreProject)
	if err != nil {
		log.Fatalf("Failed to connect to Firestore: %v", err)
	}
	defer func(firestoreClient *firestore.Client) {
		err := firestoreClient.Close()
		if err != nil {
			log.Fatalf("Failed to close Firestore client: %v", err)
		}
	}(firestoreClient)

	// Setup auth client
	authClient := authclient.NewClient(cfg.AuthClientConnectionString)

	// Setup dependencies
	likeRepo := like.NewLikeRepository(firestoreClient)
	likeService := like.NewServiceImpl(likeRepo)

	commentRepo := comment.NewCommentRepository(firestoreClient)
	commentService := comment.NewServiceImpl(commentRepo)

	// Setup routes
	mux := http.NewServeMux()

	// Like endpoints
	mux.HandleFunc("GET /trips/{tripId}/likes", like.GetTripLikesHandler(likeService))
	mux.HandleFunc("GET /comments/{commentId}/likes", like.GetCommentLikesHandler(likeService))

	mux.HandleFunc("POST /trips/{tripId}/likes", authclient.RequireAuth(authClient)(like.LikeTripHandler(likeService)))
	mux.HandleFunc("POST /comments/{commentId}/likes", authclient.RequireAuth(authClient)(like.LikeCommentHandler(likeService)))
	mux.HandleFunc("DELETE /trips/{tripId}/likes", authclient.RequireAuth(authClient)(like.UnlikeTripHandler(likeService)))
	mux.HandleFunc("DELETE /comments/{commentId}/likes", authclient.RequireAuth(authClient)(like.UnlikeCommentHandler(likeService)))

	// Comment endpoints
	mux.HandleFunc("GET /trips/{tripId}/comments", comment.ListTripCommentsHandler(commentService))
	mux.HandleFunc("GET /comments/{commentId}/replies", comment.ListRepliesHandler(commentService))

	mux.HandleFunc("POST /trips/{tripId}/comments", authclient.RequireAuth(authClient)(comment.CreateTripCommentHandler(commentService)))
	mux.HandleFunc("POST /comments/{commentId}/comments", authclient.RequireAuth(authClient)(comment.CreateReplyHandler(commentService)))
	mux.HandleFunc("PUT /comments/{commentId}", authclient.RequireAuth(authClient)(comment.UpdateCommentHandler(commentService)))
	mux.HandleFunc("DELETE /comments/{commentId}", authclient.RequireAuth(authClient)(comment.DeleteCommentHandler(commentService)))

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			return
		}
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
			log.Fatalf("Failed to shutdown server: %v", err)
			return
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
