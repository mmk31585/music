package creator

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOverview(ctx context.Context, userID string) (*CreatorStats, error) {
	uid, _ := uuid.Parse(userID)
	stats, err := s.repo.GetStats(ctx, uid)
	if err != nil {
		_ = s.repo.RefreshStats(ctx, uid)
		stats, _ = s.repo.GetStats(ctx, uid)
	}
	return stats, nil
}

func (s *Service) GetDailyStats(ctx context.Context, userID, from, to string, limit int) ([]CreatorDailyStat, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetDailyStats(ctx, uid, from, to, limit)
}

func (s *Service) GetTrackStats(ctx context.Context, userID string) ([]TrackStats, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetTrackStats(ctx, uid)
}

func (s *Service) RefreshStats(ctx context.Context, userID string) error {
	uid, _ := uuid.Parse(userID)
	return s.repo.RefreshStats(ctx, uid)
}

func (s *Service) IsCreator(ctx context.Context, userID string) (bool, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.IsCreator(ctx, uid)
}

// --- Earnings ---

func (s *Service) GetEarningsBreakdown(ctx context.Context, userID string) (*EarningsBreakdown, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetEarningsBreakdown(ctx, uid)
}

func (s *Service) GetPayoutHistory(ctx context.Context, userID string, limit int) ([]Payout, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetPayoutHistory(ctx, uid, limit)
}

func (s *Service) GetPayoutMethods(ctx context.Context, userID string) ([]PayoutMethod, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetPayoutMethods(ctx, uid)
}

// --- Audience ---

func (s *Service) GetTopListeners(ctx context.Context, userID string, limit int) ([]TopListener, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetTopListeners(ctx, uid, limit)
}

func (s *Service) GetGeographicStats(ctx context.Context, userID string) ([]GeographicStat, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetGeographicStats(ctx, uid)
}

func (s *Service) GetAudienceOverview(ctx context.Context, userID string) (*AudienceOverview, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetAudienceOverview(ctx, uid)
}

// --- Content Management ---

func (s *Service) GetCreatorContent(ctx context.Context, userID string) (*CreatorContent, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetCreatorContent(ctx, uid)
}

func (s *Service) UpdateTrack(ctx context.Context, userID, trackID string, req TrackUpdateRequest) error {
	uid, _ := uuid.Parse(userID)
	tid, _ := uuid.Parse(trackID)
	// Verify ownership
	owned, _ := s.repo.IsCreator(ctx, uid)
	if !owned {
		return fmt.Errorf("not a creator")
	}
	return s.repo.UpdateTrack(ctx, tid, req)
}

func (s *Service) UpdateAlbum(ctx context.Context, userID, albumID string, req AlbumUpdateRequest) error {
	uid, _ := uuid.Parse(userID)
	aid, _ := uuid.Parse(albumID)
	owned, _ := s.repo.IsCreator(ctx, uid)
	if !owned {
		return fmt.Errorf("not a creator")
	}
	return s.repo.UpdateAlbum(ctx, aid, req)
}

func (s *Service) DeleteTrack(ctx context.Context, userID, trackID string) error {
	uid, _ := uuid.Parse(userID)
	tid, _ := uuid.Parse(trackID)
	return s.repo.DeleteTrack(ctx, tid, uid)
}

func (s *Service) DeleteAlbum(ctx context.Context, userID, albumID string) error {
	uid, _ := uuid.Parse(userID)
	aid, _ := uuid.Parse(albumID)
	return s.repo.DeleteAlbum(ctx, aid, uid)
}
