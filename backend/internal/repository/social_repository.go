package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ShalArl/trip-manager/internal/domain"
)

const (
	collComments      = "comments"
	collReplies       = "commentReplies"
	collActivityLikes = "activityLikes"
	collCommentLikes  = "commentLikes"
	collReplyLikes    = "replyLikes"
)

type SocialRepository interface {
	ListComments(ctx context.Context, activityID string, limit, offset int) ([]*domain.Comment, int, error)
	GetComment(ctx context.Context, id string) (*domain.Comment, error)
	CreateComment(ctx context.Context, c *domain.Comment) (*domain.Comment, error)
	UpdateComment(ctx context.Context, c *domain.Comment) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id string) error

	ListReplies(ctx context.Context, commentID string, limit, offset int) ([]*domain.CommentReply, int, error)
	GetReply(ctx context.Context, id string) (*domain.CommentReply, error)
	CreateReply(ctx context.Context, r *domain.CommentReply) (*domain.CommentReply, error)
	UpdateReply(ctx context.Context, r *domain.CommentReply) (*domain.CommentReply, error)
	DeleteReply(ctx context.Context, id string) error

	LikeActivity(ctx context.Context, userID, activityID string) error
	UnlikeActivity(ctx context.Context, userID, activityID string) error
	LikeComment(ctx context.Context, userID, commentID string) error
	UnlikeComment(ctx context.Context, userID, commentID string) error
	LikeReply(ctx context.Context, userID, replyID string) error
	UnlikeReply(ctx context.Context, userID, replyID string) error

	CountActivityLikes(ctx context.Context, activityID string) (int, error)
	CountActivityComments(ctx context.Context, activityID string) (int, error)
	CountCommentLikes(ctx context.Context, commentID string) (int, error)
	CountCommentReplies(ctx context.Context, commentID string) (int, error)
	CountReplyLikes(ctx context.Context, replyID string) (int, error)
}

type SocialRepositoryImpl struct {
	client *firestore.Client
}

func NewSocialRepository(client *firestore.Client) SocialRepository {
	return &SocialRepositoryImpl{client: client}
}

// ── Firestore document shapes ─────────────────────────────────────────────

type commentDoc struct {
	ActivityID string    `firestore:"activityId"`
	UserID     string    `firestore:"userId"`
	Text       string    `firestore:"text"`
	ImageKey   string    `firestore:"imageKey,omitempty"`
	CreatedAt  time.Time `firestore:"createdAt"`
	UpdatedAt  time.Time `firestore:"updatedAt"`
}

type replyDoc struct {
	CommentID string    `firestore:"commentId"`
	UserID    string    `firestore:"userId"`
	Text      string    `firestore:"text"`
	ImageKey  string    `firestore:"imageKey,omitempty"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

type likeDoc struct {
	UserID    string    `firestore:"userId"`
	TargetID  string    `firestore:"targetId"`
	CreatedAt time.Time `firestore:"createdAt"`
}

// ── Comments ──────────────────────────────────────────────────────────────

func (r *SocialRepositoryImpl) ListComments(ctx context.Context, activityID string, limit, offset int) ([]*domain.Comment, int, error) {
	q := r.client.Collection(collComments).
		Where("activityId", "==", activityID).
		OrderBy("createdAt", firestore.Desc)

	// Total count separat — bei Firestore kein COUNT(*), nur via aggregation query
	totalAgg, err := q.NewAggregationQuery().WithCount("total").Get(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count comments: %w", err)
	}
	total := int(totalAgg["total"].(*firestore.AggregationResult).Count)

	q = q.Offset(offset).Limit(limit)

	it := q.Documents(ctx)
	defer it.Stop()

	var comments []*domain.Comment
	for {
		snap, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("iterate comments: %w", err)
		}
		c, err := commentFromSnapshot(snap)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, c)
	}

	return comments, total, nil
}

func (r *SocialRepositoryImpl) GetComment(ctx context.Context, id string) (*domain.Comment, error) {
	snap, err := r.client.Collection(collComments).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get comment: %w", err)
	}
	return commentFromSnapshot(snap)
}

func (r *SocialRepositoryImpl) CreateComment(ctx context.Context, c *domain.Comment) (*domain.Comment, error) {
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now

	doc := commentDoc{
		ActivityID: c.ActivityID,
		UserID:     c.UserID,
		Text:       c.Text,
		ImageKey:   *c.ImageKey,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}

	ref, _, err := r.client.Collection(collComments).Add(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	c.ID = ref.ID
	return c, nil
}

func (r *SocialRepositoryImpl) UpdateComment(ctx context.Context, c *domain.Comment) (*domain.Comment, error) {
	c.UpdatedAt = time.Now().UTC()
	_, err := r.client.Collection(collComments).Doc(c.ID).Update(ctx, []firestore.Update{
		{Path: "text", Value: c.Text},
		{Path: "imageKey", Value: c.ImageKey},
		{Path: "updatedAt", Value: c.UpdatedAt},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("update comment: %w", err)
	}
	return c, nil
}

func (r *SocialRepositoryImpl) DeleteComment(ctx context.Context, id string) error {
	_, err := r.client.Collection(collComments).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	return nil
}

// ── Replies ───────────────────────────────────────────────────────────────

func (r *SocialRepositoryImpl) ListReplies(ctx context.Context, commentID string, limit, offset int) ([]*domain.CommentReply, int, error) {
	q := r.client.Collection(collReplies).
		Where("commentId", "==", commentID).
		OrderBy("createdAt", firestore.Asc)

	totalAgg, err := q.NewAggregationQuery().WithCount("total").Get(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count replies: %w", err)
	}
	total := int(totalAgg["total"].(*firestore.AggregationResult).Count)

	q = q.Offset(offset).Limit(limit)
	it := q.Documents(ctx)
	defer it.Stop()

	var replies []*domain.CommentReply
	for {
		snap, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("iterate replies: %w", err)
		}
		rep, err := replyFromSnapshot(snap)
		if err != nil {
			return nil, 0, err
		}
		replies = append(replies, rep)
	}

	return replies, total, nil
}

func (r *SocialRepositoryImpl) GetReply(ctx context.Context, id string) (*domain.CommentReply, error) {
	snap, err := r.client.Collection(collReplies).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get reply: %w", err)
	}
	return replyFromSnapshot(snap)
}

func (r *SocialRepositoryImpl) CreateReply(ctx context.Context, rep *domain.CommentReply) (*domain.CommentReply, error) {
	now := time.Now().UTC()
	rep.CreatedAt = now
	rep.UpdatedAt = now

	doc := replyDoc{
		CommentID: rep.CommentID,
		UserID:    rep.UserID,
		Text:      rep.Text,
		ImageKey:  *rep.ImageKey,
		CreatedAt: rep.CreatedAt,
		UpdatedAt: rep.UpdatedAt,
	}

	ref, _, err := r.client.Collection(collReplies).Add(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("create reply: %w", err)
	}
	rep.ID = ref.ID
	return rep, nil
}

func (r *SocialRepositoryImpl) UpdateReply(ctx context.Context, rep *domain.CommentReply) (*domain.CommentReply, error) {
	rep.UpdatedAt = time.Now().UTC()
	_, err := r.client.Collection(collReplies).Doc(rep.ID).Update(ctx, []firestore.Update{
		{Path: "text", Value: rep.Text},
		{Path: "imageKey", Value: rep.ImageKey},
		{Path: "updatedAt", Value: rep.UpdatedAt},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("update reply: %w", err)
	}
	return rep, nil
}

func (r *SocialRepositoryImpl) DeleteReply(ctx context.Context, id string) error {
	_, err := r.client.Collection(collReplies).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("delete reply: %w", err)
	}
	return nil
}

// ── Likes ─────────────────────────────────────────────────────────────────

func (r *SocialRepositoryImpl) LikeActivity(ctx context.Context, userID, activityID string) error {
	return r.upsertLike(ctx, collActivityLikes, userID, activityID)
}

func (r *SocialRepositoryImpl) UnlikeActivity(ctx context.Context, userID, activityID string) error {
	return r.deleteLike(ctx, collActivityLikes, userID, activityID)
}

func (r *SocialRepositoryImpl) LikeComment(ctx context.Context, userID, commentID string) error {
	return r.upsertLike(ctx, collCommentLikes, userID, commentID)
}

func (r *SocialRepositoryImpl) UnlikeComment(ctx context.Context, userID, commentID string) error {
	return r.deleteLike(ctx, collCommentLikes, userID, commentID)
}

func (r *SocialRepositoryImpl) LikeReply(ctx context.Context, userID, replyID string) error {
	return r.upsertLike(ctx, collReplyLikes, userID, replyID)
}

func (r *SocialRepositoryImpl) UnlikeReply(ctx context.Context, userID, replyID string) error {
	return r.deleteLike(ctx, collReplyLikes, userID, replyID)
}

// upsertLike: composite key {userID}_{targetID} → max. ein Like pro User pro Target.
// Idempotent: zweimal Liken = no-op.
func (r *SocialRepositoryImpl) upsertLike(ctx context.Context, collection, userID, targetID string) error {
	docID := userID + "_" + targetID
	doc := likeDoc{
		UserID:    userID,
		TargetID:  targetID,
		CreatedAt: time.Now().UTC(),
	}
	_, err := r.client.Collection(collection).Doc(docID).Set(ctx, doc)
	if err != nil {
		return fmt.Errorf("upsert like: %w", err)
	}
	return nil
}

// deleteLike: idempotent. Wenn Like nicht existiert, kein Fehler.
func (r *SocialRepositoryImpl) deleteLike(ctx context.Context, collection, userID, targetID string) error {
	docID := userID + "_" + targetID
	_, err := r.client.Collection(collection).Doc(docID).Delete(ctx)
	if err != nil {
		// Firestore Delete is idempotent — kein NotFound-Error bei non-existent
		return fmt.Errorf("delete like: %w", err)
	}
	return nil
}

// ── Counts ────────────────────────────────────────────────────────────────

func (r *SocialRepositoryImpl) CountActivityLikes(ctx context.Context, activityID string) (int, error) {
	return r.countLikes(ctx, collActivityLikes, activityID)
}

func (r *SocialRepositoryImpl) CountActivityComments(ctx context.Context, activityID string) (int, error) {
	q := r.client.Collection(collComments).Where("activityId", "==", activityID)
	return r.runCount(ctx, q)
}

func (r *SocialRepositoryImpl) CountCommentLikes(ctx context.Context, commentID string) (int, error) {
	return r.countLikes(ctx, collCommentLikes, commentID)
}

func (r *SocialRepositoryImpl) CountCommentReplies(ctx context.Context, commentID string) (int, error) {
	q := r.client.Collection(collReplies).Where("commentId", "==", commentID)
	return r.runCount(ctx, q)
}

func (r *SocialRepositoryImpl) CountReplyLikes(ctx context.Context, replyID string) (int, error) {
	return r.countLikes(ctx, collReplyLikes, replyID)
}

func (r *SocialRepositoryImpl) countLikes(ctx context.Context, collection, targetID string) (int, error) {
	q := r.client.Collection(collection).Where("targetId", "==", targetID)
	return r.runCount(ctx, q)
}

func (r *SocialRepositoryImpl) runCount(ctx context.Context, q firestore.Query) (int, error) {
	agg, err := q.NewAggregationQuery().WithCount("total").Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("aggregation count: %w", err)
	}
	res, ok := agg["total"].(*firestore.AggregationResult)
	if !ok {
		return 0, fmt.Errorf("aggregation result type mismatch")
	}
	return int(res.Count), nil
}

// ── Snapshot mappers ──────────────────────────────────────────────────────

func commentFromSnapshot(snap *firestore.DocumentSnapshot) (*domain.Comment, error) {
	var doc commentDoc
	if err := snap.DataTo(&doc); err != nil {
		return nil, fmt.Errorf("decode comment: %w", err)
	}
	return &domain.Comment{
		ID:         snap.Ref.ID,
		ActivityID: doc.ActivityID,
		UserID:     doc.UserID,
		Text:       doc.Text,
		ImageKey:   &doc.ImageKey,
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
	}, nil
}

func replyFromSnapshot(snap *firestore.DocumentSnapshot) (*domain.CommentReply, error) {
	var doc replyDoc
	if err := snap.DataTo(&doc); err != nil {
		return nil, fmt.Errorf("decode reply: %w", err)
	}
	return &domain.CommentReply{
		ID:        snap.Ref.ID,
		CommentID: doc.CommentID,
		UserID:    doc.UserID,
		Text:      doc.Text,
		ImageKey:  &doc.ImageKey,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}, nil
}
