package recommendation

import (
	"context"

	"music/internal/platform/eventbus"
)

type EventHandler struct {
	service Service
}

func NewEventHandler(service Service) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) OnTrackPlayed(ctx context.Context, event events.Event) error {
	_, ok := event.(events.TrackPlayedEvent)
	if !ok {
		return nil
	}

	// Future hook:
	// persist recommendation signal, update embeddings, update user taste profile, etc.
	return nil
}
