package domain

import "time"

// ActivityLike represents a user liking an activity
type ActivityLike struct {
	ID         string    `json:"id"`
	ActivityID string    `json:"activityId"`
	UserID     string    `json:"userId"`
	User       *User     `json:"user,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Comment represents a comment on an activity with optional image
type Comment struct {
	ID         string    `json:"id"`
	Text       string    `json:"text"`
	ImageKey   *string   `json:"imageKey,omitempty"`
	ImageUrl   *string   `json:"imageUrl,omitempty"` // Short-lived signed URL
	ActivityID string    `json:"activityId"`
	UserID     string    `json:"userId,omitempty"`
	LikeCount  int       `json:"likeCount"`
	ReplyCount int       `json:"replyCount"` // Count of replies/nested comments
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// CommentReply represents a reply to a comment (nested comment)
type CommentReply struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	ImageKey  *string   `json:"imageKey,omitempty"`
	ImageUrl  *string   `json:"imageUrl,omitempty"` // Short-lived signed URL
	CommentID string    `json:"commentId"`          // Parent comment ID
	UserID    string    `json:"userID,omitempty"`
	LikeCount int       `json:"likeCount"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CommentLike represents a user liking a comment
type CommentLike struct {
	ID        string    `json:"id"`
	CommentID string    `json:"commentId"`
	UserID    string    `json:"userId"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// ReplyLike represents a user liking a comment reply
type ReplyLike struct {
	ID        string    `json:"id"`
	ReplyID   string    `json:"replyId"`
	UserID    string    `json:"userId"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
