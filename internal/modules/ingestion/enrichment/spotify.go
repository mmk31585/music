package enrichment

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type spotifyClient struct {
	httpClient  *http.Client
	cfg         SpotifyConfig
	mu          sync.Mutex
	token       string
	tokenExpiry time.Time
}

func NewSpotifyClient(cfg SpotifyConfig) SpotifyClient {
	return &spotifyClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cfg:        cfg,
	}
}

func (c *spotifyClient) SearchTrack(ctx context.Context, query TrackQuery) (*SpotifyResult, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, fmt.Errorf("spotify auth: %w", err)
	}

	q := query.Title
	if query.Artist != "" {
		q = fmt.Sprintf("track:%s artist:%s", query.Title, query.Artist)
	}

	params := url.Values{}
	params.Set("q", q)
	params.Set("type", "track")
	params.Set("limit", "5")

	u := fmt.Sprintf("%s/search?%s", c.cfg.BaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.mu.Lock()
		c.token = ""
		c.mu.Unlock()
		return c.SearchTrack(ctx, query)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify returned %d", resp.StatusCode)
	}

	var spotResp struct {
		Tracks struct {
			Items []struct {
				ID         string `json:"id"`
				Popularity int    `json:"popularity"`
				PreviewURL string `json:"preview_url"`
				Album      struct {
					Images []struct {
						URL string `json:"url"`
					} `json:"images"`
				} `json:"album"`
				Artists []struct {
					Name string `json:"name"`
					ID   string `json:"id"`
				} `json:"artists"`
			} `json:"items"`
		} `json:"tracks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&spotResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(spotResp.Tracks.Items) == 0 {
		return nil, nil
	}

	item := spotResp.Tracks.Items[0]
	result := &SpotifyResult{
		SpotifyID:  item.ID,
		Popularity: item.Popularity,
		PreviewURL: item.PreviewURL,
	}

	if len(item.Album.Images) > 0 {
		result.AlbumCoverURL = item.Album.Images[0].URL
	}

	if len(item.Artists) > 0 && item.Artists[0].ID != "" {
		result.ArtistImageURL = c.fetchArtistImage(ctx, item.Artists[0].ID)
	}

	return result, nil
}

func (c *spotifyClient) ensureToken(ctx context.Context) error {
	c.mu.Lock()
	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	if c.cfg.ClientID == "" || c.cfg.ClientSecret == "" {
		return fmt.Errorf("spotify client credentials not configured")
	}

	authStr := base64.StdEncoding.EncodeToString([]byte(c.cfg.ClientID + ":" + c.cfg.ClientSecret))

	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", c.cfg.AuthBaseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("create auth request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+authStr)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("spotify auth returned %d", resp.StatusCode)
	}

	var authResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("decode auth response: %w", err)
	}

	c.mu.Lock()
	c.token = authResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(authResp.ExpiresIn-60) * time.Second)
	c.mu.Unlock()

	return nil
}

// SearchArtistImage searches for an artist by name on Spotify and returns their image URL.
func (c *spotifyClient) SearchArtistImage(ctx context.Context, name string) (string, error) {
	if err := c.ensureToken(ctx); err != nil {
		return "", fmt.Errorf("spotify auth: %w", err)
	}

	cleaned := NormalizeName(name)

	params := url.Values{}
	params.Set("q", cleaned)
	params.Set("type", "artist")
	params.Set("limit", "3")

	u := fmt.Sprintf("%s/search?%s", c.cfg.BaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return "", fmt.Errorf("create spotify artist request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("spotify artist request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.mu.Lock()
		c.token = ""
		c.mu.Unlock()
		return c.SearchArtistImage(ctx, name)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("spotify artist search returned %d", resp.StatusCode)
	}

	var searchResp struct {
		Artists struct {
			Items []struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
			} `json:"items"`
		} `json:"artists"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return "", fmt.Errorf("decode spotify artist response: %w", err)
	}

	if len(searchResp.Artists.Items) == 0 {
		return "", nil
	}

	// Try to find best name match, then get high-res image
	cleanedLower := strings.ToLower(cleaned)
	for _, artist := range searchResp.Artists.Items {
		an := NormalizeName(artist.Name)
		if strings.EqualFold(an, cleaned) || strings.Contains(strings.ToLower(an), cleanedLower) {
			if img := bestSpotifyImage(artist.Images); img != "" {
				return img, nil
			}
		}
	}

	// Fallback: first result's image
	if img := bestSpotifyImage(searchResp.Artists.Items[0].Images); img != "" {
		return img, nil
	}

	return "", nil
}

func bestSpotifyImage(images []struct {
	URL string `json:"url"`
}) string {
	// Prefer the largest image (last in array — Spotify returns ascending order)
	for i := len(images) - 1; i >= 0; i-- {
		if images[i].URL != "" {
			return images[i].URL
		}
	}
	return ""
}

func (c *spotifyClient) fetchArtistImage(ctx context.Context, artistID string) string {
	u := fmt.Sprintf("%s/artists/%s", c.cfg.BaseURL, artistID)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var artistResp struct {
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&artistResp); err != nil {
		return ""
	}

	if len(artistResp.Images) > 0 {
		return artistResp.Images[0].URL
	}
	return ""
}
