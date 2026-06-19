package social

import "time"

type ClubDiscussion struct {
	ID         string    `json:"id" db:"id"`
	ClubID     string    `json:"club_id" db:"club_id"`
	AuthorID   string    `json:"author_id" db:"author_id"`
	Title      string    `json:"title" db:"title"`
	Body       string    `json:"body" db:"body"`
	ReplyCount int       `json:"reply_count" db:"reply_count"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ClubDiscussionReply struct {
	ID           string    `json:"id" db:"id"`
	DiscussionID string    `json:"discussion_id" db:"discussion_id"`
	AuthorID     string    `json:"author_id" db:"author_id"`
	Body         string    `json:"body" db:"body"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type CreateClubDiscussionRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type CreateDiscussionReplyRequest struct {
	Body string `json:"body" binding:"required"`
}
