package library

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

const favoritesPlaylistName = "Liked Songs"

var (
	ErrFavoritesNotFound = errors.New("favorites playlist not found")
)

// EnsureFavoritesPlaylist returns the ID of the user's "Liked Songs" playlist,
// creating it if it doesn't exist yet.
func (r *Repository) EnsureFavoritesPlaylist(ctx context.Context, userID string) (string, error) {
	// Try to find existing
	var id string
	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM playlists WHERE user_id = $1 AND name = $2 LIMIT 1`,
		userID, favoritesPlaylistName,
	).Scan(&id)
	if err == nil {
		return id, nil
	}

	// Create one
	newID := uuid.New().String()
	now := time.Now()
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO playlists (id, user_id, name, description, cover_url, is_public, is_collaborative, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		newID, userID, favoritesPlaylistName,
		"Auto-created playlist for liked songs", nil,
		false, false, now, now,
	)
	if err != nil {
		return "", err
	}
	return newID, nil
}

// AddTrackToFavorites adds a track to the user's Liked Songs playlist.
func (r *Repository) AddTrackToFavorites(ctx context.Context, userID, trackID string) error {
	playlistID, err := r.EnsureFavoritesPlaylist(ctx, userID)
	if err != nil {
		return err
	}

	// Get next position
	var maxPos int
	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(MAX(position), 0) FROM playlist_tracks WHERE playlist_id = $1`,
		playlistID,
	).Scan(&maxPos)
	if err != nil {
		return err
	}

	// Insert track (ignore if already exists)
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO playlist_tracks (playlist_id, track_id, position, added_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (playlist_id, track_id) DO NOTHING`,
		playlistID, trackID, maxPos+1,
	)
	return err
}

// RemoveTrackFromFavorites removes a track from the user's Liked Songs playlist.
func (r *Repository) RemoveTrackFromFavorites(ctx context.Context, userID, trackID string) error {
	playlistID, err := r.EnsureFavoritesPlaylist(ctx, userID)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx,
		`DELETE FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`,
		playlistID, trackID,
	)
	return err
}
