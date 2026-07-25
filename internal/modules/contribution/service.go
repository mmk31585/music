package contribution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrForbiddenReview = errors.New("you do not have permission to review contributions")
	ErrAlreadyDecided  = errors.New("contribution has already been decided")
	ErrNotApproved     = errors.New("contribution is not approved")
)

type AIVerifier interface {
	Verify(ctx context.Context, contributionType string, data json.RawMessage) (verdict string, confidence float64, reason string, err error)
}

type XPAwarder interface {
	AwardXP(ctx context.Context, userID string, amount int, source string) error
}

type Service struct {
	repo       *Repository
	aiVerifier AIVerifier
	xpAwarder  XPAwarder
}

func NewService(repo *Repository, aiVerifier AIVerifier, xpAwarder XPAwarder) *Service {
	return &Service{repo: repo, aiVerifier: aiVerifier, xpAwarder: xpAwarder}
}

func (s *Service) Create(ctx context.Context, req CreateContributionRequest, contributorID string) (Contribution, error) {
	req.ContributionType = strings.TrimSpace(req.ContributionType)
	req.TargetType = strings.TrimSpace(req.TargetType)

	dataJSON, err := json.Marshal(req.Data)
	if err != nil {
		return Contribution{}, fmt.Errorf("invalid data: %w", err)
	}

	var aiVerdict *string
	var aiConfidence *float64
	var aiReason *string

	if s.aiVerifier != nil {
		verdict, confidence, reason, vErr := s.aiVerifier.Verify(ctx, req.ContributionType, dataJSON)
		if vErr == nil {
			aiVerdict = &verdict
			aiConfidence = &confidence
			if reason != "" {
				aiReason = &reason
			}
		}
	}

	return s.repo.Create(ctx, req, contributorID, dataJSON, aiVerdict, aiConfidence, aiReason)
}

func (s *Service) GetByID(ctx context.Context, id string) (Contribution, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]Contribution, int, error) {
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.ListByUser(ctx, userID, pageSize, offset)
}

func (s *Service) ListPending(ctx context.Context, page, pageSize int) ([]Contribution, int, error) {
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.ListPending(ctx, pageSize, offset)
}

func (s *Service) ListByStatus(ctx context.Context, status string, page, pageSize int) ([]Contribution, int, error) {
	return s.ListByStatusFiltered(ctx, status, "", "", page, pageSize)
}

func (s *Service) ListByStatusFiltered(ctx context.Context, status, contributionType, search string, page, pageSize int) ([]Contribution, int, error) {
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.ListByStatusFiltered(ctx, status, contributionType, search, pageSize, offset)
}

func (s *Service) ListByTarget(ctx context.Context, targetType, targetID string) ([]Contribution, error) {
	return s.repo.ListByTarget(ctx, targetType, targetID)
}

func (s *Service) Review(ctx context.Context, id, moderatorID string, req ReviewContributionRequest) (Contribution, error) {
	contribution, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Contribution{}, err
	}

	if contribution.Status != "pending" && contribution.Status != "needs_review" {
		return Contribution{}, ErrAlreadyDecided
	}

	newVersion := contribution.Version + 1

	result, err := s.repo.Review(ctx, id, moderatorID, req.Action, req.Note, newVersion)
	if err != nil {
		return Contribution{}, err
	}

	if req.Action == "approve" {
		if err := s.repo.CreateContentVersion(ctx, contribution.TargetType, contribution.TargetID, moderatorID, []byte(contribution.Data), newVersion); err != nil {
			return Contribution{}, fmt.Errorf("failed to create content version: %w", err)
		}

		if s.xpAwarder != nil {
			xpAmount := 50
			switch contribution.ContributionType {
			case "translation":
				xpAmount = 30
			case "credits":
				xpAmount = 20
			case "metadata":
				xpAmount = 15
			case "album_art":
				xpAmount = 25
			case "bio":
				xpAmount = 20
			}
			if err := s.xpAwarder.AwardXP(ctx, contribution.ContributorID, xpAmount, fmt.Sprintf("contribution_%s", contribution.ContributionType)); err != nil {
				return Contribution{}, fmt.Errorf("failed to award xp: %w", err)
			}
		}
	}

	changeType := "approve"
	if req.Action == "reject" {
		changeType = "reject"
	}
	if err := s.repo.RecordHistory(ctx, id, moderatorID, changeType, []byte(contribution.Data), nil); err != nil {
		return Contribution{}, fmt.Errorf("failed to record history: %w", err)
	}

	return result, nil
}

func (s *Service) Apply(ctx context.Context, id string) (Contribution, error) {
	contribution, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Contribution{}, err
	}
	if contribution.Status != "approved" {
		return Contribution{}, ErrNotApproved
	}
	if err := s.repo.ApplyToTarget(ctx, contribution); err != nil {
		return Contribution{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetHistory(ctx context.Context, contributionID string) ([]ContributionHistory, error) {
	return s.repo.GetHistory(ctx, contributionID)
}

func (s *Service) GetContentVersions(ctx context.Context, targetType, targetID string) ([]ContentVersion, error) {
	return s.repo.GetContentVersions(ctx, targetType, targetID)
}

func (s *Service) GetLeaderboard(ctx context.Context, limit int) ([]ContributorStats, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.GetLeaderboard(ctx, limit)
}
