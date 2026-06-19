package recommendation

import (
	"context"
	"errors"
)

var ErrInvalidLimit = errors.New("invalid limit")

type Service interface {
	PopularTracks(ctx context.Context, limit int) ([]TrackItem, error)
	BestTracks(ctx context.Context, limit int) ([]TrackItem, error)
	RecentTracks(ctx context.Context, userID string, limit int) ([]TrackItem, error)
	SimilarTracks(ctx context.Context, trackID string, limit int) ([]TrackItem, error)
	TracksByArtist(ctx context.Context, artistID string, limit int) ([]TrackItem, error)
	TracksByGenre(ctx context.Context, genre string, limit int) ([]TrackItem, error)
	ForYou(ctx context.Context, userID string, limit int) ([]TrackItem, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func normalizeLimit(limit int) (int, error) {
	if limit == 0 {
		return 20, nil
	}
	if limit < 0 || limit > 100 {
		return 0, ErrInvalidLimit
	}
	return limit, nil
}

func (s *service) PopularTracks(ctx context.Context, limit int) ([]TrackItem, error) {
	limit, err := normalizeLimit(limit)
	if err != nil {
		return nil, err
	}
	return s.repo.GetPopularTracks(ctx, limit)
}

func (s *service) BestTracks(ctx context.Context, limit int) ([]TrackItem, error) {
	limit, err := normalizeLimit(limit)
	if err != nil {
		return nil, err
	}
	return s.repo.GetBestTracks(ctx, limit)
}

func (s *service) RecentTracks(ctx context.Context, userID string, limit int) ([]TrackItem, error) {
	limit, err := normalizeLimit(limit)
	if err != nil {
		return nil, err
	}
	return s.repo.GetRecentTracks(ctx, userID, limit)
}

func (s *service) SimilarTracks(ctx context.Context, trackID string, limit int) ([]TrackItem, error) {
	limit, err := normalizeLimit(limit)
	if err != nil {
		return nil, err
	}

	meta, err := s.repo.GetTrackMeta(ctx, trackID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetSimilarTracksByMeta(ctx, trackID, meta.ArtistID, meta.AlbumID, meta.Genre, limit)
}

func (s *service) TracksByArtist(ctx context.Context, artistID string, limit int) ([]TrackItem, error) {
	limit, err := normalizeLimit(limit)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTracksByArtist(ctx, artistID, limit)
}

func (s *service) TracksByGenre(ctx context.Context, genre string, limit int) ([]TrackItem, error) {
	limit, err := normalizeLimit(limit)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTracksByGenre(ctx, genre, limit)
}

func (s *service) ForYou(ctx context.Context, userID string, limit int) ([]TrackItem, error) {
	limit, err := normalizeLimit(limit)
	if err != nil {
		return nil, err
	}

	candidates := make(map[string]*scoredTrack)

	addItems := func(items []TrackItem, weight float64) {
		for _, item := range items {
			entry, exists := candidates[item.ID]
			if !exists {
				candidates[item.ID] = &scoredTrack{
					Item:  item,
					Score: weight,
				}
				continue
			}
			entry.Score += weight
		}
	}

	// 1. Followed artist
	followedArtistIDs, err := s.repo.GetFollowedArtistIDs(ctx, userID, 10)
	if err != nil {
		return nil, err
	}
	if len(followedArtistIDs) > 0 {
		items, err := s.repo.GetTracksFromArtists(ctx, followedArtistIDs, 100)
		if err != nil {
			return nil, err
		}
		addItems(items, 5)
	}

	// 2. Top artist from listening history
	topArtistIDs, err := s.repo.GetTopArtistIDs(ctx, userID, 10)
	if err != nil {
		return nil, err
	}
	if len(topArtistIDs) > 0 {
		items, err := s.repo.GetTracksFromArtists(ctx, topArtistIDs, 100)
		if err != nil {
			return nil, err
		}
		addItems(items, 4)
	}

	// 3. Top genres from history
	topGenres, err := s.repo.GetTopGenres(ctx, userID, 10)
	if err != nil {
		return nil, err
	}
	if len(topGenres) > 0 {
		items, err := s.repo.GetTracksFromGenres(ctx, topGenres, 100)
		if err != nil {
			return nil, err
		}
		addItems(items, 3)
	}

	// 4. Similar to liked tracks
	likedTrackIDs, err := s.repo.GetLikedTrackIDs(ctx, userID, 10)
	if err != nil {
		return nil, err
	}
	for _, likedTrackID := range likedTrackIDs {
		meta, err := s.repo.GetTrackMeta(ctx, likedTrackID)
		if err != nil {
			if errors.Is(err, ErrTrackNotFound) {
				continue
			}
			return nil, err
		}

		items, err := s.repo.GetSimilarTracksByMeta(ctx, likedTrackID, meta.ArtistID, meta.AlbumID, meta.Genre, 20)
		if err != nil {
			return nil, err
		}
		addItems(items, 4)
	}

	// 5. Global popularity boost
	popularIDs, err := s.repo.GetPopularTrackIDs(ctx, 20)
	if err != nil {
		return nil, err
	}
	if len(popularIDs) > 0 {
		items, err := s.repo.GetTracksByIDs(ctx, popularIDs)
		if err != nil {
			return nil, err
		}
		addItems(items, 1)
	}

	if len(candidates) == 0 {
		return s.PopularTracks(ctx, limit)
	}

	// Optionally exclude already liked tracks from results
	likedSet := map[string]struct{}{}
	for _, id := range likedTrackIDs {
		likedSet[id] = struct{}{}
	}

	scored := sortScoredTracks(candidates)

	result := make([]TrackItem, 0, limit)
	for _, item := range scored {
		if _, exists := likedSet[item.ID]; exists {
			continue
		}
		result = append(result, item)
		if len(result) >= limit {
			break
		}
	}

	if len(result) == 0 {
		return s.PopularTracks(ctx, limit)
	}

	return result, nil
}
