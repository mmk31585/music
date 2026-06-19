package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type DeezerProvider struct {
	client  *http.Client
	baseURL string
}

func NewDeezerProvider() *DeezerProvider {
	return &DeezerProvider{
		client:  &http.Client{Timeout: 8 * time.Second},
		baseURL: "https://api.deezer.com",
	}
}

func (p *DeezerProvider) Name() string { return "deezer" }

type deezerTrack struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	Preview  string `json:"preview"`
	Duration int    `json:"duration"`
	Artist   struct {
		Name string `json:"name"`
	} `json:"artist"`
	Album struct {
		Title string `json:"title"`
		Cover string `json:"cover_medium"`
	} `json:"album"`
}

type deezerResponse struct {
	Data []deezerTrack `json:"data"`
}

func (p *DeezerProvider) Search(ctx context.Context, q SearchQuery) ([]Result, error) {
	query := q.Raw
	if query == "" {
		query = q.Artist + " " + q.Title
	}

	u, _ := url.Parse(p.baseURL + "/search/track")
	u.RawQuery = url.Values{"q": {query}, "limit": {"10"}}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("deezer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deezer: status %d", resp.StatusCode)
	}

	var dr deezerResponse
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return nil, fmt.Errorf("deezer decode: %w", err)
	}

	results := make([]Result, 0, len(dr.Data))
	for _, t := range dr.Data {
		results = append(results, Result{
			Title:     t.Title,
			Artist:    t.Artist.Name,
			Album:     t.Album.Title,
			URL:       t.Link,
			Duration:  t.Duration,
			Thumbnail: t.Album.Cover,
			Source:    "deezer",
			ExternalIDs: ExternalIDs{
				DeezerID: fmt.Sprintf("%d", t.ID),
			},
		})
	}

	return results, nil
}
