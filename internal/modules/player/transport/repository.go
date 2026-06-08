package player

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

var ErrTrackNotFound = errors.New("track not found")

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetTrackForPlayback(ctx context.Context, id string) (*PlaybackTrack, error) {
	query := `
		SELECT t.id, t.title,
		       COALESCE(ar.name, '') AS artist_name,
		       al.title AS album_title,
		       COALESCE(t.audio_url, '') AS audio_url,  
		       COALESCE(t.cover_url, '') AS cover_url,   
		       t.duration_seconds,
		       COALESCE(t.is_public, true) AS is_public,
		       t.created_at
		FROM tracks t
		LEFT JOIN artists ar ON ar.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		WHERE t.id = $1
		LIMIT 1
	`

	var track PlaybackTrack

	err := r.db.QueryRowxContext(ctx, query, id).Scan(
		&track.ID,
		&track.Title,
		&track.ArtistName,
		&track.AlbumTitle,
		&track.AudioURL,
		&track.CoverURL,
		&track.DurationSeconds,
		&track.IsPublic,
		&track.CreatedAt,
	)
	log.Print(track)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTrackNotFound
		}

		return nil, err
	}

	return &track, nil
}

func (r *Repository) IncrementPlayCount(ctx context.Context, trackID string) error {
	query := `UPDATE tracks SET play_count = COALESCE(play_count, 0) + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, trackID)
	return err
}
