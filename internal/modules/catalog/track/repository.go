package track

import (
	"context"
	"database/sql"
	"music/internal/modules/catalog/common"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, req CreateRequest) (*Track, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	slug := common.Slugify(req.Title)

	explicit := false
	if req.Explicit != nil {
		explicit = *req.Explicit
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	var item Track
	err = tx.GetContext(ctx, &item, `
		INSERT INTO tracks (
			artist_id, album_id, title, slug, duration_seconds, track_number,
			explicit, audio_url, cover_url, is_public
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING
			id, artist_id, album_id, title, slug, duration_seconds, track_number,
			explicit, audio_url, cover_url, play_count, is_public, created_at, updated_at
	`,
		req.ArtistID,
		req.AlbumID,
		req.Title,
		slug,
		req.DurationSeconds,
		req.TrackNumber,
		explicit,
		req.AudioURL,
		req.CoverURL,
		isPublic,
	)
	if err != nil {
		return nil, common.MapPGError(err)
	}

	if err := r.replaceGenresTx(ctx, tx, item.ID, req.GenreIDs); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, item.ID)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Track, error) {
	var item Track
	err := r.db.GetContext(ctx, &item, `
		SELECT
			id, artist_id, album_id, title, slug, duration_seconds, track_number,
			explicit, audio_url, cover_url, play_count, is_public, created_at, updated_at
		FROM tracks
		WHERE id = $1
	`, id)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	genres, err := r.listGenresByTrack(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Genres = genres

	return &item, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int, publicOnly bool) ([]Track, error) {
	var items []Track

	query := `
		SELECT
			id, artist_id, album_id, title, slug, duration_seconds, track_number,
			explicit, audio_url, cover_url, play_count, is_public, created_at, updated_at
		FROM tracks
	`
	args := []any{}
	if publicOnly {
		query += ` WHERE is_public = TRUE`
	}
	query += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, err
	}

	for i := range items {
		genres, err := r.listGenresByTrack(ctx, items[i].ID)
		if err != nil {
			return nil, err
		}
		items[i].Genres = genres
	}

	return items, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Track, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	title := current.Title
	if req.Title != nil {
		title = *req.Title
	}
	slug := common.Slugify(title)

	var item Track
	err = tx.GetContext(ctx, &item, `
		UPDATE tracks
		SET
			album_id = COALESCE($2, album_id),
			title = COALESCE($3, title),
			slug = $4,
			duration_seconds = COALESCE($5, duration_seconds),
			track_number = COALESCE($6, track_number),
			explicit = COALESCE($7, explicit),
			audio_url = COALESCE($8, audio_url),
			cover_url = COALESCE($9, cover_url),
			is_public = COALESCE($10, is_public),
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id, artist_id, album_id, title, slug, duration_seconds, track_number,
			explicit, audio_url, cover_url, play_count, is_public, created_at, updated_at
	`,
		id,
		req.AlbumID,
		req.Title,
		slug,
		req.DurationSeconds,
		req.TrackNumber,
		req.Explicit,
		req.AudioURL,
		req.CoverURL,
		req.IsPublic,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, common.MapPGError(err)
	}

	if req.GenreIDs != nil {
		if err := r.replaceGenresTx(ctx, tx, id, req.GenreIDs); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, item.ID)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tracks WHERE id = $1`, id)
	if err != nil {
		return common.MapPGError(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return common.ErrNotFound
	}
	return nil
}

func (r *Repository) replaceGenresTx(ctx context.Context, tx *sqlx.Tx, trackID uuid.UUID, genreIDs []uuid.UUID) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM track_genres WHERE track_id = $1`, trackID); err != nil {
		return common.MapPGError(err)
	}

	for _, gid := range genreIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO track_genres (track_id, genre_id)
			VALUES ($1, $2)
		`, trackID, gid); err != nil {
			return common.MapPGError(err)
		}
	}
	return nil
}

func (r *Repository) listGenresByTrack(ctx context.Context, trackID uuid.UUID) ([]Genre, error) {
	var items []Genre
	err := r.db.SelectContext(ctx, &items, `
		SELECT g.id, g.name, g.slug, g.created_at
		FROM genres g
		INNER JOIN track_genres tg ON tg.genre_id = g.id
		WHERE tg.track_id = $1
		ORDER BY g.name ASC
	`, trackID)
	return items, err
}
