package notification

import "time"

const (
	TypeArtistPublishedTrack = "artist_published_track"
	TypePlaylistShared       = "playlist_shared"
	TypeSubscriptionRenewed  = "subscription_renewed"
)

type NotificationResponse struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Title      string         `json:"title"`
	Body       string         `json:"body"`
	EntityType *string        `json:"entityType,omitempty"`
	EntityID   *string        `json:"entityId,omitempty"`
	Payload    map[string]any `json:"payload,omitempty"`
	IsRead     bool           `json:"isRead"`
	ReadAt     *time.Time     `json:"readAt,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
}

type ListNotificationsResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int                    `json:"unreadCount"`
}

type MarkReadResponse struct {
	Message string `json:"message"`
}

type CreateNotificationInput struct {
	UserID     string
	Type       string
	Title      string
	Body       string
	EntityType *string
	EntityID   *string
	Payload    map[string]any
}

type ArtistPublishedTrackEvent struct {
	ArtistID   string
	ArtistName string
	TrackID    string
	TrackTitle string
	UserIDs    []string
}

type PlaylistSharedEvent struct {
	PlaylistID   string
	PlaylistName string
	OwnerID      string
	OwnerName    string
	TargetUserID string
}

type SubscriptionRenewedEvent struct {
	UserID      string
	PlanID      string
	PlanName    string
	AmountCents int64
	Currency    string
}
