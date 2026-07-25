package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

// escapeMBQuery escapes a value for safe use in a MusicBrainz Lucene query.
// Within quoted strings, only " and \ need escaping; Lucene treats other
// special characters as literals.
func escapeMBQuery(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `"`, `\"`)
	return replacer.Replace(s)
}

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

	// Build a proper MusicBrainz Lucene query.
	// Examples:
	//   "d4vd - Feel It" → artist:"d4vd" AND recording:"Feel It"
	//   "feel it"        → feel it (loose search across all fields)
	var luceneQuery string
	if q.Artist != "" && q.Title != "" {
		luceneQuery = fmt.Sprintf(`artist:"%s" AND recording:"%s"`,
			escapeMBQuery(q.Artist),
			escapeMBQuery(q.Title),
		)
	} else {
		// No structured info — use the raw query as-is for a loose search.
		// If the raw query looks like it has separator syntax, also try
		// a broad recording search as fallback (handled below).
		luceneQuery = escapeMBQuery(q.Raw)
	}

	u, _ := url.Parse(p.baseURL + "/recording")
	u.RawQuery = url.Values{
		"query": {luceneQuery},
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

	// If the structured search returned nothing and we had artist+title,
	// also try a loose recording search as a fallback.
	if len(mr.Recordings) == 0 && q.Artist != "" && q.Title != "" {
		// Wait for rate limiter before fallback request
		select {
		case <-p.tick:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		fallbackQuery := fmt.Sprintf(`recording:"%s"`, escapeMBQuery(q.Title))
		fallbackURL, _ := url.Parse(p.baseURL + "/recording")
		fallbackURL.RawQuery = url.Values{
			"query": {fallbackQuery},
			"fmt":   {"json"},
			"limit": {"10"},
		}.Encode()

		req2, err2 := http.NewRequestWithContext(ctx, http.MethodGet, fallbackURL.String(), nil)
		if err2 == nil {
			req2.Header.Set("User-Agent", req.UserAgent())
			resp2, err2 := p.client.Do(req2)
			if err2 == nil {
				var mr2 mbSearchResponse
				if json.NewDecoder(resp2.Body).Decode(&mr2) == nil {
					mr.Recordings = append(mr.Recordings, mr2.Recordings...)
				}
				resp2.Body.Close()
			}
		}
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
