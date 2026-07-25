package history

import (
	"context"

	"github.com/google/uuid"

	"music/internal/platform/events"
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

	sessionID := ""
	if e.SessionID != uuid.Nil {
		sessionID = e.SessionID.String()
	}

	_, err := h.service.Record(ctx, e.UserID, RecordListeningRequest{
		TrackID:         e.TrackID.String(),
		Duration:        e.Duration,
		Completed:       e.Completed,
		SessionID:       sessionID,
		TrackDurationMs: e.TrackDurationMs,
	})

	return err
}
