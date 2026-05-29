package notification

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, in CreateNotificationInput) (*Notification, error) {
	payload := json.RawMessage(`{}`)
	if in.Payload != nil {
		b, err := json.Marshal(in.Payload)
		if err != nil {
			return nil, err
		}
		payload = b
	}

	var n Notification
	err := r.db.GetContext(ctx, &n, `
		INSERT INTO notifications (
			user_id,
			type,
			title,
			body,
			entity_type,
			entity_id,
			payload,
			is_read
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, FALSE)
		RETURNING
			id,
			user_id,
			type,
			title,
			body,
			entity_type,
			entity_id,
			payload,
			is_read,
			read_at,
			created_at
	`,
		in.UserID,
		in.Type,
		in.Title,
		in.Body,
		in.EntityType,
		in.EntityID,
		payload,
	)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *Repository) CreateBulk(ctx context.Context, inputs []CreateNotificationInput) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		INSERT INTO notifications (
			user_id,
			type,
			title,
			body,
			entity_type,
			entity_id,
			payload,
			is_read
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, FALSE)
	`

	for _, in := range inputs {
		payload := json.RawMessage(`{}`)
		if in.Payload != nil {
			b, err := json.Marshal(in.Payload)
			if err != nil {
				return err
			}
			payload = b
		}

		if _, err := tx.ExecContext(
			ctx,
			query,
			in.UserID,
			in.Type,
			in.Title,
			in.Body,
			in.EntityType,
			in.EntityID,
			payload,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]Notification, error) {
	var items []Notification

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			id,
			user_id,
			type,
			title,
			body,
			entity_type,
			entity_id,
			payload,
			is_read,
			read_at,
			created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	return items, err
}

func (r *Repository) CountUnreadByUser(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1
		  AND is_read = FALSE
	`, userID)
	return count, err
}

func (r *Repository) MarkAsRead(ctx context.Context, userID, notificationID string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE notifications
		SET
			is_read = TRUE,
			read_at = NOW()
		WHERE id = $1
		  AND user_id = $2
		  AND is_read = FALSE
	`, notificationID, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		var exists bool
		err = r.db.GetContext(ctx, &exists, `
			SELECT EXISTS(
				SELECT 1 FROM notifications WHERE id = $1 AND user_id = $2
			)
		`, notificationID, userID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrNotFound
		}
	}

	return nil
}

func (r *Repository) MarkAllAsRead(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notifications
		SET
			is_read = TRUE,
			read_at = NOW()
		WHERE user_id = $1
		  AND is_read = FALSE
	`, userID)
	return err
}

func (r *Repository) GetByID(ctx context.Context, userID, notificationID string) (*Notification, error) {
	var n Notification
	err := r.db.GetContext(ctx, &n, `
		SELECT
			id,
			user_id,
			type,
			title,
			body,
			entity_type,
			entity_id,
			payload,
			is_read,
			read_at,
			created_at
		FROM notifications
		WHERE id = $1
		  AND user_id = $2
	`, notificationID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}
