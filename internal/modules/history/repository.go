package history

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrTrackNotFound = errors.New("track not found")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByUser(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]ListeningHistoryItem, error) {
	var items []ListeningHistoryItem

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			id,
			user_id,
			track_id,
			played_at,
			duration,
			completed
		FROM listening_history
		WHERE user_id = $1
		ORDER BY played_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Record is intentionally available at repository/service level,
// even if the only public endpoint right now is GET /history.
//
// Later your player/playback module can call:
// historyService.Record(ctx, userID, req)
func (r *Repository) Record(
	ctx context.Context,
	userID uuid.UUID,
	trackID uuid.UUID,
	duration int,
	completed bool,
) (*ListeningHistoryItem, error) {
	exists, err := r.trackExists(ctx, trackID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTrackNotFound
	}

	var item ListeningHistoryItem

	err = r.db.GetContext(ctx, &item, `
		INSERT INTO listening_history (
			user_id,
			track_id,
			duration,
			completed
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			user_id,
			track_id,
			played_at,
			duration,
			completed
	`, userID, trackID, duration, completed)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) trackExists(ctx context.Context, trackID uuid.UUID) (bool, error) {
	var exists bool

	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM tracks
			WHERE id = $1
		)
	`, trackID)

	return exists, err
}
