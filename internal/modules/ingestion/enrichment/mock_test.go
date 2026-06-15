package enrichment

import (
	"context"
	"sync"
)

type mockMusicBrainz struct {
	mu      sync.Mutex
	result  *MusicBrainzResult
	err     error
	callCount int
}

func (m *mockMusicBrainz) SearchRecording(ctx context.Context, query TrackQuery) (*MusicBrainzResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

type mockLastFM struct {
	mu      sync.Mutex
	result  *LastFMResult
	err     error
	callCount int
}

func (m *mockLastFM) SearchTrack(ctx context.Context, query TrackQuery) (*LastFMResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

type mockSpotify struct {
	mu      sync.Mutex
	result  *SpotifyResult
	err     error
	callCount int
}

func (m *mockSpotify) SearchTrack(ctx context.Context, query TrackQuery) (*SpotifyResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func fullMusicBrainzResult() *MusicBrainzResult {
	return &MusicBrainzResult{
		MBID:        "mb-id-123",
		Title:       "Test Song",
		ArtistName:  "Test Artist",
		ArtistMBID:  "artist-mb-id-456",
		AlbumName:   "Test Album",
		AlbumMBID:   "album-mb-id-789",
		ReleaseYear: 2020,
		Duration:    240,
		Genres:      []string{"rock", "alternative"},
	}
}

func fullLastFMResult() *LastFMResult {
	return &LastFMResult{
		PlayCount:     1500000,
		ListenerCount: 500000,
		Tags:          []string{"rock", "alternative", "indie"},
		ArtistBio:     "Test Artist is a fictional band created for testing purposes.",
		SimilarArtists: []string{"Similar Band 1", "Similar Band 2"},
	}
}

func fullSpotifyResult() *SpotifyResult {
	return &SpotifyResult{
		SpotifyID:      "spotify-id-abc",
		PreviewURL:     "https://p.scdn.co/mp3-preview/test",
		AlbumCoverURL:  "https://i.scdn.co/image/test",
		ArtistImageURL: "https://i.scdn.co/image/artist",
		Popularity:     78,
	}
}
