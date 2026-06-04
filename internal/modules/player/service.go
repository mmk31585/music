package player

import (
	"context"
	"errors"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"music/internal/platform/eventbus"
)

var (
	ErrInvalidTrackID  = errors.New("invalid track id")
	ErrAudioNotFound   = errors.New("track audio not found")
	ErrInvalidMediaURL = errors.New("invalid media url")
	ErrPrivateTrack    = errors.New("track is private")
)

type Service struct {
	repo      *Repository
	mediaRoot string
	publisher events.Publisher
}

func NewService(repo *Repository, mediaRoot string, publisher events.Publisher) *Service {
	if strings.TrimSpace(mediaRoot) == "" {
		mediaRoot = "./media"
	}

	return &Service{
		repo:      repo,
		mediaRoot: mediaRoot,
		publisher: publisher,
	}
}

func (s *Service) GetPlaybackTrack(ctx context.Context, id string) (*PlaybackTrack, error) {
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

	if !track.IsPublic {
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
		StreamURL:       "/api/player/tracks/" + track.ID + "/stream",
	}
}

func (s *Service) ResolveAudioFilePath(audioURL string) (string, error) {
	audioURL = strings.TrimSpace(audioURL)
	if audioURL == "" {
		return "", ErrAudioNotFound
	}

	parsedURL, err := url.Parse(audioURL)
	if err == nil && parsedURL.Host != "" {
		audioURL = parsedURL.Path
	}

	audioURL = strings.ReplaceAll(audioURL, "\\", "/")
	audioURL = strings.TrimPrefix(audioURL, "/")

	if strings.HasPrefix(audioURL, "media/") {
		audioURL = strings.TrimPrefix(audioURL, "media/")
	}

	if audioURL == "" {
		return "", ErrInvalidMediaURL
	}

	cleanRelative := filepath.Clean(audioURL)

	if cleanRelative == "." ||
		strings.HasPrefix(cleanRelative, "..") ||
		strings.Contains(cleanRelative, string(filepath.Separator)+".."+string(filepath.Separator)) {
		return "", ErrInvalidMediaURL
	}

	fullPath := filepath.Join(s.mediaRoot, cleanRelative)

	absRoot, err := filepath.Abs(s.mediaRoot)
	if err != nil {
		return "", err
	}

	absFile, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(absRoot, absFile)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", ErrInvalidMediaURL
	}

	return absFile, nil
}

// TrackPlayed keeps your counter update, but also emits an event.
// This is the important upgrade.
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
