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
	EventTrackPlayed           = "track.played"
	EventPlaylistCreated       = "playlist.created"
	EventSubscriptionPurchased = "subscription.purchased"
	EventUserRegistered        = "user.registered"
)

type TrackPlayedEvent struct {
	BaseEvent

	UserID    uuid.UUID `json:"user_id"`
	TrackID   uuid.UUID `json:"track_id"`
	ArtistID  uuid.UUID `json:"artist_id"`
	AlbumID   uuid.UUID `json:"album_id"`
	Duration  int       `json:"duration"`
	Completed bool      `json:"completed"`
	Source    string    `json:"source,omitempty"`
}

func (e TrackPlayedEvent) EventName() string {
	return EventTrackPlayed
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
