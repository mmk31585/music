package importcmd

import (
	"github.com/gin-gonic/gin"

	appErr "music/internal/common/errors"
	"music/internal/common/response"
)

type Handler struct {
	importSvc *ImportService
}

func NewHandler(importSvc *ImportService) *Handler {
	return &Handler{importSvc: importSvc}
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

	userID := getUserID(c)
	if userID == "" {
		response.Error(c, appErr.Unauthorized("user not authenticated", nil))
		return
	}

	resp, err := h.importSvc.Import(c.Request.Context(), req.URL, req.Source, userID)
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
