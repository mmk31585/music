package reactions

import (
	"time"

	"github.com/google/uuid"
)

type Reaction struct {
	ID         uuid.UUID `db:"id" json:"id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	TargetID   string    `db:"target_id" json:"target_id"`
	TargetType string    `db:"target_type" json:"target_type"`
	Type       string    `db:"type" json:"type"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type ReactionCount struct {
	TargetID   string `db:"target_id" json:"target_id"`
	TargetType string `db:"target_type" json:"target_type"`
	Type       string `db:"type" json:"type"`
	Count      int    `db:"count" json:"count"`
}
