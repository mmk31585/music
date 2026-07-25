package upload

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	apperrors "music/internal/common/errors"

	"github.com/jmoiron/sqlx"
)

type RepositoryInterface interface {
	// Draft management
	CreateDraft(ctx context.Context, draft Draft) (Draft, error)
	GetDraft(ctx context.Context, id string) (Draft, error)
	UpdateDraftStatus(ctx context.Context, id string, status string, reviewNotes *string, reviewedBy *string) error
	ListPendingDrafts(ctx context.Context, source *string, limit, offset int) ([]Draft, int, error)
	ListUserDrafts(ctx context.Context, userID string, limit, offset int) ([]Draft, int, error)

	// Co-uploaders
	AddCoUploader(ctx context.Context, co CoUploader) (CoUploader, error)
	RemoveCoUploader(ctx context.Context, trackID int64, userID string) error
	GetCoUploaders(ctx context.Context, trackID int64) ([]CoUploader, error)
	GetUserCoUploads(ctx context.Context, userID string) ([]CoUploader, error)

	// File attachment
	AttachFile(ctx context.Context, id, filePath, fileURL string, fileSize int64) error

	// Upload slots
	GetUploadSlots(ctx context.Context, userID string) (UploadSlot, error)
	IncrementSlots(ctx context.Context, userID string) error
	DecrementSlots(ctx context.Context, userID string) error
	UpdateMaxSlots(ctx context.Context, userID string, maxSlots int) error
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateDraft(ctx context.Context, draft Draft) (Draft, error) {
	query := `
		INSERT INTO ingestion_drafts (id, uploaded_by, original_filename, file_path, file_size,
		    duration_seconds, bitrate, format, status, extracted_metadata, upload_source, club_id, needs_review)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, uploaded_by, original_filename, file_path, file_size, duration_seconds, bitrate,
		    format, status, extracted_metadata, upload_source, club_id, needs_review, created_at, updated_at
	`

	var created Draft
	err := r.db.QueryRowContext(ctx, query,
		draft.ID, draft.UploadedBy, draft.OriginalFilename, draft.FilePath, draft.FileSize,
		draft.DurationSeconds, draft.Bitrate, draft.Format, draft.Status, draft.ExtractedMetadata,
		draft.UploadSource, draft.ClubID, draft.NeedsReview,
	).Scan(
		&created.ID, &created.UploadedBy, &created.OriginalFilename, &created.FilePath,
		&created.FileSize, &created.DurationSeconds, &created.Bitrate,
		&created.Format, &created.Status, &created.ExtractedMetadata,
		&created.UploadSource, &created.ClubID, &created.NeedsReview,
		&created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return Draft{}, err
	}
	return created, nil
}

func (r *Repository) GetDraft(ctx context.Context, id string) (Draft, error) {
	query := `
		SELECT id, uploaded_by, original_filename, file_path, file_size, duration_seconds, bitrate,
		    format, status, extracted_metadata, enriched_metadata, final_metadata,
		    upload_source, club_id, needs_review, review_notes, reviewed_by, reviewed_at,
		    created_at, updated_at
		FROM ingestion_drafts
		WHERE id = $1
	`

	var d Draft
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.UploadedBy, &d.OriginalFilename, &d.FilePath,
		&d.FileSize, &d.DurationSeconds, &d.Bitrate,
		&d.Format, &d.Status, &d.ExtractedMetadata, &d.EnrichedMetadata, &d.FinalMetadata,
		&d.UploadSource, &d.ClubID, &d.NeedsReview, &d.ReviewNotes, &d.ReviewedBy, &d.ReviewedAt,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Draft{}, apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "draft not found", nil)
		}
		return Draft{}, err
	}
	return d, nil
}

func (r *Repository) UpdateDraftStatus(ctx context.Context, id string, status string, reviewNotes *string, reviewedBy *string) error {
	query := `
		UPDATE ingestion_drafts
		SET status = $2, review_notes = $3, reviewed_by = $4, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, status, reviewNotes, reviewedBy)
	return err
}

func (r *Repository) ListPendingDrafts(ctx context.Context, source *string, limit, offset int) ([]Draft, int, error) {
	where := "WHERE needs_review = TRUE AND status = 'pending'"
	args := []any{}
	argIdx := 1

	if source != nil {
		where += " AND upload_source = $" + string(rune('0'+argIdx))
		args = append(args, *source)
		argIdx++
	}

	countQuery := "SELECT COUNT(*) FROM ingestion_drafts " + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, uploaded_by, original_filename, file_path, file_size, duration_seconds, bitrate,
		    format, status, extracted_metadata, upload_source, club_id, needs_review, created_at, updated_at
		FROM ingestion_drafts ` + where + `
		ORDER BY created_at ASC
		LIMIT $` + string(rune('0'+argIdx)) + ` OFFSET $` + string(rune('0'+argIdx+1))

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var drafts []Draft
	for rows.Next() {
		var d Draft
		if err := rows.Scan(
			&d.ID, &d.UploadedBy, &d.OriginalFilename, &d.FilePath,
			&d.FileSize, &d.DurationSeconds, &d.Bitrate,
			&d.Format, &d.Status, &d.ExtractedMetadata,
			&d.UploadSource, &d.ClubID, &d.NeedsReview,
			&d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		drafts = append(drafts, d)
	}
	return drafts, total, nil
}

func (r *Repository) ListUserDrafts(ctx context.Context, userID string, limit, offset int) ([]Draft, int, error) {
	where := "WHERE uploaded_by = $1"
	args := []any{userID}
	argIdx := 2

	countQuery := "SELECT COUNT(*) FROM ingestion_drafts " + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, uploaded_by, original_filename, file_path, file_size, duration_seconds, bitrate,
		    format, status, extracted_metadata, upload_source, club_id, needs_review, created_at, updated_at
		FROM ingestion_drafts ` + where + `
		ORDER BY created_at DESC
		LIMIT $` + string(rune('0'+argIdx)) + ` OFFSET $` + string(rune('0'+argIdx+1))

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var drafts []Draft
	for rows.Next() {
		var d Draft
		if err := rows.Scan(
			&d.ID, &d.UploadedBy, &d.OriginalFilename, &d.FilePath,
			&d.FileSize, &d.DurationSeconds, &d.Bitrate,
			&d.Format, &d.Status, &d.ExtractedMetadata,
			&d.UploadSource, &d.ClubID, &d.NeedsReview,
			&d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		drafts = append(drafts, d)
	}
	return drafts, total, nil
}

func (r *Repository) AttachFile(ctx context.Context, id, filePath, fileURL string, fileSize int64) error {
	query := `UPDATE ingestion_drafts SET file_path = $1, file_size = $2, format = $3, updated_at = NOW() WHERE id = $4`
	ext := ""
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '.' {
			ext = filePath[i+1:]
			break
		}
	}
	_, err := r.db.ExecContext(ctx, query, fileURL, fileSize, ext, id)
	return err
}

func (r *Repository) AddCoUploader(ctx context.Context, co CoUploader) (CoUploader, error) {
	query := `
		INSERT INTO co_uploaders (track_id, user_id, role, xp_share_percent, contribution_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, track_id, user_id, role, xp_share_percent, contribution_type, added_at
	`

	var created CoUploader
	err := r.db.QueryRowContext(ctx, query,
		co.TrackID, co.UserID, co.Role, co.XpSharePercent, co.ContributionType,
	).Scan(
		&created.ID, &created.TrackID, &created.UserID, &created.Role,
		&created.XpSharePercent, &created.ContributionType, &created.AddedAt,
	)
	if err != nil {
		return CoUploader{}, err
	}
	return created, nil
}

func (r *Repository) RemoveCoUploader(ctx context.Context, trackID int64, userID string) error {
	query := `DELETE FROM co_uploaders WHERE track_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, trackID, userID)
	return err
}

func (r *Repository) GetCoUploaders(ctx context.Context, trackID int64) ([]CoUploader, error) {
	query := `
		SELECT id, track_id, user_id, role, xp_share_percent, contribution_type, added_at
		FROM co_uploaders
		WHERE track_id = $1
		ORDER BY added_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cos []CoUploader
	for rows.Next() {
		var co CoUploader
		if err := rows.Scan(&co.ID, &co.TrackID, &co.UserID, &co.Role, &co.XpSharePercent, &co.ContributionType, &co.AddedAt); err != nil {
			return nil, err
		}
		cos = append(cos, co)
	}
	return cos, nil
}

func (r *Repository) GetUserCoUploads(ctx context.Context, userID string) ([]CoUploader, error) {
	query := `
		SELECT id, track_id, user_id, role, xp_share_percent, contribution_type, added_at
		FROM co_uploaders
		WHERE user_id = $1
		ORDER BY added_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cos []CoUploader
	for rows.Next() {
		var co CoUploader
		if err := rows.Scan(&co.ID, &co.TrackID, &co.UserID, &co.Role, &co.XpSharePercent, &co.ContributionType, &co.AddedAt); err != nil {
			return nil, err
		}
		cos = append(cos, co)
	}
	return cos, nil
}

func (r *Repository) GetUploadSlots(ctx context.Context, userID string) (UploadSlot, error) {
	query := `
		SELECT id, user_id, used_slots, max_slots, last_upload_at, created_at, updated_at
		FROM upload_slots
		WHERE user_id = $1
	`

	var slot UploadSlot
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&slot.ID, &slot.UserID, &slot.UsedSlots, &slot.MaxSlots,
		&slot.LastUploadAt, &slot.CreatedAt, &slot.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UploadSlot{}, apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "upload slots not found", nil)
		}
		return UploadSlot{}, err
	}
	return slot, nil
}

func (r *Repository) IncrementSlots(ctx context.Context, userID string) error {
	query := `
		UPDATE upload_slots
		SET used_slots = used_slots + 1, last_upload_at = NOW(), updated_at = NOW()
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) DecrementSlots(ctx context.Context, userID string) error {
	query := `
		UPDATE upload_slots
		SET used_slots = GREATEST(used_slots - 1, 0), updated_at = NOW()
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) UpdateMaxSlots(ctx context.Context, userID string, maxSlots int) error {
	query := `
		UPDATE upload_slots
		SET max_slots = $2, updated_at = NOW()
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, userID, maxSlots)
	return err
}
