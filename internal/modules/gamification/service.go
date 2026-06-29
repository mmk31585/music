package gamification

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*UserProfileResponse, error) {
	username, avatarURL, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalXP, err := s.repo.GetUserXP(ctx, userID)
	if err != nil {
		return nil, err
	}

	level, title, titlePersian, nextLevelXP := CalculateLevel(totalXP)

	rank, err := s.repo.GetUserRank(ctx, userID)
	if err != nil {
		rank = 0
	}

	badges, err := s.repo.GetUserBadges(ctx, userID)
	if err != nil {
		badges = []UserBadge{}
	}

	challenges, err := s.repo.GetUserChallenges(ctx, userID)
	if err != nil {
		challenges = []UserChallenge{}
	}

	return &UserProfileResponse{
		UserID:       userID,
		Username:     username,
		AvatarURL:    avatarURL,
		Level:        level,
		CurrentXP:    totalXP,
		NextLevelXP:  nextLevelXP,
		TotalXP:      totalXP,
		Title:        title,
		TitlePersian: titlePersian,
		Rank:         rank,
		Badges:       badges,
		Challenges:   challenges,
	}, nil
}

func (s *Service) AddXP(ctx context.Context, userID string, amount int, source string) (int64, error) {
	meta := fmt.Sprintf(`{"source": "%s"}`, source)
	metaPtr := &meta
	newBalance, err := s.repo.AddXP(ctx, userID, amount, source, metaPtr)
	if err != nil {
		return 0, err
	}

	s.repo.CheckAndAwardBadges(ctx, userID)

	return newBalance, nil
}

func (s *Service) GetBadges(ctx context.Context, userID string) ([]Badge, []UserBadge, error) {
	allBadges, err := s.repo.GetAllBadges(ctx)
	if err != nil {
		return nil, nil, err
	}

	userBadges, err := s.repo.GetUserBadges(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	return allBadges, userBadges, nil
}

func (s *Service) GetActiveChallenges(ctx context.Context, userID string) ([]DailyChallenge, []UserChallenge, error) {
	challenges, err := s.repo.GetActiveChallenges(ctx)
	if err != nil {
		return nil, nil, err
	}

	userChallenges, err := s.repo.GetUserChallenges(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	return challenges, userChallenges, nil
}

func (s *Service) GetLeaderboard(ctx context.Context, lbType string, limit int) (*LeaderboardResponse, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}

	entries, err := s.repo.GetLeaderboard(ctx, lbType, limit)
	if err != nil {
		return nil, err
	}

	return &LeaderboardResponse{
		Type:    lbType,
		Entries: entries,
	}, nil
}

func (s *Service) IncrementChallengeProgress(ctx context.Context, userID, challengeType string) error {
	challenges, err := s.repo.GetActiveChallengesByType(ctx, challengeType)
	if err != nil {
		return err
	}
	if len(challenges) == 0 {
		return nil
	}

	// Exclude already-completed challenges — each event should only increment
	// and reward unfinished challenges.
	activeIDs := make([]string, 0, len(challenges))
	activeChallenges := make([]DailyChallenge, 0, len(challenges))
	for _, c := range challenges {
		done, err := s.repo.IsChallengeCompleted(ctx, userID, c.ID, c.TargetCount)
		if err != nil {
			continue
		}
		if !done {
			activeIDs = append(activeIDs, c.ID)
			activeChallenges = append(activeChallenges, c)
		}
	}
	if len(activeIDs) == 0 {
		return nil
	}

	err = s.repo.IncrementChallengeProgress(ctx, userID, activeIDs)
	if err != nil {
		return err
	}

	// Check each challenge for completion and award rewards
	for _, c := range activeChallenges {
		completed, err := s.repo.IsChallengeCompleted(ctx, userID, c.ID, c.TargetCount)
		if err != nil || !completed {
			continue
		}

		// Mark as completed
		if err := s.repo.CompleteChallenge(ctx, userID, c.ID); err != nil {
			continue
		}

		// Award challenge XP
		if c.XPReward > 0 {
			meta := fmt.Sprintf(`{"challenge_id": "%s", "challenge_type": "%s"}`, c.ID, c.ChallengeType)
			metaPtr := &meta
			_, _ = s.repo.AddXP(ctx, userID, c.XPReward, fmt.Sprintf("challenge_%s", c.ChallengeType), metaPtr)
		}

		// Award linked badge
		if c.BadgeRewardID != nil && *c.BadgeRewardID != "" {
			_ = s.repo.AwardBadge(ctx, userID, *c.BadgeRewardID)
		}
	}

	return nil
}

func (s *Service) CheckBadges(ctx context.Context, userID string) ([]Badge, error) {
	return s.repo.CheckAndAwardBadges(ctx, userID)
}

func (s *Service) AwardXP(ctx context.Context, userID string, amount int, source string) error {
	_, err := s.AddXP(ctx, userID, amount, source)
	return err
}
