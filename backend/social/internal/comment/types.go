package comment

import "time"

// Comment represents a comment on an entity
type Comment struct {
	ID        string    `firestore:"id"`
	EntityID  string    `firestore:"entityId"`
	UserID    string    `firestore:"userId"`
	Text      string    `firestore:"text"`
	ImageKey  *string   `firestore:"imageKey,omitempty"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

// CommentResponse represents a comment in API responses
type CommentResponse struct {
	ID        string    `json:"id"`
	EntityID  string    `json:"entityId"`
	UserID    string    `json:"userId"`
	Text      string    `json:"text"`
	ImageKey  *string   `json:"imageKey,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CommentListResponse wraps a list of comments
type CommentListResponse struct {
	Data []*CommentResponse `json:"data"`
}

// CreateCommentRequest is the request body for creating a comment
type CreateCommentRequest struct {
	Text     string  `json:"text" binding:"required"`
	ImageKey *string `json:"imageKey,omitempty"`
}
