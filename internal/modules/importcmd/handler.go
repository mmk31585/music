package importcmd

import (
	"github.com/gin-gonic/gin"

	appErr "music/internal/common/errors"
	"music/internal/common/response"
)

type Handler struct {
	importSvc      *ImportService
	artistSearcher *ArtistSearcher
}

func NewHandler(importSvc *ImportService, artistSearcher *ArtistSearcher) *Handler {
	return &Handler{importSvc: importSvc, artistSearcher: artistSearcher}
}

func (h *Handler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Error(c, appErr.BadRequest("search query 'q' is required", nil))
		return
	}

	results, err := h.importSvc.Search(c.Request.Context(), q)
	if err != nil {
		response.Error(c, appErr.Internal("search failed: "+err.Error(), err))
		return
	}

	// Convert to DTO
	dto := make([]SearchResult, 0, len(results))
	for _, r := range results {
		extIDs := make(map[string]string)
		if r.ExternalIDs.SpotifyID != "" {
			extIDs["spotify_id"] = r.ExternalIDs.SpotifyID
		}
		if r.ExternalIDs.DeezerID != "" {
			extIDs["deezer_id"] = r.ExternalIDs.DeezerID
		}
		if r.ExternalIDs.MBID != "" {
			extIDs["mbid"] = r.ExternalIDs.MBID
		}

		dto = append(dto, SearchResult{
			Title:       r.Title,
			Artist:      r.Artist,
			URL:         r.URL,
			Duration:    r.Duration,
			Thumbnail:   r.Thumbnail,
			Source:      r.Source,
			Score:       r.Score,
			ISRC:        r.ISRC,
			ExternalIDs: extIDs,
		})
	}

	response.OK(c, "search results retrieved", dto)
}

func (h *Handler) Import(c *gin.Context) {
	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", err))
		return
	}

	// Validate: need either a URL or at least title+artist
	if req.URL == "" && (req.Title == "" || req.Artist == "") {
		response.Error(c, appErr.BadRequest("either 'url' or 'title'+'artist' is required", nil))
		return
	}

	userID := getUserID(c)
	if userID == "" {
		response.Error(c, appErr.Unauthorized("user not authenticated", nil))
		return
	}

	resp, err := h.importSvc.Import(c.Request.Context(), req.URL, req.Source, req.Title, req.Artist, req.Album, req.ExternalIDs, userID)
	if err != nil {
		response.Error(c, appErr.Internal("import failed: "+err.Error(), err))
		return
	}

	response.Created(c, "track import initiated", resp)
}

func (h *Handler) GetProgress(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		response.Error(c, appErr.BadRequest("jobId is required", nil))
		return
	}

	pu, err := h.importSvc.GetProgress(c.Request.Context(), jobID)
	if err != nil {
		response.Error(c, appErr.Internal("failed to get progress: "+err.Error(), err))
		return
	}
	if pu == nil {
		response.Error(c, appErr.NotFound("job not found", nil))
		return
	}

	resp := ProgressResponse{
		JobID:    pu.JobID,
		Status:   string(pu.Status),
		Progress: pu.Progress,
		Stage:    pu.Stage,
		Error:    pu.Error,
		DraftID:  pu.DraftID,
	}

	response.OK(c, "import progress", resp)
}

func (h *Handler) SearchArtist(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		response.Error(c, appErr.BadRequest("artist name 'name' is required", nil))
		return
	}

	result, err := h.artistSearcher.SearchArtist(c.Request.Context(), name)
	if err != nil {
		response.Error(c, appErr.Internal("artist search failed: "+err.Error(), err))
		return
	}

	// Convert to DTO
	dto := ArtistDiscographyResult{
		ArtistInfo: ArtistInfoDTO{
			Name:  result.ArtistInfo.Name,
			Image: result.ArtistInfo.Image,
		},
		Albums: make([]AlbumGroupDTO, 0, len(result.Albums)),
	}

	for _, album := range result.Albums {
		tracks := make([]TrackResultDTO, 0, len(album.Tracks))
		for _, t := range album.Tracks {
			tracks = append(tracks, TrackResultDTO{
				Title:       t.Title,
				Duration:    t.Duration,
				Source:      t.Source,
				Album:       t.Album,
				ExternalIDs: t.ExternalIDs,
			})
		}

		dto.Albums = append(dto.Albums, AlbumGroupDTO{
			Title:  album.Title,
			Cover:  album.Cover,
			Source: album.Source,
			Tracks: tracks,
		})
	}

	response.OK(c, "artist discography retrieved", dto)
}

func (h *Handler) BatchImport(c *gin.Context) {
	var req BatchImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", err))
		return
	}

	userID := getUserID(c)
	if userID == "" {
		response.Error(c, appErr.Unauthorized("user not authenticated", nil))
		return
	}

	resp, err := h.importSvc.BatchImport(c.Request.Context(), req.Tracks, userID)
	if err != nil {
		response.Error(c, appErr.Internal("batch import failed: "+err.Error(), err))
		return
	}

	response.Created(c, "batch import initiated", resp)
}

func (h *Handler) GetBatchProgress(c *gin.Context) {
	batchID := c.Param("batchId")
	if batchID == "" {
		response.Error(c, appErr.BadRequest("batchId is required", nil))
		return
	}

	bp, err := h.importSvc.GetBatchProgress(c.Request.Context(), batchID)
	if err != nil {
		response.Error(c, appErr.Internal("failed to get batch progress: "+err.Error(), err))
		return
	}
	if bp == nil {
		response.Error(c, appErr.NotFound("batch not found", nil))
		return
	}

	response.OK(c, "batch progress", bp)
}

func getUserID(c *gin.Context) string {
	for _, key := range []string{"auth_user_id", "user_id", "userID", "userId", "sub"} {
		raw, exists := c.Get(key)
		if !exists || raw == nil {
			continue
		}
		switch v := raw.(type) {
		case string:
			if v != "" {
				return v
			}
		}
	}
	return ""
}
