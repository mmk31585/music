package search

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SearchTracks(ctx context.Context, query string, limit int) ([]TrackResult, error) {
	sql := `
		SELECT
			t.id::text AS id,
			t.title,
			t.artist_id::text AS artist_id,
			a.name AS artist_name,
			t.album_id::text AS album_id,
			al.title AS album_title,
			t.cover_url,
			t.audio_url,
			t.duration_seconds
		FROM tracks t
		LEFT JOIN artists a ON a.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		WHERE t.title ILIKE $1 OR a.name ILIKE $1 OR al.title ILIKE $1
		ORDER BY t.play_count DESC
		LIMIT $2
	`
	pattern := fmt.Sprintf("%%%s%%", query)
	var results []TrackResult
	err := r.db.SelectContext(ctx, &results, sql, pattern, limit)
	return results, err
}

func (r *Repository) SearchAlbums(ctx context.Context, query string, limit int) ([]AlbumResult, error) {
	sql := `
		SELECT
			al.id::text AS id,
			al.title,
			al.artist_id::text AS artist_id,
			a.name AS artist_name,
			al.cover_url,
			al.release_date::text AS release_date
		FROM albums al
		LEFT JOIN artists a ON a.id = al.artist_id
		WHERE al.title ILIKE $1 OR a.name ILIKE $1
		ORDER BY al.release_date DESC NULLS LAST
		LIMIT $2
	`
	pattern := fmt.Sprintf("%%%s%%", query)
	var results []AlbumResult
	err := r.db.SelectContext(ctx, &results, sql, pattern, limit)
	return results, err
}

func (r *Repository) SearchArtists(ctx context.Context, query string, limit int) ([]ArtistResult, error) {
	sql := `
		SELECT
			id::text AS id,
			name,
			image_url AS cover_url
		FROM artists
		WHERE name ILIKE $1
		ORDER BY monthly_listeners DESC
		LIMIT $2
	`
	pattern := fmt.Sprintf("%%%s%%", query)
	var results []ArtistResult
	err := r.db.SelectContext(ctx, &results, sql, pattern, limit)
	return results, err
}

func (r *Repository) GetUserTopGenres(ctx context.Context, userID string, limit int) ([]string, error) {
	var genres []string
	err := r.db.SelectContext(ctx, &genres, `
		SELECT g.name FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		JOIN track_genres tg ON tg.track_id = t.id
		JOIN genres g ON g.id = tg.genre_id
		WHERE lh.user_id = $1
		GROUP BY g.name
		ORDER BY COUNT(*) DESC
		LIMIT $2
	`, userID, limit)
	return genres, err
}

func (r *Repository) GetUserTopArtistIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	var ids []string
	err := r.db.SelectContext(ctx, &ids, `
		SELECT a.id::text FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		JOIN artists a ON a.id = t.artist_id
		WHERE lh.user_id = $1
		GROUP BY a.id
		ORDER BY COUNT(*) DESC
		LIMIT $2
	`, userID, limit)
	return ids, err
}

func (r *Repository) SearchPlaylists(ctx context.Context, query string, limit int) ([]PlaylistResult, error) {
	sql := `
		SELECT
			id::text AS id,
			name,
			description,
			cover_url,
			user_id::text AS user_id,
			is_public
		FROM playlists
		WHERE name ILIKE $1 AND is_public = true
		ORDER BY name ASC
		LIMIT $2
	`
	pattern := fmt.Sprintf("%%%s%%", query)
	var results []PlaylistResult
	err := r.db.SelectContext(ctx, &results, sql, pattern, limit)
	return results, err
}
