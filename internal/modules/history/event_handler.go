package history

import (
	"context"

	"music/internal/platform/eventbus"
)

type EventHandler struct {
	service *Service
}

func NewEventHandler(service *Service) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) OnTrackPlayed(ctx context.Context, event events.Event) error {
	e, ok := event.(events.TrackPlayedEvent)
	if !ok {
		return nil
	}

	_, err := h.service.Record(ctx, e.UserID, RecordListeningRequest{
		TrackID:   e.TrackID.String(),
		Duration:  e.Duration,
		Completed: e.Completed,
	})

	return err
}
