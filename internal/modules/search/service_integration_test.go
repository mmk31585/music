package search

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	osclient "music/internal/platform/opensearch"
)

func TestSearch_WithMockOS_HappyPath(t *testing.T) {
	ts := newMockOSServer(t)
	defer ts.Close()

	client := newOSClientForTest(t, ts.URL)
	svc := NewService(client, nil)

	resp, err := svc.Search(context.Background(), "test query", "all", 10)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "test query", resp.Query)

	require.Len(t, resp.Tracks, 1)
	assert.Equal(t, "track-1", resp.Tracks[0].ID)
	assert.Equal(t, "Test Track", resp.Tracks[0].Title)
	assert.Equal(t, strPtr("artist-1"), resp.Tracks[0].ArtistID)
	assert.Equal(t, strPtr("Test Artist"), resp.Tracks[0].ArtistName)

	require.Len(t, resp.Albums, 1)
	assert.Equal(t, "album-1", resp.Albums[0].ID)
	assert.Equal(t, "Test Album", resp.Albums[0].Title)

	require.Len(t, resp.Artists, 1)
	assert.Equal(t, "artist-1", resp.Artists[0].ID)
	assert.Equal(t, "Test Artist", resp.Artists[0].Name)

	require.Len(t, resp.Playlists, 1)
	assert.Equal(t, "playlist-1", resp.Playlists[0].ID)
	assert.Equal(t, "Test Playlist", resp.Playlists[0].Name)
}

func TestSearch_WithMockOS_TracksOnly(t *testing.T) {
	ts := newMockOSServer(t)
	defer ts.Close()

	client := newOSClientForTest(t, ts.URL)
	svc := NewService(client, nil)

	resp, err := svc.Search(context.Background(), "test", "tracks", 10)
	require.NoError(t, err)
	require.Len(t, resp.Tracks, 1)
	assert.Empty(t, resp.Albums)
	assert.Empty(t, resp.Artists)
	assert.Empty(t, resp.Playlists)
}

func TestSearch_WithMockOS_AlbumsOnly(t *testing.T) {
	ts := newMockOSServer(t)
	defer ts.Close()

	client := newOSClientForTest(t, ts.URL)
	svc := NewService(client, nil)

	resp, err := svc.Search(context.Background(), "test", "albums", 10)
	require.NoError(t, err)
	assert.Empty(t, resp.Tracks)
	require.Len(t, resp.Albums, 1)
	assert.Empty(t, resp.Artists)
	assert.Empty(t, resp.Playlists)
}

func TestSearch_WithMockOS_ArtistsOnly(t *testing.T) {
	ts := newMockOSServer(t)
	defer ts.Close()

	client := newOSClientForTest(t, ts.URL)
	svc := NewService(client, nil)

	resp, err := svc.Search(context.Background(), "test", "artists", 10)
	require.NoError(t, err)
	assert.Empty(t, resp.Tracks)
	assert.Empty(t, resp.Albums)
	require.Len(t, resp.Artists, 1)
	assert.Empty(t, resp.Playlists)
}

func TestSearch_WithMockOS_PlaylistsOnly(t *testing.T) {
	ts := newMockOSServer(t)
	defer ts.Close()

	client := newOSClientForTest(t, ts.URL)
	svc := NewService(client, nil)

	resp, err := svc.Search(context.Background(), "test", "playlists", 10)
	require.NoError(t, err)
	assert.Empty(t, resp.Tracks)
	assert.Empty(t, resp.Albums)
	assert.Empty(t, resp.Artists)
	require.Len(t, resp.Playlists, 1)
}

func TestSearch_WithMockOS_EmptyResults(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			fmt.Fprint(w, `{"name":"test-cluster","version":{"number":"7.10.0","distribution":"opensearch"}}`)
			return
		}
		fmt.Fprint(w, `{"took":1,"timed_out":false,"hits":{"total":{"value":0},"hits":[]}}`)
	}))
	defer ts.Close()

	client := newOSClientForTest(t, ts.URL)
	svc := NewService(client, nil)

	resp, err := svc.Search(context.Background(), "nonexistent", "all", 10)
	require.NoError(t, err)
	assert.Empty(t, resp.Tracks)
	assert.Empty(t, resp.Albums)
	assert.Empty(t, resp.Artists)
	assert.Empty(t, resp.Playlists)
}

func TestSearch_WithMockOS_RankedQueryBuilt(t *testing.T) {
	var requestBodies []string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			fmt.Fprint(w, `{"name":"test-cluster","version":{"number":"7.10.0","distribution":"opensearch"}}`)
			return
		}
		bodyBytes, _ := io.ReadAll(r.Body)
		requestBodies = append(requestBodies, string(bodyBytes))
		fmt.Fprint(w, `{"took":1,"timed_out":false,"hits":{"total":{"value":0},"hits":[]}}`)
	}))
	defer ts.Close()

	client := newOSClientForTest(t, ts.URL)
	svc := NewService(client, &Repository{})

	_, err := svc.Search(context.Background(), "hello world", "all", 10)
	require.NoError(t, err)

	// Find the track search body (contains "title^3" field or "function_score")
	var trackBody string
	for _, b := range requestBodies {
		if strings.Contains(b, "title") && strings.Contains(b, "function_score") {
			trackBody = b
			break
		}
	}
	require.NotEmpty(t, trackBody, "should have found a track search request. bodies: %v", requestBodies)

	assert.Contains(t, trackBody, "function_score")
	assert.Contains(t, trackBody, "field_value_factor")
	assert.Contains(t, trackBody, "popularity")
	assert.Contains(t, trackBody, "gauss")
	assert.Contains(t, trackBody, "multi_match")
	assert.Contains(t, trackBody, "fuzziness")
}

func TestSearch_WithMockOS_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			fmt.Fprint(w, `{"name":"test-cluster","version":{"number":"7.10.0","distribution":"opensearch"}}`)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":"internal"}`)
	}))
	defer ts.Close()

	client := newOSClientForTest(t, ts.URL)
	svc := NewService(client, &Repository{})

	resp, err := svc.Search(context.Background(), "test", "tracks", 10)
	require.NoError(t, err)
	assert.Empty(t, resp.Tracks)
}

func TestSuggestions_WithMockOS(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			fmt.Fprint(w, `{"name":"test-cluster","version":{"number":"7.10.0","distribution":"opensearch"}}`)
			return
		}
		fmt.Fprint(w, `{
			"took":1,"timed_out":false,
			"hits":{"total":{"value":2},"hits":[
				{"_source":{"title":"Test Track"}},
				{"_source":{"name":"Test Artist"}}
			]}
		}`)
	}))
	defer ts.Close()

	client := newOSClientForTest(t, ts.URL)
	svc := NewService(client, nil)

	suggestions, err := svc.Suggestions(context.Background(), "test", 5)
	require.NoError(t, err)
	require.Len(t, suggestions, 2)
	assert.Equal(t, "Test Track", suggestions[0])
	assert.Equal(t, "Test Artist", suggestions[1])
}

// -- helpers --

func newMockOSServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			fmt.Fprint(w, `{"name":"test-cluster","version":{"number":"7.10.0","distribution":"opensearch"}}`)
			return
		}

		index := extractIndex(r.URL.Path)
		fmt.Fprint(w, buildSearchResponse(index))
	}))
}

func newOSClientForTest(t *testing.T, url string) *osclient.Client {
	t.Helper()

	client, err := osclient.NewClient(osclient.Config{
		URL: url,
	}, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, client)
	return client
}

func extractIndex(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for _, p := range parts {
		if p != "_search" {
			return p
		}
	}
	return ""
}

func buildSearchResponse(index string) string {
	switch index {
	case "tracks":
		return `{
			"took":2,"timed_out":false,
			"hits":{"total":{"value":1},"hits":[
				{"_source":{"id":"track-1","title":"Test Track","artist_id":"artist-1","artist_name":"Test Artist","album_id":"album-1","album_title":"Test Album"}}
			]}
		}`
	case "albums":
		return `{
			"took":1,"timed_out":false,
			"hits":{"total":{"value":1},"hits":[
				{"_source":{"id":"album-1","title":"Test Album","artist_id":"artist-1","artist_name":"Test Artist"}}
			]}
		}`
	case "artists":
		return `{
			"took":1,"timed_out":false,
			"hits":{"total":{"value":1},"hits":[
				{"_source":{"id":"artist-1","name":"Test Artist"}}
			]}
		}`
	case "playlists":
		return `{
			"took":1,"timed_out":false,
			"hits":{"total":{"value":1},"hits":[
				{"_source":{"id":"playlist-1","name":"Test Playlist","is_public":true}}
			]}
		}`
	default:
		return `{"took":0,"timed_out":false,"hits":{"total":{"value":0},"hits":[]}}`
	}
}

func strPtr(s string) *string {
	return &s
}
