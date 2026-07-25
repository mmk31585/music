package enrichment

import (
	"context"
	"encoding/json"
	"sync"
)

type mockMusicBrainz struct {
	mu        sync.Mutex
	result    *MusicBrainzResult
	err       error
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
	mu        sync.Mutex
	result    *LastFMResult
	err       error
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

func (m *mockLastFM) SearchArtist(ctx context.Context, name string) (*LastFMResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

type mockSpotify struct {
	mu        sync.Mutex
	result    *SpotifyResult
	err       error
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

func (m *mockSpotify) SearchArtistImage(ctx context.Context, name string) (string, error) {
	// For tests, just return empty — we don't test artist image enrichment here
	return "", nil
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
		PlayCount:      1500000,
		ListenerCount:  500000,
		Tags:           []string{"rock", "alternative", "indie"},
		ArtistBio:      "Test Artist is a fictional band created for testing purposes.",
		SimilarArtists: []string{"Similar Band 1", "Similar Band 2"},
	}
}

type mockLRCLib struct {
	mu        sync.Mutex
	result    *LRCLibResult
	err       error
	callCount int
}

func (m *mockLRCLib) SearchLyrics(ctx context.Context, query TrackQuery) (*LRCLibResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func fullLRCLibResult() *LRCLibResult {
	return &LRCLibResult{
		ID:           1,
		TrackName:    "Test Song",
		ArtistName:   "Test Artist",
		AlbumName:    "Test Album",
		Duration:     json.Number("240"),
		Synced:       true,
		SyncedLyrics: "[00:00.00]Test lyric line\n[00:05.00]Another test line",
	}
}

type mockMLClient struct {
	mu     sync.Mutex
	result *MLResult
	err    error
}

func (m *mockMLClient) FetchLyrics(ctx context.Context, query TrackQuery) (*MLResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func (m *mockMLClient) FetchCover(ctx context.Context, query TrackQuery) (*MLResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
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
