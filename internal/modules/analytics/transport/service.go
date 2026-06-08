package analytics

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidEventType    = errors.New("invalid event type")
	ErrInvalidTrackID      = errors.New("invalid track id")
	ErrInvalidArtistID     = errors.New("invalid artist id")
	ErrInvalidAlbumID      = errors.New("invalid album id")
	ErrInvalidPlaylistID   = errors.New("invalid playlist id")
	ErrSearchQueryRequired = errors.New("search query is required")
	ErrTrackIDRequired     = errors.New("track id is required for this event")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) TrackEvent(
	ctx context.Context,
	userID *uuid.UUID,
	req TrackEventRequest,
) error {
	eventType, err := parseEventType(req.EventType)
	if err != nil {
		return err
	}

	params := CreateEventParams{
		UserID:    userID,
		EventType: eventType,
		Metadata:  req.Metadata,
	}

	if req.TrackID != "" {
		id, err := uuid.Parse(req.TrackID)
		if err != nil {
			return ErrInvalidTrackID
		}
		params.TrackID = &id
	}

	if req.ArtistID != "" {
		id, err := uuid.Parse(req.ArtistID)
		if err != nil {
			return ErrInvalidArtistID
		}
		params.ArtistID = &id
	}

	if req.AlbumID != "" {
		id, err := uuid.Parse(req.AlbumID)
		if err != nil {
			return ErrInvalidAlbumID
		}
		params.AlbumID = &id
	}

	if req.PlaylistID != "" {
		id, err := uuid.Parse(req.PlaylistID)
		if err != nil {
			return ErrInvalidPlaylistID
		}
		params.PlaylistID = &id
	}

	if trimmed := strings.TrimSpace(req.Query); trimmed != "" {
		params.Query = &trimmed
	}

	switch eventType {
	case EventSearch:
		if params.Query == nil {
			return ErrSearchQueryRequired
		}

	case EventPlay, EventPause, EventSkip, EventCompletion:
		if params.TrackID == nil {
			return ErrTrackIDRequired
		}
	}

	return s.repo.CreateEvent(ctx, params)
}

func parseEventType(value string) (EventType, error) {
	switch EventType(strings.TrimSpace(strings.ToLower(value))) {
	case EventPlay:
		return EventPlay, nil
	case EventPause:
		return EventPause, nil
	case EventSkip:
		return EventSkip, nil
	case EventCompletion:
		return EventCompletion, nil
	case EventSearch:
		return EventSearch, nil
	case EventPlaylistOpen:
		return EventPlaylistOpen, nil
	default:
		return "", ErrInvalidEventType
	}
}
