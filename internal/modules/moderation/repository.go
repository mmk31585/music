package moderation

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateReport(ctx context.Context, report ContentReport) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO content_reports (id, reporter_id, target_id, target_type, reason, description)
		VALUES (:id, :reporter_id, :target_id, :target_type, :reason, :description)
	`, report)
	return err
}

func (r *Repository) GetReportByID(ctx context.Context, id uuid.UUID) (*ContentReport, error) {
	var report ContentReport
	err := r.db.GetContext(ctx, &report, `
		SELECT id, reporter_id, target_id, target_type, reason, description, status, moderator_id, resolved_at, created_at, resolution_note
		FROM content_reports WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *Repository) GetPendingReports(ctx context.Context, limit, offset int) ([]ContentReport, error) {
	var items []ContentReport
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, reporter_id, target_id, target_type, reason, description, status, moderator_id, resolved_at, created_at, resolution_note
		FROM content_reports
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *Repository) GetReportsByStatus(ctx context.Context, status string, limit, offset int) ([]ContentReport, error) {
	var items []ContentReport
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, reporter_id, target_id, target_type, reason, description, status, moderator_id, resolved_at, created_at, resolution_note
		FROM content_reports
		WHERE status = $1
		ORDER BY resolved_at DESC NULLS LAST, created_at DESC
		LIMIT $2 OFFSET $3
	`, status, limit, offset)
	return items, err
}

func (r *Repository) ResolveReport(ctx context.Context, reportID uuid.UUID, moderatorID uuid.UUID, status, note string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE content_reports
		SET status = $1, moderator_id = $2, resolved_at = NOW(), resolution_note = $3
		WHERE id = $4
	`, status, moderatorID, note, reportID)
	return err
}

func (r *Repository) BulkResolve(ctx context.Context, ids []uuid.UUID, moderatorID uuid.UUID, status, note string) error {
	if len(ids) == 0 {
		return nil
	}
	qry := `UPDATE content_reports SET status = $1, moderator_id = $2, resolved_at = NOW(), resolution_note = $3 WHERE id = ANY($4)`
	_, err := r.db.ExecContext(ctx, qry, status, moderatorID, note, ids)
	return err
}

func (r *Repository) FlagContent(ctx context.Context, flag ContentFlag) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO content_flags (id, target_id, target_type, flag_type, expires_at)
		VALUES (:id, :target_id, :target_type, :flag_type, :expires_at)
	`, flag)
	return err
}

func (r *Repository) GetFlags(ctx context.Context, targetID, targetType string) ([]ContentFlag, error) {
	var items []ContentFlag
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, target_id, target_type, flag_type, flagged_at, expires_at
		FROM content_flags
		WHERE target_id = $1 AND target_type = $2
	`, targetID, targetType)
	return items, err
}

func (r *Repository) GetAllFlags(ctx context.Context, limit, offset int, expired bool) ([]ContentFlag, error) {
	var items []ContentFlag
	var err error
	if expired {
		err = r.db.SelectContext(ctx, &items, `
			SELECT id, target_id, target_type, flag_type, flagged_at, expires_at
			FROM content_flags ORDER BY flagged_at DESC LIMIT $1 OFFSET $2
		`, limit, offset)
	} else {
		err = r.db.SelectContext(ctx, &items, `
			SELECT id, target_id, target_type, flag_type, flagged_at, expires_at
			FROM content_flags WHERE expires_at IS NULL OR expires_at > NOW()
			ORDER BY flagged_at DESC LIMIT $1 OFFSET $2
		`, limit, offset)
	}
	return items, err
}

func (r *Repository) RemoveExpiredFlags(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM content_flags WHERE expires_at < NOW()`)
	return err
}

func (r *Repository) CreateAction(ctx context.Context, a ModerationAction) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO moderation_actions (id, report_id, moderator_id, action, target_id, target_type, previous_status, new_status, note, metadata)
		VALUES (:id, :report_id, :moderator_id, :action, :target_id, :target_type, :previous_status, :new_status, :note, :metadata)
	`, a)
	return err
}

func (r *Repository) GetActions(ctx context.Context, reportID *uuid.UUID, limit, offset int) ([]ModerationAction, error) {
	var items []ModerationAction
	var err error
	if reportID != nil {
		err = r.db.SelectContext(ctx, &items, `
			SELECT id, report_id, moderator_id, action, target_id, target_type, previous_status, new_status, note, metadata, created_at
			FROM moderation_actions WHERE report_id = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3
		`, *reportID, limit, offset)
	} else {
		err = r.db.SelectContext(ctx, &items, `
			SELECT id, report_id, moderator_id, action, target_id, target_type, previous_status, new_status, note, metadata, created_at
			FROM moderation_actions
			ORDER BY created_at DESC LIMIT $1 OFFSET $2
		`, limit, offset)
	}
	return items, err
}

func (r *Repository) GetStats(ctx context.Context) (*ModerationStats, error) {
	var s ModerationStats
	err := r.db.GetContext(ctx, &s, `
		SELECT
			COALESCE((SELECT COUNT(*) FROM content_reports), 0) AS total_reports,
			COALESCE((SELECT COUNT(*) FROM content_reports WHERE status = 'pending'), 0) AS pending_reports,
			COALESCE((SELECT COUNT(*) FROM content_reports WHERE resolved_at::date = CURRENT_DATE), 0) AS resolved_today,
			COALESCE((SELECT COUNT(*) FROM content_flags WHERE expires_at IS NULL OR expires_at > NOW()), 0) AS flagged_content,
			COALESCE((SELECT COUNT(DISTINCT reporter_id) FROM content_reports), 0) AS unique_reporters
	`)
	if err != nil {
		return nil, err
	}

	var reasons []struct {
		Reason string `db:"reason"`
		Count  int64  `db:"cnt"`
	}
	err = r.db.SelectContext(ctx, &reasons, `
		SELECT reason, COUNT(*) AS cnt FROM content_reports GROUP BY reason ORDER BY cnt DESC
	`)
	if err == nil {
		s.ByReason = make(map[string]int64, len(reasons))
		for _, r := range reasons {
			s.ByReason[r.Reason] = r.Count
		}
	}

	var types []struct {
		TargetType string `db:"target_type"`
		Count      int64  `db:"cnt"`
	}
	err = r.db.SelectContext(ctx, &types, `
		SELECT target_type, COUNT(*) AS cnt FROM content_reports GROUP BY target_type ORDER BY cnt DESC
	`)
	if err == nil {
		s.ByTargetType = make(map[string]int64, len(types))
		for _, t := range types {
			s.ByTargetType[t.TargetType] = t.Count
		}
	}

	var avg struct {
		Hours float64 `db:"hours"`
	}
	err = r.db.GetContext(ctx, &avg, `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (resolved_at - created_at)) / 3600), 0) AS hours
		FROM content_reports WHERE status != 'pending' AND resolved_at IS NOT NULL
	`)
	if err == nil {
		s.AvgResolutionHrs = avg.Hours
	}

	return &s, nil
}
