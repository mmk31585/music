package recommendation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"music/internal/modules/history"
)

const (
	debounceDuration  = 1 * time.Hour
	debounceBatchSize = 10
)

type TasteProfileService struct {
	repo        TasteProfileRepository
	historyRepo *history.Repository
	mlClient    *SimilarityMLClient // optional; if nil, avg_tempo_preference won't be populated
}

func NewTasteProfileService(repo TasteProfileRepository, historyRepo *history.Repository, mlClient *SimilarityMLClient) *TasteProfileService {
	return &TasteProfileService{repo: repo, historyRepo: historyRepo, mlClient: mlClient}
}

// RecomputeProfile rebuilds the user's taste profile from recent signal
// history using exponential decay weighting.
//
// Debouncing: skips recompute if profile was computed less than 1 hour ago
// AND fewer than 10 new signals have accumulated.
func (s *TasteProfileService) RecomputeProfile(ctx context.Context, userID string) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	existing, err := s.repo.GetProfile(ctx, userID)
	if err == nil && existing != nil && s.historyRepo != nil {
		elapsed := time.Since(existing.LastComputedAt)
		if elapsed < debounceDuration {
			signals, err := s.historyRepo.ListSignalsByUser(ctx, parsedUserID, MaxSignalAgeDays)
			if err == nil {
				newCount := 0
				for _, sig := range signals {
					if sig.PlayedAt.After(existing.LastComputedAt) {
						newCount++
					}
				}
				if newCount < debounceBatchSize {
					return nil
				}
			}
		}
	}

	since := time.Now().AddDate(0, 0, -MaxSignalAgeDays)
	rows, err := s.repo.GetSignalRows(ctx, userID, since)
	if err != nil {
		return fmt.Errorf("get signal rows: %w", err)
	}

	now := time.Now()
	genreScores := make(map[string]float64)
	artistScores := make(map[string]float64)
	seenReplays := make(map[string]bool)
	var positiveTrackIDs []string

	for _, row := range rows {
		weight := signalBaseWeight(row.SignalType, row.IsExplicitLike)
		if weight == 0 {
			continue
		}

		daysAgo := now.Sub(row.PlayedAt).Hours() / 24
		decay := math.Exp(-math.Ln2 * daysAgo / ProfileDecayHalfLifeDays)
		effectiveWeight := weight * decay

		if row.ArtistID != nil && *row.ArtistID != "" {
			artistScores[*row.ArtistID] += effectiveWeight
		}
		if row.GenreID != nil && *row.GenreID != "" {
			for _, g := range strings.Split(*row.GenreID, ",") {
				g = strings.TrimSpace(g)
				if g != "" {
					genreScores[g] += effectiveWeight
				}
			}
		}

		if weight > 0 && !seenReplays[row.TrackID] {
			seenReplays[row.TrackID] = true
			positiveTrackIDs = append(positiveTrackIDs, row.TrackID)
		}
	}

	seedTrackIDs := uniqueFirstN(positiveTrackIDs, MaxSeedTracks)

	// Populate avg_tempo_preference from ml-service if available.
	var avgTempo *float64
	if s.mlClient != nil && len(seedTrackIDs) > 0 {
		tempos, err := s.mlClient.GetBulkTempo(ctx, seedTrackIDs)
		if err != nil {
			log.Printf("failed to fetch bulk tempos for profile (non-fatal): %v", err)
		} else {
			var sum float64
			var count int
			for _, bpm := range tempos {
				if bpm != nil {
					sum += *bpm
					count++
				}
			}
			if count > 0 {
				v := sum / float64(count)
				avgTempo = &v
			}
		}
	}

	profile := &TasteProfile{
		UserID:             userID,
		TopGenreIDs:        topNByScore(genreScores, TopGenreLimit),
		TopArtistIDs:       topNByScore(artistScores, TopArtistLimit),
		SeedTrackIDs:       seedTrackIDs,
		LastComputedAt:     now,
		AvgTempoPreference: avgTempo,
	}

	// Preserve onboarding genres
	if existing != nil && len(existing.OnboardingGenreIDs) > 0 {
		profile.OnboardingGenreIDs = existing.OnboardingGenreIDs
	}

	if err := s.repo.UpsertProfile(ctx, profile); err != nil {
		return fmt.Errorf("upsert profile: %w", err)
	}

	return nil
}

// GetOrInitProfile returns the user's taste profile. If none exists,
// returns a minimal profile seeded with defaults.
func (s *TasteProfileService) GetOrInitProfile(ctx context.Context, userID string) (*TasteProfile, error) {
	profile, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &TasteProfile{UserID: userID}, nil
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return profile, nil
}

// SetOnboardingGenres stores the user's genre preferences from onboarding.
func (s *TasteProfileService) SetOnboardingGenres(ctx context.Context, userID string, genreIDs []string) error {
	if len(genreIDs) == 0 {
		return errors.New("at least one genre required")
	}
	return s.repo.SetOnboardingGenres(ctx, userID, genreIDs)
}

// ── Helpers ────────────────────────────────────────────────────────────

func signalBaseWeight(signalType string, isExplicitLike bool) float64 {
	if isExplicitLike {
		return SignalWeightExplicitLike
	}
	switch signalType {
	case "complete_positive":
		return SignalWeightCompletePositive
	case "replay_strong":
		return SignalWeightReplayStrong
	case "skip_negative":
		return SignalWeightSkipNegative
	case "partial_neutral":
		return SignalWeightPartialNeutral
	default:
		return 0
	}
}

func topNByScore(scores map[string]float64, n int) []string {
	type kv struct {
		Key   string
		Value float64
	}
	var sorted []kv
	for k, v := range scores {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	if len(sorted) > n {
		sorted = sorted[:n]
	}
	result := make([]string, len(sorted))
	for i, kv := range sorted {
		result[i] = kv.Key
	}
	return result
}

func uniqueFirstN(items []string, n int) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
		if len(result) >= n {
			break
		}
	}
	return result
}
