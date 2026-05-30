package queue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrQueueItemNotFound = errors.New("queue item not found")
	ErrTrackNotFound     = errors.New("track not found")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, userID uuid.UUID) ([]QueueItem, error) {
	var items []QueueItem

	err := r.db.SelectContext(ctx, &items, `
		SELECT id, user_id, track_id, position, created_at, updated_at
		FROM queue_items
		WHERE user_id = $1
		ORDER BY position ASC, created_at ASC
	`, userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) AddLater(ctx context.Context, userID, trackID uuid.UUID) (*QueueItem, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer rollback(tx)

	if err := lockUserQueue(ctx, tx, userID); err != nil {
		return nil, err
	}

	exists, err := r.trackExists(ctx, tx, trackID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTrackNotFound
	}

	var position int
	err = tx.GetContext(ctx, &position, `
		SELECT COALESCE(MAX(position), 0) + 1
		FROM queue_items
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	var item QueueItem
	err = tx.GetContext(ctx, &item, `
		INSERT INTO queue_items (
			user_id,
			track_id,
			position
		)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, track_id, position, created_at, updated_at
	`, userID, trackID, position)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) AddNext(ctx context.Context, userID, trackID uuid.UUID) (*QueueItem, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer rollback(tx)

	if err := lockUserQueue(ctx, tx, userID); err != nil {
		return nil, err
	}

	exists, err := r.trackExists(ctx, tx, trackID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTrackNotFound
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE queue_items
		SET position = position + 1,
		    updated_at = NOW()
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	var item QueueItem
	err = tx.GetContext(ctx, &item, `
		INSERT INTO queue_items (
			user_id,
			track_id,
			position
		)
		VALUES ($1, $2, 1)
		RETURNING id, user_id, track_id, position, created_at, updated_at
	`, userID, trackID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) Remove(ctx context.Context, userID, itemID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)

	if err := lockUserQueue(ctx, tx, userID); err != nil {
		return err
	}

	var removedPosition int
	err = tx.GetContext(ctx, &removedPosition, `
		DELETE FROM queue_items
		WHERE id = $1 AND user_id = $2
		RETURNING position
	`, itemID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrQueueItemNotFound
	}
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE queue_items
		SET position = position - 1,
		    updated_at = NOW()
		WHERE user_id = $1
		  AND position > $2
	`, userID, removedPosition)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) Clear(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM queue_items
		WHERE user_id = $1
	`, userID)

	return err
}

func (r *Repository) Reorder(ctx context.Context, userID uuid.UUID, items []ReorderItem) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)

	if err := lockUserQueue(ctx, tx, userID); err != nil {
		return err
	}

	for _, item := range items {
		result, err := tx.ExecContext(ctx, `
			UPDATE queue_items
			SET position = $1,
			    updated_at = NOW()
			WHERE id = $2 AND user_id = $3
		`, item.Position, item.ID, userID)
		if err != nil {
			return err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if affected == 0 {
			return ErrQueueItemNotFound
		}
	}

	if err := normalizePositions(ctx, tx, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) Move(ctx context.Context, userID, itemID uuid.UUID, newPosition int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)

	if err := lockUserQueue(ctx, tx, userID); err != nil {
		return err
	}

	var currentPosition int
	err = tx.GetContext(ctx, &currentPosition, `
		SELECT position
		FROM queue_items
		WHERE id = $1 AND user_id = $2
	`, itemID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrQueueItemNotFound
	}
	if err != nil {
		return err
	}

	var maxPosition int
	err = tx.GetContext(ctx, &maxPosition, `
		SELECT COALESCE(MAX(position), 0)
		FROM queue_items
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return err
	}

	if maxPosition == 0 {
		return ErrQueueItemNotFound
	}

	if newPosition < 1 {
		newPosition = 1
	}

	if newPosition > maxPosition {
		newPosition = maxPosition
	}

	if currentPosition == newPosition {
		return tx.Commit()
	}

	if newPosition < currentPosition {
		_, err = tx.ExecContext(ctx, `
			UPDATE queue_items
			SET position = position + 1,
			    updated_at = NOW()
			WHERE user_id = $1
			  AND position >= $2
			  AND position < $3
		`, userID, newPosition, currentPosition)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.ExecContext(ctx, `
			UPDATE queue_items
			SET position = position - 1,
			    updated_at = NOW()
			WHERE user_id = $1
			  AND position <= $2
			  AND position > $3
		`, userID, newPosition, currentPosition)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE queue_items
		SET position = $1,
		    updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`, newPosition, itemID, userID)
	if err != nil {
		return err
	}

	if err := normalizePositions(ctx, tx, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) trackExists(ctx context.Context, tx *sqlx.Tx, trackID uuid.UUID) (bool, error) {
	var exists bool

	err := tx.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM tracks
			WHERE id = $1
		)
	`, trackID)

	return exists, err
}

func lockUserQueue(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) error {
	_, err := tx.ExecContext(ctx, `
		SELECT pg_advisory_xact_lock(hashtext($1))
	`, fmt.Sprintf("queue:%s", userID.String()))

	return err
}

func normalizePositions(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) error {
	_, err := tx.ExecContext(ctx, `
		WITH ordered AS (
			SELECT
				id,
				ROW_NUMBER() OVER (
					ORDER BY position ASC, created_at ASC, id ASC
				) AS new_position
			FROM queue_items
			WHERE user_id = $1
		)
		UPDATE queue_items qi
		SET position = ordered.new_position,
		    updated_at = NOW()
		FROM ordered
		WHERE qi.id = ordered.id
	`, userID)

	return err
}

func rollback(tx *sqlx.Tx) {
	_ = tx.Rollback()
}

type ReorderItem struct {
	ID       uuid.UUID
	Position int
}
