package importcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ── Deezer API types ──────────────────────────────────────────────

type deezerArtist struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture_medium"`
}

type deezerArtistSearch struct {
	Data []deezerArtist `json:"data"`
}

type deezerAlbum struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Cover string `json:"cover_medium"`
}

type deezerAlbums struct {
	Data []deezerAlbum `json:"data"`
}

type deezerTrack struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

type deezerAlbumTracks struct {
	Data []deezerTrack `json:"data"`
}

// ── MusicBrainz API types ─────────────────────────────────────────

type mbArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type mbArtistSearch struct {
	Artists []mbArtist `json:"artists"`
}

type mbRecording struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Length   *int   `json:"length"`
	Releases []struct {
		Title string `json:"title"`
		ID    string `json:"id"`
	} `json:"releases"`
}

type mbArtistReleases struct {
	Recordings []mbRecording `json:"recordings"`
	Releases   []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"releases"`
}

// ── Public types ──────────────────────────────────────────────────

type ArtistDiscography struct {
	ArtistInfo ArtistInfo   `json:"artist_info"`
	Albums     []AlbumGroup `json:"albums"`
}

type ArtistInfo struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type AlbumGroup struct {
	Title  string        `json:"title"`
	Cover  string        `json:"cover"`
	Source string        `json:"source"`
	Tracks []TrackResult `json:"tracks"`
}

type TrackResult struct {
	Title       string            `json:"title"`
	Duration    int               `json:"duration"`
	Source      string            `json:"source"`
	Album       string            `json:"album,omitempty"`
	ExternalIDs map[string]string `json:"external_ids,omitempty"`
}

// ── Artist search service ─────────────────────────────────────────

type ArtistSearcher struct {
	logger *zap.Logger
	client *http.Client
	mbTick <-chan time.Time
}

func NewArtistSearcher(logger *zap.Logger) *ArtistSearcher {
	return &ArtistSearcher{
		logger: logger,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
			},
		},
		mbTick: time.Tick(time.Second), // MusicBrainz rate limit: 1 req/s
	}
}

// SearchArtist finds all albums and tracks for an artist using
// Deezer (free, no key) and MusicBrainz (free, no key) in parallel.
// Falls back to iTunes if both fail.
// Returns whichever returns first with results.
//
// Each provider gets its own short timeout so a slow provider
// (e.g. Deezer when unreachable) doesn't block the other.
func (as *ArtistSearcher) SearchArtist(ctx context.Context, name string) (*ArtistDiscography, error) {
	start := time.Now()

	type searchResult struct {
		result *ArtistDiscography
		err    error
	}

	deezerCh := make(chan searchResult, 1)
	mbCh := make(chan searchResult, 1)

	// Use independent timeouts per provider so Deezer hanging
	// doesn't block MusicBrainz from returning results quickly.
	deezerCtx, deezerCancel := context.WithTimeout(ctx, 8*time.Second)
	defer deezerCancel()
	mbCtx, mbCancel := context.WithTimeout(ctx, 25*time.Second)
	defer mbCancel()

	// Try both providers in parallel
	go func() {
		r, err := as.searchDeezer(deezerCtx, name)
		deezerCh <- searchResult{r, err}
	}()
	go func() {
		r, err := as.searchMusicBrainz(mbCtx, name)
		mbCh <- searchResult{r, err}
	}()

	// Wait for whichever completes first with results
	var lastErr error

	// MusicBrainz is the primary — wait for it unconditionally,
	// but use Deezer if it returns first with results.
	for pending := 2; pending > 0; pending-- {
		select {
		case dr := <-deezerCh:
			if dr.err == nil && dr.result != nil && len(dr.result.Albums) > 0 {
				as.logger.Info("artist discography found via Deezer",
					zap.String("artist", name),
					zap.Int("albums", len(dr.result.Albums)),
					zap.Duration("elapsed", time.Since(start)),
				)
				return dr.result, nil
			}
			if dr.err != nil {
				lastErr = dr.err
				as.logger.Warn("Deezer search failed",
					zap.String("artist", name),
					zap.Error(dr.err),
					zap.Duration("elapsed", time.Since(start)),
				)
			}
		case mr := <-mbCh:
			if mr.err == nil && mr.result != nil && len(mr.result.Albums) > 0 {
				as.logger.Info("artist discography found via MusicBrainz",
					zap.String("artist", name),
					zap.Int("albums", len(mr.result.Albums)),
					zap.Duration("elapsed", time.Since(start)),
				)
				return mr.result, nil
			}
			if mr.err != nil {
				lastErr = mr.err
				as.logger.Warn("MusicBrainz search failed",
					zap.String("artist", name),
					zap.Error(mr.err),
					zap.Duration("elapsed", time.Since(start)),
				)
			}
		case <-ctx.Done():
			as.logger.Warn("artist search cancelled",
				zap.String("artist", name),
				zap.Duration("elapsed", time.Since(start)),
			)
			if lastErr != nil {
				return nil, fmt.Errorf("search timed out: %w", lastErr)
			}
			return nil, fmt.Errorf("search timed out for artist '%s'", name)
		}
	}

	// Both failed — try iTunes as fallback
	as.logger.Info("both Deezer and MusicBrainz failed, trying iTunes fallback",
		zap.String("artist", name),
		zap.Duration("elapsed", time.Since(start)),
	)

	itunesCtx, itunesCancel := context.WithTimeout(ctx, 10*time.Second)
	defer itunesCancel()

	itunesResult, itunesErr := as.searchiTunes(itunesCtx, name)
	if itunesErr == nil && itunesResult != nil && len(itunesResult.Albums) > 0 {
		as.logger.Info("artist discography found via iTunes",
			zap.String("artist", name),
			zap.Int("albums", len(itunesResult.Albums)),
			zap.Duration("elapsed", time.Since(start)),
		)
		return itunesResult, nil
	}

	if itunesErr != nil {
		as.logger.Warn("iTunes fallback also failed",
			zap.String("artist", name),
			zap.Error(itunesErr),
		)
	}

	// All providers failed
	if lastErr != nil {
		return nil, fmt.Errorf("all providers failed for '%s': %w", name, lastErr)
	}
	return nil, fmt.Errorf("no results found for artist '%s' on any provider", name)
}

// ── Deezer implementation (parallel album track fetching) ─────────

func (as *ArtistSearcher) searchDeezer(ctx context.Context, name string) (*ArtistDiscography, error) {
	start := time.Now()

	// Step 1: Find artist
	artist, err := as.findDeezerArtist(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("find deezer artist: %w", err)
	}
	if artist == nil {
		return nil, fmt.Errorf("artist '%s' not found on Deezer", name)
	}

	as.logger.Info("found deezer artist",
		zap.String("name", artist.Name),
		zap.Int("id", artist.ID),
	)

	info := ArtistInfo{
		Name:  artist.Name,
		Image: artist.Picture,
	}

	// Step 2: Get albums
	albums, err := as.getDeezerAlbums(ctx, artist.ID)
	if err != nil {
		return nil, fmt.Errorf("get deezer albums: %w", err)
	}

	as.logger.Info("got deezer albums",
		zap.Int("count", len(albums)),
		zap.Duration("elapsed", time.Since(start)),
	)

	if len(albums) == 0 {
		return &ArtistDiscography{
			ArtistInfo: info,
			Albums:     []AlbumGroup{},
		}, nil
	}

	// Limit albums to prevent excessive API calls
	maxAlbums := 20
	if len(albums) > maxAlbums {
		as.logger.Info("limiting albums for performance",
			zap.Int("total", len(albums)),
			zap.Int("max", maxAlbums),
		)
		albums = albums[:maxAlbums]
	}

	// Step 3: Get tracks for all albums in parallel using goroutines
	type albumResult struct {
		title  string
		cover  string
		tracks []TrackResult
		err    error
	}

	resultCh := make(chan albumResult, len(albums))
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for _, album := range albums {
		album := album
		go func() {
			tracks, err := as.getDeezerAlbumTracks(ctx, album.ID)
			if err != nil {
				as.logger.Warn("failed to get deezer album tracks, skipping",
					zap.String("album", album.Title),
					zap.Int("album_id", album.ID),
					zap.Error(err),
				)
				resultCh <- albumResult{title: album.Title, err: err}
				return
			}

			trackResults := make([]TrackResult, 0, len(tracks))
			for _, t := range tracks {
				trackResults = append(trackResults, TrackResult{
					Title:    t.Title,
					Duration: t.Duration,
					Source:   "deezer",
					Album:    album.Title,
					ExternalIDs: map[string]string{
						"deezer_id": fmt.Sprintf("%d", t.ID),
					},
				})
			}

			resultCh <- albumResult{
				title:  album.Title,
				cover:  album.Cover,
				tracks: trackResults,
			}
		}()
	}

	// Collect results
	albumGroups := make([]AlbumGroup, 0, len(albums))
	for i := 0; i < len(albums); i++ {
		select {
		case res := <-resultCh:
			if res.err == nil && len(res.tracks) > 0 {
				albumGroups = append(albumGroups, AlbumGroup{
					Title:  res.title,
					Cover:  res.cover,
					Source: "deezer",
					Tracks: res.tracks,
				})
			}
		case <-ctx.Done():
			as.logger.Warn("deezer album track fetching cancelled",
				zap.Duration("elapsed", time.Since(start)),
			)
			break
		}
	}

	as.logger.Info("deezer search completed",
		zap.Int("albums", len(albumGroups)),
		zap.Duration("elapsed", time.Since(start)),
	)

	return &ArtistDiscography{
		ArtistInfo: info,
		Albums:     albumGroups,
	}, nil
}

func (as *ArtistSearcher) findDeezerArtist(ctx context.Context, name string) (*deezerArtist, error) {
	u, _ := url.Parse("https://api.deezer.com/search/artist")
	u.RawQuery = url.Values{"q": {name}, "limit": {"5"}}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := as.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deezer: status %d", resp.StatusCode)
	}

	var sr deezerArtistSearch
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("deezer decode: %w", err)
	}

	if len(sr.Data) == 0 {
		return nil, nil
	}

	// Find best match by name (case-insensitive)
	nameLower := strings.ToLower(name)
	for _, a := range sr.Data {
		if strings.EqualFold(a.Name, name) || strings.Contains(strings.ToLower(a.Name), nameLower) {
			return &a, nil
		}
	}

	// Return first result if no exact match
	return &sr.Data[0], nil
}

func (as *ArtistSearcher) getDeezerAlbums(ctx context.Context, artistID int) ([]deezerAlbum, error) {
	u := fmt.Sprintf("https://api.deezer.com/artist/%d/albums?limit=50", artistID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := as.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deezer albums: status %d", resp.StatusCode)
	}

	var result deezerAlbums
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("deezer albums decode: %w", err)
	}

	return result.Data, nil
}

func (as *ArtistSearcher) getDeezerAlbumTracks(ctx context.Context, albumID int) ([]deezerTrack, error) {
	u := fmt.Sprintf("https://api.deezer.com/album/%d/tracks?limit=50", albumID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := as.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deezer tracks: status %d", resp.StatusCode)
	}

	var result deezerAlbumTracks
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("deezer tracks decode: %w", err)
	}

	return result.Data, nil
}

// ── MusicBrainz implementation (fallback) ─────────────────────────

func (as *ArtistSearcher) searchMusicBrainz(ctx context.Context, name string) (*ArtistDiscography, error) {
	// Step 1: Find artist
	artist, err := as.findMBArtist(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("find musicbrainz artist: %w", err)
	}
	if artist == nil {
		return nil, fmt.Errorf("artist '%s' not found on MusicBrainz", name)
	}

	info := ArtistInfo{
		Name:  artist.Name,
		Image: "", // MusicBrainz doesn't provide images directly
	}

	// Step 2: Get recordings with releases
	recordings, err := as.getMBRecordings(ctx, artist.ID)
	if err != nil {
		return nil, fmt.Errorf("get musicbrainz recordings: %w", err)
	}

	// Group by album — collect album IDs for parallel cover art fetch
	albumMap := make(map[string]*AlbumGroup)
	albumOrder := make([]string, 0)
	albumIDs := make([]string, 0)

	for _, rec := range recordings {
		albumName := ""
		albumID := ""
		if len(rec.Releases) > 0 {
			albumName = rec.Releases[0].Title
			albumID = rec.Releases[0].ID
		}
		if albumName == "" {
			albumName = "Singles & EPs"
		}

		duration := 0
		if rec.Length != nil {
			duration = *rec.Length / 1000
		}

		if _, exists := albumMap[albumName]; !exists {
			albumOrder = append(albumOrder, albumName)
			albumIDs = append(albumIDs, albumID)
			albumMap[albumName] = &AlbumGroup{
				Title:  albumName,
				Source: "musicbrainz",
				Tracks: make([]TrackResult, 0),
			}
		}

		albumMap[albumName].Tracks = append(albumMap[albumName].Tracks, TrackResult{
			Title:    rec.Title,
			Duration: duration,
			Source:   "musicbrainz",
			Album:    albumName,
			ExternalIDs: map[string]string{
				"mbid": rec.ID,
			},
		})
	}

	// Fetch cover art URLs in parallel (bounded at 10 goroutines)
	var wg sync.WaitGroup
	coverCtx, coverCancel := context.WithTimeout(ctx, 10*time.Second)
	defer coverCancel()
	sem := make(chan struct{}, 10)
	mu := sync.Mutex{}

	for i, albumID := range albumIDs {
		if albumID == "" {
			continue
		}
		i, albumID := i, albumID
		wg.Add(1)
		sem <- struct{}{} // acquire semaphore
		go func() {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore
			url := as.getCoverArtURL(coverCtx, albumID)
			if url != "" {
				mu.Lock()
				if i < len(albumOrder) {
					if group := albumMap[albumOrder[i]]; group != nil {
						group.Cover = url
					}
				}
				mu.Unlock()
			}
		}()
	}

	// Wait for all cover fetches to complete
	wg.Wait()
	close(sem)

	albums := make([]AlbumGroup, 0, len(albumOrder))
	for _, albumName := range albumOrder {
		if group := albumMap[albumName]; group != nil {
			albums = append(albums, *group)
		}
	}

	return &ArtistDiscography{
		ArtistInfo: info,
		Albums:     albums,
	}, nil
}

func (as *ArtistSearcher) findMBArtist(ctx context.Context, name string) (*mbArtist, error) {
	// Rate limit
	select {
	case <-as.mbTick:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	u, _ := url.Parse("https://musicbrainz.org/ws/2/artist")
	u.RawQuery = url.Values{
		"query": {fmt.Sprintf(`artist:"%s"`, name)},
		"fmt":   {"json"},
		"limit": {"5"},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "MuseMusic/1.0 (music@muse.app)")

	resp, err := as.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("musicbrainz: status %d", resp.StatusCode)
	}

	var sr mbArtistSearch
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("musicbrainz decode: %w", err)
	}

	if len(sr.Artists) == 0 {
		return nil, nil
	}

	// Find best match
	nameLower := strings.ToLower(name)
	for _, a := range sr.Artists {
		if strings.EqualFold(a.Name, name) || strings.Contains(strings.ToLower(a.Name), nameLower) {
			return &a, nil
		}
	}

	return &sr.Artists[0], nil
}

func (as *ArtistSearcher) getMBRecordings(ctx context.Context, artistID string) ([]mbRecording, error) {
	// Rate limit
	select {
	case <-as.mbTick:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Use the release-group browse endpoint to get all releases, then recordings
	// First get releases
	releaseURL := fmt.Sprintf(
		"https://musicbrainz.org/ws/2/release?artist=%s&inc=recordings&fmt=json&limit=100",
		artistID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releaseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "MuseMusic/1.0 (music@muse.app)")

	resp, err := as.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("musicbrainz releases: status %d", resp.StatusCode)
	}

	var releaseResp struct {
		Releases []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Media []struct {
				Tracks []struct {
					ID     string `json:"id"`
					Title  string `json:"title"`
					Length *int   `json:"length"`
					Number string `json:"number"`
				} `json:"tracks"`
			} `json:"media"`
		} `json:"releases"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&releaseResp); err != nil {
		return nil, fmt.Errorf("musicbrainz releases decode: %w", err)
	}

	// Flatten to recording-like entries with release info
	recordings := make([]mbRecording, 0)
	seen := make(map[string]bool) // dedup by recording ID

	for _, release := range releaseResp.Releases {
		for _, media := range release.Media {
			for _, track := range media.Tracks {
				if seen[track.ID] {
					continue
				}
				seen[track.ID] = true
				recordings = append(recordings, mbRecording{
					ID:     track.ID,
					Title:  track.Title,
					Length: track.Length,
					Releases: []struct {
						Title string `json:"title"`
						ID    string `json:"id"`
					}{
						{Title: release.Title, ID: release.ID},
					},
				})
			}
		}
	}

	return recordings, nil
}

func (as *ArtistSearcher) getCoverArtURL(ctx context.Context, releaseID string) string {
	u := fmt.Sprintf("https://coverartarchive.org/release/%s/front", releaseID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return ""
	}

	// Cover Art Archive redirects to the actual image
	client := &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // Don't follow redirect
	}}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusFound {
		if loc := resp.Header.Get("Location"); loc != "" {
			return loc
		}
	}

	// Direct 200 means the image was served
	if resp.StatusCode == http.StatusOK {
		return u
	}

	return ""
}

// ── iTunes / Apple Music fallback ──────────────────────────────────

func (as *ArtistSearcher) searchiTunes(ctx context.Context, name string) (*ArtistDiscography, error) {
	start := time.Now()

	// Step 1: Find artist on iTunes
	u, _ := url.Parse("https://itunes.apple.com/search")
	u.RawQuery = url.Values{
		"term":   {name},
		"entity": {"musicArtist"},
		"limit":  {"5"},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create itunes request: %w", err)
	}

	resp, err := as.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("itunes request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("itunes: status %d", resp.StatusCode)
	}

	var searchResp struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			ArtistID   int    `json:"artistId"`
			ArtistName string `json:"artistName"`
			ArtworkURL string `json:"artworkUrl100"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("itunes decode: %w", err)
	}

	if searchResp.ResultCount == 0 {
		return nil, fmt.Errorf("artist '%s' not found on iTunes", name)
	}

	// Find best match
	bestIdx := -1
	nameLower := strings.ToLower(name)
	for i, a := range searchResp.Results {
		if strings.EqualFold(a.ArtistName, name) || strings.Contains(strings.ToLower(a.ArtistName), nameLower) {
			bestIdx = i
			break
		}
	}
	if bestIdx == -1 {
		bestIdx = 0
	}
	bestArtist := searchResp.Results[bestIdx]

	// Upgrade artwork URL: replace 100x100 with 600x600
	artwork := strings.Replace(bestArtist.ArtworkURL, "100x100bb", "600x600bb", 1)

	info := ArtistInfo{
		Name:  bestArtist.ArtistName,
		Image: artwork,
	}

	as.logger.Info("found itunes artist",
		zap.String("name", bestArtist.ArtistName),
		zap.Int("id", bestArtist.ArtistID),
		zap.Duration("elapsed", time.Since(start)),
	)

	// Step 2: Get artist's albums via iTunes lookup
	lookupURL := fmt.Sprintf("https://itunes.apple.com/lookup?id=%d&entity=album&limit=50", bestArtist.ArtistID)
	lookupReq, err := http.NewRequestWithContext(ctx, http.MethodGet, lookupURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create itunes lookup request: %w", err)
	}

	lookupResp, err := as.client.Do(lookupReq)
	if err != nil {
		return nil, fmt.Errorf("itunes lookup request: %w", err)
	}
	defer lookupResp.Body.Close()

	if lookupResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("itunes lookup: status %d", lookupResp.StatusCode)
	}

	var albumResp struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			CollectionID   int    `json:"collectionId"`
			CollectionName string `json:"collectionName"`
			ArtworkURL     string `json:"artworkUrl100"`
			TrackCount     int    `json:"trackCount"`
			ReleaseDate    string `json:"releaseDate"`
		} `json:"results"`
	}

	if err := json.NewDecoder(lookupResp.Body).Decode(&albumResp); err != nil {
		return nil, fmt.Errorf("itunes lookup decode: %w", err)
	}

	// Filter to only albums (not singles or compilations) and get tracks
	albumGroups := make([]AlbumGroup, 0)
	for _, alb := range albumResp.Results {
		if alb.TrackCount == 0 {
			continue
		}

		albumCover := strings.Replace(alb.ArtworkURL, "100x100bb", "600x600bb", 1)

		// Get tracks for this album
		tracks, err := as.getiTunesAlbumTracks(ctx, alb.CollectionID)
		if err != nil {
			as.logger.Warn("failed to get itunes album tracks, skipping",
				zap.String("album", alb.CollectionName),
				zap.Error(err),
			)
			continue
		}

		albumGroups = append(albumGroups, AlbumGroup{
			Title:  alb.CollectionName,
			Cover:  albumCover,
			Source: "itunes",
			Tracks: tracks,
		})
	}

	as.logger.Info("itunes search completed",
		zap.Int("albums", len(albumGroups)),
		zap.Duration("elapsed", time.Since(start)),
	)

	return &ArtistDiscography{
		ArtistInfo: info,
		Albums:     albumGroups,
	}, nil
}

func (as *ArtistSearcher) getiTunesAlbumTracks(ctx context.Context, albumID int) ([]TrackResult, error) {
	u := fmt.Sprintf("https://itunes.apple.com/lookup?id=%d&entity=song&limit=200", albumID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := as.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("itunes tracks: status %d", resp.StatusCode)
	}

	var trackResp struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			TrackName   string `json:"trackName"`
			Duration    int    `json:"trackTimeMillis"`
			TrackID     int64  `json:"trackId"`
			PreviewURL  string `json:"previewUrl"`
			ReleaseDate string `json:"releaseDate"`
			DiscNumber  int    `json:"discNumber"`
			TrackNumber int    `json:"trackNumber"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&trackResp); err != nil {
		return nil, fmt.Errorf("itunes tracks decode: %w", err)
	}

	tracks := make([]TrackResult, 0, len(trackResp.Results))
	for _, t := range trackResp.Results {
		if t.TrackName == "" {
			continue
		}

		duration := 0
		if t.Duration > 0 {
			duration = t.Duration / 1000
		}

		tracks = append(tracks, TrackResult{
			Title:    t.TrackName,
			Duration: duration,
			Source:   "itunes",
			ExternalIDs: map[string]string{
				"itunes_id": fmt.Sprintf("%d", t.TrackID),
			},
		})
	}

	return tracks, nil
}
