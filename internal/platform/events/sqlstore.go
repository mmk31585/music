package events

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type SQLOutboxStore struct {
	db *sqlx.DB
}

func NewSQLOutboxStore(db *sqlx.DB) *SQLOutboxStore {
	return &SQLOutboxStore{db: db}
}

func (s *SQLOutboxStore) Insert(ctx context.Context, event OutboxEvent) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO event_outbox (id, event_type, payload, status, created_at, retry_count)
		 VALUES ($1, $2, $3::jsonb, 'pending', $4, 0)
		 ON CONFLICT (id) DO NOTHING`,
		event.ID, event.EventType, string(event.Payload), event.CreatedAt,
	)
	return err
}

func (s *SQLOutboxStore) FetchPending(ctx context.Context, limit int) ([]OutboxEvent, error) {
	var events []OutboxEvent
	err := s.db.SelectContext(ctx, &events,
		`SELECT id, event_type, payload, status, created_at, processed_at, retry_count, last_error
		 FROM event_outbox WHERE status = 'pending'
		 ORDER BY created_at ASC LIMIT $1 FOR UPDATE SKIP LOCKED`, limit)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLOutboxStore) MarkProcessed(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE event_outbox SET status = 'processed', processed_at = NOW()
		 WHERE id = $1`, id)
	return err
}

func (s *SQLOutboxStore) MarkFailed(ctx context.Context, id string, errMsg string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE event_outbox SET status = 'failed', retry_count = retry_count + 1, last_error = $2
		 WHERE id = $1`, id, errMsg)
	return err
}

func (s *SQLOutboxStore) DeleteProcessed(ctx context.Context, before time.Time) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM event_outbox WHERE status = 'processed' AND processed_at < $1`, before)
	return err
}
