package video

import (
	"time"

	"github.com/google/uuid"
)

// VideoComment represents a comment on a video.
type VideoComment struct {
	ID        uuid.UUID `db:"id"`
	VideoID   uuid.UUID `db:"video_id"`
	UserID    uuid.UUID `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// CommentResponse is the public API response for a comment.
type CommentResponse struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// CreateCommentRequest is the request body for creating a comment.
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// AdminCommentItem is the admin view of a comment (includes user info).
type AdminCommentItem struct {
	ID         uuid.UUID `db:"id"          json:"id"`
	VideoID    uuid.UUID `db:"video_id"    json:"video_id"`
	UserID     uuid.UUID `db:"user_id"     json:"user_id"`
	Content    string    `db:"content"     json:"content"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"  json:"updated_at"`
	AuthorName string    `db:"author_name" json:"author_name"`
}

// AdminUpdateCommentRequest is the request body for updating a comment as admin.
type AdminUpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}
