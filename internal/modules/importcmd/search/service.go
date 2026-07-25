package search

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"music/internal/platform/cache"
	"music/internal/platform/circuitbreaker"
	"music/internal/platform/fuzzy"
	"music/internal/platform/metrics"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Metrics struct {
	SearchesTotal   int64
	CacheHits       int64
	CacheMisses     int64
	ProviderLatency map[string]time.Duration
	ProviderErrors  map[string]int64
	ProviderSuccess map[string]int64
	mu              sync.Mutex
}

func NewMetrics() *Metrics {
	return &Metrics{
		ProviderLatency: make(map[string]time.Duration),
		ProviderErrors:  make(map[string]int64),
		ProviderSuccess: make(map[string]int64),
	}
}

type Service struct {
	providers []Provider
	cache     *cache.SearchCache
	metrics   *Metrics
	logger    *zap.Logger
	breakers  map[string]*circuitbreaker.Breaker
}

func NewService(logger *zap.Logger, rdb redis.UniversalClient, providers []Provider) *Service {
	svc := &Service{
		providers: providers,
		cache:     cache.NewSearchCache(rdb, 5*time.Minute),
		metrics:   NewMetrics(),
		logger:    logger,
		breakers:  make(map[string]*circuitbreaker.Breaker),
	}

	for _, p := range providers {
		svc.breakers[p.Name()] = circuitbreaker.New(circuitbreaker.Config{
			Threshold:    3,
			RecoveryTime: 30 * time.Second,
			HalfOpenMax:  1,
		})
	}

	return svc
}

func (s *Service) Metrics() *Metrics { return s.metrics }

func (s *Service) Search(ctx context.Context, rawQuery string) ([]Result, error) {
	start := time.Now()
	q := SearchQuery{Raw: rawQuery}

	if parts := splitQuery(rawQuery); len(parts) >= 2 {
		q.Artist = parts[0]
		q.Title = parts[1]
	}

	cacheKey := normalizeCacheKey(rawQuery)
	if cached, hit, err := s.cache.Get(ctx, cacheKey); err == nil && hit {
		s.metrics.mu.Lock()
		s.metrics.CacheHits++
		s.metrics.mu.Unlock()

		metrics.IncImportSearchTotal("cache", true)
		metrics.IncImportCacheOp("hit")
		metrics.ObserveImportSearch("cache", start)

		var results []Result
		if err := json.Unmarshal(cached, &results); err == nil {
			s.logger.Debug("search cache hit",
				zap.String("query", rawQuery),
				zap.Int("results", len(results)),
				zap.Duration("elapsed", time.Since(start)),
			)
			return results, nil
		}
	}

	metrics.IncImportCacheOp("miss")

	s.metrics.mu.Lock()
	s.metrics.CacheMisses++
	s.metrics.SearchesTotal++
	s.metrics.mu.Unlock()

	type providerResult struct {
		name    string
		results []Result
		err     error
		latency time.Duration
	}

	g, gCtx := errgroup.WithContext(ctx)
	ch := make(chan providerResult, len(s.providers))

	for _, p := range s.providers {
		p := p
		g.Go(func() error {
			pStart := time.Now()
			breaker := s.breakers[p.Name()]

			var results []Result
			err := breaker.Call(gCtx, func(ctx context.Context) error {
				var innerErr error
				results, innerErr = p.Search(ctx, q)
				return innerErr
			})

			latency := time.Since(pStart)

			metrics.ObserveImportSearch(p.Name(), pStart)
			if err != nil {
				metrics.IncImportSearchError(p.Name())
			}
			metrics.SetCircuitBreakerState(p.Name(), int(breaker.State()))

			s.metrics.mu.Lock()
			s.metrics.ProviderLatency[p.Name()] = latency
			if err != nil {
				s.metrics.ProviderErrors[p.Name()]++
			} else {
				s.metrics.ProviderSuccess[p.Name()]++
			}
			s.metrics.mu.Unlock()

			select {
			case ch <- providerResult{name: p.Name(), results: results, err: err, latency: latency}:
			case <-gCtx.Done():
			}
			return nil
		})
	}

	_ = g.Wait()
	close(ch)

	var allResults []Result
	for pr := range ch {
		if pr.err != nil {
			s.logger.Warn("provider search failed",
				zap.String("provider", pr.name),
				zap.Duration("latency", pr.latency),
				zap.Error(pr.err),
			)
			continue
		}
		allResults = append(allResults, pr.results...)
	}

	if len(allResults) == 0 {
		s.logger.Warn("all providers returned no results",
			zap.String("query", rawQuery),
			zap.Duration("elapsed", time.Since(start)),
		)
		return []Result{}, nil
	}

	deduped := Deduplicate(allResults)

	for i, r := range deduped {
		score := fuzzy.CombinedScore(
			r.Title, r.Artist, r.Album, r.Duration,
			q.Title, q.Artist, "", 0,
		)
		deduped[i].Score = score
	}

	Rank(deduped)

	elapsed := time.Since(start)
	s.logger.Info("search completed",
		zap.String("query", rawQuery),
		zap.Int("raw_results", len(allResults)),
		zap.Int("deduped", len(deduped)),
		zap.Duration("elapsed", elapsed),
	)

	if len(deduped) > 20 {
		deduped = deduped[:20]
	}

	if cacheData, err := json.Marshal(deduped); err == nil {
		if setErr := s.cache.Set(ctx, cacheKey, cacheData); setErr != nil {
			s.logger.Warn("failed to cache search results", zap.Error(setErr))
		}
	}

	return deduped, nil
}

func splitQuery(raw string) []string {
	parts := strings.SplitN(raw, " - ", 2)
	if len(parts) == 2 {
		return parts
	}
	parts = strings.SplitN(raw, " by ", 2)
	if len(parts) == 2 {
		return parts
	}
	parts = strings.SplitN(raw, " from ", 2)
	if len(parts) == 2 {
		return parts
	}
	return []string{raw}
}

func normalizeCacheKey(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}
