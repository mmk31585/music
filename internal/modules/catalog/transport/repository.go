package album

import (
	"context"
	"database/sql"
	"music/internal/modules/catalog/transport"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func parseDate(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, common.ErrInvalidInput
	}
	return &t, nil
}

func (r *Repository) Create(ctx context.Context, req CreateRequest) (*Album, error) {
	slug := common.Slugify(req.Title)
	releaseDate, err := parseDate(req.ReleaseDate)
	if err != nil {
		return nil, err
	}

	albumType := "album"
	if req.AlbumType != nil && *req.AlbumType != "" {
		albumType = *req.AlbumType
	}

	var item Album
	err = r.db.GetContext(ctx, &item, `
		INSERT INTO albums (artist_id, title, slug, cover_url, release_date, album_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, artist_id, title, slug, cover_url, release_date, album_type, created_at, updated_at
	`,
		req.ArtistID,
		req.Title,
		slug,
		req.CoverURL,
		releaseDate,
		albumType,
	)
	if err != nil {
		return nil, common.MapPGError(err)
	}
	return &item, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Album, error) {
	var item Album
	err := r.db.GetContext(ctx, &item, `
		SELECT id, artist_id, title, slug, cover_url, release_date, album_type, created_at, updated_at
		FROM albums
		WHERE id = $1
	`, id)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Album, error) {
	var items []Album
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, artist_id, title, slug, cover_url, release_date, album_type, created_at, updated_at
		FROM albums
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Album, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	title := current.Title
	if req.Title != nil {
		title = *req.Title
	}
	slug := common.Slugify(title)

	releaseDate, err := parseDate(req.ReleaseDate)
	if err != nil && req.ReleaseDate != nil {
		return nil, err
	}

	albumType := current.AlbumType
	if req.AlbumType != nil && *req.AlbumType != "" {
		albumType = *req.AlbumType
	}

	var item Album
	err = r.db.GetContext(ctx, &item, `
		UPDATE albums
		SET
			title = COALESCE($2, title),
			slug = $3,
			cover_url = COALESCE($4, cover_url),
			release_date = COALESCE($5, release_date),
			album_type = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, artist_id, title, slug, cover_url, release_date, album_type, created_at, updated_at
	`,
		id,
		req.Title,
		slug,
		req.CoverURL,
		releaseDate,
		albumType,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, common.MapPGError(err)
	}
	return &item, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM albums WHERE id = $1`, id)
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
