package analytics

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrTrackNotFound    = errors.New("track not found")
	ErrArtistNotFound   = errors.New("artist not found")
	ErrAlbumNotFound    = errors.New("album not found")
	ErrPlaylistNotFound = errors.New("playlist not found")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type CreateEventParams struct {
	UserID     *uuid.UUID
	EventType  EventType
	TrackID    *uuid.UUID
	ArtistID   *uuid.UUID
	AlbumID    *uuid.UUID
	PlaylistID *uuid.UUID
	Query      *string
	Metadata   map[string]any
}

func (r *Repository) CreateEvent(ctx context.Context, params CreateEventParams) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)

	if err := r.validateReferences(ctx, tx, params); err != nil {
		return err
	}

	if params.Metadata == nil {
		params.Metadata = map[string]any{}
	}

	metadataJSON, err := json.Marshal(params.Metadata)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO analytics_events (
			user_id,
			event_type,
			track_id,
			artist_id,
			album_id,
			playlist_id,
			query,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb)
	`,
		params.UserID,
		string(params.EventType),
		params.TrackID,
		params.ArtistID,
		params.AlbumID,
		params.PlaylistID,
		params.Query,
		string(metadataJSON),
	)
	if err != nil {
		return err
	}

	switch params.EventType {
	case EventPlay, EventPause, EventSkip, EventCompletion:
		if params.TrackID != nil {
			if err := r.upsertTrackAnalytics(ctx, tx, *params.TrackID, params.EventType); err != nil {
				return err
			}
		}

		if params.ArtistID != nil {
			if err := r.upsertArtistAnalytics(ctx, tx, *params.ArtistID, params.EventType); err != nil {
				return err
			}
		}

		if params.AlbumID != nil {
			if err := r.upsertAlbumAnalytics(ctx, tx, *params.AlbumID, params.EventType); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *Repository) validateReferences(ctx context.Context, tx *sqlx.Tx, params CreateEventParams) error {
	if params.TrackID != nil {
		ok, err := existsByID(ctx, tx, "tracks", *params.TrackID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrTrackNotFound
		}
	}

	if params.ArtistID != nil {
		ok, err := existsByID(ctx, tx, "artists", *params.ArtistID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrArtistNotFound
		}
	}

	if params.AlbumID != nil {
		ok, err := existsByID(ctx, tx, "albums", *params.AlbumID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrAlbumNotFound
		}
	}

	if params.PlaylistID != nil {
		ok, err := existsByID(ctx, tx, "playlists", *params.PlaylistID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrPlaylistNotFound
		}
	}

	return nil
}

func (r *Repository) upsertTrackAnalytics(ctx context.Context, tx *sqlx.Tx, trackID uuid.UUID, eventType EventType) error {
	playInc, pauseInc, skipInc, completionInc := counterFlags(eventType)

	_, err := tx.ExecContext(ctx, `
		INSERT INTO track_analytics (
			track_id,
			play_count,
			pause_count,
			skip_count,
			completion_count,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (track_id)
		DO UPDATE SET
			play_count = track_analytics.play_count + EXCLUDED.play_count,
			pause_count = track_analytics.pause_count + EXCLUDED.pause_count,
			skip_count = track_analytics.skip_count + EXCLUDED.skip_count,
			completion_count = track_analytics.completion_count + EXCLUDED.completion_count,
			updated_at = NOW()
	`, trackID, playInc, pauseInc, skipInc, completionInc)

	return err
}

func (r *Repository) upsertArtistAnalytics(ctx context.Context, tx *sqlx.Tx, artistID uuid.UUID, eventType EventType) error {
	playInc, pauseInc, skipInc, completionInc := counterFlags(eventType)

	_, err := tx.ExecContext(ctx, `
		INSERT INTO artist_analytics (
			artist_id,
			play_count,
			pause_count,
			skip_count,
			completion_count,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (artist_id)
		DO UPDATE SET
			play_count = artist_analytics.play_count + EXCLUDED.play_count,
			pause_count = artist_analytics.pause_count + EXCLUDED.pause_count,
			skip_count = artist_analytics.skip_count + EXCLUDED.skip_count,
			completion_count = artist_analytics.completion_count + EXCLUDED.completion_count,
			updated_at = NOW()
	`, artistID, playInc, pauseInc, skipInc, completionInc)

	return err
}

func (r *Repository) upsertAlbumAnalytics(ctx context.Context, tx *sqlx.Tx, albumID uuid.UUID, eventType EventType) error {
	playInc, pauseInc, skipInc, completionInc := counterFlags(eventType)

	_, err := tx.ExecContext(ctx, `
		INSERT INTO album_analytics (
			album_id,
			play_count,
			pause_count,
			skip_count,
			completion_count,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (album_id)
		DO UPDATE SET
			play_count = album_analytics.play_count + EXCLUDED.play_count,
			pause_count = album_analytics.pause_count + EXCLUDED.pause_count,
			skip_count = album_analytics.skip_count + EXCLUDED.skip_count,
			completion_count = album_analytics.completion_count + EXCLUDED.completion_count,
			updated_at = NOW()
	`, albumID, playInc, pauseInc, skipInc, completionInc)

	return err
}

func existsByID(ctx context.Context, tx *sqlx.Tx, table string, id uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM ` + table + `
			WHERE id = $1
		)
	`

	var exists bool
	err := tx.GetContext(ctx, &exists, query, id)
	return exists, err
}

func counterFlags(eventType EventType) (play int64, pause int64, skip int64, completion int64) {
	switch eventType {
	case EventPlay:
		play = 1
	case EventPause:
		pause = 1
	case EventSkip:
		skip = 1
	case EventCompletion:
		completion = 1
	}
	return
}

func rollback(tx *sqlx.Tx) {
	_ = tx.Rollback()
}

type Overview struct {
	TotalTracks    int     `db:"total_tracks"`
	TotalUsers     int     `db:"total_users"`
	TotalAlbums    int     `db:"total_albums"`
	TotalPlays     int     `db:"total_plays"`
	ActiveUsers24h int     `db:"active_users_last_24h"`
	StorageUsedMB  float64 `db:"storage_used_mb"`
}

func (r *Repository) GetOverview(ctx context.Context) (Overview, error) {
	var ov Overview
	err := r.db.GetContext(ctx, &ov, `
		SELECT
			(SELECT COUNT(*) FROM tracks) AS total_tracks,
			(SELECT COUNT(*) FROM users) AS total_users,
			(SELECT COUNT(*) FROM albums) AS total_albums,
			(SELECT COUNT(*) FROM listening_history WHERE played_at >= NOW() - INTERVAL '24 hours') AS active_users_last_24h,
			COALESCE((
				SELECT COUNT(DISTINCT user_id)
				FROM listening_history
				WHERE played_at >= NOW() - INTERVAL '24 hours'
			), 0) AS active_users_last_24h,
			COALESCE((
				SELECT SUM(file_size)::numeric / 1048576.0
				FROM media_assets
				WHERE deleted_at IS NULL
			), 0) AS storage_used_mb
	`)
	if err != nil {
		return Overview{}, err
	}

	// TotalPlays — get count from the partitioned listening_history table
	err = r.db.GetContext(ctx, &ov.TotalPlays, `SELECT COUNT(*) FROM listening_history`)
	if err != nil {
		return Overview{}, err
	}

	return ov, nil
}

func nullableUUID(id *uuid.UUID) any {
	if id == nil {
		return sql.NullString{}
	}
	return *id
}
