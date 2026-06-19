package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MusicBrainzProvider struct {
	client  *http.Client
	baseURL string
	appName string
	tick    <-chan time.Time
}

func NewMusicBrainzProvider() *MusicBrainzProvider {
	return &MusicBrainzProvider{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: "https://musicbrainz.org/ws/2",
		appName: "musicapp/1.0",
		tick:    time.Tick(time.Second),
	}
}

func (p *MusicBrainzProvider) Name() string { return "musicbrainz" }

type mbRecording struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Length       *int   `json:"length"`
	Score        int    `json:"score"`
	ArtistCredit []struct {
		Name string `json:"name"`
	} `json:"artist-credit"`
	Releases []struct {
		Title           string `json:"title"`
		ID              string `json:"id"`
		CoverArtArchive struct {
			Front bool `json:"front"`
		} `json:"cover-art-archive"`
	} `json:"releases"`
	ISRCs []struct {
		ISRC string `json:"isrc"`
	} `json:"isrcs"`
}

type mbSearchResponse struct {
	Recordings []mbRecording `json:"recordings"`
}

func (p *MusicBrainzProvider) Search(ctx context.Context, q SearchQuery) ([]Result, error) {
	select {
	case <-p.tick:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	query := q.Raw
	if query == "" {
		query = q.Artist + " " + q.Title
	}

	u, _ := url.Parse(p.baseURL + "/recording")
	u.RawQuery = url.Values{
		"query": {fmt.Sprintf(`"%s"`, query)},
		"fmt":   {"json"},
		"limit": {"10"},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", p.appName+" (music@example.com)")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("musicbrainz: status %d", resp.StatusCode)
	}

	var mr mbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
		return nil, fmt.Errorf("musicbrainz decode: %w", err)
	}

	results := make([]Result, 0, len(mr.Recordings))
	for _, rec := range mr.Recordings {
		artistName := ""
		if len(rec.ArtistCredit) > 0 {
			artistName = rec.ArtistCredit[0].Name
		}

		albumName := ""
		if len(rec.Releases) > 0 {
			albumName = rec.Releases[0].Title
		}

		duration := 0
		if rec.Length != nil {
			duration = *rec.Length / 1000
		}

		isrc := ""
		if len(rec.ISRCs) > 0 {
			isrc = rec.ISRCs[0].ISRC
		}

		results = append(results, Result{
			Title:    rec.Title,
			Artist:   artistName,
			Album:    albumName,
			Duration: duration,
			Source:   "musicbrainz",
			ISRC:     isrc,
			ExternalIDs: ExternalIDs{
				MBID: rec.ID,
				ISRC: isrc,
			},
		})
	}

	return results, nil
}
