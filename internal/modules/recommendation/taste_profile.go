package recommendation

import (
	"time"
)

// ── Taste profile constants ────────────────────────────────────────────
const (
	ProfileDecayHalfLifeDays = 30.0
	MaxSignalAgeDays         = 90
	MaxSeedTracks            = 15
	TopGenreLimit            = 10
	TopArtistLimit           = 20

	SignalWeightCompletePositive = 1.0
	SignalWeightReplayStrong     = 1.5
	SignalWeightExplicitLike     = 2.0
	SignalWeightSkipNegative     = -0.5
	SignalWeightPartialNeutral   = 0.2
)

// ── Models ─────────────────────────────────────────────────────────────

type TasteProfile struct {
	UserID             string    `json:"user_id" db:"user_id"`
	TopGenreIDs        []string  `json:"top_genre_ids" db:"top_genre_ids"`
	TopArtistIDs       []string  `json:"top_artist_ids" db:"top_artist_ids"`
	SeedTrackIDs       []string  `json:"seed_track_ids" db:"seed_track_ids"`
	AvgTempoPreference *float64  `json:"avg_tempo_preference" db:"avg_tempo_preference"`
	LastComputedAt     time.Time `json:"last_computed_at" db:"last_computed_at"`
	OnboardingGenreIDs []string  `json:"onboarding_genre_ids" db:"onboarding_genre_ids"`
}

// signalRow is a raw row from listening_history joined with catalog data,
// used as input for the profile recomputation algorithm.
type signalRow struct {
	ID                string    `db:"id"`
	UserID            string    `db:"user_id"`
	TrackID           string    `db:"track_id"`
	PlayedAt          time.Time `db:"played_at"`
	PlayedDurationMs  int64     `db:"played_duration_ms"`
	TrackDurationMs   int64     `db:"track_duration_ms"`
	CompletionPercent float64   `db:"completion_percent"`
	SignalType        string    `db:"signal_type"`
	IsExplicitLike    bool      `db:"is_explicit_like"`
	ArtistID          *string   `db:"artist_id"`
	GenreID           *string   `db:"genre_id"`
}
