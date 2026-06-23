package recommendation

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type TrackStat struct {
	TrackID    string `json:"track_id"`
	Title      string `json:"title"`
	ArtistName string `json:"artist_name"`
	PlayCount  int    `json:"play_count"`
}

type ArtistStat struct {
	ArtistID   string `json:"artist_id"`
	ArtistName string `json:"artist_name"`
	PlayCount  int    `json:"play_count"`
}

type GenreStat struct {
	GenreName string `json:"genre_name"`
	PlayCount int    `json:"play_count"`
}

type ListeningStats struct {
	PeriodLabel          string       `json:"period_label"`
	TotalMinutesListened int          `json:"total_minutes_listened"`
	TotalTracksPlayed    int          `json:"total_tracks_played"`
	TopTracks            []TrackStat  `json:"top_tracks"`
	TopArtists           []ArtistStat `json:"top_artists"`
	TopGenres            []GenreStat  `json:"top_genres"`
	UniqueArtistsCount   int          `json:"unique_artists_count"`
	LongestStreak        int          `json:"longest_listening_streak_days"`
	DiscoveryScore       float64      `json:"discovery_score"`
}

type ListeningStatsService struct {
	db *sqlx.DB
}

func NewListeningStatsService(db *sqlx.DB) *ListeningStatsService {
	return &ListeningStatsService{db: db}
}

func (s *ListeningStatsService) GetStats(ctx context.Context, userID string, period string) (*ListeningStats, error) {
	since, label, err := periodBounds(period)
	if err != nil {
		return nil, err
	}

	var totalMinutes sql.NullFloat64
	err = s.db.GetContext(ctx, &totalMinutes, `
		SELECT COALESCE(SUM(COALESCE(played_duration_ms, duration * 1000)), 0) / 60000.0
		FROM listening_history
		WHERE user_id = $1 AND played_at >= $2
	`, userID, since)
	if err != nil {
		return nil, fmt.Errorf("total minutes: %w", err)
	}

	var totalTracks int
	err = s.db.GetContext(ctx, &totalTracks, `
		SELECT COUNT(*)
		FROM listening_history
		WHERE user_id = $1 AND played_at >= $2
	`, userID, since)
	if err != nil {
		return nil, fmt.Errorf("total tracks: %w", err)
	}

	var topTracks []TrackStat
	err = s.db.SelectContext(ctx, &topTracks, `
		SELECT lh.track_id::text, t.title AS title,
			COALESCE(a.name, '') AS artist_name,
			COUNT(*) AS play_count
		FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		LEFT JOIN artists a ON a.id = t.artist_id
		WHERE lh.user_id = $1 AND lh.played_at >= $2
		GROUP BY lh.track_id, t.title, a.name
		ORDER BY COUNT(*) DESC
		LIMIT 10
	`, userID, since)
	if err != nil {
		return nil, fmt.Errorf("top tracks: %w", err)
	}
	if topTracks == nil {
		topTracks = []TrackStat{}
	}

	var topArtists []ArtistStat
	err = s.db.SelectContext(ctx, &topArtists, `
		SELECT a.id::text AS artist_id, a.name AS artist_name, COUNT(*) AS play_count
		FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		JOIN artists a ON a.id = t.artist_id
		WHERE lh.user_id = $1 AND lh.played_at >= $2
		GROUP BY a.id, a.name
		ORDER BY COUNT(*) DESC
		LIMIT 10
	`, userID, since)
	if err != nil {
		return nil, fmt.Errorf("top artists: %w", err)
	}
	if topArtists == nil {
		topArtists = []ArtistStat{}
	}

	var topGenres []GenreStat
	err = s.db.SelectContext(ctx, &topGenres, `
		SELECT g.name AS genre_name, COUNT(*) AS play_count
		FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		JOIN track_genres tg ON tg.track_id = t.id
		JOIN genres g ON g.id = tg.genre_id
		WHERE lh.user_id = $1 AND lh.played_at >= $2
		GROUP BY g.name
		ORDER BY COUNT(*) DESC
		LIMIT 10
	`, userID, since)
	if err != nil {
		return nil, fmt.Errorf("top genres: %w", err)
	}
	if topGenres == nil {
		topGenres = []GenreStat{}
	}

	var uniqueArtists int
	err = s.db.GetContext(ctx, &uniqueArtists, `
		SELECT COUNT(DISTINCT t.artist_id)
		FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		WHERE lh.user_id = $1 AND lh.played_at >= $2 AND t.artist_id IS NOT NULL
	`, userID, since)
	if err != nil {
		return nil, fmt.Errorf("unique artists: %w", err)
	}

	// Longest streak: consecutive days with at least one play
	longestStreak, err := s.calcLongestStreak(ctx, userID, since)
	if err != nil {
		longestStreak = 0
	}

	// Discovery score: ratio of unique tracks to total plays
	discoveryScore := s.calcDiscoveryScore(ctx, userID, since, totalTracks)

	minutes := 0
	if totalMinutes.Valid {
		minutes = int(totalMinutes.Float64)
	}

	return &ListeningStats{
		PeriodLabel:          label,
		TotalMinutesListened: minutes,
		TotalTracksPlayed:    totalTracks,
		TopTracks:            topTracks,
		TopArtists:           topArtists,
		TopGenres:            topGenres,
		UniqueArtistsCount:   uniqueArtists,
		LongestStreak:        longestStreak,
		DiscoveryScore:       discoveryScore,
	}, nil
}

func (s *ListeningStatsService) calcLongestStreak(ctx context.Context, userID string, since time.Time) (int, error) {
	var dates []time.Time
	err := s.db.SelectContext(ctx, &dates, `
		SELECT DISTINCT played_at::date
		FROM listening_history
		WHERE user_id = $1 AND played_at >= $2
		ORDER BY played_at::date
	`, userID, since)
	if err != nil || len(dates) == 0 {
		return 0, err
	}

	longest := 1
	current := 1
	for i := 1; i < len(dates); i++ {
		diff := dates[i].Sub(dates[i-1]).Hours() / 24
		if diff <= 1.5 { // within ~36h (accounting for timezone edge cases)
			current++
			if current > longest {
				longest = current
			}
		} else {
			current = 1
		}
	}
	return longest, nil
}

func (s *ListeningStatsService) calcDiscoveryScore(ctx context.Context, userID string, since time.Time, totalPlays int) float64 {
	if totalPlays == 0 {
		return 0
	}

	var uniqueTracks int
	err := s.db.GetContext(ctx, &uniqueTracks, `
		SELECT COUNT(DISTINCT track_id)
		FROM listening_history
		WHERE user_id = $1 AND played_at >= $2
	`, userID, since)
	if err != nil {
		return 0
	}

	score := float64(uniqueTracks) / float64(totalPlays)
	if score > 1.0 {
		score = 1.0
	}
	return score
}

func periodBounds(period string) (time.Time, string, error) {
	now := time.Now()
	switch period {
	case "month":
		return now.AddDate(0, -1, 0), "این ماه", nil
	case "year":
		return now.AddDate(-1, 0, 0), "امسال", nil
	case "all_time":
		return time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location()), "همیشه", nil
	default:
		return time.Time{}, "", fmt.Errorf("invalid period: %s", period)
	}
}
