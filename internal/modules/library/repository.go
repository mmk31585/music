package library

import (
	"context"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	ErrTrackNotFound          = errors.New("track not found")
	ErrAlbumNotFound          = errors.New("album not found")
	ErrArtistNotFound         = errors.New("artist not found")
	ErrTrackAlreadyLiked      = errors.New("track already liked")
	ErrAlbumAlreadyLiked      = errors.New("album already liked")
	ErrArtistAlreadyFollowed  = errors.New("artist already followed")
	ErrLikedTrackNotFound     = errors.New("liked track not found")
	ErrLikedAlbumNotFound     = errors.New("liked album not found")
	ErrFollowedArtistNotFound = errors.New("followed artist not found")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) TrackExists(ctx context.Context, trackID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM tracks WHERE id = $1)`,
		trackID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) AlbumExists(ctx context.Context, albumID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM albums WHERE id = $1)`,
		albumID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) ArtistExists(ctx context.Context, artistID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM artists WHERE id = $1)`,
		artistID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) LikeTrack(ctx context.Context, userID, trackID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO liked_tracks (user_id, track_id) VALUES ($1, $2)`,
		userID, trackID,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrTrackAlreadyLiked
		}
		return err
	}
	return nil
}

func (r *Repository) UnlikeTrack(ctx context.Context, userID, trackID int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM liked_tracks WHERE user_id = $1 AND track_id = $2`,
		userID, trackID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrLikedTrackNotFound
	}
	return nil
}

func (r *Repository) LikeAlbum(ctx context.Context, userID, albumID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO liked_albums (user_id, album_id) VALUES ($1, $2)`,
		userID, albumID,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrAlbumAlreadyLiked
		}
		return err
	}
	return nil
}

func (r *Repository) UnlikeAlbum(ctx context.Context, userID, albumID int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM liked_albums WHERE user_id = $1 AND album_id = $2`,
		userID, albumID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrLikedAlbumNotFound
	}
	return nil
}

func (r *Repository) FollowArtist(ctx context.Context, userID, artistID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO followed_artists (user_id, artist_id) VALUES ($1, $2)`,
		userID, artistID,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrArtistAlreadyFollowed
		}
		return err
	}
	return nil
}

func (r *Repository) UnfollowArtist(ctx context.Context, userID, artistID int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM followed_artists WHERE user_id = $1 AND artist_id = $2`,
		userID, artistID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrFollowedArtistNotFound
	}
	return nil
}

func (r *Repository) AddPlayHistory(ctx context.Context, userID, trackID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO play_history (user_id, track_id) VALUES ($1, $2)`,
		userID, trackID,
	)
	return err
}

func (r *Repository) ListLikedTracks(ctx context.Context, userID int64) ([]LibraryTrackItem, error) {
	query := `
		SELECT
			t.id,
			t.title,
			ar.id,
			ar.name,
			al.id,
			al.title,
			t.cover_url,
			t.audio_url,
			t.duration_seconds,
			lt.created_at,
			NULL
		FROM liked_tracks lt
		JOIN tracks t ON t.id = lt.track_id
		LEFT JOIN artist ar ON ar.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		WHERE lt.user_id = $1
		ORDER BY lt.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]LibraryTrackItem, 0)
	for rows.Next() {
		var item LibraryTrackItem
		if err := rows.Scan(
			&item.TrackID,
			&item.Title,
			&item.ArtistID,
			&item.ArtistName,
			&item.AlbumID,
			&item.AlbumTitle,
			&item.CoverURL,
			&item.AudioURL,
			&item.DurationSeconds,
			&item.AddedAt,
			&item.LastPlayedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Repository) ListLikedAlbums(ctx context.Context, userID int64) ([]LibraryAlbumItem, error) {
	query := `
		SELECT
			al.id,
			al.title,
			ar.id,
			ar.name,
			al.cover_url,
			al.release_date,
			la.created_at
		FROM liked_albums la
		JOIN albums al ON al.id = la.album_id
		LEFT JOIN artist ar ON ar.id = al.artist_id
		WHERE la.user_id = $1
		ORDER BY la.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]LibraryAlbumItem, 0)
	for rows.Next() {
		var item LibraryAlbumItem
		if err := rows.Scan(
			&item.AlbumID,
			&item.Title,
			&item.ArtistID,
			&item.ArtistName,
			&item.CoverURL,
			&item.ReleaseDate,
			&item.AddedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Repository) ListFollowedArtists(ctx context.Context, userID int64) ([]LibraryArtistItem, error) {
	query := `
		SELECT
			ar.id,
			ar.name,
			ar.cover_url,
			fa.created_at
		FROM followed_artists fa
		JOIN artist ar ON ar.id = fa.artist_id
		WHERE fa.user_id = $1
		ORDER BY fa.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]LibraryArtistItem, 0)
	for rows.Next() {
		var item LibraryArtistItem
		if err := rows.Scan(
			&item.ArtistID,
			&item.Name,
			&item.CoverURL,
			&item.FollowedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Repository) ListPlayHistory(ctx context.Context, userID int64, limit int) ([]LibraryTrackItem, error) {
	query := `
		SELECT
			t.id,
			t.title,
			ar.id,
			ar.name,
			al.id,
			al.title,
			t.cover_url,
			t.audio_url,
			t.duration_seconds,
			ph.played_at,
			ph.played_at
		FROM play_history ph
		JOIN tracks t ON t.id = ph.track_id
		LEFT JOIN artist ar ON ar.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		WHERE ph.user_id = $1
		ORDER BY ph.played_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]LibraryTrackItem, 0)
	for rows.Next() {
		var item LibraryTrackItem
		if err := rows.Scan(
			&item.TrackID,
			&item.Title,
			&item.ArtistID,
			&item.ArtistName,
			&item.AlbumID,
			&item.AlbumTitle,
			&item.CoverURL,
			&item.AudioURL,
			&item.DurationSeconds,
			&item.AddedAt,
			&item.LastPlayedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Repository) ListRecentlyPlayed(ctx context.Context, userID int64, limit int) ([]LibraryTrackItem, error) {
	query := `
		SELECT
			t.id,
			t.title,
			ar.id,
			ar.name,
			al.id,
			al.title,
			t.cover_url,
			t.audio_url,
			t.duration_seconds,
			MAX(ph.played_at) AS added_at,
			MAX(ph.played_at) AS last_played_at
		FROM play_history ph
		JOIN tracks t ON t.id = ph.track_id
		LEFT JOIN artist ar ON ar.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		WHERE ph.user_id = $1
		GROUP BY
			t.id, t.title,
			ar.id, ar.name,
			al.id, al.title,
			t.cover_url, t.audio_url, t.duration_seconds
		ORDER BY MAX(ph.played_at) DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]LibraryTrackItem, 0)
	for rows.Next() {
		var item LibraryTrackItem
		if err := rows.Scan(
			&item.TrackID,
			&item.Title,
			&item.ArtistID,
			&item.ArtistName,
			&item.AlbumID,
			&item.AlbumTitle,
			&item.CoverURL,
			&item.AudioURL,
			&item.DurationSeconds,
			&item.AddedAt,
			&item.LastPlayedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}
