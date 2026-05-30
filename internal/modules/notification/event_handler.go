package notification

import (
	"context"
	"fmt"

	"music/internal/platform/events"
)

type EventHandler struct {
	service *Service
}

func NewEventHandler(service *Service) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) OnUserRegistered(ctx context.Context, event events.Event) error {
	e, ok := event.(events.UserRegisteredEvent)
	if !ok {
		return nil
	}

	_, err := h.service.repo.Create(ctx, CreateNotificationInput{
		UserID: e.UserID.String(),
		Type:   "welcome",
		Title:  "Welcome",
		Body:   fmt.Sprintf("Welcome %s! Your account has been created successfully.", e.Name),
		Payload: map[string]any{
			"userId": e.UserID.String(),
			"email":  e.Email,
			"name":   e.Name,
		},
	})

	return err
}

func (h *EventHandler) OnSubscriptionPurchased(ctx context.Context, event events.Event) error {
	e, ok := event.(events.SubscriptionPurchasedEvent)
	if !ok {
		return nil
	}

	entityType := "subscription"
	entityID := e.SubscriptionID.String()

	_, err := h.service.repo.Create(ctx, CreateNotificationInput{
		UserID:     e.UserID.String(),
		Type:       "subscription_purchased",
		Title:      "Subscription activated",
		Body:       fmt.Sprintf("Your %s subscription has been activated.", e.Plan),
		EntityType: &entityType,
		EntityID:   &entityID,
		Payload: map[string]any{
			"subscriptionId": e.SubscriptionID.String(),
			"plan":           e.Plan,
			"amount":         e.Amount,
			"currency":       e.Currency,
		},
	})

	return err
}

func (h *EventHandler) OnPlaylistCreated(ctx context.Context, event events.Event) error {
	e, ok := event.(events.PlaylistCreatedEvent)
	if !ok {
		return nil
	}

	entityType := "playlist"
	entityID := e.PlaylistID.String()

	_, err := h.service.repo.Create(ctx, CreateNotificationInput{
		UserID:     e.UserID.String(),
		Type:       "playlist_created",
		Title:      "Playlist created",
		Body:       fmt.Sprintf("Your playlist \"%s\" was created successfully.", e.Name),
		EntityType: &entityType,
		EntityID:   &entityID,
		Payload: map[string]any{
			"playlistId": e.PlaylistID.String(),
			"name":       e.Name,
		},
	})

	return err
}
