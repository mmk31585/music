package notification

import (
	"context"
	"encoding/json"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListNotifications(ctx context.Context, userID string, limit, offset int) (*ListNotificationsResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	items, err := s.repo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.repo.CountUnreadByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]NotificationResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapNotification(item))
	}

	return &ListNotificationsResponse{
		Notifications: result,
		UnreadCount:   unreadCount,
	}, nil
}

func (s *Service) MarkAsRead(ctx context.Context, userID, notificationID string) error {
	return s.repo.MarkAsRead(ctx, userID, notificationID)
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *Service) NotifyArtistPublishedTrack(ctx context.Context, event ArtistPublishedTrackEvent) error {
	entityType := "track"
	entityID := event.TrackID

	inputs := make([]CreateNotificationInput, 0, len(event.UserIDs))
	for _, userID := range event.UserIDs {
		if userID == "" {
			continue
		}

		inputs = append(inputs, CreateNotificationInput{
			UserID:     userID,
			Type:       TypeArtistPublishedTrack,
			Title:      "New release",
			Body:       fmt.Sprintf("%s published a new track: %s", event.ArtistName, event.TrackTitle),
			EntityType: strPtr(entityType),
			EntityID:   strPtr(entityID),
			Payload: map[string]any{
				"artistId":   event.ArtistID,
				"artistName": event.ArtistName,
				"trackId":    event.TrackID,
				"trackTitle": event.TrackTitle,
			},
		})
	}

	if len(inputs) == 0 {
		return nil
	}

	return s.repo.CreateBulk(ctx, inputs)
}

func (s *Service) NotifyPlaylistShared(ctx context.Context, event PlaylistSharedEvent) error {
	entityType := "playlist"
	entityID := event.PlaylistID

	_, err := s.repo.Create(ctx, CreateNotificationInput{
		UserID:     event.TargetUserID,
		Type:       TypePlaylistShared,
		Title:      "Playlist shared with you",
		Body:       fmt.Sprintf("%s shared playlist \"%s\" with you", event.OwnerName, event.PlaylistName),
		EntityType: strPtr(entityType),
		EntityID:   strPtr(entityID),
		Payload: map[string]any{
			"playlistId":   event.PlaylistID,
			"playlistName": event.PlaylistName,
			"ownerId":      event.OwnerID,
			"ownerName":    event.OwnerName,
		},
	})
	return err
}

func (s *Service) NotifySubscriptionRenewed(ctx context.Context, event SubscriptionRenewedEvent) error {
	entityType := "subscription"
	entityID := event.PlanID

	_, err := s.repo.Create(ctx, CreateNotificationInput{
		UserID:     event.UserID,
		Type:       TypeSubscriptionRenewed,
		Title:      "Subscription renewed",
		Body:       fmt.Sprintf("Your %s subscription has been renewed", event.PlanName),
		EntityType: strPtr(entityType),
		EntityID:   strPtr(entityID),
		Payload: map[string]any{
			"planId":      event.PlanID,
			"planName":    event.PlanName,
			"amountCents": event.AmountCents,
			"currency":    event.Currency,
		},
	})
	return err
}

func mapNotification(n Notification) NotificationResponse {
	var payload map[string]any
	if len(n.Payload) > 0 {
		_ = json.Unmarshal(n.Payload, &payload)
	}
	if payload == nil {
		payload = map[string]any{}
	}

	return NotificationResponse{
		ID:         n.ID,
		Type:       n.Type,
		Title:      n.Title,
		Body:       n.Body,
		EntityType: n.EntityType,
		EntityID:   n.EntityID,
		Payload:    payload,
		IsRead:     n.IsRead,
		ReadAt:     n.ReadAt,
		CreatedAt:  n.CreatedAt,
	}
}

func strPtr(v string) *string {
	return &v
}
