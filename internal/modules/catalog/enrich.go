package catalog

import (
	"context"
	"errors"
	"net/http"
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
	enricher       *enrichment.Enricher
	trackSvc       *track.Service
	artistSvc      *artistpkg.Service
	albumSvc       *album.Service
	lyricsService  *lyricsMod.Service
	deezerClient   enrichment.DeezerClient
	coverArtClient enrichment.CoverArtClient
	lastFMClient   enrichment.LastFMClient
	spotifyClient  enrichment.SpotifyClient
	logger         *zap.Logger
	enrichMutex    sync.Mutex
}

func NewEnrichHandler(
	enricher *enrichment.Enricher,
	trackSvc *track.Service,
	artistSvc *artistpkg.Service,
	albumSvc *album.Service,
	lyricsService *lyricsMod.Service,
	deezerClient enrichment.DeezerClient,
	coverArtClient enrichment.CoverArtClient,
	lastFMClient enrichment.LastFMClient,
	spotifyClient enrichment.SpotifyClient,
	logger *zap.Logger,
) *EnrichHandler {
	return &EnrichHandler{
		enricher:       enricher,
		trackSvc:       trackSvc,
		artistSvc:      artistSvc,
		albumSvc:       albumSvc,
		lyricsService:  lyricsService,
		deezerClient:   deezerClient,
		coverArtClient: coverArtClient,
		lastFMClient:   lastFMClient,
		spotifyClient:  spotifyClient,
		logger:         logger,
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

	// Normalize name to strip invisible characters (soft hyphens, zero-width spaces)
	// that can break matching with external APIs
	artistName := enrichment.NormalizeName(artist.Name)

	// We enrich artist by looking up their info from multiple external sources.
	// Run lookups in parallel:
	// 1. Deezer search for artist image
	// 2. Last.fm artist.getInfo for bio, image, tags, similar artists
	// 3. Spotify search for artist image (high quality)

	type imgResult struct {
		url string
		err error
	}
	type lastFMResult struct {
		bio *string
		img *string
		err error
	}

	deezerCh := make(chan imgResult, 1)
	spotifyCh := make(chan imgResult, 1)
	lastfmCh := make(chan lastFMResult, 1)

	// 1. Deezer lookup for image (with 15s timeout)
	go func() {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()
		imgURL, err := h.deezerClient.SearchArtistImage(ctx, artistName)
		deezerCh <- imgResult{url: imgURL, err: err}
	}()

	// 2. Spotify lookup for high-res artist image
	go func() {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()
		imgURL, err := h.spotifyClient.SearchArtistImage(ctx, artistName)
		spotifyCh <- imgResult{url: imgURL, err: err}
	}()

	// 3. Last.fm artist.getInfo for bio and image
	go func() {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()
		result, err := h.lastFMClient.SearchArtist(ctx, artistName)
		if err != nil || result == nil {
			lastfmCh <- lastFMResult{err: err}
			return
		}

		var bio *string
		var img *string
		if result.ArtistBio != "" {
			bio = &result.ArtistBio
		}
		if result.ArtistImageURL != "" {
			img = &result.ArtistImageURL
		}
		lastfmCh <- lastFMResult{bio: bio, img: img}
	}()

	// Wait for all three
	var deezerResp, spotifyResp imgResult
	var lastfmResp lastFMResult
	for i := 0; i < 3; i++ {
		select {
		case dr := <-deezerCh:
			deezerResp = dr
		case sr := <-spotifyCh:
			spotifyResp = sr
		case lr := <-lastfmCh:
			lastfmResp = lr
		case <-c.Request.Context().Done():
			c.JSON(http.StatusGatewayTimeout, gin.H{"success": false, "message": "enrichment timed out"})
			return
		}
	}

	// Merge results: Spotify > Deezer > Last.fm > existing
	// Log what each source returned for debugging
	if spotifyResp.err != nil {
		h.logger.Warn("spotify enrichment failed",
			zap.String("artist", artistName),
			zap.Error(spotifyResp.err),
		)
	} else {
		h.logger.Info("spotify enrichment result",
			zap.String("artist", artistName),
			zap.String("image_url", spotifyResp.url),
		)
	}

	if deezerResp.err != nil {
		h.logger.Warn("deezer enrichment failed",
			zap.String("artist", artistName),
			zap.Error(deezerResp.err),
		)
	} else {
		h.logger.Info("deezer enrichment result",
			zap.String("artist", artistName),
			zap.String("image_url", deezerResp.url),
		)
	}

	if lastfmResp.err != nil {
		h.logger.Warn("lastfm enrichment failed",
			zap.String("artist", artistName),
			zap.Error(lastfmResp.err),
		)
	} else {
		h.logger.Info("lastfm enrichment result",
			zap.String("artist", artistName),
			zap.Any("image_url", lastfmResp.img),
			zap.Any("bio", lastfmResp.bio != nil),
		)
	}

	imageURL := spotifyResp.url
	if imageURL == "" {
		imageURL = deezerResp.url
	}
	if imageURL == "" && lastfmResp.img != nil {
		imageURL = *lastfmResp.img
	}

	h.logger.Info("final image selection",
		zap.String("artist", artistName),
		zap.String("selected_url", imageURL),
		zap.String("current_url", derefStr(artist.ImageURL)),
	)

	// Build update request — only set ImageURL when we have a new image
	// to avoid overwriting an existing image with empty string
	updateReq := artistpkg.UpdateRequest{
		IsVerified:       &artist.IsVerified,
		MonthlyListeners: &artist.MonthlyListeners,
	}
	if imageURL != "" {
		updateReq.ImageURL = &imageURL
	}

	// Apply bio from Last.fm
	if lastfmResp.bio != nil {
		updateReq.Bio = lastfmResp.bio
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
		zap.String("new_image_url", derefStr(updated.ImageURL)),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    artistpkg.ArtistToResponse(updated),
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

	// Get artist name for external lookups, normalize to strip invisible characters
	artistName := ""
	artists, listErr := h.albumSvc.ListArtists(c.Request.Context(), albumID)
	if listErr == nil && len(artists) > 0 {
		artistName = enrichment.NormalizeName(artists[0].Name)
	} else {
		if al.ArtistID != uuidNil {
			art, artErr := h.artistSvc.GetByID(c.Request.Context(), al.ArtistID.String())
			if artErr == nil && art != nil {
				artistName = enrichment.NormalizeName(art.Name)
			}
		}
	}
	albumTitle := enrichment.NormalizeName(al.Title)

	// Look for cover art from multiple sources in parallel
	type coverResult struct {
		url string
		err error
	}

	lastfmCh := make(chan coverResult, 1)
	mbCh := make(chan coverResult, 1)
	deezerCh := make(chan coverResult, 1)

	// 1. Deezer album search for cover (fast, free, no API key)
	go func() {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()
		coverURL, err := h.deezerClient.SearchAlbumCover(ctx, albumTitle, artistName)
		deezerCh <- coverResult{url: coverURL, err: err}
	}()

	// 2. Last.fm album search for cover (via track enrichment)
	go func() {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()
		result, err := h.enricher.Enrich(ctx, albumTitle, artistName, albumTitle, 0)
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

	// 3. MusicBrainz / Cover Art Archive (via CoverArtClient)
	go func() {
		coverURL, err := h.coverArtClient.SearchAlbumCover(c.Request.Context(), albumTitle, artistName)
		mbCh <- coverResult{url: coverURL, err: err}
	}()

	// Wait for all three
	var deezerRes, lfmRes, mbRes coverResult
	for i := 0; i < 3; i++ {
		select {
		case r := <-deezerCh:
			deezerRes = r
		case r := <-lastfmCh:
			lfmRes = r
		case r := <-mbCh:
			mbRes = r
		case <-c.Request.Context().Done():
			c.JSON(http.StatusGatewayTimeout, gin.H{"success": false, "message": "enrichment timed out"})
			return
		}
	}

	// Pick best cover URL (Deezer > CAA > Last.fm > existing)
	coverURL := deezerRes.url
	if coverURL == "" {
		coverURL = mbRes.url
	}
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
			"data":    album.AlbumToResponse(al),
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
		"data":    album.AlbumToResponse(updated),
	})
}

// uuidNil is a helper for zero UUID comparison
var uuidNil = [16]byte{}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

var (
	ErrEnrichInProgress = errors.New("enrichment already in progress")
)
