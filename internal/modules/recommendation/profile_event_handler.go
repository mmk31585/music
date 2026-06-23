package recommendation

import (
	"context"
	"log"

	"music/internal/modules/history"
	"music/internal/platform/events"
)

// ProfileEventHandler consumes signal and like events to asynchronously
// update the user's taste profile. It follows the same pattern as
// notification/event_handler.go and history/event_handler.go.
type ProfileEventHandler struct {
	tasteService   *TasteProfileService
	historyService *history.Service
}

func NewProfileEventHandler(tasteService *TasteProfileService, historyService *history.Service) *ProfileEventHandler {
	return &ProfileEventHandler{
		tasteService:   tasteService,
		historyService: historyService,
	}
}

// OnPlaybackSignalRecorded triggers a taste profile recompute when a new
// playback signal is recorded. Async via the event bus — never blocks the
// request path.
func (h *ProfileEventHandler) OnPlaybackSignalRecorded(ctx context.Context, event events.Event) error {
	e, ok := event.(events.PlaybackSignalRecordedEvent)
	if !ok {
		return nil
	}

	if err := h.tasteService.RecomputeProfile(ctx, e.UserID.String()); err != nil {
		log.Printf("profile recompute error for user %s: %v", e.UserID.String(), err)
	}

	return nil
}

// OnTrackLiked marks the latest playback signal for the liked track as
// explicitly liked and triggers a profile recompute.
func (h *ProfileEventHandler) OnTrackLiked(ctx context.Context, event events.Event) error {
	e, ok := event.(events.TrackLikedEvent)
	if !ok {
		return nil
	}

	// Mark the most recent listening history entry for this (user, track) as explicitly liked
	if err := h.historyService.GetRepo().UpdateSignalExplicitLike(ctx, e.UserID, e.TrackID); err != nil {
		log.Printf("update explicit like error for user %s track %s: %v", e.UserID, e.TrackID, err)
	}

	if err := h.tasteService.RecomputeProfile(ctx, e.UserID.String()); err != nil {
		log.Printf("profile recompute error (like) for user %s: %v", e.UserID.String(), err)
	}

	return nil
}
