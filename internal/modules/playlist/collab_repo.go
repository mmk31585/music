package playlist

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotCollaborator     = errors.New("user is not a collaborator")
	ErrAlreadyCollaborator = errors.New("user is already a collaborator")
	ErrCannotRemoveOwner   = errors.New("cannot remove playlist owner as collaborator")
	ErrPlaylistNotCollab   = errors.New("playlist is not collaborative")
)

func (r *repository) IsCollaborator(ctx context.Context, playlistID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM playlist_collaborators WHERE playlist_id = $1 AND user_id = $2)`,
		playlistID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *repository) IsCollaborativePlaylist(ctx context.Context, playlistID string) (bool, error) {
	var collab bool
	err := r.db.QueryRowContext(ctx,
		`SELECT is_collaborative FROM playlists WHERE id = $1`,
		playlistID,
	).Scan(&collab)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrPlaylistNotFound
	}
	return collab, err
}

func (r *repository) SetCollaborative(ctx context.Context, playlistID string, collab bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE playlists SET is_collaborative = $1, updated_at = NOW() WHERE id = $2`,
		collab, playlistID,
	)
	return err
}

func (r *repository) AddCollaborator(ctx context.Context, playlistID, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO playlist_collaborators (playlist_id, user_id, added_by) VALUES ($1, $2, $2)`,
		playlistID, userID,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrAlreadyCollaborator
		}
		return err
	}
	return nil
}

func (r *repository) RemoveCollaborator(ctx context.Context, playlistID, userID string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM playlist_collaborators WHERE playlist_id = $1 AND user_id = $2`,
		playlistID, userID,
	)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotCollaborator
	}
	return nil
}

func (r *repository) ListCollaborators(ctx context.Context, playlistID string) ([]CollaboratorResponse, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT pc.user_id, pc.created_at::text, false AS is_creator
		FROM playlist_collaborators pc
		WHERE pc.playlist_id = $1
		ORDER BY pc.created_at ASC
	`, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]CollaboratorResponse, 0)
	for rows.Next() {
		var item CollaboratorResponse
		if err := rows.Scan(&item.UserID, &item.AddedAt, &item.IsCreator); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
