package notification

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID         string          `db:"id"`
	UserID     string          `db:"user_id"`
	Type       string          `db:"type"`
	Title      string          `db:"title"`
	Body       string          `db:"body"`
	EntityType *string         `db:"entity_type"`
	EntityID   *string         `db:"entity_id"`
	Payload    json.RawMessage `db:"payload"`
	IsRead     bool            `db:"is_read"`
	ReadAt     *time.Time      `db:"read_at"`
	CreatedAt  time.Time       `db:"created_at"`
}
