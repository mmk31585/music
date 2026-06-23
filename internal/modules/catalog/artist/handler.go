package artist

import (
	"errors"
	"music/internal/modules/catalog/common"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List godoc
// @Summary List artists
// @Description Returns a paginated list of artists.
// @Tags artists
// @Produce json
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Param q query string false "Search query to filter artists by name"
// @Success 200 {array} ArtistResponse
// @Failure 500 {object} map[string]interface{}
// @Router /artists [get]
func (h *Handler) List(c *gin.Context) {
	p := common.ParsePagination(c)
	q := strings.TrimSpace(c.Query("q"))

	var items []Artist
	var err error

	if q != "" {
		items, err = h.service.Search(c.Request.Context(), q, p.Limit, p.Offset)
	} else {
		items, err = h.service.List(c.Request.Context(), p.Limit, p.Offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list artists"})
		return
	}

	c.JSON(http.StatusOK, ArtistListToResponse(items))
}

// Get godoc
// @Summary Get artist by ID
// @Description Returns a single artist by its ID.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Success 200 {object} ArtistResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /artists/{artistID} [get]
func (h *Handler) Get(c *gin.Context) {
	item, err := h.service.GetByID(c.Request.Context(), c.Param("artistID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artist"})
		return
	}
	c.JSON(http.StatusOK, ArtistToResponse(item))
}

// Create godoc
// @Summary Create artist
// @Description Creates a new artist.
// @Tags artists
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Artist creation payload"
// @Success 201 {object} ArtistResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /artists [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Create(c.Request.Context(), req)
	if errors.Is(err, common.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "artist already exists"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create artist"})
		return
	}
	c.JSON(http.StatusCreated, ArtistToResponse(item))
}

// Update godoc
// @Summary Update artist
// @Description Updates an existing artist by ID.
// @Tags artists
// @Accept json
// @Produce json
// @Param artistID path string true "Artist ID"
// @Param request body UpdateRequest true "Artist update payload"
// @Success 200 {object} ArtistResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /artists/{artistID} [put]
func (h *Handler) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Update(c.Request.Context(), c.Param("artistID"), req)
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update artist"})
		return
	}
	c.JSON(http.StatusOK, ArtistToResponse(item))
}

// Delete godoc
// @Summary Delete artist
// @Description Deletes an artist by ID.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /artists/{artistID} [delete]
func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(c.Request.Context(), c.Param("artistID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete artist"})
		return
	}
	c.Status(http.StatusNoContent)
}

// Overview godoc
// @Summary Get artist overview
// @Description Returns a comprehensive overview of an artist including top tracks, albums, singles, and appearances.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Success 200 {object} ArtistOverviewResponse
// @Router /artists/{artistID}/overview [get]
func (h *Handler) Overview(c *gin.Context) {
	overview, err := h.service.GetOverview(c.Request.Context(), c.Param("artistID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artist overview"})
		return
	}

	resp := &ArtistOverviewResponse{
		Artist:    ArtistToResponse(overview.Artist),
		TopTracks: make([]ArtistTrackResponse, len(overview.TopTracks)),
		Albums:    make([]ArtistAlbumResponse, len(overview.Albums)),
		Singles:   make([]ArtistAlbumResponse, len(overview.Singles)),
		AppearsOn: make([]ArtistTrackResponse, len(overview.AppearsOn)),
	}
	for i, t := range overview.TopTracks {
		resp.TopTracks[i] = artistTrackToResponse(&t)
	}
	for i, a := range overview.Albums {
		resp.Albums[i] = artistAlbumToResponse(&a)
	}
	for i, a := range overview.Singles {
		resp.Singles[i] = artistAlbumToResponse(&a)
	}
	for i, t := range overview.AppearsOn {
		resp.AppearsOn[i] = artistTrackToResponse(&t)
	}
	c.JSON(http.StatusOK, resp)
}

func artistTrackToResponse(t *ArtistTrack) ArtistTrackResponse {
	return ArtistTrackResponse{
		ID:              t.ID,
		Title:           t.Title,
		Slug:            t.Slug,
		AlbumID:         t.AlbumID,
		DurationSeconds: t.DurationSeconds,
		TrackNumber:     t.TrackNumber,
		Explicit:        t.Explicit,
		AudioURL:        t.AudioURL,
		CoverURL:        t.CoverURL,
		PlayCount:       t.PlayCount,
		ArtistRole:      t.ArtistRole,
		ArtistPosition:  t.ArtistPosition,
		CreatedAt:       t.CreatedAt,
	}
}

func artistAlbumToResponse(a *ArtistAlbum) ArtistAlbumResponse {
	return ArtistAlbumResponse{
		ID:          a.ID,
		Title:       a.Title,
		Slug:        a.Slug,
		CoverURL:    a.CoverURL,
		ReleaseDate: a.ReleaseDate,
		AlbumType:   a.AlbumType,
		ArtistRole:  a.ArtistRole,
		CreatedAt:   a.CreatedAt,
	}
}

// Tracks godoc
// @Summary Get artist tracks
// @Description Returns a paginated list of tracks by the artist.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} ArtistTrackResponse
// @Router /artists/{artistID}/tracks [get]
func (h *Handler) Tracks(c *gin.Context) {
	p := common.ParsePagination(c)
	items, err := h.service.ListTracks(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tracks"})
		return
	}
	if items == nil {
		c.JSON(http.StatusOK, []ArtistTrackResponse{})
		return
	}
	resp := make([]ArtistTrackResponse, len(items))
	for i, t := range items {
		resp[i] = artistTrackToResponse(&t)
	}
	c.JSON(http.StatusOK, resp)
}

// Albums godoc
// @Summary Get artist albums
// @Description Returns a paginated list of albums by the artist.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} ArtistAlbumResponse
// @Router /artists/{artistID}/albums [get]
func (h *Handler) Albums(c *gin.Context) {
	p := common.ParsePagination(c)
	items, err := h.service.ListAlbums(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list albums"})
		return
	}
	if items == nil {
		c.JSON(http.StatusOK, []ArtistAlbumResponse{})
		return
	}
	resp := make([]ArtistAlbumResponse, len(items))
	for i, a := range items {
		resp[i] = artistAlbumToResponse(&a)
	}
	c.JSON(http.StatusOK, resp)
}

// Singles godoc
// @Summary Get artist singles
// @Description Returns a paginated list of singles by the artist.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} ArtistAlbumResponse
// @Router /artists/{artistID}/singles [get]
func (h *Handler) Singles(c *gin.Context) {
	p := common.ParsePagination(c)
	items, err := h.service.ListSingles(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list singles"})
		return
	}
	if items == nil {
		c.JSON(http.StatusOK, []ArtistAlbumResponse{})
		return
	}
	resp := make([]ArtistAlbumResponse, len(items))
	for i, a := range items {
		resp[i] = artistAlbumToResponse(&a)
	}
	c.JSON(http.StatusOK, resp)
}

// AppearsOn godoc
// @Summary Get artist appearances
// @Description Returns a paginated list of tracks the artist appears on but did not create.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Success 200 {array} ArtistTrackResponse
// @Router /artists/{artistID}/appears-on [get]
func (h *Handler) AppearsOn(c *gin.Context) {
	p := common.ParsePagination(c)
	items, err := h.service.ListAppearsOn(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list appearances"})
		return
	}
	if items == nil {
		c.JSON(http.StatusOK, []ArtistTrackResponse{})
		return
	}
	resp := make([]ArtistTrackResponse, len(items))
	for i, t := range items {
		resp[i] = artistTrackToResponse(&t)
	}
	c.JSON(http.StatusOK, resp)
}

// TopTracks godoc
// @Summary Get artist top tracks
// @Description Returns the top tracks for an artist.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Success 200 {array} ArtistTrackResponse
// @Router /artists/{artistID}/top-tracks [get]
func (h *Handler) TopTracks(c *gin.Context) {
	p := common.ParsePagination(c)
	items, err := h.service.ListTopTracks(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list top tracks"})
		return
	}
	if items == nil {
		c.JSON(http.StatusOK, []ArtistTrackResponse{})
		return
	}
	resp := make([]ArtistTrackResponse, len(items))
	for i, t := range items {
		resp[i] = artistTrackToResponse(&t)
	}
	c.JSON(http.StatusOK, resp)
}

// Related godoc
// @Summary Get related artists
// @Description Returns a list of artists related to the given artist.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Success 200 {array} RelatedArtistResponse
// @Router /artists/{artistID}/related [get]
func (h *Handler) Related(c *gin.Context) {
	p := common.ParsePagination(c)
	items, err := h.service.ListRelated(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list related artists"})
		return
	}
	if items == nil {
		c.JSON(http.StatusOK, []RelatedArtistResponse{})
		return
	}
	resp := make([]RelatedArtistResponse, len(items))
	for i, ra := range items {
		resp[i] = RelatedArtistResponse{
			ID:               ra.ID,
			Name:             ra.Name,
			Slug:             ra.Slug,
			Bio:              ra.Bio,
			ImageURL:         ra.ImageURL,
			AvatarMediaID:    ra.AvatarMediaID,
			BannerMediaID:    ra.BannerMediaID,
			IsVerified:       ra.IsVerified,
			MonthlyListeners: ra.MonthlyListeners,
			RelatedScore:     ra.RelatedScore,
			RelatedSource:    ra.RelatedSource,
			RelatedAt:        ra.RelatedAt,
		}
	}
	c.JSON(http.StatusOK, resp)
}

// ReplaceRelated godoc
// @Summary Replace related artists
// @Description Replaces the list of related artists for the given artist.
// @Tags artists
// @Accept json
// @Produce json
// @Param artistID path string true "Artist ID"
// @Param request body []RelatedArtistRequest true "Related artists list"
// @Success 200 {array} RelatedArtistResponse
// @Router /artists/{artistID}/related [put]
func (h *Handler) ReplaceRelated(c *gin.Context) {
	var req []RelatedArtistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	items, err := h.service.ReplaceRelated(c.Request.Context(), c.Param("artistID"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to replace related artists"})
		return
	}
	resp := make([]RelatedArtistResponse, len(items))
	for i, ra := range items {
		resp[i] = RelatedArtistResponse{
			ID:               ra.ID,
			Name:             ra.Name,
			Slug:             ra.Slug,
			Bio:              ra.Bio,
			ImageURL:         ra.ImageURL,
			AvatarMediaID:    ra.AvatarMediaID,
			BannerMediaID:    ra.BannerMediaID,
			IsVerified:       ra.IsVerified,
			MonthlyListeners: ra.MonthlyListeners,
			RelatedScore:     ra.RelatedScore,
			RelatedSource:    ra.RelatedSource,
			RelatedAt:        ra.RelatedAt,
		}
	}
	c.JSON(http.StatusOK, resp)
}

// ReplaceTopTracks godoc
// @Summary Replace top tracks
// @Description Replaces the top tracks list for the given artist.
// @Tags artists
// @Accept json
// @Produce json
// @Param artistID path string true "Artist ID"
// @Param request body []ArtistTopTrackRequest true "Top tracks list"
// @Success 200 {object} map[string]interface{}
// @Router /artists/{artistID}/top-tracks [put]
func (h *Handler) ReplaceTopTracks(c *gin.Context) {
	var req []ArtistTopTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	items, err := h.service.ReplaceTopTracks(c.Request.Context(), c.Param("artistID"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to replace top tracks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}
