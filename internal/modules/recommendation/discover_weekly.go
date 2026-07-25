package recommendation

import (
	"context"
	"fmt"
	"sort"
	"time"
)

type DiscoverWeeklyPlaylist struct {
	UserID      string    `json:"user_id" db:"user_id"`
	TrackIDs    []string  `json:"track_ids" db:"track_ids"`
	GeneratedAt time.Time `json:"generated_at" db:"generated_at"`
	WeekOf      time.Time `json:"week_of" db:"week_of"`
}

type DiscoverWeeklyService struct {
	profileService *TasteProfileService
	mlClient       *SimilarityMLClient
	repo           Repository
	dwRepo         DiscoverWeeklyRepository
}

func NewDiscoverWeeklyService(profileService *TasteProfileService, mlClient *SimilarityMLClient, repo Repository, dwRepo DiscoverWeeklyRepository) *DiscoverWeeklyService {
	return &DiscoverWeeklyService{profileService: profileService, mlClient: mlClient, repo: repo, dwRepo: dwRepo}
}

// GetOrGenerateWeeklyPlaylist returns the current week's playlist, generating
// a new one if none exists or if the week has rolled over.
func (s *DiscoverWeeklyService) GetOrGenerateWeeklyPlaylist(ctx context.Context, userID string) (*DiscoverWeeklyPlaylist, []TrackItem, error) {
	weekOf := currentWeekStart()

	// Check for existing playlist this week
	existing, err := s.dwRepo.GetByWeek(ctx, userID, weekOf)
	if err == nil && existing != nil {
		tracks, err := s.repo.GetTracksByIDs(ctx, existing.TrackIDs)
		if err != nil {
			return nil, nil, fmt.Errorf("get existing playlist tracks: %w", err)
		}
		return existing, tracks, nil
	}

	// Generate new playlist
	profile, err := s.profileService.GetOrInitProfile(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get profile: %w", err)
	}
	if len(profile.SeedTrackIDs) == 0 {
		return nil, nil, fmt.Errorf("insufficient taste profile data")
	}

	allHeard, err := s.repo.GetAllUserPlayedTrackIDs(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get played tracks: %w", err)
	}
	alreadyHeard := make(map[string]struct{}, len(allHeard))
	for _, id := range allHeard {
		alreadyHeard[id] = struct{}{}
	}

	likedIDs, _ := s.repo.GetLikedTrackIDs(ctx, userID, 200)
	for _, id := range likedIDs {
		alreadyHeard[id] = struct{}{}
	}

	desiredSize := 30
	candidates := s.fetchNoveltyCandidates(ctx, profile.SeedTrackIDs, desiredSize, alreadyHeard)

	// Add adjacent-genre tracks (~20% of total)
	adjacentCount := int(float64(desiredSize) * 0.2)
	if adjacentCount > 0 && len(profile.TopGenreIDs) > 0 {
		adjacent := s.getAdjacentGenreTracks(ctx, profile.TopGenreIDs, adjacentCount, alreadyHeard)
		candidates = append(candidates, adjacent...)
	}

	// Deduplicate and apply exclusion
	trackIDs := s.deduplicateAndLimit(candidates, alreadyHeard, desiredSize)

	if len(trackIDs) == 0 {
		return nil, nil, fmt.Errorf("no tracks could be generated")
	}

	playlist := &DiscoverWeeklyPlaylist{
		UserID:      userID,
		TrackIDs:    trackIDs,
		GeneratedAt: time.Now(),
		WeekOf:      weekOf,
	}

	if err := s.dwRepo.Save(ctx, playlist); err != nil {
		return nil, nil, fmt.Errorf("save weekly playlist: %w", err)
	}

	tracks, err := s.repo.GetTracksByIDs(ctx, trackIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("get playlist tracks: %w", err)
	}

	return playlist, tracks, nil
}

func (s *DiscoverWeeklyService) fetchNoveltyCandidates(ctx context.Context, seedIDs []string, desiredSize int, exclude map[string]struct{}) []TrackItem {
	candidateSize := desiredSize * 3

	if s.mlClient != nil {
		excludeSlice := make([]string, 0, len(exclude))
		for id := range exclude {
			excludeSlice = append(excludeSlice, id)
		}

		results, err := s.mlClient.GetSimilarToTrackSet(ctx, seedIDs, candidateSize, excludeSlice)
		if err == nil && len(results) > 0 {
			// Novelty filter: sort by similarity, skip top 25% (near-duplicates),
			// take from the remaining middle band.
			// If the candidate pool is too small (< desiredSize/2), skip the
			// cutoff to avoid shrinking the pool to nothing.
			sort.Slice(results, func(i, j int) bool {
				return results[i].SimilarityScore > results[j].SimilarityScore
			})

			useAll := len(results) < desiredSize/2
			var cutoff int
			if !useAll {
				cutoff = int(float64(len(results)) * 0.25)
				if cutoff >= len(results) {
					cutoff = 0
				}
			}

			nCandidates := len(results) - cutoff
			if nCandidates > desiredSize {
				nCandidates = desiredSize
			}

			middleBand := results[cutoff : cutoff+nCandidates]
			ids := make([]string, len(middleBand))
			for i, r := range middleBand {
				ids[i] = r.TrackID
			}

			tracks, err := s.repo.GetTracksByIDs(ctx, ids)
			if err != nil {
				return nil
			}

			idOrder := make(map[string]int, len(ids))
			for i, id := range ids {
				idOrder[id] = int(i)
			}
			sort.Slice(tracks, func(i, j int) bool {
				oi, oki := idOrder[tracks[i].ID]
				oj, okj := idOrder[tracks[j].ID]
				if !oki {
					return false
				}
				if !okj {
					return true
				}
				return oi < oj
			})

			return tracks
		}
	}

	// Fallback: meta similarity with novelty adjustment
	return s.similarToAnySeed(ctx, seedIDs, candidateSize, exclude)
}

func (s *DiscoverWeeklyService) similarToAnySeed(ctx context.Context, seedIDs []string, limit int, exclude map[string]struct{}) []TrackItem {
	skipSeeds := len(seedIDs) / 4
	if skipSeeds >= len(seedIDs) {
		skipSeeds = 0
	}
	// Skip the most-used seed to encourage novelty
	novelSeeds := seedIDs[skipSeeds:]

	seen := make(map[string]struct{})
	var results []TrackItem

	for _, seedID := range novelSeeds {
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

	return results
}

func (s *DiscoverWeeklyService) getAdjacentGenreTracks(ctx context.Context, topGenreIDs []string, count int, exclude map[string]struct{}) []TrackItem {
	if count <= 0 {
		return nil
	}

	topGenreSet := make(map[string]struct{}, len(topGenreIDs))
	for _, id := range topGenreIDs {
		topGenreSet[id] = struct{}{}
	}

	// Get globally popular genres
	popularGenres, err := s.repo.GetPopularGenreIDs(ctx, 20)
	if err != nil || len(popularGenres) == 0 {
		return nil
	}

	var adjacentGenres []string
	for _, gid := range popularGenres {
		if _, inTop := topGenreSet[gid]; !inTop {
			adjacentGenres = append(adjacentGenres, gid)
			if len(adjacentGenres) >= count {
				break
			}
		}
	}

	if len(adjacentGenres) == 0 {
		return nil
	}

	genreNames, err := s.repo.GetGenreNames(ctx, adjacentGenres)
	if err != nil || len(genreNames) == 0 {
		return nil
	}

	tracks, err := s.repo.GetTracksFromGenres(ctx, genreNames, count*2)
	if err != nil {
		return nil
	}

	var filtered []TrackItem
	for _, t := range tracks {
		if _, ok := exclude[t.ID]; ok {
			continue
		}
		filtered = append(filtered, t)
		if len(filtered) >= count {
			break
		}
	}

	return filtered
}

func (s *DiscoverWeeklyService) deduplicateAndLimit(items []TrackItem, exclude map[string]struct{}, limit int) []string {
	seen := make(map[string]struct{})
	// Also add exclude items as seen
	for id := range exclude {
		seen[id] = struct{}{}
	}

	var result []string
	for _, t := range items {
		if _, ok := seen[t.ID]; ok {
			continue
		}
		seen[t.ID] = struct{}{}
		result = append(result, t.ID)
		if len(result) >= limit {
			break
		}
	}
	return result
}

func currentWeekStart() time.Time {
	now := time.Now().UTC()
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	start := now.AddDate(0, 0, -int(weekday-time.Monday))
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
}
