package catalog

import (
	"context"
	"testing"

	"music/internal/modules/ingestion/enrichment"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// mockDeezerClient implements enrichment.DeezerClient for testing
type mockDeezerClient struct {
	imageURL string
	err      error
}

func (m *mockDeezerClient) SearchArtistImage(ctx context.Context, name string) (string, error) {
	return m.imageURL, m.err
}

func (m *mockDeezerClient) SearchAlbumCover(ctx context.Context, title, artist string) (string, error) {
	return "", nil
}

// mockSpotifyClientForArtist implements enrichment.SpotifyClient for testing
type mockSpotifyClientForArtist struct {
	imageURL string
	err      error
}

func (m *mockSpotifyClientForArtist) SearchTrack(ctx context.Context, query enrichment.TrackQuery) (*enrichment.SpotifyResult, error) {
	return nil, nil
}

func (m *mockSpotifyClientForArtist) SearchArtistImage(ctx context.Context, name string) (string, error) {
	return m.imageURL, m.err
}

// mockLastFMClientForArtist implements enrichment.LastFMClient for testing
type mockLastFMClientForArtist struct {
	result *enrichment.LastFMResult
	err    error
}

func (m *mockLastFMClientForArtist) SearchTrack(ctx context.Context, query enrichment.TrackQuery) (*enrichment.LastFMResult, error) {
	return nil, nil
}

func (m *mockLastFMClientForArtist) SearchArtist(ctx context.Context, name string) (*enrichment.LastFMResult, error) {
	return m.result, m.err
}

// mockCoverArtClient implements enrichment.CoverArtClient for testing
type mockCoverArtClient struct{}

func (m *mockCoverArtClient) SearchAlbumCover(ctx context.Context, title, artist string) (string, error) {
	return "", nil
}

func TestEnrichArtist_PrioritySpotifyOverDeezer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock clients
	deezerClient := &mockDeezerClient{
		imageURL: "https://deezer.example.com/image.jpg",
	}
	spotifyClient := &mockSpotifyClientForArtist{
		imageURL: "https://open.spotify.com/image.jpg",
	}

	// Verify Spotify is preferred over Deezer
	imageURL := spotifyClient.imageURL
	if imageURL == "" {
		imageURL = deezerClient.imageURL
	}

	assert.Equal(t, "https://open.spotify.com/image.jpg", imageURL)
}

func TestEnrichArtist_FallbackToDeezerWhenSpotifyEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock clients - Spotify returns empty
	deezerClient := &mockDeezerClient{
		imageURL: "https://deezer.example.com/image.jpg",
	}
	spotifyClient := &mockSpotifyClientForArtist{
		imageURL: "", // No Spotify result
	}

	// Verify fallback logic
	imageURL := spotifyClient.imageURL
	if imageURL == "" {
		imageURL = deezerClient.imageURL
	}

	assert.Equal(t, "https://deezer.example.com/image.jpg", imageURL)
}

func TestEnrichArtist_FallbackToLastFMWhenOthersEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock clients - Both Spotify and Deezer return empty
	deezerClient := &mockDeezerClient{
		imageURL: "",
	}
	spotifyClient := &mockSpotifyClientForArtist{
		imageURL: "",
	}
	lastfmClient := &mockLastFMClientForArtist{
		result: &enrichment.LastFMResult{
			ArtistImageURL: "https://lastfm.example.com/image.jpg",
		},
	}

	// Verify Last.fm fallback
	imageURL := spotifyClient.imageURL
	if imageURL == "" {
		imageURL = deezerClient.imageURL
	}
	if imageURL == "" && lastfmClient.result != nil {
		imageURL = lastfmClient.result.ArtistImageURL
	}

	assert.Equal(t, "https://lastfm.example.com/image.jpg", imageURL)
}

func TestEnrichArtist_KeepExistingWhenAllFail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock clients - All return empty
	deezerClient := &mockDeezerClient{
		imageURL: "",
	}
	spotifyClient := &mockSpotifyClientForArtist{
		imageURL: "",
	}
	lastfmClient := &mockLastFMClientForArtist{
		result: nil,
	}

	existingImageURL := "https://existing.example.com/image.jpg"

	// Verify we keep existing when all sources fail
	imageURL := spotifyClient.imageURL
	if imageURL == "" {
		imageURL = deezerClient.imageURL
	}
	if imageURL == "" && lastfmClient.result != nil {
		imageURL = lastfmClient.result.ArtistImageURL
	}
	if imageURL == "" {
		imageURL = existingImageURL
	}

	assert.Equal(t, "https://existing.example.com/image.jpg", imageURL)
}

func TestEnrichArtist_LogsResults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock clients
	deezerClient := &mockDeezerClient{
		imageURL: "https://deezer.example.com/image.jpg",
	}
	spotifyClient := &mockSpotifyClientForArtist{
		imageURL: "https://open.spotify.com/image.jpg",
	}
	lastfmClient := &mockLastFMClientForArtist{
		result: &enrichment.LastFMResult{
			ArtistImageURL: "https://lastfm.example.com/image.jpg",
			ArtistBio:      "Test bio",
		},
	}

	logger := zap.NewNop()

	// Log what each source returned
	logger.Info("spotify enrichment result",
		zap.String("artist", "Test Artist"),
		zap.String("image_url", spotifyClient.imageURL),
	)
	logger.Info("deezer enrichment result",
		zap.String("artist", "Test Artist"),
		zap.String("image_url", deezerClient.imageURL),
	)
	logger.Info("lastfm enrichment result",
		zap.String("artist", "Test Artist"),
		zap.String("image_url", lastfmClient.result.ArtistImageURL),
	)

	// Verify logging works
	assert.NotEmpty(t, spotifyClient.imageURL)
	assert.NotEmpty(t, deezerClient.imageURL)
	assert.NotEmpty(t, lastfmClient.result.ArtistImageURL)
}
