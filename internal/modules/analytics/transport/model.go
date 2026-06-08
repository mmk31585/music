package analytics

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventPlay         EventType = "play"
	EventPause        EventType = "pause"
	EventSkip         EventType = "skip"
	EventCompletion   EventType = "completion"
	EventSearch       EventType = "search"
	EventPlaylistOpen EventType = "playlist_open"
)

type AnalyticsEvent struct {
	ID         uuid.UUID      `db:"id" json:"id"`
	UserID     *uuid.UUID     `db:"user_id" json:"user_id,omitempty"`
	EventType  EventType      `db:"event_type" json:"event_type"`
	TrackID    *uuid.UUID     `db:"track_id" json:"track_id,omitempty"`
	ArtistID   *uuid.UUID     `db:"artist_id" json:"artist_id,omitempty"`
	AlbumID    *uuid.UUID     `db:"album_id" json:"album_id,omitempty"`
	PlaylistID *uuid.UUID     `db:"playlist_id" json:"playlist_id,omitempty"`
	Query      *string        `db:"query" json:"query,omitempty"`
	Metadata   map[string]any `db:"metadata" json:"metadata"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
}

type TrackAnalytics struct {
	TrackID         uuid.UUID `db:"track_id" json:"track_id"`
	PlayCount       int64     `db:"play_count" json:"play_count"`
	PauseCount      int64     `db:"pause_count" json:"pause_count"`
	SkipCount       int64     `db:"skip_count" json:"skip_count"`
	CompletionCount int64     `db:"completion_count" json:"completion_count"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type ArtistAnalytics struct {
	ArtistID        uuid.UUID `db:"artist_id" json:"artist_id"`
	PlayCount       int64     `db:"play_count" json:"play_count"`
	PauseCount      int64     `db:"pause_count" json:"pause_count"`
	SkipCount       int64     `db:"skip_count" json:"skip_count"`
	CompletionCount int64     `db:"completion_count" json:"completion_count"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type AlbumAnalytics struct {
	AlbumID         uuid.UUID `db:"album_id" json:"album_id"`
	PlayCount       int64     `db:"play_count" json:"play_count"`
	PauseCount      int64     `db:"pause_count" json:"pause_count"`
	SkipCount       int64     `db:"skip_count" json:"skip_count"`
	CompletionCount int64     `db:"completion_count" json:"completion_count"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
