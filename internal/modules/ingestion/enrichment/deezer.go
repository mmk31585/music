package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
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

type deezerAlbum struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Cover  string `json:"cover_medium"`
	Artist struct {
		Name string `json:"name"`
	} `json:"artist"`
}

type deezerSearchData[T any] struct {
	Data []T `json:"data"`
}

// NormalizeName strips invisible/format characters (soft hyphens, zero-width spaces, etc.)
// that can be embedded in artist/album names and break external API matching.
func NormalizeName(s string) string {
	return strings.Map(func(r rune) rune {
		// Remove format characters (Cf): soft hyphen, zero-width spaces, etc.
		if unicode.Is(unicode.Cf, r) {
			return -1
		}
		return r
	}, strings.TrimSpace(s))
}

func (c *deezerClient) SearchArtistImage(ctx context.Context, name string) (string, error) {
	cleaned := NormalizeName(name)
	searchURL := fmt.Sprintf("%s/search/artist?q=%s&limit=3", c.baseURL, url.QueryEscape(cleaned))

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

	var sr deezerSearchData[deezerArtist]
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", fmt.Errorf("deezer decode: %w", err)
	}

	if len(sr.Data) == 0 {
		return "", nil
	}

	// Find best match using normalized names
	cleanedLower := strings.ToLower(cleaned)
	for _, a := range sr.Data {
		an := NormalizeName(a.Name)
		if strings.EqualFold(an, cleaned) || strings.Contains(strings.ToLower(an), cleanedLower) {
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

// SearchAlbumCover searches for an album cover via Deezer's album search API.
func (c *deezerClient) SearchAlbumCover(ctx context.Context, albumTitle, artistName string) (string, error) {
	cleanedTitle := NormalizeName(albumTitle)
	cleanedArtist := NormalizeName(artistName)
	q := cleanedTitle
	if cleanedArtist != "" {
		q = fmt.Sprintf(`%s "%s"`, cleanedTitle, cleanedArtist)
	}

	searchURL := fmt.Sprintf("%s/search/album?q=%s&limit=3", c.baseURL, url.QueryEscape(q))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return "", fmt.Errorf("create deezer album request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("deezer album request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deezer album: status %d", resp.StatusCode)
	}

	var sr deezerSearchData[deezerAlbum]
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", fmt.Errorf("deezer album decode: %w", err)
	}

	if len(sr.Data) == 0 {
		return "", nil
	}

	// Try to find best match by comparing title + artist
	titleLower := strings.ToLower(cleanedTitle)
	artistLower := strings.ToLower(cleanedArtist)
	for _, al := range sr.Data {
		alTitle := NormalizeName(al.Title)
		alArtist := NormalizeName(al.Artist.Name)
		if strings.EqualFold(alTitle, cleanedTitle) ||
			strings.Contains(strings.ToLower(alTitle), titleLower) {
			// If we have an artist name, prefer matching artist too
			if cleanedArtist == "" || strings.EqualFold(alArtist, cleanedArtist) ||
				strings.Contains(strings.ToLower(alArtist), artistLower) {
				if al.Cover != "" {
					return al.Cover, nil
				}
			}
		}
	}

	// Fallback: first result with a cover
	for _, al := range sr.Data {
		if al.Cover != "" {
			return al.Cover, nil
		}
	}

	return "", nil
}
