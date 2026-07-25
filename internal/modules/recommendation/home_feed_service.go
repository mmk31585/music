package recommendation

import (
	"context"
	"fmt"
	"sort"
)

type HomeFeedService struct {
	repo       Repository
	profileSvc *TasteProfileService
	mlClient   *SimilarityMLClient
	svc        Service
}

func NewHomeFeedService(repo Repository, profileSvc *TasteProfileService, mlClient *SimilarityMLClient, svc Service) *HomeFeedService {
	return &HomeFeedService{repo: repo, profileSvc: profileSvc, mlClient: mlClient, svc: svc}
}

func (h *HomeFeedService) BuildFeed(ctx context.Context, userID string) (*HomeFeedResponse, error) {
	profile, _ := h.profileSvc.GetOrInitProfile(ctx, userID)

	var sections []HomeFeedSection

	// 1. Recently played
	if recent, err := h.repo.GetRecentTracks(ctx, userID, 10); err == nil && len(recent) > 0 {
		sections = append(sections, HomeFeedSection{
			ID: "recently_played", Title: "Recently Played", Subtitle: "Continue where you left off",
			Type: "track", Items: recent,
		})
	}

	// 2. Trending
	if popular, err := h.repo.GetPopularTracks(ctx, 10); err == nil && len(popular) > 0 {
		sections = append(sections, HomeFeedSection{
			ID: "trending", Title: "Trending", Subtitle: "Popular right now",
			Type: "track", Items: popular,
		})
	}

	// 3. For You
	if forYou, err := h.svc.ForYou(ctx, userID, 10); err == nil && len(forYou) > 0 {
		sections = append(sections, HomeFeedSection{
			ID: "for_you", Title: "For You", Subtitle: "Personalized picks",
			Type: "track", Items: forYou,
		})
	}

	// 4. From your favorite artists (taste profile TopArtistIDs)
	if len(profile.TopArtistIDs) > 0 {
		if artistTracks, err := h.repo.GetTracksFromArtists(ctx, profile.TopArtistIDs, 10); err == nil && len(artistTracks) > 0 {
			sections = append(sections, HomeFeedSection{
				ID: "from_your_artists", Title: "More from your favorite artists", Subtitle: "Based on your listening",
				Type: "track", Items: artistTracks,
			})
		}
	}

	// 5. Because you listened to X (seed track similarity)
	maxSeeds := 2
	for i, seedID := range profile.SeedTrackIDs {
		if i >= maxSeeds {
			break
		}

		seedItems, err := h.repo.GetTracksByIDs(ctx, []string{seedID})
		if err != nil || len(seedItems) == 0 {
			continue
		}
		seedTrack := seedItems[0]

		var similar []TrackItem
		if h.mlClient != nil {
			if results, mlErr := h.mlClient.GetSimilarTracks(ctx, seedID, 10, nil); mlErr == nil && len(results) > 0 {
				ids := make([]string, len(results))
				for j, r := range results {
					ids[j] = r.TrackID
				}
				if tracks, err := h.repo.GetTracksByIDs(ctx, ids); err == nil {
					idOrder := make(map[string]int, len(ids))
					for j, id := range ids {
						idOrder[id] = j
					}
					sort.Slice(tracks, func(a, b int) bool {
						return idOrder[tracks[a].ID] < idOrder[tracks[b].ID]
					})
					similar = tracks
				}
			}
		}

		if len(similar) == 0 {
			if meta, metaErr := h.repo.GetTrackMeta(ctx, seedID); metaErr == nil {
				if tracks, simErr := h.repo.GetSimilarTracksByMeta(ctx, seedID, meta.ArtistID, meta.AlbumID, meta.Genre, 10); simErr == nil {
					similar = tracks
				}
			}
		}

		if len(similar) > 0 {
			shortID := seedID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}
			sections = append(sections, HomeFeedSection{
				ID:        fmt.Sprintf("because_of_%s", shortID),
				Title:     fmt.Sprintf("Because you listened to %s", seedTrack.Title),
				Type:      "track",
				Items:     similar,
				SeedTrack: &seedTrack,
			})
		}
	}

	// 6. Popular in your genres
	if genres, err := h.repo.GetTopGenres(ctx, userID, 5); err == nil && len(genres) > 0 {
		if genreTracks, err := h.repo.GetTracksFromGenres(ctx, genres, 10); err == nil && len(genreTracks) > 0 {
			sections = append(sections, HomeFeedSection{
				ID: "your_genres", Title: "Popular in your genres",
				Type: "track", Items: genreTracks,
			})
		}
	}

	if sections == nil {
		sections = []HomeFeedSection{}
	}

	return &HomeFeedResponse{Sections: sections}, nil
}
