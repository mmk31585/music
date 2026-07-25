package gamification

import (
	"context"

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

	_, err := h.service.AddXP(ctx, e.UserID.String(), XPForStream, "stream")
	if err != nil {
		return err
	}

	// Track streaming challenge progress
	return h.service.IncrementChallengeProgress(ctx, e.UserID.String(), "streams")
}

func (h *EventHandler) OnTrackLiked(ctx context.Context, event events.Event) error {
	e, ok := event.(events.TrackLikedEvent)
	if !ok {
		return nil
	}

	_, err := h.service.AddXP(ctx, e.UserID.String(), XPForLike, "like")
	if err != nil {
		return err
	}

	return h.service.IncrementChallengeProgress(ctx, e.UserID.String(), "likes")
}

func (h *EventHandler) OnTrackUploaded(ctx context.Context, event events.Event) error {
	e, ok := event.(events.TrackUploadedEvent)
	if !ok {
		return nil
	}

	_, err := h.service.AddXP(ctx, e.UserID.String(), XPForUpload, "upload")
	if err != nil {
		return err
	}

	return h.service.IncrementChallengeProgress(ctx, e.UserID.String(), "upload")
}

func (h *EventHandler) OnUploadPublished(ctx context.Context, event events.Event) error {
	e, ok := event.(events.UploadPublishedEvent)
	if !ok {
		return nil
	}

	_, err := h.service.AddXP(ctx, e.UserID.String(), XPForUploadPublished, "upload_published")
	if err != nil {
		return err
	}

	return h.service.IncrementChallengeProgress(ctx, e.UserID.String(), "publish")
}

func (h *EventHandler) OnContributionAccepted(ctx context.Context, event events.Event) error {
	e, ok := event.(events.ContributionAcceptedEvent)
	if !ok {
		return nil
	}

	_, err := h.service.AddXP(ctx, e.UserID.String(), XPForContributionAccepted, "contribution_accepted")
	if err != nil {
		return err
	}

	return h.service.IncrementChallengeProgress(ctx, e.UserID.String(), "contribute")
}
