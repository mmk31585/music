package history

import (
	"time"

	"github.com/google/uuid"
)

// ── Signal type constants ──────────────────────────────────────────────
const (
	SkipThresholdMs     = 30_000 // under 30s played = skip
	CompletionThreshold = 0.80   // 80%+ played = strong positive

	SignalTypeSkipNegative     = "skip_negative"
	SignalTypeCompletePositive = "complete_positive"
	SignalTypeReplayStrong     = "replay_strong"
	SignalTypePartialNeutral   = "partial_neutral"
)

// PlaybackSignal is the enriched listening history entry with signal
// classification. The underlying storage is the listening_history table
// with the additional columns added by migration 000036.
type PlaybackSignal struct {
	ID                uuid.UUID `db:"id" json:"id"`
	UserID            uuid.UUID `db:"user_id" json:"user_id"`
	TrackID           uuid.UUID `db:"track_id" json:"track_id"`
	SessionID         uuid.UUID `db:"session_id" json:"session_id,omitempty"`
	PlayedDurationMs  int64     `db:"played_duration_ms" json:"played_duration_ms"`
	TrackDurationMs   int64     `db:"track_duration_ms" json:"track_duration_ms"`
	CompletionPercent float64   `db:"completion_percent" json:"completion_percent"`
	SignalType        string    `db:"signal_type" json:"signal_type"`
	IsExplicitLike    bool      `db:"is_explicit_like" json:"is_explicit_like"`
	PlayedAt          time.Time `db:"played_at" json:"played_at"`
}

type ListeningHistoryItem struct {
	ID            uuid.UUID `db:"id" json:"id"`
	UserID        uuid.UUID `db:"user_id" json:"user_id"`
	TrackID       uuid.UUID `db:"track_id" json:"track_id"`
	PlayedAt      time.Time `db:"played_at" json:"played_at"`
	Duration      int       `db:"duration" json:"duration"`
	Completed     bool      `db:"completed" json:"completed"`
	TrackTitle    string    `db:"track_title" json:"track_title,omitempty"`
	TrackDuration *int      `db:"track_duration" json:"track_duration,omitempty"`
	TrackCoverURL string    `db:"track_cover_url" json:"track_cover_url,omitempty"`
	ArtistName    string    `db:"artist_name" json:"artist_name,omitempty"`
}
