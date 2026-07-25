package reputation

import (
	"context"
)

type ServiceInterface interface {
	GetUserReputation(ctx context.Context, userID string) (ReputationSummary, error)
	RecordContribution(ctx context.Context, userID string, contribType string, contribID *int64, accepted bool, reviewedBy *string) error
	GetTrustTiers(ctx context.Context) ([]TrustTier, error)
	GetTopContributors(ctx context.Context, limit int) ([]ReputationSummary, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserReputation(ctx context.Context, userID string) (ReputationSummary, error) {
	rep, err := s.repo.GetUserReputation(ctx, userID)
	if err != nil {
		return ReputationSummary{}, err
	}

	tier, err := s.repo.GetTrustTier(ctx, rep.Tier)
	if err != nil {
		// Fallback: use minimal info from reputation
		return ReputationSummary{
			UserID:                rep.UserID,
			TrustScore:            rep.TrustScore,
			Tier:                  rep.Tier,
			TierLabel:             rep.Tier,
			AcceptedContributions: rep.AcceptedContributions,
			TotalContributions:    rep.TotalContributions,
			UploadSlots:           rep.UploadSlots,
			AutoPublish:           rep.AutoPublish,
			CanReview:             rep.CanReview,
		}, nil
	}

	return ReputationSummary{
		UserID:                rep.UserID,
		TrustScore:            rep.TrustScore,
		Tier:                  rep.Tier,
		TierLabel:             tier.Label,
		AcceptedContributions: rep.AcceptedContributions,
		TotalContributions:    rep.TotalContributions,
		UploadSlots:           rep.UploadSlots,
		AutoPublish:           rep.AutoPublish,
		CanReview:             rep.CanReview,
	}, nil
}

func (s *Service) RecordContribution(ctx context.Context, userID string, contribType string, contribID *int64, accepted bool, reviewedBy *string) error {
	// 1. Update stats
	if err := s.repo.UpdateContributionStats(ctx, userID, accepted); err != nil {
		return err
	}

	// 2. Calculate score delta
	var delta float64
	var reason string
	if accepted {
		switch contribType {
		case "track_upload":
			delta = 5.0
			reason = "track upload accepted"
		case "metadata_edit":
			delta = 2.0
			reason = "metadata edit accepted"
		case "lyrics_edit":
			delta = 3.0
			reason = "lyrics edit accepted"
		case "review":
			delta = 1.0
			reason = "review completed"
		default:
			delta = 1.0
			reason = contribType + " accepted"
		}
	} else {
		delta = -1.0
		reason = contribType + " rejected"
	}

	// 3. Record score change
	if err := s.repo.RecordScoreChange(ctx, userID, contribType, contribID, delta, reason, reviewedBy); err != nil {
		return err
	}

	// 4. Recalculate tier
	rep, err := s.repo.GetUserReputation(ctx, userID)
	if err != nil {
		return err
	}

	newTier, err := s.repo.GetTierForScore(ctx, rep.TrustScore, rep.AcceptedContributions)
	if err != nil {
		return err
	}

	// 5. Update tier and capabilities if changed
	if newTier.Slug != rep.Tier {
		// Tier changed — update reputation record
		// This is handled by a trigger or manual update in production
		// For now, we just record the change
		if err := s.repo.RecordScoreChange(ctx, userID, "tier_change", nil, 0, "tier changed to "+newTier.Slug, nil); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetTrustTiers(ctx context.Context) ([]TrustTier, error) {
	return s.repo.GetTrustTiers(ctx)
}

func (s *Service) GetTopContributors(ctx context.Context, limit int) ([]ReputationSummary, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetTopContributors(ctx, limit)
}
