package recommendation

import (
	"context"
	"fmt"
	"sort"
)

type HomeRecommendationService struct {
	profileService *TasteProfileService
	mlClient       *SimilarityMLClient
	repo           Repository
}

func NewHomeRecommendationService(profileService *TasteProfileService, mlClient *SimilarityMLClient, repo Repository) *HomeRecommendationService {
	return &HomeRecommendationService{profileService: profileService, mlClient: mlClient, repo: repo}
}

func (s *HomeRecommendationService) GetPersonalizedTracks(ctx context.Context, userID string, limit int) ([]TrackItem, error) {
	profile, err := s.profileService.GetOrInitProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	excludeIDs := s.collectExcludeIDs(ctx, userID)

	if len(profile.SeedTrackIDs) > 0 {
		return s.fromSeedTracks(ctx, profile.SeedTrackIDs, limit, excludeIDs)
	}

	if len(profile.OnboardingGenreIDs) > 0 {
		return s.fromOnboardingGenres(ctx, profile.OnboardingGenreIDs, limit, excludeIDs)
	}

	return []TrackItem{}, nil
}

func (s *HomeRecommendationService) collectExcludeIDs(ctx context.Context, userID string) []string {
	seen := make(map[string]struct{})

	recent, err := s.repo.GetRecentTracks(ctx, userID, 20)
	if err == nil {
		for _, t := range recent {
			seen[t.ID] = struct{}{}
		}
	}

	liked, err := s.repo.GetLikedTrackIDs(ctx, userID, 100)
	if err == nil {
		for _, id := range liked {
			seen[id] = struct{}{}
		}
	}

	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

func (s *HomeRecommendationService) fromSeedTracks(ctx context.Context, seedIDs []string, limit int, excludeIDs []string) ([]TrackItem, error) {
	if s.mlClient != nil {
		results, err := s.mlClient.GetSimilarToTrackSet(ctx, seedIDs, limit, excludeIDs)
		if err == nil && len(results) > 0 {
			ids := make([]string, len(results))
			for i, r := range results {
				ids[i] = r.TrackID
			}
			tracks, err := s.repo.GetTracksByIDs(ctx, ids)
			if err != nil {
				return nil, err
			}
			idOrder := make(map[string]int, len(ids))
			for i, id := range ids {
				idOrder[id] = i
			}
			sortByOrder(tracks, idOrder)
			if len(tracks) > limit {
				tracks = tracks[:limit]
			}
			return tracks, nil
		}
	}

	return s.similarToAnySeed(ctx, seedIDs, limit, excludeIDs)
}

func (s *HomeRecommendationService) similarToAnySeed(ctx context.Context, seedIDs []string, limit int, excludeIDs []string) ([]TrackItem, error) {
	exclude := make(map[string]struct{}, len(excludeIDs))
	for _, id := range excludeIDs {
		exclude[id] = struct{}{}
	}

	seen := make(map[string]struct{})
	var results []TrackItem

	for _, seedID := range seedIDs {
		if len(results) >= limit {
			break
		}

		meta, err := s.repo.GetTrackMeta(ctx, seedID)
		if err != nil {
			continue
		}

		tracks, err := s.repo.GetSimilarTracksByMeta(ctx, seedID, meta.ArtistID, meta.AlbumID, meta.Genre, limit)
		if err != nil {
			continue
		}

		for _, t := range tracks {
			if _, ok := exclude[t.ID]; ok {
				continue
			}
			if _, ok := seen[t.ID]; ok {
				continue
			}
			seen[t.ID] = struct{}{}
			results = append(results, t)
			if len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

func (s *HomeRecommendationService) fromOnboardingGenres(ctx context.Context, genreIDs []string, limit int, excludeIDs []string) ([]TrackItem, error) {
	genreNames, err := s.repo.GetGenreNames(ctx, genreIDs)
	if err != nil || len(genreNames) == 0 {
		return []TrackItem{}, nil
	}

	tracks, err := s.repo.GetTracksFromGenres(ctx, genreNames, limit*2)
	if err != nil {
		return []TrackItem{}, nil
	}

	exclude := make(map[string]struct{}, len(excludeIDs))
	for _, id := range excludeIDs {
		exclude[id] = struct{}{}
	}

	var filtered []TrackItem
	for _, t := range tracks {
		if _, ok := exclude[t.ID]; ok {
			continue
		}
		filtered = append(filtered, t)
		if len(filtered) >= limit {
			break
		}
	}

	return filtered, nil
}

func sortByOrder(items []TrackItem, order map[string]int) {
	sort.Slice(items, func(i, j int) bool {
		return order[items[i].ID] < order[items[j].ID]
	})
}
