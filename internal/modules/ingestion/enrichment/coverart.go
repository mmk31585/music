package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type coverArtClient struct {
	httpClient *http.Client
	mbBaseURL  string
	caaBaseURL string
	userAgent  string
	caaClient  *http.Client
}

func NewCoverArtClient() CoverArtClient {
	return &coverArtClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		mbBaseURL:  "https://musicbrainz.org/ws/2",
		caaBaseURL: "https://coverartarchive.org",
		userAgent:  "MuseMusic/1.0 (music@muse.app)",
		caaClient: &http.Client{
			Timeout: 5 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *coverArtClient) SearchAlbumCover(ctx context.Context, albumTitle, artistName string) (string, error) {
	if albumTitle == "" {
		return "", nil
	}

	// Step 1: Search MusicBrainz for the release
	releaseID, err := c.searchRelease(ctx, albumTitle, artistName)
	if err != nil || releaseID == "" {
		return "", err
	}

	// Step 2: Try Cover Art Archive for the release
	return c.fetchCoverFromCAA(ctx, releaseID)
}

func (c *coverArtClient) searchRelease(ctx context.Context, albumTitle, artistName string) (string, error) {
	query := fmt.Sprintf("release:\"%s\" AND artist:\"%s\"",
		url.QueryEscape(albumTitle),
		url.QueryEscape(artistName),
	)
	u := fmt.Sprintf("%s/release?query=%s&fmt=json&limit=5", c.mbBaseURL, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("create musicbrainz request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("musicbrainz request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	var mbResp struct {
		Releases []struct {
			ID string `json:"id"`
		} `json:"releases"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mbResp); err != nil {
		return "", fmt.Errorf("musicbrainz decode: %w", err)
	}

	if len(mbResp.Releases) == 0 {
		return "", nil
	}

	return mbResp.Releases[0].ID, nil
}

func (c *coverArtClient) fetchCoverFromCAA(ctx context.Context, releaseID string) (string, error) {
	// Try first result from CAA
	caaURL := fmt.Sprintf("%s/release/%s/front", c.caaBaseURL, releaseID)
	caaReq, err := http.NewRequestWithContext(ctx, http.MethodGet, caaURL, nil)
	if err != nil {
		return "", nil
	}
	caaResp, err := c.caaClient.Do(caaReq)
	if err != nil {
		return "", nil
	}
	caaResp.Body.Close()

	if loc := redirectLocation(caaResp); loc != "" {
		return loc, nil
	}
	if caaResp.StatusCode == http.StatusOK {
		return caaURL, nil
	}

	return "", nil
}

func redirectLocation(resp *http.Response) string {
	if resp.StatusCode == http.StatusSeeOther ||
		resp.StatusCode == http.StatusTemporaryRedirect ||
		resp.StatusCode == http.StatusFound {
		return resp.Header.Get("Location")
	}
	return ""
}
