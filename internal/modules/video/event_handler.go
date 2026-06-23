package video

import (
	"context"

	"music/internal/platform/events"
)

// EventHandler subscribes to domain events that affect the video/music-status
// domain. Currently handles:
//   - EventTrackPlayed: upsert the user's music status so the "now playing"
//     indicator shows their current track.
type EventHandler struct {
	service *Service
}

func NewEventHandler(service *Service) *EventHandler {
	return &EventHandler{service: service}
}

// OnTrackPlayed upserts the user's music status with the track they just
// started playing. Visibility defaults to "private" — the user can change
// it later via the SetMusicStatus API. The 30-minute staleness TTL in
// GetMusicStatus prevents stale "now playing" indicators from persisting.
func (h *EventHandler) OnTrackPlayed(ctx context.Context, event events.Event) error {
	e, ok := event.(events.TrackPlayedEvent)
	if !ok {
		return nil
	}

	trackIDStr := e.TrackID.String()
	_, err := h.service.SetMusicStatus(ctx, e.UserID, &trackIDStr, "private")
	return err
}
