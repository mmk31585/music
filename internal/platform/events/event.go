package events

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	EventName() string
	OccurredAt() time.Time
}

type BaseEvent struct {
	At time.Time `json:"occurred_at"`
}

func NewBaseEvent() BaseEvent {
	return BaseEvent{
		At: time.Now().UTC(),
	}
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.At
}

const (
	EventTrackPlayed            = "track.played"
	EventPlaylistCreated        = "playlist.created"
	EventSubscriptionPurchased  = "subscription.purchased"
	EventUserRegistered         = "user.registered"
	EventPlaybackSignalRecorded = "playback.signal_recorded"
	EventTrackLiked             = "track.liked"
	EventTrackUploaded          = "track.uploaded"
	EventUploadPublished        = "upload.published"
	EventContributionAccepted   = "contribution.accepted"
)

type TrackPlayedEvent struct {
	BaseEvent

	UserID          uuid.UUID `json:"user_id"`
	TrackID         uuid.UUID `json:"track_id"`
	ArtistID        uuid.UUID `json:"artist_id"`
	AlbumID         uuid.UUID `json:"album_id"`
	Duration        int       `json:"duration"`  // seconds listened
	Completed       bool      `json:"completed"` // frontend-reported completion
	Source          string    `json:"source,omitempty"`
	SessionID       uuid.UUID `json:"session_id,omitempty"`
	TrackDurationMs int64     `json:"track_duration_ms,omitempty"`
}

func (e TrackPlayedEvent) EventName() string {
	return EventTrackPlayed
}

type PlaybackSignalRecordedEvent struct {
	BaseEvent

	UserID            uuid.UUID `json:"user_id"`
	TrackID           uuid.UUID `json:"track_id"`
	SessionID         uuid.UUID `json:"session_id,omitempty"`
	PlayedDurationMs  int64     `json:"played_duration_ms"`
	TrackDurationMs   int64     `json:"track_duration_ms"`
	CompletionPercent float64   `json:"completion_percent"`
	SignalType        string    `json:"signal_type"`
}

func (e PlaybackSignalRecordedEvent) EventName() string {
	return EventPlaybackSignalRecorded
}

type TrackLikedEvent struct {
	BaseEvent

	UserID  uuid.UUID `json:"user_id"`
	TrackID uuid.UUID `json:"track_id"`
}

func (e TrackLikedEvent) EventName() string {
	return EventTrackLiked
}

type PlaylistCreatedEvent struct {
	BaseEvent

	UserID     uuid.UUID `json:"user_id"`
	PlaylistID uuid.UUID `json:"playlist_id"`
	Name       string    `json:"name"`
	IsPublic   bool      `json:"is_public"`
}

func (e PlaylistCreatedEvent) EventName() string {
	return EventPlaylistCreated
}

type SubscriptionPurchasedEvent struct {
	BaseEvent

	UserID         uuid.UUID `json:"user_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Plan           string    `json:"plan"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
}

func (e SubscriptionPurchasedEvent) EventName() string {
	return EventSubscriptionPurchased
}

type UserRegisteredEvent struct {
	BaseEvent

	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
}

func (e UserRegisteredEvent) EventName() string {
	return EventUserRegistered
}

type TrackUploadedEvent struct {
	BaseEvent

	UserID  uuid.UUID `json:"user_id"`
	DraftID string    `json:"draft_id"`
	Title   string    `json:"title"`
}

func (e TrackUploadedEvent) EventName() string {
	return EventTrackUploaded
}

type UploadPublishedEvent struct {
	BaseEvent

	UserID  uuid.UUID `json:"user_id"`
	TrackID uuid.UUID `json:"track_id"`
	DraftID string    `json:"draft_id"`
	Title   string    `json:"title"`
}

func (e UploadPublishedEvent) EventName() string {
	return EventUploadPublished
}

type ContributionAcceptedEvent struct {
	BaseEvent

	UserID           uuid.UUID `json:"user_id"`
	ContributionID   int64     `json:"contribution_id"`
	ContributionType string    `json:"contribution_type"`
}

func (e ContributionAcceptedEvent) EventName() string {
	return EventContributionAccepted
}
