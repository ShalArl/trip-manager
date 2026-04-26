package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ShalArl/trip-manager/internal/domain"
)

const (
	collTripLikes    = "tripLikes"
	collTripComments = "tripComments"
)

// ── Firestore document shapes ──────────────────────────────────────────────

type tripLikeDoc struct {
	TripID    string    `firestore:"tripId"`
	UserID    string    `firestore:"userId"`
	CreatedAt time.Time `firestore:"createdAt"`
}

type tripCommentDoc struct {
	TripID    string    `firestore:"tripId"`
	UserID    string    `firestore:"userId"`
	Text      string    `firestore:"text"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

// ── Interface ──────────────────────────────────────────────────────────────

type SocialRepository interface {
	// Likes
	LikeTrip(ctx context.Context, userID, tripID string) error
	UnlikeTrip(ctx context.Context, userID, tripID string) error
	HasLiked(ctx context.Context, userID, tripID string) (bool, error)
	CountTripLikes(ctx context.Context, tripID string) (int, error)

	// Comments
	CreateTripComment(ctx context.Context, c *domain.TripComment) (*domain.TripComment, error)
	ListTripComments(ctx context.Context, tripID string) ([]*domain.TripComment, error)
	DeleteTripComment(ctx context.Context, commentID, userID string) error
}

// ── Implementation ─────────────────────────────────────────────────────────

type SocialRepositoryImpl struct {
	client *firestore.Client
}

func NewSocialRepository(client *firestore.Client) SocialRepository {
	return &SocialRepositoryImpl{client: client}
}

// ── Likes ──────────────────────────────────────────────────────────────────

func (r *SocialRepositoryImpl) LikeTrip(ctx context.Context, userID, tripID string) error {
	docID := userID + "_" + tripID
	_, err := r.client.Collection(collTripLikes).Doc(docID).Create(ctx, tripLikeDoc{
		TripID:    tripID,
		UserID:    userID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return domain.ErrConflict
		}
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return nil
}

func (r *SocialRepositoryImpl) UnlikeTrip(ctx context.Context, userID, tripID string) error {
	docID := userID + "_" + tripID
	_, err := r.client.Collection(collTripLikes).Doc(docID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return nil
}

func (r *SocialRepositoryImpl) HasLiked(ctx context.Context, userID, tripID string) (bool, error) {
	docID := userID + "_" + tripID
	_, err := r.client.Collection(collTripLikes).Doc(docID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return true, nil
}

func (r *SocialRepositoryImpl) CountTripLikes(ctx context.Context, tripID string) (int, error) {
	iter := r.client.Collection(collTripLikes).
		Where("tripId", "==", tripID).
		Documents(ctx)
	defer iter.Stop()

	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("%w: %v", domain.ErrInternal, err)
		}
		count++
	}
	return count, nil
}

// ── Comments ───────────────────────────────────────────────────────────────

func (r *SocialRepositoryImpl) CreateTripComment(ctx context.Context, c *domain.TripComment) (*domain.TripComment, error) {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	ref, _, err := r.client.Collection(collTripComments).Add(ctx, tripCommentDoc{
		TripID:    c.TripID,
		UserID:    c.UserID,
		Text:      c.Text,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	c.ID = ref.ID
	return c, nil
}

func (r *SocialRepositoryImpl) ListTripComments(ctx context.Context, tripID string) ([]*domain.TripComment, error) {
	iter := r.client.Collection(collTripComments).
		Where("tripId", "==", tripID).
		OrderBy("createdAt", firestore.Asc).
		Documents(ctx)
	defer iter.Stop()

	var comments []*domain.TripComment
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
		}

		var d tripCommentDoc
		if err := doc.DataTo(&d); err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
		}

		comments = append(comments, &domain.TripComment{
			ID:        doc.Ref.ID,
			TripID:    d.TripID,
			UserID:    d.UserID,
			Text:      d.Text,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}
	return comments, nil
}

func (r *SocialRepositoryImpl) DeleteTripComment(ctx context.Context, commentID, userID string) error {
	doc, err := r.client.Collection(collTripComments).Doc(commentID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	// Nur der Ersteller darf löschen
	var d tripCommentDoc
	if err := doc.DataTo(&d); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	if d.UserID != userID {
		return domain.ErrForbidden
	}

	_, err = r.client.Collection(collTripComments).Doc(commentID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return nil
}
