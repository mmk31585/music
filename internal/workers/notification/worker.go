package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Worker struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewWorker(db *pgxpool.Pool, logger *zap.Logger) *Worker {
	return &Worker{
		db:     db,
		logger: logger,
	}
}

func (w *Worker) Name() string {
	return "notification"
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.logger.Info("processing notifications")
			if err := w.process(ctx); err != nil {
				w.logger.Error("notification processing failed", zap.Error(err))
			}
		}
	}
}

type outboxEvent struct {
	ID        string          `db:"id"`
	EventType string          `db:"event_type"`
	Payload   json.RawMessage `db:"payload"`
}

var notificationEventTypes = []string{
	"user.registered",
	"playlist.created",
	"subscription.purchased",
}

func (w *Worker) process(ctx context.Context) error {
	rows, err := w.db.Query(ctx, `
		SELECT id, event_type, payload
		FROM event_outbox
		WHERE status = 'pending'
		  AND event_type = ANY($1)
		ORDER BY created_at ASC
		LIMIT 50
		FOR UPDATE SKIP LOCKED
	`, notificationEventTypes)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var e outboxEvent
		if err := rows.Scan(&e.ID, &e.EventType, &e.Payload); err != nil {
			continue
		}

		if err := w.handleEvent(ctx, e); err != nil {
			w.logger.Error("failed to handle event for notification",
				zap.String("event_id", e.ID),
				zap.String("event_type", e.EventType),
				zap.Error(err))
			w.markFailed(ctx, e.ID, err.Error())
			continue
		}

		if _, err := w.db.Exec(ctx,
			`UPDATE event_outbox SET status = 'processed', processed_at = NOW() WHERE id = $1`,
			e.ID,
		); err != nil {
			w.logger.Error("failed to mark event processed",
				zap.String("event_id", e.ID), zap.Error(err))
		}
	}

	return rows.Err()
}

func (w *Worker) handleEvent(ctx context.Context, e outboxEvent) error {
	switch e.EventType {
	case "user.registered":
		return w.handleUserRegistered(ctx, e.Payload)
	case "playlist.created":
		return w.handlePlaylistCreated(ctx, e.Payload)
	case "subscription.purchased":
		return w.handleSubscriptionPurchased(ctx, e.Payload)
	}
	return nil
}

type userRegisteredPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

func (w *Worker) handleUserRegistered(ctx context.Context, raw json.RawMessage) error {
	var p userRegisteredPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	_, err := w.db.Exec(ctx, `
		INSERT INTO notifications (user_id, type, title, body, payload, is_read)
		VALUES ($1, 'welcome', 'Welcome', $2, $3::jsonb, FALSE)
	`, p.UserID,
		fmt.Sprintf("Welcome %s! Your account has been created successfully.", p.Name),
		fmt.Sprintf(`{"userId":"%s","email":"%s","name":"%s"}`, p.UserID, p.Email, p.Name))
	return err
}

type playlistCreatedPayload struct {
	UserID     string `json:"user_id"`
	PlaylistID string `json:"playlist_id"`
	Name       string `json:"name"`
	IsPublic   bool   `json:"is_public"`
}

func (w *Worker) handlePlaylistCreated(ctx context.Context, raw json.RawMessage) error {
	var p playlistCreatedPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	entityType := "playlist"
	_, err := w.db.Exec(ctx, `
		INSERT INTO notifications (user_id, type, title, body, entity_type, entity_id, payload, is_read)
		VALUES ($1, 'playlist_created', 'Playlist created', $2, $3, $4, $5::jsonb, FALSE)
	`, p.UserID,
		fmt.Sprintf("Your playlist \"%s\" was created successfully.", p.Name),
		entityType, p.PlaylistID,
		fmt.Sprintf(`{"playlistId":"%s","name":"%s"}`, p.PlaylistID, p.Name))
	return err
}

type subscriptionPurchasedPayload struct {
	UserID         string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
	Plan           string `json:"plan"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
}

func (w *Worker) handleSubscriptionPurchased(ctx context.Context, raw json.RawMessage) error {
	var p subscriptionPurchasedPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	entityType := "subscription"
	entityID := p.SubscriptionID
	_, err := w.db.Exec(ctx, `
		INSERT INTO notifications (user_id, type, title, body, entity_type, entity_id, payload, is_read)
		VALUES ($1, 'subscription_purchased', 'Subscription activated', $2, $3, $4, $5::jsonb, FALSE)
	`, p.UserID,
		fmt.Sprintf("Your %s subscription has been activated.", p.Plan),
		entityType, entityID,
		fmt.Sprintf(`{"subscriptionId":"%s","plan":"%s","amount":%d,"currency":"%s"}`,
			p.SubscriptionID, p.Plan, p.Amount, p.Currency))
	return err
}

func (w *Worker) markFailed(ctx context.Context, id, errMsg string) {
	w.db.Exec(ctx,
		`UPDATE event_outbox SET status = 'failed', retry_count = retry_count + 1, last_error = $2
		 WHERE id = $1`, id, errMsg)
}
