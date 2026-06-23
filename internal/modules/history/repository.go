package history

import (
	"context"
	"errors"
	"fmt"

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
			lh.id,
			lh.user_id,
			lh.track_id,
			lh.played_at,
			lh.duration,
			lh.completed,
			COALESCE(t.title, '') AS track_title,
			t.duration_seconds AS track_duration,
			COALESCE(t.cover_url, '') AS track_cover_url,
			COALESCE(a.name, '') AS artist_name
		FROM listening_history lh
		LEFT JOIN tracks t ON t.id = lh.track_id
		LEFT JOIN track_artists ta ON ta.track_id = t.id AND ta.role = 'primary'
		LEFT JOIN artists a ON a.id = ta.artist_id
		WHERE lh.user_id = $1
		ORDER BY lh.played_at DESC, lh.id DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Record inserts a basic listening history entry (backward-compatible).
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

// RecordSignal inserts a signal-enriched listening history entry.
func (r *Repository) RecordSignal(
	ctx context.Context,
	userID, trackID uuid.UUID,
	playedDurationMs, trackDurationMs int64,
	completionPercent float64,
	signalType string,
	isExplicitLike bool,
	sessionID uuid.UUID,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO listening_history (
			user_id,
			track_id,
			duration,
			completed,
			session_id,
			played_duration_ms,
			track_duration_ms,
			completion_percent,
			signal_type,
			is_explicit_like
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		userID,
		trackID,
		playedDurationMs/1000, // duration in seconds for backward compat
		signalType == SignalTypeCompletePositive || signalType == SignalTypeReplayStrong,
		sessionID,
		playedDurationMs,
		trackDurationMs,
		completionPercent,
		signalType,
		isExplicitLike,
	)
	return err
}

// ListSignalsByUser returns signal-enriched history for a user within
// the last N days. Used by the taste profile recomputation.
func (r *Repository) ListSignalsByUser(ctx context.Context, userID uuid.UUID, days int) ([]PlaybackSignal, error) {
	var signals []PlaybackSignal
	err := r.db.SelectContext(ctx, &signals, `
		SELECT
			id,
			user_id,
			track_id,
			COALESCE(session_id, '00000000-0000-0000-0000-000000000000') AS session_id,
			COALESCE(played_duration_ms, duration * 1000) AS played_duration_ms,
			COALESCE(track_duration_ms, 0) AS track_duration_ms,
			COALESCE(completion_percent, 0) AS completion_percent,
			COALESCE(signal_type, '') AS signal_type,
			COALESCE(is_explicit_like, false) AS is_explicit_like,
			played_at
		FROM listening_history
		WHERE user_id = $1
		  AND played_at >= NOW() - ($2 || ' days')::INTERVAL
		ORDER BY played_at DESC
	`, userID, fmt.Sprintf("%d", days))
	if err != nil {
		return nil, err
	}
	return signals, nil
}

// UpdateSignalExplicitLike marks the most recent listening_history entry
// for the given (user, track) as explicitly liked.
func (r *Repository) UpdateSignalExplicitLike(ctx context.Context, userID, trackID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE listening_history
		SET is_explicit_like = true
		WHERE id = (
			SELECT id FROM listening_history
			WHERE user_id = $1 AND track_id = $2
			ORDER BY played_at DESC
			LIMIT 1
		)
	`)
	return err
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
