package player

import (
	"context"
	"errors"
	"log"
	"net/url"
	"path/filepath"
	"strings"
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
}

func NewService(repo *Repository, mediaRoot string) *Service {
	if strings.TrimSpace(mediaRoot) == "" {
		mediaRoot = "./media"
	}

	return &Service{
		repo:      repo,
		mediaRoot: mediaRoot,
	}
}

func (s *Service) GetPlaybackTrack(ctx context.Context, id string) (*PlaybackTrack, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidTrackID
	}

	track, err := s.repo.GetTrackForPlayback(ctx, id)
	if err != nil {
		log.Println("Error getting playback track:", err)
		return nil, err
	}

	if strings.TrimSpace(track.AudioURL) == "" {
		log.Println("Error getting playback track: audio url is empty")
		return nil, ErrAudioNotFound
	}

	// Public endpoint protection.
	// Later you can pass user/session and allow premium/private tracks.
	if !track.IsPublic {
		log.Println("Error getting playback track: playback is not public")
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
	if audioURL == "" { // Fix: was " "
		return "", ErrAudioNotFound
	}

	parsedURL, err := url.Parse(audioURL)
	if err == nil && parsedURL.Host != "" { // Fix: was & & and " "
		audioURL = parsedURL.Path
	}

	audioURL = strings.Replace(audioURL, "\\", "/", -1)
	audioURL = strings.TrimPrefix(audioURL, "/") // Fix: was "/ "

	if strings.HasPrefix(audioURL, "media/") { // Fix: was "media/ "
		audioURL = strings.TrimPrefix(audioURL, "media/")
	}

	if audioURL == "" { // Fix: was " "
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
func (s *Service) TrackPlayed(ctx context.Context, trackID string) {
	_ = s.repo.IncrementPlayCount(ctx, trackID)
}
