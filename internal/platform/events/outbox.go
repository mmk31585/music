package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID          string          `db:"id" json:"id"`
	EventType   string          `db:"event_type" json:"event_type"`
	Payload     json.RawMessage `db:"payload" json:"payload"`
	Status      string          `db:"status" json:"status"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	ProcessedAt *time.Time      `db:"processed_at" json:"processed_at,omitempty"`
	RetryCount  int             `db:"retry_count" json:"retry_count"`
	LastError   *string         `db:"last_error" json:"last_error,omitempty"`
	TraceID     *string         `db:"trace_id" json:"trace_id,omitempty"`
}

type OutboxStore interface {
	Insert(ctx context.Context, event OutboxEvent) error
	FetchPending(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkProcessed(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string, errMsg string) error
	DeleteProcessed(ctx context.Context, before time.Time) error
}

func NewOutboxEvent(event Event) OutboxEvent {
	payload := map[string]interface{}{
		"event_name":  event.EventName(),
		"occurred_at": event.OccurredAt(),
	}

	switch e := event.(type) {
	case TrackPlayedEvent:
		payload["user_id"] = e.UserID.String()
		payload["track_id"] = e.TrackID.String()
		payload["artist_id"] = e.ArtistID.String()
		payload["album_id"] = e.AlbumID.String()
		payload["duration"] = e.Duration
		payload["completed"] = e.Completed
		payload["source"] = e.Source
	case PlaylistCreatedEvent:
		payload["user_id"] = e.UserID.String()
		payload["playlist_id"] = e.PlaylistID.String()
		payload["name"] = e.Name
		payload["is_public"] = e.IsPublic
	case SubscriptionPurchasedEvent:
		payload["user_id"] = e.UserID.String()
		payload["subscription_id"] = e.SubscriptionID.String()
		payload["plan"] = e.Plan
		payload["amount"] = e.Amount
		payload["currency"] = e.Currency
	case UserRegisteredEvent:
		payload["user_id"] = e.UserID.String()
		payload["email"] = e.Email
		payload["name"] = e.Name
	}

	raw, _ := json.Marshal(payload)

	return OutboxEvent{
		ID:        uuid.New().String(),
		EventType: event.EventName(),
		Payload:   raw,
		Status:    "pending",
		CreatedAt: time.Now().UTC(),
	}
}
