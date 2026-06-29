package acquisition

import (
	"context"
	"fmt"
	"time"

	"music/internal/platform/metrics"

	"go.uber.org/zap"
)

type Service struct {
	logger    *zap.Logger
	resolvers []Resolver
}

func NewService(logger *zap.Logger, proxy string) *Service {
	resolvers := []Resolver{
		newSoundCloudResolver(proxy),
		newYTMusicResolver(proxy),
		// iTunes resolver intentionally excluded — it only returns 30-second
		// preview clips, not full tracks. Full tracks come from SoundCloud/YouTube.
	}

	return &Service{
		logger:    logger,
		resolvers: resolvers,
	}
}

func (s *Service) Resolve(ctx context.Context, q ResolveQuery) (*Candidate, error) {
	start := time.Now()

	for _, r := range sortByPriority(s.resolvers) {
		select {
		case <-ctx.Done():
			s.logger.Warn("resolve cancelled",
				zap.String("query", q.Artist+" - "+q.Title),
				zap.Duration("elapsed", time.Since(start)),
			)
			return nil, ctx.Err()
		default:
		}

		rStart := time.Now()
		s.logger.Debug("trying resolver",
			zap.String("resolver", r.Name()),
			zap.String("query", q.Artist+" - "+q.Title),
		)

		resolveCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		candidate, err := r.Resolve(resolveCtx, q)
		cancel()

		if err != nil {
			metrics.IncImportAcquisition(r.Name(), "fail", rStart)
			s.logger.Debug("resolver failed, trying next",
				zap.String("resolver", r.Name()),
				zap.Error(err),
			)
			continue
		}
		if candidate == nil {
			metrics.IncImportAcquisition(r.Name(), "fail", rStart)
			continue
		}

		metrics.IncImportAcquisition(r.Name(), "success", rStart)
		s.logger.Info("track resolved",
			zap.String("resolver", r.Name()),
			zap.String("query", q.Artist+" - "+q.Title),
			zap.String("url", candidate.URL),
			zap.String("quality", candidate.Quality),
			zap.Duration("elapsed", time.Since(start)),
		)
		return candidate, nil
	}

	return nil, fmt.Errorf("Could not find a downloadable copy of '%s - %s'. Try a different search term or use a direct YouTube/SoundCloud URL.", q.Artist, q.Title)
}

func sortByPriority(resolvers []Resolver) []Resolver {
	sorted := make([]Resolver, len(resolvers))
	copy(sorted, resolvers)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].Priority() > sorted[j].Priority() {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}
