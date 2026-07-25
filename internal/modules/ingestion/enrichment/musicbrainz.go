package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type musicBrainzClient struct {
	httpClient *http.Client
	cfg        MusicBrainzConfig
	rateLimit  <-chan time.Time
}

func NewMusicBrainzClient(cfg MusicBrainzConfig) MusicBrainzClient {
	return &musicBrainzClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		cfg:        cfg,
		rateLimit:  time.Tick(time.Second),
	}
}

func (c *musicBrainzClient) SearchRecording(ctx context.Context, query TrackQuery) (*MusicBrainzResult, error) {
	select {
	case <-c.rateLimit:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	params := url.Values{}
	queryParts := []string{}
	if query.Title != "" {
		queryParts = append(queryParts, fmt.Sprintf(`"%s"`, query.Title))
	}
	if query.Artist != "" {
		queryParts = append(queryParts, fmt.Sprintf(`artist:"%s"`, query.Artist))
	}
	if len(queryParts) == 0 {
		return nil, fmt.Errorf("empty query")
	}

	params.Set("query", strings.Join(queryParts, " AND "))
	params.Set("fmt", "json")
	params.Set("limit", "5")

	u := fmt.Sprintf("%s/recording?%s", c.cfg.BaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("musicbrainz returned %d", resp.StatusCode)
	}

	var mbResp struct {
		Recordings []struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			Length *int   `json:"length,omitempty"`
			ArtistCredit []struct {
				Name string `json:"name"`
				Artist struct {
					ID string `json:"id"`
				} `json:"artist"`
			} `json:"artist-credit"`
			Releases []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
				Date  string `json:"date"`
			} `json:"releases"`
			Tags []struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"tags"`
			ArtistMBID string `json:"-"`
		} `json:"recordings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&mbResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(mbResp.Recordings) == 0 {
		return nil, nil
	}

	r := mbResp.Recordings[0]
	result := &MusicBrainzResult{
		MBID:   r.ID,
		Title:  r.Title,
		Genres: extractTags(r.Tags),
	}

	if r.Length != nil {
		result.Duration = *r.Length / 1000
	}

	if len(r.ArtistCredit) > 0 {
		result.ArtistName = r.ArtistCredit[0].Name
		result.ArtistMBID = r.ArtistCredit[0].Artist.ID
	}

	if len(r.Releases) > 0 {
		result.AlbumName = r.Releases[0].Title
		result.AlbumMBID = r.Releases[0].ID
		if year := extractYear(r.Releases[0].Date); year > 0 {
			result.ReleaseYear = year
		}
	}

	return result, nil
}

func extractTags(tags []struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}) []string {
	var result []string
	seen := map[string]bool{}
	for _, t := range tags {
		name := strings.TrimSpace(t.Name)
		if name != "" && !seen[name] {
			seen[name] = true
			result = append(result, name)
		}
	}
	if len(result) > 10 {
		result = result[:10]
	}
	return result
}

func extractYear(date string) int {
	if date == "" {
		return 0
	}
	parts := strings.SplitN(date, "-", 2)
	if len(parts) == 0 {
		return 0
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return year
}
