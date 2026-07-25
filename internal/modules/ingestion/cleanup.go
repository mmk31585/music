package ingestion

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type CleanupService struct {
	repo  *Repository
	logger *zap.Logger
}

func NewCleanupService(repo *Repository, logger *zap.Logger) *CleanupService {
	return &CleanupService{repo: repo, logger: logger}
}

func (s *CleanupService) FlagStaleDrafts(ctx context.Context, maxAge time.Duration) (int, error) {
	cutoff := time.Now().Add(-maxAge)
	drafts, err := s.repo.ListStaleDrafts(ctx, cutoff)
	if err != nil {
		return 0, err
	}
	flagged := 0
	for _, d := range drafts {
		if err := s.repo.MarkDraftStale(ctx, d.ID); err != nil {
			s.logger.Warn("failed to mark draft stale", zap.String("draft_id", d.ID), zap.Error(err))
			continue
		}
		flagged++
	}
	s.logger.Info("stale drafts flagged", zap.Int("count", flagged), zap.Duration("max_age", maxAge))
	return flagged, nil
}

func (s *CleanupService) StartPeriodicCleanup(ctx context.Context, interval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if _, err := s.FlagStaleDrafts(ctx, maxAge); err != nil {
					s.logger.Error("periodic stale draft cleanup failed", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	s.logger.Info("periodic stale draft cleanup started", zap.Duration("interval", interval), zap.Duration("max_age", maxAge))
}
