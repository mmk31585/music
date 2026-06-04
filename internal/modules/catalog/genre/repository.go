package genre

import (
	"context"
	"database/sql"
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

func (r *Repository) Create(ctx context.Context, req CreateRequest) (*Genre, error) {
	slug := common.Slugify(req.Name)

	var item Genre
	err := r.db.GetContext(ctx, &item, `
		INSERT INTO genres (name, slug)
		VALUES ($1, $2)
		RETURNING id, name, slug, created_at
	`, req.Name, slug)
	if err != nil {
		return nil, common.MapPGError(err)
	}
	return &item, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Genre, error) {
	var item Genre
	err := r.db.GetContext(ctx, &item, `
		SELECT id, name, slug, created_at
		FROM genres
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

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Genre, error) {
	var items []Genre
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, name, slug, created_at
		FROM genres
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Genre, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := current.Name
	if req.Name != nil {
		name = *req.Name
	}
	slug := common.Slugify(name)

	var item Genre
	err = r.db.GetContext(ctx, &item, `
		UPDATE genres
		SET
			name = COALESCE($2, name),
			slug = $3
		WHERE id = $1
		RETURNING id, name, slug, created_at
	`, id, req.Name, slug)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, common.MapPGError(err)
	}
	return &item, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM genres WHERE id = $1`, id)
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
