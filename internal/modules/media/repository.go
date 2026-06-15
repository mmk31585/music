package media

import (
	"context"
	"database/sql"

	catalogCommon "music/internal/modules/catalog/common"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func normalizeStorageProvider(v string) string {
	if v == "" {
		return "local"
	}

	return v
}

func normalizeMetadata(v string) string {
	if v == "" {
		return "{}"
	}

	return v
}

func (r *Repository) Create(ctx context.Context, req CreateMediaRequest) (*Media, error) {
	var item Media

	err := r.db.GetContext(ctx, &item, `
		INSERT INTO media (
			media_type,
			storage_provider,
			bucket,
			object_key,
			public_url,
			mime_type,
			file_size,
			checksum_sha256,
			duration_seconds,
			width,
			height,
			original_filename,
			metadata,
			created_by
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13::jsonb,
			$14
		)
		RETURNING
			id,
			media_type,
			storage_provider,
			bucket,
			object_key,
			public_url,
			mime_type,
			file_size,
			checksum_sha256,
			duration_seconds,
			width,
			height,
			original_filename,
			metadata::text AS metadata,
			created_by,
			created_at,
			updated_at
	`,
		req.MediaType,
		normalizeStorageProvider(req.StorageProvider),
		req.Bucket,
		req.ObjectKey,
		req.PublicURL,
		req.MimeType,
		req.FileSize,
		req.ChecksumSHA256,
		req.DurationSeconds,
		req.Width,
		req.Height,
		req.OriginalFilename,
		normalizeMetadata(req.Metadata),
		req.CreatedBy,
	)
	if err != nil {
		return nil, catalogCommon.MapPGError(err)
	}

	return &item, nil
}

func (r *Repository) FindByChecksum(ctx context.Context, checksum string) (*Media, error) {
	var item Media

	err := r.db.GetContext(ctx, &item, `
		SELECT
			id,
			media_type,
			storage_provider,
			bucket,
			object_key,
			public_url,
			mime_type,
			file_size,
			checksum_sha256,
			duration_seconds,
			width,
			height,
			original_filename,
			metadata::text AS metadata,
			created_by,
			created_at,
			updated_at
		FROM media
		WHERE checksum_sha256 = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, checksum)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) List(ctx context.Context) ([]Media, error) {
	var items []Media

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			id,
			media_type,
			storage_provider,
			bucket,
			object_key,
			public_url,
			mime_type,
			file_size,
			checksum_sha256,
			duration_seconds,
			width,
			height,
			original_filename,
			metadata::text AS metadata,
			created_by,
			created_at,
			updated_at
		FROM media
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM media WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return catalogCommon.ErrNotFound
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Media, error) {
	var item Media

	err := r.db.GetContext(ctx, &item, `
		SELECT
			id,
			media_type,
			storage_provider,
			bucket,
			object_key,
			public_url,
			mime_type,
			file_size,
			checksum_sha256,
			duration_seconds,
			width,
			height,
			original_filename,
			metadata::text AS metadata,
			created_by,
			created_at,
			updated_at
		FROM media
		WHERE id = $1
	`, id)

	if err == sql.ErrNoRows {
		return nil, catalogCommon.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}
