package artist

import (
	"context"
	"database/sql"
	"fmt"
	"music/internal/modules/catalog/transport"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, req CreateRequest) (*Artist, error) {
	slug := common.Slugify(req.Name)

	var item Artist
	err := r.db.GetContext(ctx, &item, `
		INSERT INTO artist (
			name, slug, bio, image_url, is_verified, monthly_listeners
		)
		VALUES ($1, $2, $3, $4, COALESCE($5, FALSE), COALESCE($6, 0))
		RETURNING id, name, slug, bio, image_url, is_verified, monthly_listeners, created_at, updated_at
	`,
		req.Name,
		slug,
		req.Bio,
		req.ImageURL,
		req.IsVerified,
		req.MonthlyListeners,
	)
	if err != nil {
		return nil, common.MapPGError(err)
	}
	return &item, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Artist, error) {
	var item Artist
	err := r.db.GetContext(ctx, &item, `
		SELECT id, name, slug, bio, image_url, is_verified, monthly_listeners, created_at, updated_at
		FROM artist
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

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Artist, error) {
	var items []Artist
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, name, slug, bio, image_url, is_verified, monthly_listeners, created_at, updated_at
		FROM artist
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Artist, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := current.Name
	if req.Name != nil {
		name = *req.Name
	}
	slug := common.Slugify(name)

	var item Artist
	err = r.db.GetContext(ctx, &item, `
		UPDATE artist
		SET
			name = COALESCE($2, name),
			slug = $3,
			bio = COALESCE($4, bio),
			image_url = COALESCE($5, image_url),
			is_verified = COALESCE($6, is_verified),
			monthly_listeners = COALESCE($7, monthly_listeners),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, slug, bio, image_url, is_verified, monthly_listeners, created_at, updated_at
	`,
		id,
		req.Name,
		slug,
		req.Bio,
		req.ImageURL,
		req.IsVerified,
		req.MonthlyListeners,
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
	res, err := r.db.ExecContext(ctx, `DELETE FROM artist WHERE id = $1`, id)
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

func (r *Repository) Search(ctx context.Context, q string, limit, offset int) ([]Artist, error) {
	var items []Artist
	pattern := fmt.Sprintf("%%%s%%", q)
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, name, slug, bio, image_url, is_verified, monthly_listeners, created_at, updated_at
		FROM artist
		WHERE name ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, pattern, limit, offset)
	return items, err
}
