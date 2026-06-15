package playlist

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	ErrPlaylistNotFound       = errors.New("playlist not found")
	ErrPlaylistTrackNotFound  = errors.New("playlist track not found")
	ErrTrackAlreadyInPlaylist = errors.New("track already in playlist")
	ErrTrackNotFound          = errors.New("track not found")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePlaylist(ctx context.Context, req CreatePlaylistRequest, userID int64) (Playlist, error) {
	query := `
		INSERT INTO playlists (user_id, name, description, cover_url, is_public)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, description, cover_url, is_public, created_at, updated_at
	`

	var p Playlist
	err := r.db.QueryRowContext(ctx, query,
		userID,
		req.Name,
		req.Description,
		req.CoverURL,
		req.IsPublic,
	).Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.Description,
		&p.CoverURL,
		&p.IsPublic,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	return p, err
}

func (r *Repository) UpdatePlaylist(ctx context.Context, playlistID, userID int64, req UpdatePlaylistRequest) (Playlist, error) {
	query := `
		UPDATE playlists
		SET
			name = $1,
			description = $2,
			cover_url = $3,
			is_public = $4,
			updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING id, user_id, name, description, cover_url, is_public, created_at, updated_at
	`

	var p Playlist
	err := r.db.QueryRowContext(ctx, query,
		req.Name,
		req.Description,
		req.CoverURL,
		req.IsPublic,
		playlistID,
		userID,
	).Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.Description,
		&p.CoverURL,
		&p.IsPublic,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Playlist{}, ErrPlaylistNotFound
	}

	return p, err
}

func (r *Repository) DeletePlaylist(ctx context.Context, playlistID, userID int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM playlists WHERE id = $1 AND user_id = $2`,
		playlistID, userID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrPlaylistNotFound
	}

	return nil
}

func (r *Repository) GetPlaylistByID(ctx context.Context, playlistID int64) (Playlist, error) {
	query := `
		SELECT id, user_id, name, description, cover_url, is_public, created_at, updated_at
		FROM playlists
		WHERE id = $1
	`

	var p Playlist
	err := r.db.QueryRowContext(ctx, query, playlistID).Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.Description,
		&p.CoverURL,
		&p.IsPublic,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Playlist{}, ErrPlaylistNotFound
	}

	return p, err
}

func (r *Repository) ListPlaylistTracks(ctx context.Context, playlistID int64) ([]PlaylistTrackItem, error) {
	query := `
		SELECT
			pt.id,
			pt.track_id,
			pt.position,
			t.title,
			ar.name,
			al.title,
			t.cover_url,
			t.audio_url,
			t.duration_seconds
		FROM playlist_tracks pt
		JOIN tracks t ON t.id = pt.track_id
		LEFT JOIN artists ar ON ar.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		WHERE pt.playlist_id = $1
		ORDER BY pt.position ASC
	`

	rows, err := r.db.QueryContext(ctx, query, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]PlaylistTrackItem, 0)
	for rows.Next() {
		var item PlaylistTrackItem
		if err := rows.Scan(
			&item.PlaylistTrackID,
			&item.TrackID,
			&item.Position,
			&item.Title,
			&item.ArtistName,
			&item.AlbumTitle,
			&item.CoverURL,
			&item.AudioURL,
			&item.DurationSeconds,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Repository) ListPublicPlaylists(ctx context.Context) ([]PlaylistListItemResponse, error) {
	query := `
		SELECT
			p.id,
			p.user_id,
			p.name,
			p.description,
			p.cover_url,
			p.is_public,
			COUNT(pt.id) AS track_count,
			p.created_at,
			p.updated_at
		FROM playlists p
		LEFT JOIN playlist_tracks pt ON pt.playlist_id = p.id
		WHERE p.is_public = TRUE
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]PlaylistListItemResponse, 0)
	for rows.Next() {
		var item PlaylistListItemResponse
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Description,
			&item.CoverURL,
			&item.IsPublic,
			&item.TrackCount,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Repository) ListUserPlaylists(ctx context.Context, userID int64) ([]PlaylistListItemResponse, error) {
	query := `
		SELECT
			p.id,
			p.user_id,
			p.name,
			p.description,
			p.cover_url,
			p.is_public,
			COUNT(pt.id) AS track_count,
			p.created_at,
			p.updated_at
		FROM playlists p
		LEFT JOIN playlist_tracks pt ON pt.playlist_id = p.id
		WHERE p.user_id = $1
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]PlaylistListItemResponse, 0)
	for rows.Next() {
		var item PlaylistListItemResponse
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Description,
			&item.CoverURL,
			&item.IsPublic,
			&item.TrackCount,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Repository) TrackExists(ctx context.Context, trackID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM tracks WHERE id = $1)`,
		trackID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) GetPlaylistTrackCount(ctx context.Context, playlistID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM playlist_tracks WHERE playlist_id = $1`,
		playlistID,
	).Scan(&count)
	return count, err
}

func (r *Repository) GetTrackPosition(ctx context.Context, playlistID, trackID int64) (int, error) {
	var pos int
	err := r.db.QueryRowContext(ctx,
		`SELECT position FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`,
		playlistID, trackID,
	).Scan(&pos)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrPlaylistTrackNotFound
	}

	return pos, err
}

func (r *Repository) AddTrack(ctx context.Context, playlistID, trackID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)

	var exists bool
	if err := tx.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM tracks WHERE id = $1)`,
		trackID,
	).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrTrackNotFound
	}

	var position int
	if err := tx.QueryRowContext(ctx,
		`SELECT COALESCE(MAX(position), 0) + 1 FROM playlist_tracks WHERE playlist_id = $1`,
		playlistID,
	).Scan(&position); err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES ($1, $2, $3)`,
		playlistID, trackID, position,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrTrackAlreadyInPlaylist
		}
		return err
	}

	return tx.Commit()
}

func (r *Repository) RemoveTrack(ctx context.Context, playlistID, trackID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)

	var oldPosition int
	err = tx.QueryRowContext(ctx,
		`SELECT position FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`,
		playlistID, trackID,
	).Scan(&oldPosition)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPlaylistTrackNotFound
	}
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx,
		`DELETE FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`,
		playlistID, trackID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrPlaylistTrackNotFound
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE playlist_tracks SET position = position - 1 WHERE playlist_id = $1 AND position > $2`,
		playlistID, oldPosition,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) ReorderTrack(ctx context.Context, playlistID, trackID int64, newPosition int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)

	var oldPosition int
	err = tx.QueryRowContext(ctx,
		`SELECT position FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`,
		playlistID, trackID,
	).Scan(&oldPosition)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPlaylistTrackNotFound
	}
	if err != nil {
		return err
	}

	var total int
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM playlist_tracks WHERE playlist_id = $1`,
		playlistID,
	).Scan(&total); err != nil {
		return err
	}

	if newPosition < 1 || newPosition > total {
		return ErrInvalidTrackPosition
	}

	if newPosition == oldPosition {
		return tx.Commit()
	}

	if newPosition < oldPosition {
		_, err = tx.ExecContext(ctx, `
			UPDATE playlist_tracks
			SET position = position + 1
			WHERE playlist_id = $1
			  AND position >= $2
			  AND position < $3
		`, playlistID, newPosition, oldPosition)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.ExecContext(ctx, `
			UPDATE playlist_tracks
			SET position = position - 1
			WHERE playlist_id = $1
			  AND position <= $2
			  AND position > $3
		`, playlistID, newPosition, oldPosition)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE playlist_tracks
		SET position = $1
		WHERE playlist_id = $2 AND track_id = $3
	`, newPosition, playlistID, trackID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func rollback(tx *sql.Tx) {
	_ = tx.Rollback()
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}

func (r *Repository) DebugError(err error) error {
	return fmt.Errorf("playlist repository error: %w", err)
}
