package history

import (
	"time"

	"github.com/google/uuid"
)

type ListeningHistoryItem struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	TrackID   uuid.UUID `db:"track_id" json:"track_id"`
	PlayedAt  time.Time `db:"played_at" json:"played_at"`
	Duration  int       `db:"duration" json:"duration"`
	Completed bool      `db:"completed" json:"completed"`
}
