package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type SpotifyProvider struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client

	mu        sync.RWMutex
	token     string
	tokenTime time.Time
}

func NewSpotifyProvider(clientID, clientSecret string) *SpotifyProvider {
	return &SpotifyProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *SpotifyProvider) Name() string { return "spotify" }

type spotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type spotifyTrack struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URI         string `json:"uri"`
	Explicit    bool   `json:"explicit"`
	DurationMs  int    `json:"duration_ms"`
	TrackNumber int    `json:"track_number"`
	Artists     []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"artists"`
	Album struct {
		Name   string `json:"name"`
		ID     string `json:"id"`
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	} `json:"album"`
	ExternalIDs struct {
		ISRC string `json:"isrc"`
	} `json:"external_ids"`
}

type spotifySearchResponse struct {
	Tracks struct {
		Items []spotifyTrack `json:"items"`
	} `json:"tracks"`
}

func (p *SpotifyProvider) getToken(ctx context.Context) (string, error) {
	p.mu.RLock()
	if p.token != "" && time.Since(p.tokenTime) < 50*time.Minute {
		t := p.token
		p.mu.RUnlock()
		return t, nil
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != "" && time.Since(p.tokenTime) < 50*time.Minute {
		return p.token, nil
	}

	u := "https://accounts.spotify.com/api/token"
	form := url.Values{"grant_type": {"client_credentials"}}.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, io.NopCloser(strings.NewReader(form)))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(p.clientID, p.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tr spotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}

	p.token = tr.AccessToken
	p.tokenTime = time.Now()
	return p.token, nil
}

func (p *SpotifyProvider) Search(ctx context.Context, q SearchQuery) ([]Result, error) {
	query := q.Raw
	if query == "" {
		query = q.Artist + " " + q.Title
	}

	token, err := p.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("spotify auth: %w", err)
	}

	u, _ := url.Parse("https://api.spotify.com/v1/search")
	u.RawQuery = url.Values{"q": {query}, "type": {"track"}, "limit": {"10"}}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("spotify search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify: status %d", resp.StatusCode)
	}

	var sr spotifySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("spotify decode: %w", err)
	}

	results := make([]Result, 0, len(sr.Tracks.Items))
	for _, t := range sr.Tracks.Items {
		thumbnail := ""
		if len(t.Album.Images) > 0 {
			thumbnail = t.Album.Images[0].URL
		}
		artistName := ""
		if len(t.Artists) > 0 {
			artistName = t.Artists[0].Name
		}

		results = append(results, Result{
			Title:     t.Name,
			Artist:    artistName,
			Album:     t.Album.Name,
			URL:       "https://open.spotify.com/track/" + t.ID,
			Duration:  t.DurationMs / 1000,
			Thumbnail: thumbnail,
			Source:    "spotify",
			ISRC:      t.ExternalIDs.ISRC,
			ExternalIDs: ExternalIDs{
				SpotifyID: t.ID,
				ISRC:      t.ExternalIDs.ISRC,
			},
		})
	}

	return results, nil
}
