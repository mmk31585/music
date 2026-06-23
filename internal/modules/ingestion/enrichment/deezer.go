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

type deezerClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewDeezerClient() DeezerClient {
	return &deezerClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		baseURL:    "https://api.deezer.com",
	}
}

type deezerArtist struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture_medium"`
}

type deezerSearch struct {
	Data []deezerArtist `json:"data"`
}

func (c *deezerClient) SearchArtistImage(ctx context.Context, name string) (string, error) {
	searchURL := fmt.Sprintf("%s/search/artist?q=%s&limit=3", c.baseURL, url.QueryEscape(name))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return "", fmt.Errorf("create deezer request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("deezer request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deezer: status %d", resp.StatusCode)
	}

	var sr deezerSearch
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", fmt.Errorf("deezer decode: %w", err)
	}

	if len(sr.Data) == 0 {
		return "", nil
	}

	// Find best match
	nameLower := strings.ToLower(name)
	for _, a := range sr.Data {
		if strings.EqualFold(a.Name, name) || strings.Contains(strings.ToLower(a.Name), nameLower) {
			if a.Picture != "" {
				return a.Picture, nil
			}
		}
	}

	// Fallback to first result
	if sr.Data[0].Picture != "" {
		return sr.Data[0].Picture, nil
	}

	return "", nil
}
