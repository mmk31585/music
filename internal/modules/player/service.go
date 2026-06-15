package player

import (
	"context"
	"errors"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"music/internal/platform/events"
	platformstorage "music/internal/platform/storage"
)

var (
	ErrInvalidTrackID  = errors.New("invalid track id")
	ErrAudioNotFound   = errors.New("track audio not found")
	ErrInvalidMediaURL = errors.New("invalid media url")
	ErrPrivateTrack    = errors.New("track is private")
)

type Service struct {
	repo      *Repository
	storage   platformstorage.Storage
	publisher events.Publisher
}

func NewService(
	repo *Repository,
	storage platformstorage.Storage,
	publisher events.Publisher,
) *Service {
	return &Service{
		repo:      repo,
		storage:   storage,
		publisher: publisher,
	}
}

func (s *Service) GetPlaybackTrack(ctx context.Context, id string) (*PlaybackTrack, error) {
	return s.getPlaybackTrack(ctx, id, true)
}

func (s *Service) GetAdminPlaybackTrack(ctx context.Context, id string) (*PlaybackTrack, error) {
	return s.getPlaybackTrack(ctx, id, false)
}

func (s *Service) getPlaybackTrack(ctx context.Context, id string, requirePublic bool) (*PlaybackTrack, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidTrackID
	}

	track, err := s.repo.GetTrackForPlayback(ctx, id)
	if err != nil {
		log.Println("error getting playback track:", err)
		return nil, err
	}

	if strings.TrimSpace(track.AudioURL) == "" {
		log.Println("error getting playback track: audio url is empty")
		return nil, ErrAudioNotFound
	}

	if requirePublic && !track.IsPublic {
		log.Println("error getting playback track: playback is not public")
		return nil, ErrPrivateTrack
	}

	return track, nil
}

func (s *Service) BuildPlaybackResponse(track *PlaybackTrack) PlaybackTrackResponse {
	return PlaybackTrackResponse{
		ID:              track.ID,
		Title:           track.Title,
		ArtistName:      track.ArtistName,
		AlbumTitle:      track.AlbumTitle,
		CoverURL:        track.CoverURL,
		DurationSeconds: track.DurationSeconds,

		// Important: your registered route is /api/v1/player/...
		StreamURL: "/api/v1/player/tracks/" + track.ID + "/stream",
	}
}

func (s *Service) BuildAdminPlaybackResponse(track *PlaybackTrack) PlaybackTrackResponse {
	return PlaybackTrackResponse{
		ID:              track.ID,
		Title:           track.Title,
		ArtistName:      track.ArtistName,
		AlbumTitle:      track.AlbumTitle,
		CoverURL:        track.CoverURL,
		DurationSeconds: track.DurationSeconds,

		// Private/admin stream URL.
		StreamURL: "/api/v1/admin/player/tracks/" + track.ID + "/stream",
	}
}

// ResolveAudioStorageKey converts stored audio_url into your storage key.
//
// Examples:
//
// http://localhost:8080/uploads/track-audio/file.mp3
// -> track-audio/file.mp3
//
// /uploads/track-audio/file.mp3
// -> track-audio/file.mp3
//
// track-audio/file.mp3
// -> track-audio/file.mp3
func (s *Service) ResolveAudioStorageKey(audioURL string) (string, error) {
	audioURL = strings.TrimSpace(audioURL)
	if audioURL == "" {
		return "", ErrAudioNotFound
	}

	parsedURL, err := url.Parse(audioURL)
	if err == nil && parsedURL.Scheme != "" {
		audioURL = parsedURL.Path
	}

	audioURL = strings.ReplaceAll(audioURL, "\\", "/")
	audioURL = strings.TrimSpace(audioURL)
	audioURL = strings.TrimPrefix(audioURL, "/")

	audioURL = strings.TrimPrefix(audioURL, "uploads/")
	audioURL = strings.TrimPrefix(audioURL, "media/")

	if audioURL == "" {
		return "", ErrInvalidMediaURL
	}

	cleanKey := filepath.Clean(audioURL)
	cleanKey = strings.ReplaceAll(cleanKey, "\\", "/")

	if cleanKey == "." ||
		strings.HasPrefix(cleanKey, "../") ||
		strings.Contains(cleanKey, "/../") ||
		strings.HasPrefix(cleanKey, "/") {
		return "", ErrInvalidMediaURL
	}

	return cleanKey, nil
}

func (s *Service) ResolveAudioURL(ctx context.Context, audioURL string) (string, error) {
	log.Println("ResolveAudioURL input:", audioURL)

	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	key, err := s.ResolveAudioStorageKey(audioURL)
	if err != nil {
		log.Println("ResolveAudioStorageKey error:", err)
		return "", err
	}

	log.Println("resolved storage key:", key)

	exists, err := s.storage.Exists(ctx, key)
	if err != nil {
		log.Println("storage.Exists error:", err)
		return "", ErrAudioNotFound
	}

	log.Println("storage.Exists result:", exists)

	if !exists {
		return "", ErrAudioNotFound
	}

	resolvedURL, err := s.storage.GetURL(ctx, key)
	if err != nil {
		log.Println("storage.GetURL error:", err)
		return "", ErrAudioNotFound
	}

	log.Println("storage.GetURL result:", resolvedURL)

	return resolvedURL, nil
}

func (s *Service) TrackPlayed(
	ctx context.Context,
	userID string,
	track *PlaybackTrack,
	duration int,
	completed bool,
	source string,
) {
	if track == nil {
		return
	}

	_ = s.repo.IncrementPlayCount(ctx, track.ID)

	if s.publisher == nil {
		return
	}

	parsedUserID, err := uuid.Parse(strings.TrimSpace(userID))
	if err != nil {
		return
	}

	trackUUID, err := uuid.Parse(strings.TrimSpace(track.ID))
	if err != nil {
		return
	}

	var artistUUID uuid.UUID
	if strings.TrimSpace(track.ArtistID) != "" {
		artistUUID, _ = uuid.Parse(track.ArtistID)
	}

	var albumUUID uuid.UUID
	if strings.TrimSpace(track.AlbumID) != "" {
		albumUUID, _ = uuid.Parse(track.AlbumID)
	}

	_ = s.publisher.Publish(ctx, events.TrackPlayedEvent{
		BaseEvent: events.NewBaseEvent(),
		UserID:    parsedUserID,
		TrackID:   trackUUID,
		ArtistID:  artistUUID,
		AlbumID:   albumUUID,
		Duration:  duration,
		Completed: completed,
		Source:    source,
	})
}
