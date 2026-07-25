package reputation

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	apperrors "music/internal/common/errors"

	"github.com/jmoiron/sqlx"
)

type RepositoryInterface interface {
	GetUserReputation(ctx context.Context, userID string) (UserReputation, error)
	UpsertUserReputation(ctx context.Context, userID string) error
	UpdateContributionStats(ctx context.Context, userID string, accepted bool) error
	RecordScoreChange(ctx context.Context, userID string, contribType string, contribID *int64, delta float64, reason string, reviewedBy *string) error
	GetTrustTiers(ctx context.Context) ([]TrustTier, error)
	GetTrustTier(ctx context.Context, slug string) (TrustTier, error)
	GetTierForScore(ctx context.Context, score float64, accepted int) (TrustTier, error)
	GetTopContributors(ctx context.Context, limit int) ([]ReputationSummary, error)
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserReputation(ctx context.Context, userID string) (UserReputation, error) {
	query := `
		SELECT id, user_id, total_contributions, accepted_contributions, rejected_contributions,
		       pending_contributions, trust_score, tier, upload_slots, auto_publish, can_review,
		       last_contribution_at, contribution_streak_days, created_at, updated_at
		FROM user_reputation
		WHERE user_id = $1
	`

	var rep UserReputation
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&rep.ID, &rep.UserID, &rep.TotalContributions, &rep.AcceptedContributions,
		&rep.RejectedContributions, &rep.PendingContributions, &rep.TrustScore, &rep.Tier,
		&rep.UploadSlots, &rep.AutoPublish, &rep.CanReview,
		&rep.LastContributionAt, &rep.ContributionStreakDays, &rep.CreatedAt, &rep.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserReputation{}, apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "user reputation not found", nil)
		}
		return UserReputation{}, err
	}
	return rep, nil
}

func (r *Repository) UpsertUserReputation(ctx context.Context, userID string) error {
	query := `
		INSERT INTO user_reputation (user_id) VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) UpdateContributionStats(ctx context.Context, userID string, accepted bool) error {
	var query string
	if accepted {
		query = `
			UPDATE user_reputation
			SET total_contributions = total_contributions + 1,
			    accepted_contributions = accepted_contributions + 1,
			    last_contribution_at = NOW(),
			    updated_at = NOW()
			WHERE user_id = $1
		`
	} else {
		query = `
			UPDATE user_reputation
			SET total_contributions = total_contributions + 1,
			    rejected_contributions = rejected_contributions + 1,
			    updated_at = NOW()
			WHERE user_id = $1
		`
	}
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) RecordScoreChange(ctx context.Context, userID string, contribType string, contribID *int64, delta float64, reason string, reviewedBy *string) error {
	query := `
		INSERT INTO contribution_scores (user_id, contribution_type, contribution_id, score_delta, reason, reviewed_by)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, userID, contribType, contribID, delta, reason, reviewedBy)
	return err
}

func (r *Repository) GetTrustTiers(ctx context.Context) ([]TrustTier, error) {
	query := `
		SELECT id, slug, label, description, min_score, min_accepted, upload_slots, auto_publish, can_review, hierarchy_level, created_at
		FROM trust_tiers
		ORDER BY hierarchy_level ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tiers []TrustTier
	for rows.Next() {
		var t TrustTier
		if err := rows.Scan(&t.ID, &t.Slug, &t.Label, &t.Description, &t.MinScore, &t.MinAccepted, &t.UploadSlots, &t.AutoPublish, &t.CanReview, &t.HierarchyLevel, &t.CreatedAt); err != nil {
			return nil, err
		}
		tiers = append(tiers, t)
	}
	return tiers, nil
}

func (r *Repository) GetTrustTier(ctx context.Context, slug string) (TrustTier, error) {
	query := `
		SELECT id, slug, label, description, min_score, min_accepted, upload_slots, auto_publish, can_review, hierarchy_level, created_at
		FROM trust_tiers
		WHERE slug = $1
	`

	var t TrustTier
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&t.ID, &t.Slug, &t.Label, &t.Description, &t.MinScore, &t.MinAccepted, &t.UploadSlots, &t.AutoPublish, &t.CanReview, &t.HierarchyLevel, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TrustTier{}, apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "trust tier not found", nil)
		}
		return TrustTier{}, err
	}
	return t, nil
}

func (r *Repository) GetTierForScore(ctx context.Context, score float64, accepted int) (TrustTier, error) {
	query := `
		SELECT id, slug, label, description, min_score, min_accepted, upload_slots, auto_publish, can_review, hierarchy_level, created_at
		FROM trust_tiers
		WHERE min_score <= $1 AND min_accepted <= $2
		ORDER BY hierarchy_level DESC
		LIMIT 1
	`

	var t TrustTier
	err := r.db.QueryRowContext(ctx, query, score, accepted).Scan(
		&t.ID, &t.Slug, &t.Label, &t.Description, &t.MinScore, &t.MinAccepted, &t.UploadSlots, &t.AutoPublish, &t.CanReview, &t.HierarchyLevel, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Default to newcomer
			return r.GetTrustTier(ctx, "newcomer")
		}
		return TrustTier{}, err
	}
	return t, nil
}

func (r *Repository) GetTopContributors(ctx context.Context, limit int) ([]ReputationSummary, error) {
	query := `
		SELECT ur.user_id, ur.trust_score, ur.tier, t.label, ur.accepted_contributions,
		       ur.total_contributions, ur.upload_slots, ur.auto_publish, ur.can_review
		FROM user_reputation ur
		JOIN trust_tiers t ON t.slug = ur.tier
		WHERE ur.accepted_contributions > 0
		ORDER BY ur.trust_score DESC, ur.accepted_contributions DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ReputationSummary
	for rows.Next() {
		var s ReputationSummary
		if err := rows.Scan(&s.UserID, &s.TrustScore, &s.Tier, &s.TierLabel, &s.AcceptedContributions, &s.TotalContributions, &s.UploadSlots, &s.AutoPublish, &s.CanReview); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, nil
}
