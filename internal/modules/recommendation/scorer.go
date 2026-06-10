package recommendation

import (
	"context"
	"math"
	"sort"
	"time"
)

type Scorer struct {
	repo  Repository
	cache CacheStore
}

type ScoredTrack struct {
	Track TrackItem
	Score float64
}

func NewScorer(repo Repository, cache CacheStore) *Scorer {
	return &Scorer{repo: repo, cache: cache}
}

func (s *Scorer) Score(ctx context.Context, userID string, candidates []TrackItem, weights ScoringWeights) []ScoredTrack {
	if len(candidates) == 0 {
		return nil
	}

	userAffinities, _ := s.loadUserAffinities(ctx, userID)
	likedIDs, _ := s.repo.GetLikedTrackIDs(ctx, userID, 100)
	likedSet := make(map[string]struct{}, len(likedIDs))
	for _, id := range likedIDs {
		likedSet[id] = struct{}{}
	}

	scored := make([]ScoredTrack, 0, len(candidates))

	for _, track := range candidates {
		score := 0.0
		weight := 1.0

		if aff, ok := userAffinities[track.ID]; ok {
			weight += aff * weights.Affinity
		}

		if _, liked := likedSet[track.ID]; liked {
			weight += weights.Liked
		}

		if track.Score != nil {
			score += *track.Score * weights.Popularity
		}

		if track.DurationSeconds != nil && *track.DurationSeconds > 0 {
			dwellScore := math.Min(float64(*track.DurationSeconds)/180.0, 1.0)
			score += dwellScore * weights.Duration
		}

		if track.Genre != nil && *track.Genre != "" {
			score += weights.GenreBonus
		}

		score *= weight

		scored = append(scored, ScoredTrack{
			Track: track,
			Score: score,
		})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	return scored
}

func (s *Scorer) DiversityRerank(tracks []ScoredTrack, maxPerArtist int) []ScoredTrack {
	if len(tracks) <= maxPerArtist {
		return tracks
	}

	artistCount := make(map[string]int)
	result := make([]ScoredTrack, 0, len(tracks))

	for _, t := range tracks {
		if t.Track.ArtistID != nil {
			artistCount[*t.Track.ArtistID]++
			if artistCount[*t.Track.ArtistID] > maxPerArtist {
				continue
			}
		}
		result = append(result, t)
	}

	return result
}

type ScoringWeights struct {
	Affinity   float64
	Liked      float64
	Popularity float64
	Duration   float64
	GenreBonus float64
}

func DefaultScoringWeights() ScoringWeights {
	return ScoringWeights{
		Affinity:   1.5,
		Liked:      2.0,
		Popularity: 1.0,
		Duration:   0.5,
		GenreBonus: 0.3,
	}
}

func (s *Scorer) loadUserAffinities(ctx context.Context, userID string) (map[string]float64, error) {
	affinities := make(map[string]float64)

	rows, err := s.repo.GetUserAffinities(ctx, userID, "track", 500)
	if err != nil {
		return affinities, err
	}
	for _, aff := range rows {
		affinities[aff.TargetID] = aff.Score
	}
	return affinities, nil
}

type UserAffinity struct {
	TargetID string
	Score    float64
	Decay    float64
}

func (r *repository) GetUserAffinities(ctx context.Context, userID string, targetType string, limit int) ([]UserAffinity, error) {
	query := `
		SELECT target_id, score, decay
		FROM user_affinities
		WHERE user_id = $1 AND target_type = $2
		ORDER BY score DESC
		LIMIT $3
	`
	var affinities []UserAffinity
	if err := r.db.SelectContext(ctx, &affinities, query, userID, targetType, limit); err != nil {
		return nil, err
	}
	return affinities, nil
}

func (r *repository) GetSimilarTracksByCooccurrence(ctx context.Context, trackID string, limit int) ([]TrackItem, error) {
	query := `
		SELECT
			t.id, t.title, t.artist_id, a.name AS artist_name,
			t.album_id, al.title AS album_title,
			(SELECT g.name FROM track_genres tg JOIN genres g ON g.id = tg.genre_id WHERE tg.track_id = t.id LIMIT 1) AS genre,
			t.cover_url, t.audio_url, t.duration_seconds,
			ts.similarity AS score
		FROM track_similarities ts
		JOIN tracks t ON t.id = ts.similar_track_id
		LEFT JOIN artists a ON a.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		WHERE ts.track_id = $1 AND ts.method = 'cooccurrence'
		ORDER BY ts.similarity DESC
		LIMIT $2
	`
	var items []TrackItem
	if err := r.db.SelectContext(ctx, &items, query, trackID, limit); err != nil {
		return nil, err
	}
	return items, nil
}

type CacheStore interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Exists(ctx context.Context, keys ...string) (bool, error)
}
