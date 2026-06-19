package search

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type trackScanRow struct {
	ID              string          `db:"id"`
	Title           string          `db:"title"`
	Slug            string          `db:"slug"`
	ArtistID        *string         `db:"artist_id"`
	AlbumID         *string         `db:"album_id"`
	CoverURL        *string         `db:"cover_url"`
	AudioURL        *string         `db:"audio_url"`
	DurationSeconds int             `db:"duration_seconds"`
	TrackNumber     *int            `db:"track_number"`
	Explicit        bool            `db:"explicit"`
	PlayCount       int64           `db:"play_count"`
	IsPublic        bool            `db:"is_public"`
	CreatedAt       string          `db:"created_at"`
	ArtistsJSON     json.RawMessage `db:"artists"`
	GenresJSON      json.RawMessage `db:"genres"`
}

type albumScanRow struct {
	ID          string          `db:"id"`
	Title       string          `db:"title"`
	Slug        string          `db:"slug"`
	ArtistID    *string         `db:"artist_id"`
	CoverURL    *string         `db:"cover_url"`
	ReleaseDate *string         `db:"release_date"`
	ArtistsJSON json.RawMessage `db:"artists"`
}

type artistScanRow struct {
	ID               string  `db:"id"`
	Name             string  `db:"name"`
	Slug             string  `db:"slug"`
	Bio              *string `db:"bio"`
	ImageURL         *string `db:"image_url"`
	IsVerified       bool    `db:"is_verified"`
	MonthlyListeners int     `db:"monthly_listeners"`
}

type playlistScanRow struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	CoverURL    *string `db:"cover_url"`
	UserID      *string `db:"user_id"`
	IsPublic    *bool   `db:"is_public"`
}

func (r *Repository) SearchTracks(ctx context.Context, query string, limit int) ([]TrackResult, error) {
	sql := `
		SELECT
			t.id::text AS id,
			t.title,
			t.slug,
			t.artist_id::text AS artist_id,
			t.album_id::text AS album_id,
			t.cover_url,
			t.audio_url,
			t.duration_seconds,
			t.track_number,
			t.explicit,
			t.play_count,
			t.is_public,
			t.created_at::text AS created_at,
			COALESCE(
				(SELECT json_agg(jsonb_build_object(
					'artist_id', ta.artist_id::text,
					'name', a.name,
					'slug', a.slug,
					'role', ta.role,
					'position', ta.position
				) ORDER BY ta.position)
				FROM track_artists ta
				JOIN artists a ON a.id = ta.artist_id
				WHERE ta.track_id = t.id),
				'[]'::json
			) AS artists,
			COALESCE(
				(SELECT json_agg(jsonb_build_object(
					'id', g.id::text,
					'name', g.name,
					'slug', g.slug
				))
				FROM track_genres tg
				JOIN genres g ON g.id = tg.genre_id
				WHERE tg.track_id = t.id),
				'[]'::json
			) AS genres
		FROM tracks t
		LEFT JOIN artists artist ON artist.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		WHERE t.is_public = true AND (t.title ILIKE $1 OR artist.name ILIKE $1 OR al.title ILIKE $1)
		ORDER BY t.play_count DESC
		LIMIT $2
	`
	pattern := fmt.Sprintf("%%%s%%", query)
	var rows []trackScanRow
	err := r.db.SelectContext(ctx, &rows, sql, pattern, limit)
	if err != nil {
		return nil, err
	}
	return toTrackResults(rows), nil
}

func (r *Repository) SearchAlbums(ctx context.Context, query string, limit int) ([]AlbumResult, error) {
	sql := `
		SELECT
			al.id::text AS id,
			al.title,
			al.slug,
			al.artist_id::text AS artist_id,
			al.cover_url,
			al.release_date::text AS release_date,
			COALESCE(
				(SELECT json_agg(jsonb_build_object(
					'artist_id', aa.artist_id::text,
					'name', a.name,
					'slug', a.slug,
					'role', aa.role,
					'position', aa.position
				) ORDER BY aa.position)
				FROM album_artists aa
				JOIN artists a ON a.id = aa.artist_id
				WHERE aa.album_id = al.id),
				'[]'::json
			) AS artists
		FROM albums al
		LEFT JOIN artists artist ON artist.id = al.artist_id
		WHERE al.title ILIKE $1 OR artist.name ILIKE $1
		ORDER BY al.release_date DESC NULLS LAST
		LIMIT $2
	`
	pattern := fmt.Sprintf("%%%s%%", query)
	var rows []albumScanRow
	err := r.db.SelectContext(ctx, &rows, sql, pattern, limit)
	if err != nil {
		return nil, err
	}
	return toAlbumResults(rows), nil
}

func (r *Repository) SearchArtists(ctx context.Context, query string, limit int) ([]ArtistResult, error) {
	sql := `
		SELECT
			id::text AS id,
			name,
			slug,
			bio,
			image_url,
			is_verified,
			monthly_listeners
		FROM artists
		WHERE name ILIKE $1
		ORDER BY monthly_listeners DESC
		LIMIT $2
	`
	pattern := fmt.Sprintf("%%%s%%", query)
	var rows []artistScanRow
	err := r.db.SelectContext(ctx, &rows, sql, pattern, limit)
	if err != nil {
		return nil, err
	}
	return toArtistResults(rows), nil
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
	var rows []playlistScanRow
	err := r.db.SelectContext(ctx, &rows, sql, pattern, limit)
	if err != nil {
		return nil, err
	}
	return toPlaylistResults(rows), nil
}

func (r *Repository) SearchTracksByIDs(ctx context.Context, ids []string) ([]TrackResult, error) {
	if len(ids) == 0 {
		return []TrackResult{}, nil
	}
	query, args, err := sqlx.In(`
		SELECT
			t.id::text AS id,
			t.title,
			t.slug,
			t.artist_id::text AS artist_id,
			t.album_id::text AS album_id,
			t.cover_url,
			t.audio_url,
			t.duration_seconds,
			t.track_number,
			t.explicit,
			t.play_count,
			t.is_public,
			t.created_at::text AS created_at,
			COALESCE(
				(SELECT json_agg(jsonb_build_object(
					'artist_id', ta.artist_id::text,
					'name', a.name,
					'slug', a.slug,
					'role', ta.role,
					'position', ta.position
				) ORDER BY ta.position)
				FROM track_artists ta
				JOIN artists a ON a.id = ta.artist_id
				WHERE ta.track_id = t.id),
				'[]'::json
			) AS artists,
			COALESCE(
				(SELECT json_agg(jsonb_build_object(
					'id', g.id::text,
					'name', g.name,
					'slug', g.slug
				))
				FROM track_genres tg
				JOIN genres g ON g.id = tg.genre_id
				WHERE tg.track_id = t.id),
				'[]'::json
			) AS genres
		FROM tracks t
		WHERE t.id::text IN (?)
		ORDER BY t.play_count DESC
	`, ids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var rows []trackScanRow
	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return toTrackResults(rows), nil
}

func (r *Repository) SearchAlbumsByIDs(ctx context.Context, ids []string) ([]AlbumResult, error) {
	if len(ids) == 0 {
		return []AlbumResult{}, nil
	}
	query, args, err := sqlx.In(`
		SELECT
			al.id::text AS id,
			al.title,
			al.slug,
			al.artist_id::text AS artist_id,
			al.cover_url,
			al.release_date::text AS release_date,
			COALESCE(
				(SELECT json_agg(jsonb_build_object(
					'artist_id', aa.artist_id::text,
					'name', a.name,
					'slug', a.slug,
					'role', aa.role,
					'position', aa.position
				) ORDER BY aa.position)
				FROM album_artists aa
				JOIN artists a ON a.id = aa.artist_id
				WHERE aa.album_id = al.id),
				'[]'::json
			) AS artists
		FROM albums al
		WHERE al.id::text IN (?)
		ORDER BY al.release_date DESC NULLS LAST
	`, ids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var rows []albumScanRow
	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return toAlbumResults(rows), nil
}

func (r *Repository) SearchArtistsByIDs(ctx context.Context, ids []string) ([]ArtistResult, error) {
	if len(ids) == 0 {
		return []ArtistResult{}, nil
	}
	query, args, err := sqlx.In(`
		SELECT
			id::text AS id,
			name,
			slug,
			bio,
			image_url,
			is_verified,
			monthly_listeners
		FROM artists
		WHERE id::text IN (?)
		ORDER BY monthly_listeners DESC
	`, ids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var rows []artistScanRow
	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return toArtistResults(rows), nil
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

func toTrackResults(rows []trackScanRow) []TrackResult {
	results := make([]TrackResult, 0, len(rows))
	for _, row := range rows {
		r := TrackResult{
			ID:              row.ID,
			Title:           row.Title,
			Slug:            row.Slug,
			ArtistID:        row.ArtistID,
			AlbumID:         row.AlbumID,
			CoverURL:        row.CoverURL,
			AudioURL:        row.AudioURL,
			DurationSeconds: row.DurationSeconds,
			TrackNumber:     row.TrackNumber,
			Explicit:        row.Explicit,
			PlayCount:       row.PlayCount,
			IsPublic:        row.IsPublic,
			CreatedAt:       row.CreatedAt,
		}
		if len(row.ArtistsJSON) > 0 {
			json.Unmarshal(row.ArtistsJSON, &r.Artists)
		}
		if r.Artists == nil {
			r.Artists = []TrackArtistResult{}
		}
		if len(row.GenresJSON) > 0 {
			json.Unmarshal(row.GenresJSON, &r.Genres)
		}
		if r.Genres == nil {
			r.Genres = []GenreResult{}
		}
		results = append(results, r)
	}
	return results
}

func toAlbumResults(rows []albumScanRow) []AlbumResult {
	results := make([]AlbumResult, 0, len(rows))
	for _, row := range rows {
		r := AlbumResult{
			ID:          row.ID,
			Title:       row.Title,
			Slug:        row.Slug,
			ArtistID:    row.ArtistID,
			CoverURL:    row.CoverURL,
			ReleaseDate: row.ReleaseDate,
		}
		if len(row.ArtistsJSON) > 0 {
			json.Unmarshal(row.ArtistsJSON, &r.Artists)
		}
		if r.Artists == nil {
			r.Artists = []TrackArtistResult{}
		}
		results = append(results, r)
	}
	return results
}

func toArtistResults(rows []artistScanRow) []ArtistResult {
	results := make([]ArtistResult, 0, len(rows))
	for _, row := range rows {
		r := ArtistResult{
			ID:               row.ID,
			Name:             row.Name,
			Slug:             row.Slug,
			Bio:              row.Bio,
			CoverURL:         row.ImageURL,
			IsVerified:       row.IsVerified,
			MonthlyListeners: row.MonthlyListeners,
		}
		results = append(results, r)
	}
	return results
}

func toPlaylistResults(rows []playlistScanRow) []PlaylistResult {
	results := make([]PlaylistResult, 0, len(rows))
	for _, row := range rows {
		r := PlaylistResult{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			CoverURL:    row.CoverURL,
			UserID:      row.UserID,
			IsPublic:    row.IsPublic,
		}
		results = append(results, r)
	}
	return results
}
