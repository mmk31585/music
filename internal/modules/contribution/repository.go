package contribution

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrContributionNotFound  = errors.New("contribution not found")
	ErrInvalidTarget         = errors.New("invalid target type or id")
	ErrDuplicateContribution = errors.New("duplicate contribution for this target")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, req CreateContributionRequest, contributorID string, dataJSON []byte, aiVerdict *string, aiConfidence *float64, aiReason *string) (Contribution, error) {
	exists, err := r.targetExists(ctx, req.TargetType, req.TargetID)
	if err != nil {
		return Contribution{}, err
	}
	if !exists {
		return Contribution{}, ErrInvalidTarget
	}

	var locale *string
	if req.Locale != "" {
		locale = &req.Locale
	}

	var summary *string
	if req.Summary != "" {
		summary = &req.Summary
	}

	status := "pending"
	if aiVerdict != nil && *aiVerdict == "pass" {
		status = "approved"
	} else if aiVerdict != nil && *aiVerdict == "fail" {
		status = "rejected"
	} else if aiVerdict != nil && *aiVerdict == "review" {
		status = "needs_review"
	}

	now := time.Now()
	var decidedAt, appliedAt *time.Time
	if status == "approved" {
		appliedAt = &now
		decidedAt = &now
	} else if status == "rejected" {
		decidedAt = &now
	}

	query := `
		INSERT INTO contributions (contributor_id, contribution_type, target_type, target_id, locale, data, summary, status, ai_verdict, ai_confidence, ai_reason, decided_at, applied_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, contributor_id, contribution_type, target_type, target_id, locale, data, summary, status, ai_verdict, ai_confidence, ai_reason, moderator_id, moderator_note, version, is_minor, xp_awarded, created_at, decided_at, applied_at, updated_at
	`
	var c Contribution
	err = r.db.QueryRowContext(ctx, query,
		contributorID, req.ContributionType, req.TargetType, req.TargetID, locale, dataJSON, summary,
		status, aiVerdict, aiConfidence, aiReason, decidedAt, appliedAt,
	).Scan(&c.ID, &c.ContributorID, &c.ContributionType, &c.TargetType, &c.TargetID,
		&c.Locale, &c.Data, &c.Summary, &c.Status, &c.AIVerdict, &c.AIConfidence, &c.AIReason,
		&c.ModeratorID, &c.ModeratorNote, &c.Version, &c.IsMinor, &c.XPAwarded,
		&c.CreatedAt, &c.DecidedAt, &c.AppliedAt, &c.UpdatedAt,
	)
	if err != nil {
		return Contribution{}, fmt.Errorf("failed to create contribution: %w", err)
	}
	return c, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (Contribution, error) {
	query := `SELECT id, contributor_id, contribution_type, target_type, target_id, locale, data, summary, status, ai_verdict, ai_confidence, ai_reason, moderator_id, moderator_note, version, is_minor, xp_awarded, created_at, decided_at, applied_at, updated_at FROM contributions WHERE id = $1`
	var c Contribution
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.ContributorID, &c.ContributionType, &c.TargetType, &c.TargetID,
		&c.Locale, &c.Data, &c.Summary, &c.Status, &c.AIVerdict, &c.AIConfidence, &c.AIReason,
		&c.ModeratorID, &c.ModeratorNote, &c.Version, &c.IsMinor, &c.XPAwarded,
		&c.CreatedAt, &c.DecidedAt, &c.AppliedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Contribution{}, ErrContributionNotFound
	}
	return c, err
}

func (r *Repository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]Contribution, int, error) {
	countQuery := `SELECT COUNT(*) FROM contributions WHERE contributor_id = $1`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, contributor_id, contribution_type, target_type, target_id, locale, data, summary, status, ai_verdict, ai_confidence, ai_reason, moderator_id, moderator_note, version, is_minor, xp_awarded, created_at, decided_at, applied_at, updated_at FROM contributions WHERE contributor_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []Contribution
	for rows.Next() {
		var c Contribution
		if err := rows.Scan(&c.ID, &c.ContributorID, &c.ContributionType, &c.TargetType, &c.TargetID,
			&c.Locale, &c.Data, &c.Summary, &c.Status, &c.AIVerdict, &c.AIConfidence, &c.AIReason,
			&c.ModeratorID, &c.ModeratorNote, &c.Version, &c.IsMinor, &c.XPAwarded,
			&c.CreatedAt, &c.DecidedAt, &c.AppliedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, c)
	}
	return items, total, rows.Err()
}

func (r *Repository) ListPending(ctx context.Context, limit, offset int) ([]Contribution, int, error) {
	countQuery := `SELECT COUNT(*) FROM contributions WHERE status = 'pending' OR status = 'needs_review'`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, contributor_id, contribution_type, target_type, target_id, locale, data, summary, status, ai_verdict, ai_confidence, ai_reason, moderator_id, moderator_note, version, is_minor, xp_awarded, created_at, decided_at, applied_at, updated_at FROM contributions WHERE status = 'pending' OR status = 'needs_review' ORDER BY created_at ASC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []Contribution
	for rows.Next() {
		var c Contribution
		if err := rows.Scan(&c.ID, &c.ContributorID, &c.ContributionType, &c.TargetType, &c.TargetID,
			&c.Locale, &c.Data, &c.Summary, &c.Status, &c.AIVerdict, &c.AIConfidence, &c.AIReason,
			&c.ModeratorID, &c.ModeratorNote, &c.Version, &c.IsMinor, &c.XPAwarded,
			&c.CreatedAt, &c.DecidedAt, &c.AppliedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, c)
	}
	return items, total, rows.Err()
}

func (r *Repository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]Contribution, int, error) {
	return r.ListByStatusFiltered(ctx, status, "", "", limit, offset)
}

func (r *Repository) ListByStatusFiltered(ctx context.Context, status, contributionType, search string, limit, offset int) ([]Contribution, int, error) {
	args := []any{}
	where := "c.status = $" + strconv.Itoa(len(args)+1)
	args = append(args, status)

	if contributionType != "" {
		where += " AND c.contribution_type = $" + strconv.Itoa(len(args)+1)
		args = append(args, contributionType)
	}
	if search != "" {
		where += " AND (c.summary ILIKE $" + strconv.Itoa(len(args)+1) + " OR u.username ILIKE $" + strconv.Itoa(len(args)+1) + ")"
		args = append(args, "%"+search+"%")
	}

	countQuery := `SELECT COUNT(*) FROM contributions c LEFT JOIN users u ON u.id = c.contributor_id WHERE ` + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT c.id, c.contributor_id, c.contribution_type, c.target_type, c.target_id, c.locale, c.data, c.summary, c.status, c.ai_verdict, c.ai_confidence, c.ai_reason, c.moderator_id, c.moderator_note, c.version, c.is_minor, c.xp_awarded, c.created_at, c.decided_at, c.applied_at, c.updated_at, u.username, u.avatar_url FROM contributions c LEFT JOIN users u ON u.id = c.contributor_id WHERE ` + where + ` ORDER BY c.created_at DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []Contribution
	for rows.Next() {
		var c Contribution
		if err := rows.Scan(&c.ID, &c.ContributorID, &c.ContributionType, &c.TargetType, &c.TargetID,
			&c.Locale, &c.Data, &c.Summary, &c.Status, &c.AIVerdict, &c.AIConfidence, &c.AIReason,
			&c.ModeratorID, &c.ModeratorNote, &c.Version, &c.IsMinor, &c.XPAwarded,
			&c.CreatedAt, &c.DecidedAt, &c.AppliedAt, &c.UpdatedAt,
			&c.ContributorUsername, &c.ContributorAvatarURL,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, c)
	}
	return items, total, rows.Err()
}

func (r *Repository) ListByTarget(ctx context.Context, targetType, targetID string) ([]Contribution, error) {
	query := `SELECT id, contributor_id, contribution_type, target_type, target_id, locale, data, summary, status, ai_verdict, ai_confidence, ai_reason, moderator_id, moderator_note, version, is_minor, xp_awarded, created_at, decided_at, applied_at, updated_at FROM contributions WHERE target_type = $1 AND target_id = $2 AND status = 'approved' ORDER BY version DESC`
	rows, err := r.db.QueryContext(ctx, query, targetType, targetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Contribution
	for rows.Next() {
		var c Contribution
		if err := rows.Scan(&c.ID, &c.ContributorID, &c.ContributionType, &c.TargetType, &c.TargetID,
			&c.Locale, &c.Data, &c.Summary, &c.Status, &c.AIVerdict, &c.AIConfidence, &c.AIReason,
			&c.ModeratorID, &c.ModeratorNote, &c.Version, &c.IsMinor, &c.XPAwarded,
			&c.CreatedAt, &c.DecidedAt, &c.AppliedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (r *Repository) Review(ctx context.Context, id, moderatorID string, action, note string, version int) (Contribution, error) {
	var c Contribution

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return Contribution{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	var currentStatus string
	err = tx.QueryRowContext(ctx, `SELECT status FROM contributions WHERE id = $1 FOR UPDATE`, id).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Contribution{}, ErrContributionNotFound
		}
		return Contribution{}, err
	}

	if currentStatus != "pending" && currentStatus != "needs_review" {
		return Contribution{}, fmt.Errorf("contribution already %s", currentStatus)
	}

	now := time.Now()
	var newStatus string
	if action == "approve" {
		newStatus = "approved"
	} else {
		newStatus = "rejected"
	}

	var notePtr *string
	if note != "" {
		notePtr = &note
	}

	xpReward := 0
	if action == "approve" {
		switch currentStatus {
		case "lyrics":
			xpReward = 50
		case "translation":
			xpReward = 30
		case "credits":
			xpReward = 20
		default:
			xpReward = 10
		}
	}

	err = tx.QueryRowContext(ctx, `
		UPDATE contributions SET status = $1, moderator_id = $2, moderator_note = $3, decided_at = $4, applied_at = $5, version = $6, xp_awarded = $7, updated_at = $4
		WHERE id = $8
		RETURNING id, contributor_id, contribution_type, target_type, target_id, locale, data, summary, status, ai_verdict, ai_confidence, ai_reason, moderator_id, moderator_note, version, is_minor, xp_awarded, created_at, decided_at, applied_at, updated_at
	`, newStatus, moderatorID, notePtr, now, now, version, xpReward, id,
	).Scan(&c.ID, &c.ContributorID, &c.ContributionType, &c.TargetType, &c.TargetID,
		&c.Locale, &c.Data, &c.Summary, &c.Status, &c.AIVerdict, &c.AIConfidence, &c.AIReason,
		&c.ModeratorID, &c.ModeratorNote, &c.Version, &c.IsMinor, &c.XPAwarded,
		&c.CreatedAt, &c.DecidedAt, &c.AppliedAt, &c.UpdatedAt,
	)
	if err != nil {
		return Contribution{}, err
	}

	return c, tx.Commit()
}

func (r *Repository) GetHistory(ctx context.Context, contributionID string) ([]ContributionHistory, error) {
	query := `SELECT id, contribution_id, data, previous_data, changed_by, change_type, created_at FROM contribution_history WHERE contribution_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, contributionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ContributionHistory
	for rows.Next() {
		var h ContributionHistory
		if err := rows.Scan(&h.ID, &h.ContributionID, &h.Data, &h.PreviousData, &h.ChangedBy, &h.ChangeType, &h.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, h)
	}
	return items, rows.Err()
}

func (r *Repository) RecordHistory(ctx context.Context, contributionID, changedBy, changeType string, data, previousData []byte) error {
	query := `INSERT INTO contribution_history (contribution_id, data, previous_data, changed_by, change_type) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, contributionID, data, previousData, changedBy, changeType)
	return err
}

func (r *Repository) GetContentVersions(ctx context.Context, targetType, targetID string) ([]ContentVersion, error) {
	query := `SELECT id, target_type, target_id, version, data, applied_by, created_at FROM content_versions WHERE target_type = $1 AND target_id = $2 ORDER BY version DESC`
	rows, err := r.db.QueryContext(ctx, query, targetType, targetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ContentVersion
	for rows.Next() {
		var v ContentVersion
		if err := rows.Scan(&v.ID, &v.TargetType, &v.TargetID, &v.Version, &v.Data, &v.AppliedBy, &v.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	return items, rows.Err()
}

func (r *Repository) CreateContentVersion(ctx context.Context, targetType, targetID, appliedBy string, data []byte, version int) error {
	query := `INSERT INTO content_versions (target_type, target_id, version, data, applied_by) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (target_type, target_id, version) DO UPDATE SET data = $4, applied_by = $5, created_at = NOW()`
	_, err := r.db.ExecContext(ctx, query, targetType, targetID, version, data, appliedBy)
	return err
}

func (r *Repository) GetLeaderboard(ctx context.Context, limit int) ([]ContributorStats, error) {
	query := `
		SELECT u.id AS user_id, u.username, u.avatar_url,
			COUNT(c.id) AS total_contributions,
			COUNT(CASE WHEN c.status = 'approved' THEN 1 END) AS approved_count,
			COUNT(CASE WHEN c.status = 'pending' THEN 1 END) AS pending_count,
			COUNT(CASE WHEN c.status = 'rejected' THEN 1 END) AS rejected_count,
			COALESCE(SUM(c.xp_awarded), 0) AS xp_earned
		FROM contributions c
		JOIN users u ON u.id = c.contributor_id
		GROUP BY u.id, u.username, u.avatar_url
		ORDER BY xp_earned DESC, approved_count DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ContributorStats
	rank := 1
	for rows.Next() {
		var s ContributorStats
		if err := rows.Scan(&s.UserID, &s.Username, &s.AvatarURL, &s.TotalContributions, &s.ApprovedCount, &s.PendingCount, &s.RejectedCount, &s.XPEarned); err != nil {
			return nil, err
		}
		s.Rank = rank
		rank++
		items = append(items, s)
	}
	return items, rows.Err()
}

var ErrCannotApply = errors.New("cannot apply this contribution type to the target")

func (r *Repository) ApplyToTarget(ctx context.Context, contribution Contribution) error {
	if contribution.Status != "approved" {
		return ErrNotApproved
	}

	switch ContributionType(contribution.ContributionType) {
	case ContributionTypeLyrics:
		_, err := r.db.ExecContext(ctx, `UPDATE tracks SET lyrics = $1::jsonb, updated_at = NOW() WHERE id = $2`, contribution.Data, contribution.TargetID)
		if err != nil {
			return fmt.Errorf("failed to apply lyrics: %w", err)
		}
	case ContributionTypeTranslation:
		_, err := r.db.ExecContext(ctx, `UPDATE tracks SET translated_lyrics = $1::jsonb, updated_at = NOW() WHERE id = $2`, contribution.Data, contribution.TargetID)
		if err != nil {
			return fmt.Errorf("failed to apply translation: %w", err)
		}
	case ContributionTypeCredits:
		_, err := r.db.ExecContext(ctx, `UPDATE tracks SET credits = $1::jsonb, updated_at = NOW() WHERE id = $2`, contribution.Data, contribution.TargetID)
		if err != nil {
			return fmt.Errorf("failed to apply credits: %w", err)
		}
	case ContributionTypeMetadata:
		var table string
		switch TargetType(contribution.TargetType) {
		case TargetTypeTrack:
			table = "tracks"
		case TargetTypeAlbum:
			table = "albums"
		case TargetTypeArtist:
			table = "artists"
		default:
			return ErrCannotApply
		}
		_, err := r.db.ExecContext(ctx, fmt.Sprintf(`UPDATE %s SET metadata = $1::jsonb, updated_at = NOW() WHERE id = $2`, table), contribution.Data, contribution.TargetID)
		if err != nil {
			return fmt.Errorf("failed to apply metadata: %w", err)
		}
	case ContributionTypeAlbumArt:
		var data map[string]any
		if err := json.Unmarshal([]byte(contribution.Data), &data); err != nil {
			return fmt.Errorf("invalid album_art data: %w", err)
		}
		coverURL, ok := data["cover_url"].(string)
		if !ok || coverURL == "" {
			return fmt.Errorf("cover_url required in album_art data")
		}
		var table string
		var column string
		switch TargetType(contribution.TargetType) {
		case TargetTypeTrack:
			table, column = "tracks", "cover_media_id"
		case TargetTypeAlbum:
			table, column = "albums", "cover_media_id"
		case TargetTypeArtist:
			table, column = "artists", "avatar_media_id"
		default:
			return ErrCannotApply
		}
		_, err := r.db.ExecContext(ctx, fmt.Sprintf(`UPDATE %s SET %s = $1, updated_at = NOW() WHERE id = $2`, table, column), coverURL, contribution.TargetID)
		if err != nil {
			return fmt.Errorf("failed to apply album art: %w", err)
		}
	case ContributionTypeBio:
		_, err := r.db.ExecContext(ctx, `UPDATE artists SET bio = $1, updated_at = NOW() WHERE id = $2`, contribution.Data, contribution.TargetID)
		if err != nil {
			return fmt.Errorf("failed to apply bio: %w", err)
		}
	default:
		return ErrCannotApply
	}

	// Mark the contribution as applied
	_, err := r.db.ExecContext(ctx, `UPDATE contributions SET applied_at = NOW(), updated_at = NOW() WHERE id = $1`, contribution.ID)
	return err
}

func (r *Repository) targetExists(ctx context.Context, targetType, targetID string) (bool, error) {
	var table string
	switch TargetType(targetType) {
	case TargetTypeTrack:
		table = "tracks"
	case TargetTypeAlbum:
		table = "albums"
	case TargetTypeArtist:
		table = "artists"
	default:
		return false, ErrInvalidTarget
	}
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", table)
	var exists bool
	err := r.db.QueryRowContext(ctx, query, targetID).Scan(&exists)
	return exists, err
}
