package comment

import "time"

type UserSummary struct {
	ID        string `firestore:"userId" json:"id"`
	Name      string `firestore:"userName" json:"name"`
	Email     string `firestore:"userEmail" json:"email"`
	AvatarUrl string `firestore:"userAvatarKey" json:"avatarUrl"`
}

type Comment struct {
	ID          string      `firestore:"id"`
	EntityID    string      `firestore:"entityId"`
	FirebaseUID string      `firestore:"firebaseUid"`
	User        UserSummary `firestore:"user"`
	Text        string      `firestore:"text"`
	ImageKey    *string     `firestore:"imageKey,omitempty"`
	CreatedAt   time.Time   `firestore:"createdAt"`
	UpdatedAt   time.Time   `firestore:"updatedAt"`
}

type CommentResponse struct {
	ID        string      `json:"id"`
	EntityID  string      `json:"entityId"`
	User      UserSummary `json:"user"`
	Text      string      `json:"text"`
	ImageKey  *string     `json:"imageKey,omitempty"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type CommentListResponse struct {
	Data []*CommentResponse `json:"data"`
}

type CreateCommentRequest struct {
	Text     string  `json:"text" binding:"required"`
	ImageKey *string `json:"imageKey,omitempty"`
}
