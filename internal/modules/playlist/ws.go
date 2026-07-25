package playlist

import (
	"time"

	"music/internal/platform/ws"

	"github.com/google/uuid"
)

type PlaylistEvent struct {
	Type       string      `json:"type"`
	PlaylistID string      `json:"playlist_id"`
	TrackID    string      `json:"track_id,omitempty"`
	Position   int         `json:"position,omitempty"`
	UserID     string      `json:"user_id"`
	Payload    interface{} `json:"payload,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
}

type PlaylistBroadcaster struct {
	hub *ws.Hub
}

func NewPlaylistBroadcaster(hub *ws.Hub) *PlaylistBroadcaster {
	return &PlaylistBroadcaster{hub: hub}
}

func (b *PlaylistBroadcaster) Broadcast(playlistID string, event PlaylistEvent) {
	if b.hub == nil {
		return
	}
	channel := "playlist:" + playlistID
	b.hub.BroadcastToChannel(channel, ws.Message{
		Type:      "playlist." + event.Type,
		Payload:   event,
		Timestamp: time.Now(),
	}, uuid.UUID{})
}

func (b *PlaylistBroadcaster) TrackAdded(playlistID, trackID, userID string) {
	b.Broadcast(playlistID, PlaylistEvent{
		Type: "track_added", PlaylistID: playlistID, TrackID: trackID, UserID: userID, Timestamp: time.Now(),
	})
}

func (b *PlaylistBroadcaster) TrackRemoved(playlistID, trackID, userID string) {
	b.Broadcast(playlistID, PlaylistEvent{
		Type: "track_removed", PlaylistID: playlistID, TrackID: trackID, UserID: userID, Timestamp: time.Now(),
	})
}

func (b *PlaylistBroadcaster) TrackReordered(playlistID, trackID, userID string, position int) {
	b.Broadcast(playlistID, PlaylistEvent{
		Type: "track_reordered", PlaylistID: playlistID, TrackID: trackID, Position: position, UserID: userID, Timestamp: time.Now(),
	})
}

func (b *PlaylistBroadcaster) PlaylistUpdated(playlistID, userID string) {
	b.Broadcast(playlistID, PlaylistEvent{
		Type: "updated", PlaylistID: playlistID, UserID: userID, Timestamp: time.Now(),
	})
}
