package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type lastFMClient struct {
	httpClient *http.Client
	cfg        LastFMConfig
}

func NewLastFMClient(cfg LastFMConfig) LastFMClient {
	return &lastFMClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cfg:        cfg,
	}
}

// SearchArtist looks up an artist by name using Last.fm artist.getInfo.
// It returns the artist's bio, image, tags, similar artists, and stats.
func (c *lastFMClient) SearchArtist(ctx context.Context, name string) (*LastFMResult, error) {
	if c.cfg.APIKey == "" {
		return nil, fmt.Errorf("last.fm API key not configured")
	}
	if name == "" {
		return nil, fmt.Errorf("artist name is required")
	}

	params := url.Values{}
	params.Set("method", "artist.getInfo")
	params.Set("api_key", c.cfg.APIKey)
	params.Set("format", "json")
	params.Set("autocorrect", "1")
	params.Set("artist", name)

	u := fmt.Sprintf("%s?%s", c.cfg.BaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("last.fm returned %d", resp.StatusCode)
	}

	var lfmResp struct {
		Artist *struct {
			Name       string `json:"name"`
			PlayCount  string `json:"playcount"`
			Listeners  string `json:"listeners"`
			Streamable string `json:"streamable"`
			URL        string `json:"url"`
			Image      []struct {
				URL  string `json:"#text"`
				Size string `json:"size"`
			} `json:"image"`
			Bio struct {
				Summary string `json:"summary"`
			} `json:"bio"`
			Tags struct {
				Tag []struct {
					Name string `json:"name"`
				} `json:"tag"`
			} `json:"tags"`
			Similar struct {
				Artist []struct {
					Name string `json:"name"`
				} `json:"artist"`
			} `json:"similar"`
		} `json:"artist"`
		Error   *int    `json:"error,omitempty"`
		Message *string `json:"message,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&lfmResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if lfmResp.Error != nil {
		return nil, fmt.Errorf("last.fm error %d: %s", *lfmResp.Error, stringPtr(lfmResp.Message))
	}

	if lfmResp.Artist == nil {
		return nil, nil
	}

	a := lfmResp.Artist
	result := &LastFMResult{
		PlayCount:     parseInt(a.PlayCount),
		ListenerCount: parseInt(a.Listeners),
	}

	for _, tag := range a.Tags.Tag {
		name := strings.TrimSpace(tag.Name)
		if name != "" {
			result.Tags = append(result.Tags, name)
		}
	}

	for _, img := range a.Image {
		if img.Size == "extralarge" || img.Size == "large" {
			if img.URL != "" {
				result.ArtistImageURL = img.URL
				break
			}
		}
	}
	if result.ArtistImageURL == "" && len(a.Image) > 0 {
		result.ArtistImageURL = a.Image[0].URL
	}

	if a.Bio.Summary != "" {
		result.ArtistBio = truncateBio(a.Bio.Summary)
	}

	for _, sa := range a.Similar.Artist {
		saname := strings.TrimSpace(sa.Name)
		if saname != "" {
			result.SimilarArtists = append(result.SimilarArtists, saname)
		}
	}
	if len(result.SimilarArtists) > 10 {
		result.SimilarArtists = result.SimilarArtists[:10]
	}

	return result, nil
}

func (c *lastFMClient) SearchTrack(ctx context.Context, query TrackQuery) (*LastFMResult, error) {
	if c.cfg.APIKey == "" {
		return nil, fmt.Errorf("last.fm API key not configured")
	}

	params := url.Values{}
	params.Set("method", "track.getInfo")
	params.Set("api_key", c.cfg.APIKey)
	params.Set("format", "json")
	params.Set("autocorrect", "1")

	if query.Artist != "" {
		params.Set("artist", query.Artist)
	}
	if query.Title != "" {
		params.Set("track", query.Title)
	}

	u := fmt.Sprintf("%s?%s", c.cfg.BaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("last.fm returned %d", resp.StatusCode)
	}

	var lfmResp struct {
		Track *struct {
			PlayCount string `json:"playcount"`
			Listeners string `json:"listeners"`
			Album     *struct {
				Title string `json:"title"`
				Image []struct {
					URL  string `json:"#text"`
					Size string `json:"size"`
				} `json:"image"`
			} `json:"album"`
			Tags struct {
				Tag []struct {
					Name string `json:"name"`
				} `json:"tag"`
			} `json:"toptags"`
			Artist struct {
				Bio struct {
					Summary string `json:"summary"`
				} `json:"bio"`
				Image []struct {
					URL  string `json:"#text"`
					Size string `json:"size"`
				} `json:"image"`
				Similar struct {
					Artist []struct {
						Name string `json:"name"`
					} `json:"artist"`
				} `json:"similar"`
			} `json:"artist"`
		} `json:"track"`
		Error   *int    `json:"error,omitempty"`
		Message *string `json:"message,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&lfmResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if lfmResp.Error != nil {
		return nil, fmt.Errorf("last.fm error %d: %s", *lfmResp.Error, stringPtr(lfmResp.Message))
	}

	if lfmResp.Track == nil {
		return nil, nil
	}

	t := lfmResp.Track
	result := &LastFMResult{
		PlayCount:     parseInt(t.PlayCount),
		ListenerCount: parseInt(t.Listeners),
	}

	for _, tag := range t.Tags.Tag {
		name := strings.TrimSpace(tag.Name)
		if name != "" {
			result.Tags = append(result.Tags, name)
		}
	}

	if t.Album != nil {
		for _, img := range t.Album.Image {
			if img.Size == "extralarge" || img.Size == "large" {
				if img.URL != "" {
					result.AlbumCoverURL = img.URL
					break
				}
			}
		}
		if result.AlbumCoverURL == "" && len(t.Album.Image) > 0 {
			result.AlbumCoverURL = t.Album.Image[0].URL
		}
	}

	if t.Artist.Bio.Summary != "" {
		result.ArtistBio = truncateBio(t.Artist.Bio.Summary)
	}

	for _, img := range t.Artist.Image {
		if img.Size == "extralarge" || img.Size == "large" {
			if img.URL != "" {
				result.ArtistImageURL = img.URL
				break
			}
		}
	}
	if result.ArtistImageURL == "" && len(t.Artist.Image) > 0 {
		result.ArtistImageURL = t.Artist.Image[0].URL
	}

	for _, a := range t.Artist.Similar.Artist {
		name := strings.TrimSpace(a.Name)
		if name != "" {
			result.SimilarArtists = append(result.SimilarArtists, name)
		}
	}
	if len(result.SimilarArtists) > 10 {
		result.SimilarArtists = result.SimilarArtists[:10]
	}

	return result, nil
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}

func truncateBio(bio string) string {
	idx := strings.Index(bio, "<a")
	if idx > 0 {
		bio = strings.TrimSpace(bio[:idx])
	}
	runes := []rune(bio)
	if len(runes) > 500 {
		runes = runes[:500]
	}
	return strings.TrimSpace(string(runes))
}

func stringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
