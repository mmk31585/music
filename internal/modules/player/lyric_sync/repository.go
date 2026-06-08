package lyrics

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrLyricsNotFound      = errors.New("lyrics not found")
	ErrLyricsAlreadyExists = errors.New("lyrics already exists for this track and language")
	ErrTrackNotFound       = errors.New("track not found")
	ErrInvalidLyricsType   = errors.New("invalid lyrics type")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLyrics(ctx context.Context, req CreateLyricsRequest) (Lyrics, error) {
	// Validate track exists
	var trackExists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM tracks WHERE id = $1)`,
		req.TrackID,
	).Scan(&trackExists)
	if err != nil {
		return Lyrics{}, fmt.Errorf("failed to check track existence: %w", err)
	}
	if !trackExists {
		return Lyrics{}, ErrTrackNotFound
	}

	// Check if lyrics already exists for this track and language
	var exists bool
	err = r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM lyrics WHERE track_id = $1 AND language = $2)`,
		req.TrackID, req.Language,
	).Scan(&exists)
	if err != nil {
		return Lyrics{}, fmt.Errorf("failed to check existing lyrics: %w", err)
	}
	if exists {
		return Lyrics{}, ErrLyricsAlreadyExists
	}

	// Validate lyrics type
	if req.Type != "plain" && req.Type != "lrc" {
		return Lyrics{}, ErrInvalidLyricsType
	}

	query := `
		INSERT INTO lyrics (track_id, language, type, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, track_id, language, type, content, created_at, updated_at
	`

	var l Lyrics
	err = r.db.QueryRowContext(ctx, query,
		req.TrackID,
		req.Language,
		req.Type,
		req.Content,
	).Scan(
		&l.ID,
		&l.TrackID,
		&l.Language,
		&l.Type,
		&l.Content,
		&l.CreatedAt,
		&l.UpdatedAt,
	)

	if err != nil {
		return Lyrics{}, fmt.Errorf("failed to create lyrics: %w", err)
	}

	return l, nil
}

func (r *Repository) UpdateLyrics(ctx context.Context, lyricsID string, req UpdateLyricsRequest) (Lyrics, error) {
	query := `
		UPDATE lyrics
		SET
			language = COALESCE($1, language),
			type = COALESCE($2, type),
			content = COALESCE($3, content),
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, track_id, language, type, content, created_at, updated_at
	`

	var l Lyrics
	err := r.db.QueryRowContext(ctx, query,
		req.Language,
		req.Type,
		req.Content,
		lyricsID,
	).Scan(
		&l.ID,
		&l.TrackID,
		&l.Language,
		&l.Type,
		&l.Content,
		&l.CreatedAt,
		&l.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Lyrics{}, ErrLyricsNotFound
	}

	if err != nil {
		return Lyrics{}, fmt.Errorf("failed to update lyrics: %w", err)
	}

	return l, nil
}

func (r *Repository) DeleteLyrics(ctx context.Context, lyricsID string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM lyrics WHERE id = $1`,
		lyricsID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete lyrics: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrLyricsNotFound
	}

	return nil
}

func (r *Repository) GetLyricsByTrackID(ctx context.Context, trackID string) ([]Lyrics, error) {
	query := `
		SELECT id, track_id, language, type, content, created_at, updated_at
		FROM lyrics
		WHERE track_id = $1
		ORDER BY language ASC
	`

	rows, err := r.db.QueryContext(ctx, query, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lyrics by track id: %w", err)
	}
	defer rows.Close()

	lyricsList := make([]Lyrics, 0)
	for rows.Next() {
		var l Lyrics
		if err := rows.Scan(
			&l.ID,
			&l.TrackID,
			&l.Language,
			&l.Type,
			&l.Content,
			&l.CreatedAt,
			&l.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan lyrics: %w", err)
		}
		lyricsList = append(lyricsList, l)
	}

	return lyricsList, nil
}

func (r *Repository) GetLyricsByTrackAndLanguage(ctx context.Context, trackID, language string) (Lyrics, error) {
	query := `
		SELECT id, track_id, language, type, content, created_at, updated_at
		FROM lyrics
		WHERE track_id = $1 AND language = $2
	`

	var l Lyrics
	err := r.db.QueryRowContext(ctx, query, trackID, language).Scan(
		&l.ID,
		&l.TrackID,
		&l.Language,
		&l.Type,
		&l.Content,
		&l.CreatedAt,
		&l.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Lyrics{}, ErrLyricsNotFound
	}

	if err != nil {
		return Lyrics{}, fmt.Errorf("failed to get lyrics by track and language: %w", err)
	}

	return l, nil
}

func (r *Repository) LyricsExists(ctx context.Context, trackID, language string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM lyrics WHERE track_id = $1 AND language = $2)`,
		trackID, language,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check lyrics existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) DebugError(err error) error {
	return fmt.Errorf("lyrics repository error: %w", err)
}
