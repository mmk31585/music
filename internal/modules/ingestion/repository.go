package ingestion

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	catalogCommon "music/internal/modules/catalog/common"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type CreateDraftParams struct {
	ID                string
	UploadedBy        string
	OriginalFilename  string
	FilePath          string
	FileSize          int64
	DurationSeconds   *float64
	Bitrate           *int
	Format            string
	ExtractedMetadata string
}

func (r *Repository) CreateDraft(ctx context.Context, params CreateDraftParams) (*IngestionDraft, error) {
	var draft IngestionDraft

	err := r.db.GetContext(ctx, &draft, `
		INSERT INTO ingestion_drafts (
			id, uploaded_by, original_filename, file_path, file_size,
			duration_seconds, bitrate, format, extracted_metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb)
		RETURNING
			id, uploaded_by, original_filename, file_path, file_size,
			duration_seconds, bitrate, format, status,
			extracted_metadata::text AS extracted_metadata,
			enriched_metadata::text AS enriched_metadata,
			final_metadata::text AS final_metadata,
			created_at, updated_at
	`,
		params.ID,
		params.UploadedBy,
		params.OriginalFilename,
		params.FilePath,
		params.FileSize,
		params.DurationSeconds,
		params.Bitrate,
		params.Format,
		params.ExtractedMetadata,
	)
	if err != nil {
		return nil, catalogCommon.MapPGError(err)
	}

	return &draft, nil
}

type CreateAssetParams struct {
	ID        string
	DraftID   string
	AssetType string
	URL       string
	Source    string
}

func (r *Repository) CreateAsset(ctx context.Context, params CreateAssetParams) (*IngestionDraftAsset, error) {
	var asset IngestionDraftAsset

	err := r.db.GetContext(ctx, &asset, `
		INSERT INTO ingestion_draft_assets (id, draft_id, asset_type, url, source)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, draft_id, asset_type, url, source, created_at
	`,
		params.ID,
		params.DraftID,
		params.AssetType,
		params.URL,
		params.Source,
	)
	if err != nil {
		return nil, catalogCommon.MapPGError(err)
	}

	return &asset, nil
}

func (r *Repository) GetDraftByID(ctx context.Context, id string) (*IngestionDraft, error) {
	var draft IngestionDraft

	err := r.db.GetContext(ctx, &draft, `
		SELECT
			id, uploaded_by, original_filename, file_path, file_size,
			duration_seconds, bitrate, format, status,
			extracted_metadata::text AS extracted_metadata,
			enriched_metadata::text AS enriched_metadata,
			final_metadata::text AS final_metadata,
			created_at, updated_at
		FROM ingestion_drafts
		WHERE id = $1
	`, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &draft, nil
}

func (r *Repository) GetAssetsByDraftID(ctx context.Context, draftID string) ([]IngestionDraftAsset, error) {
	var assets []IngestionDraftAsset

	err := r.db.SelectContext(ctx, &assets, `
		SELECT id, draft_id, asset_type, url, source, created_at
		FROM ingestion_draft_assets
		WHERE draft_id = $1
		ORDER BY created_at ASC
	`, draftID)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

type ListDraftsParams struct {
	Status string
	Limit  int
	Offset int
}

type ListDraftsRow struct {
	ID               string     `db:"id"`
	OriginalFilename string     `db:"original_filename"`
	FileSize         int64      `db:"file_size"`
	Format           string     `db:"format"`
	DurationSeconds  *float64   `db:"duration_seconds"`
	Status           DraftStatus `db:"status"`
	ExtractedMetadata string   `db:"extracted_metadata"`
	CreatedAt        time.Time  `db:"created_at"`
	CoverArtURL      *string    `db:"cover_art_url"`
}

func (r *Repository) ListDrafts(ctx context.Context, params ListDraftsParams) ([]ListDraftsRow, int, error) {
	countWhere := ""
	where := ""
	args := []interface{}{}

	if params.Status != "" {
		countWhere = " WHERE status = $1"
		where = " WHERE d.status = $1"
		args = append(args, params.Status)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM ingestion_drafts" + countWhere
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			d.id, d.original_filename, d.file_size, d.format,
			d.duration_seconds, d.status,
			d.extracted_metadata::text AS extracted_metadata,
			d.created_at,
			a.url AS cover_art_url
		FROM ingestion_drafts d
		LEFT JOIN ingestion_draft_assets a ON a.draft_id = d.id AND a.asset_type = 'cover'
	` + where + `
		ORDER BY d.created_at DESC
		LIMIT $` + fmt.Sprintf("%d", len(args)+1) + `
		OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	args = append(args, params.Limit, params.Offset)

	var rows []ListDraftsRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *Repository) UpdateDraftStatus(ctx context.Context, id string, status DraftStatus) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ingestion_drafts SET status = $1, updated_at = NOW() WHERE id = $2
	`, string(status), id)
	return err
}

func (r *Repository) UpdateDraftEnrichedMetadata(ctx context.Context, id string, enrichedJSON string, status DraftStatus) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ingestion_drafts SET enriched_metadata = $1::jsonb, status = $2, updated_at = NOW() WHERE id = $3
	`, enrichedJSON, string(status), id)
	return err
}

func (r *Repository) UpdateDraftFinalMetadata(ctx context.Context, id string, finalJSON string, status DraftStatus) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ingestion_drafts SET final_metadata = $1::jsonb, status = $2, updated_at = NOW() WHERE id = $3
	`, finalJSON, string(status), id)
	return err
}

func (r *Repository) SearchArtists(ctx context.Context, q string, limit int) ([]ArtistSearchResult, error) {
	pattern := "%" + q + "%"
	var results []ArtistSearchResult

	err := r.db.SelectContext(ctx, &results, `
		SELECT id, name, slug, COALESCE(bio, '') AS bio,
			COALESCE(image_url, '') AS image_url,
			COALESCE(country, '') AS country
		FROM artists
		WHERE name ILIKE $1
		ORDER BY monthly_listeners DESC
		LIMIT $2
	`, pattern, limit)

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *Repository) SearchAlbums(ctx context.Context, q string, limit int) ([]AlbumSearchResult, error) {
	pattern := "%" + q + "%"
	var results []AlbumSearchResult

	err := r.db.SelectContext(ctx, &results, `
		SELECT a.id, a.title, a.slug,
			COALESCE(ar.name, '') AS artist_name,
			COALESCE(a.cover_url, '') AS cover_url,
			EXTRACT(YEAR FROM a.release_date)::int AS release_year
		FROM albums a
		LEFT JOIN album_artists aa ON aa.album_id = a.id AND aa.role = 'primary'
		LEFT JOIN artists ar ON ar.id = aa.artist_id
		WHERE a.title ILIKE $1
		ORDER BY a.release_date DESC NULLS LAST
		LIMIT $2
	`, pattern, limit)

	if err != nil {
		return nil, err
	}
	return results, nil
}
