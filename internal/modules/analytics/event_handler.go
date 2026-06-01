package analytics

import (
	"context"

	"github.com/google/uuid"

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

	req := TrackEventRequest{
		EventType: string(EventPlay),
		TrackID:   e.TrackID.String(),
		Metadata: map[string]any{
			"source":    e.Source,
			"duration":  e.Duration,
			"completed": e.Completed,
		},
	}

	if e.ArtistID != uuid.Nil {
		req.ArtistID = e.ArtistID.String()
	}

	if e.AlbumID != uuid.Nil {
		req.AlbumID = e.AlbumID.String()
	}

	if err := h.service.TrackEvent(ctx, &e.UserID, req); err != nil {
		return err
	}

	if !e.Completed {
		return nil
	}

	completionReq := TrackEventRequest{
		EventType: string(EventCompletion),
		TrackID:   e.TrackID.String(),
		Metadata: map[string]any{
			"source":   e.Source,
			"duration": e.Duration,
		},
	}

	if e.ArtistID != uuid.Nil {
		completionReq.ArtistID = e.ArtistID.String()
	}

	if e.AlbumID != uuid.Nil {
		completionReq.AlbumID = e.AlbumID.String()
	}

	return h.service.TrackEvent(ctx, &e.UserID, completionReq)
}

func (h *EventHandler) OnPlaylistCreated(ctx context.Context, event events.Event) error {
	_, ok := event.(events.PlaylistCreatedEvent)
	if !ok {
		return nil
	}

	return nil
}

func (h *EventHandler) OnSubscriptionPurchased(ctx context.Context, event events.Event) error {
	_, ok := event.(events.SubscriptionPurchasedEvent)
	if !ok {
		return nil
	}

	return nil
}

func (h *EventHandler) OnUserRegistered(ctx context.Context, event events.Event) error {
	_, ok := event.(events.UserRegisteredEvent)
	if !ok {
		return nil
	}

	return nil
}
