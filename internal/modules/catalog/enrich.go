package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	album "music/internal/modules/catalog/album"
	artistpkg "music/internal/modules/catalog/artist"
	"music/internal/modules/catalog/track"
	"music/internal/modules/ingestion/enrichment"
	lyricsMod "music/internal/modules/lyrics"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EnrichHandler handles enrichment requests for published catalog items.
type EnrichHandler struct {
	enricher      *enrichment.Enricher
	trackSvc      *track.Service
	artistSvc     *artistpkg.Service
	albumSvc      *album.Service
	lyricsService *lyricsMod.Service
	logger        *zap.Logger
	enrichMutex   sync.Mutex
}

func NewEnrichHandler(
	enricher *enrichment.Enricher,
	trackSvc *track.Service,
	artistSvc *artistpkg.Service,
	albumSvc *album.Service,
	lyricsService *lyricsMod.Service,
	logger *zap.Logger,
) *EnrichHandler {
	return &EnrichHandler{
		enricher:      enricher,
		trackSvc:      trackSvc,
		artistSvc:     artistSvc,
		albumSvc:      albumSvc,
		lyricsService: lyricsService,
		logger:        logger,
	}
}

// EnrichTrack runs full enrichment for a single catalog track.
// Scope query param: "all" (default), "lyrics", "cover", "artist"
func (h *EnrichHandler) EnrichTrack(c *gin.Context) {
	trackID := c.Param("trackID")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "track id is required"})
		return
	}

	scope := c.DefaultQuery("scope", "all")

	tr, err := h.trackSvc.GetByID(c.Request.Context(), trackID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "track not found"})
		return
	}

	// Get primary artist name
	artistName := h.getArtistName(c.Request.Context(), tr)

	// Run enrichment
	result, err := h.enricher.Enrich(c.Request.Context(), tr.Title, artistName, "", tr.DurationSeconds)
	if err != nil {
		h.logger.Error("enrichment failed", zap.String("track_id", trackID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "enrichment failed: " + err.Error()})
		return
	}

	enriched := map[string]interface{}{
		"track_id":    trackID,
		"title":       tr.Title,
		"suggestions": result.Suggestions,
	}

	// Apply enrichment based on scope
	if scope == "all" || scope == "lyrics" {
		h.applyLyrics(c.Request.Context(), trackID, result)
	}

	// In future: apply cover, artist bio, etc. based on scope

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    enriched,
	})
}

// EnrichAllTracks starts background enrichment for ALL catalog tracks.
func (h *EnrichHandler) EnrichAllTracks(c *gin.Context) {
	go h.EnrichAllBackground(c.Request.Context())

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "enrichment started for all tracks in background",
	})
}

func (h *EnrichHandler) getArtistName(ctx context.Context, tr *track.Track) string {
	if tr.ArtistID != uuidNil {
		artist, err := h.artistSvc.GetByID(ctx, tr.ArtistID.String())
		if err == nil && artist != nil {
			return artist.Name
		}
	}
	// Try from Artists slice
	if len(tr.Artists) > 0 {
		for _, a := range tr.Artists {
			if a.Role == "primary" || a.Role == "" {
				artist, err := h.artistSvc.GetByID(ctx, a.ArtistID.String())
				if err == nil && artist != nil {
					return artist.Name
				}
			}
		}
		// Fallback to first artist
		artist, err := h.artistSvc.GetByID(ctx, tr.Artists[0].ArtistID.String())
		if err == nil && artist != nil {
			return artist.Name
		}
	}
	return ""
}

func (h *EnrichHandler) applyLyrics(ctx context.Context, trackID string, result *enrichment.EnrichmentResult) {
	if result == nil || result.LRCLib == nil {
		return
	}
	lrc := result.LRCLib
	content := lrc.SyncedLyrics
	lrcType := "lrc"
	if content == "" {
		content = lrc.PlainLyrics
		lrcType = "plain"
	}
	if content == "" {
		return
	}

	// Remove existing lyrics for "en", then create new
	existing, lookupErr := h.lyricsService.GetLyricsByTrackAndLanguage(ctx, trackID, "en")
	if lookupErr == nil {
		_ = h.lyricsService.DeleteLyrics(ctx, existing.ID)
	}

	_, err := h.lyricsService.CreateLyrics(ctx, lyricsMod.CreateLyricsRequest{
		TrackID:  trackID,
		Language: "en",
		Type:     lrcType,
		Content:  content,
	})
	if err != nil {
		h.logger.Warn("failed to save enriched lyrics",
			zap.String("track_id", trackID),
			zap.Error(err),
		)
	}
}

func (h *EnrichHandler) EnrichAllBackground(ctx context.Context) {
	if !h.enrichMutex.TryLock() {
		h.logger.Warn("enrich-all already in progress, skipping")
		return
	}
	defer h.enrichMutex.Unlock()

	tracks, err := h.trackSvc.List(ctx, 10000, 0, false)
	if err != nil {
		h.logger.Error("failed to list tracks for enrich-all", zap.Error(err))
		return
	}

	h.logger.Info("starting enrich-all", zap.Int("count", len(tracks)))

	for _, tr := range tracks {
		select {
		case <-ctx.Done():
			return
		default:
		}

		func(track track.Track) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			artistName := h.getArtistName(ctx, &track)
			result, err := h.enricher.Enrich(ctx, track.Title, artistName, "", track.DurationSeconds)
			if err != nil {
				h.logger.Warn("enrichment failed for track",
					zap.String("track_id", track.ID.String()),
					zap.Error(err),
				)
				return
			}

			h.applyLyrics(ctx, track.ID.String(), result)
			h.logger.Info("enriched track",
				zap.String("track_id", track.ID.String()),
				zap.String("title", track.Title),
				zap.Int("suggestions", len(result.Suggestions)),
			)
		}(tr)

		time.Sleep(1 * time.Second) // rate limit
	}

	h.logger.Info("enrich-all completed")
}

// EnrichArtist looks up artist info from external sources and updates the record.
func (h *EnrichHandler) EnrichArtist(c *gin.Context) {
	artistID := c.Param("artistID")
	if artistID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "artist id is required"})
		return
	}

	artist, err := h.artistSvc.GetByID(c.Request.Context(), artistID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "artist not found"})
		return
	}

	h.logger.Info("enriching artist",
		zap.String("artist_id", artistID),
		zap.String("name", artist.Name),
	)

	// We enrich artist by looking up their discography via ArtistSearcher (Deezer/MusicBrainz)
	// which returns artist info (image). For bio, we run a track enrichment as a proxy
	// since the Last.fm client returns artist bio in track enrichment results.
	//
	// Run both lookups in parallel:
	// 1. Deezer search for artist image
	// 2. Last.fm track search (via enricher) for bio

	type deezerResult struct {
		imageURL string
		err      error
	}
	type enrichmentResult struct {
		bio *string
		img *string
		err error
	}

	deezerCh := make(chan deezerResult, 1)
	enrichCh := make(chan enrichmentResult, 1)

	// 1. Deezer lookup for image
	go func() {
		imgURL, err := h.searchDeezerArtistImage(c.Request.Context(), artist.Name)
		deezerCh <- deezerResult{imageURL: imgURL, err: err}
	}()

	// 2. Enricher with Last.fm for bio
	go func() {
		// Use a dummy track query — Last.fm returns artist info in track responses
		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()
		result, err := h.enricher.Enrich(ctx, "", artist.Name, "", 0)
		if err != nil || result == nil || !result.Attempted {
			enrichCh <- enrichmentResult{err: err}
			return
		}

		var bio *string
		var img *string
		for _, s := range result.Suggestions {
			if s.Field == "artist_bio" {
				if v, ok := s.Value.(string); ok && v != "" {
					bio = &v
				}
			}
			if s.Field == "artist_image_url" {
				if v, ok := s.Value.(string); ok && v != "" {
					img = &v
				}
			}
		}
		enrichCh <- enrichmentResult{bio: bio, img: img}
	}()

	// Wait for both
	var deezerResp deezerResult
	var enrichResp enrichmentResult
	for i := 0; i < 2; i++ {
		select {
		case dr := <-deezerCh:
			deezerResp = dr
		case er := <-enrichCh:
			enrichResp = er
		case <-c.Request.Context().Done():
			c.JSON(http.StatusGatewayTimeout, gin.H{"success": false, "message": "enrichment timed out"})
			return
		}
	}

	// Merge results
	imageURL := deezerResp.imageURL
	if imageURL == "" && enrichResp.img != nil {
		imageURL = *enrichResp.img
	}
	if imageURL == "" && artist.ImageURL != nil {
		imageURL = *artist.ImageURL
	}

	// Build update request
	updateReq := artistpkg.UpdateRequest{
		ImageURL:         &imageURL,
		IsVerified:       &artist.IsVerified,
		MonthlyListeners: &artist.MonthlyListeners,
	}

	// Apply bio from enrichment
	if enrichResp.bio != nil {
		updateReq.Bio = enrichResp.bio
	}

	updated, err := h.artistSvc.Update(c.Request.Context(), artistID, updateReq)
	if err != nil {
		h.logger.Error("failed to update artist after enrichment",
			zap.String("artist_id", artistID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to save enriched data"})
		return
	}

	h.logger.Info("artist enriched successfully",
		zap.String("artist_id", artistID),
		zap.String("name", artist.Name),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// EnrichAlbum looks up album cover and info from external sources and updates the record.
func (h *EnrichHandler) EnrichAlbum(c *gin.Context) {
	albumID := c.Param("albumID")
	if albumID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "album id is required"})
		return
	}

	al, err := h.albumSvc.GetByID(c.Request.Context(), albumID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "album not found"})
		return
	}

	h.logger.Info("enriching album",
		zap.String("album_id", albumID),
		zap.String("title", al.Title),
	)

	// Get artist name for external lookups
	artistName := ""
	artists, listErr := h.albumSvc.ListArtists(c.Request.Context(), albumID)
	if listErr == nil && len(artists) > 0 {
		artistName = artists[0].Name
	} else {
		// Fallback: try to get the artist by the old artist_id field
		if al.ArtistID != uuidNil {
			art, artErr := h.artistSvc.GetByID(c.Request.Context(), al.ArtistID.String())
			if artErr == nil && art != nil {
				artistName = art.Name
			}
		}
	}

	// Look for cover art from multiple sources in parallel
	type coverResult struct {
		url string
		err error
	}

	lastfmCh := make(chan coverResult, 1)
	mbCh := make(chan coverResult, 1)

	// 1. Last.fm album search for cover
	go func() {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()
		// Run enrichment with track info — Last.fm may return album cover
		result, err := h.enricher.Enrich(ctx, al.Title, artistName, al.Title, 0)
		if err != nil || result == nil {
			lastfmCh <- coverResult{err: err}
			return
		}
		for _, s := range result.Suggestions {
			if s.Field == "album_cover_url" || s.Field == "album_cover_url_lastfm" {
				if v, ok := s.Value.(string); ok && v != "" {
					lastfmCh <- coverResult{url: v}
					return
				}
			}
		}
		lastfmCh <- coverResult{}
	}()

	// 2. MusicBrainz / Cover Art Archive (via ArtistSearcher-like approach)
	go func() {
		coverURL, err := h.searchAlbumCoverCAA(c.Request.Context(), al.Title, artistName)
		mbCh <- coverResult{url: coverURL, err: err}
	}()

	// Wait for both
	var lfmRes, mbRes coverResult
	for i := 0; i < 2; i++ {
		select {
		case r := <-lastfmCh:
			lfmRes = r
		case r := <-mbCh:
			mbRes = r
		case <-c.Request.Context().Done():
			c.JSON(http.StatusGatewayTimeout, gin.H{"success": false, "message": "enrichment timed out"})
			return
		}
	}

	// Pick best cover URL (CAA > Last.fm > existing)
	coverURL := mbRes.url
	if coverURL == "" {
		coverURL = lfmRes.url
	}
	if coverURL == "" && al.CoverURL != nil {
		coverURL = *al.CoverURL
	}

	if coverURL == "" || (al.CoverURL != nil && *al.CoverURL == coverURL) {
		h.logger.Info("no new cover data for album",
			zap.String("album_id", albumID),
		)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    al,
			"message": "no new cover found",
		})
		return
	}

	// Update album
	updateReq := album.UpdateRequest{
		CoverURL: &coverURL,
	}

	updated, err := h.albumSvc.Update(c.Request.Context(), albumID, updateReq)
	if err != nil {
		h.logger.Error("failed to update album after enrichment",
			zap.String("album_id", albumID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to save enriched data"})
		return
	}

	h.logger.Info("album enriched successfully",
		zap.String("album_id", albumID),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// searchDeezerArtistImage searches Deezer for an artist and returns their image URL.
func (h *EnrichHandler) searchDeezerArtistImage(ctx context.Context, name string) (string, error) {
	type deezerArtist struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Picture string `json:"picture_medium"`
	}
	type deezerSearch struct {
		Data []deezerArtist `json:"data"`
	}

	searchURL := fmt.Sprintf("https://api.deezer.com/search/artist?q=%s&limit=3", urlQueryEscape(name))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deezer: status %d", resp.StatusCode)
	}

	var sr deezerSearch
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", err
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

// searchAlbumCoverCAA searches Cover Art Archive for album cover via MusicBrainz.
func (h *EnrichHandler) searchAlbumCoverCAA(ctx context.Context, albumTitle, artistName string) (string, error) {
	if albumTitle == "" {
		return "", nil
	}

	// Search MusicBrainz for the release
	mbURL := fmt.Sprintf(
		"https://musicbrainz.org/ws/2/release?query=release:\"%s\"%%20AND%%20artist:\"%s\"&fmt=json&limit=5",
		urlQueryEscape(albumTitle),
		urlQueryEscape(artistName),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mbURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "MuseMusic/1.0 (music@muse.app)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
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
		return "", err
	}

	if len(mbResp.Releases) == 0 {
		return "", nil
	}

	// Try CAA for each release
	caaClient := &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	for _, rel := range mbResp.Releases {
		caaURL := fmt.Sprintf("https://coverartarchive.org/release/%s/front", rel.ID)
		caaReq, err := http.NewRequestWithContext(ctx, http.MethodGet, caaURL, nil)
		if err != nil {
			continue
		}
		caaResp, err := caaClient.Do(caaReq)
		if err != nil {
			continue
		}
		caaResp.Body.Close()

		if caaResp.StatusCode == http.StatusSeeOther || caaResp.StatusCode == http.StatusTemporaryRedirect || caaResp.StatusCode == http.StatusFound {
			if loc := caaResp.Header.Get("Location"); loc != "" {
				return loc, nil
			}
		}
		if caaResp.StatusCode == http.StatusOK {
			return caaURL, nil
		}
	}

	return "", nil
}

// urlQueryEscape escapes a string for use in a URL query parameter.
func urlQueryEscape(s string) string {
	return url.QueryEscape(s)
}

// uuidNil is a helper for zero UUID comparison
var uuidNil = [16]byte{}

var (
	ErrEnrichInProgress = errors.New("enrichment already in progress")
)
