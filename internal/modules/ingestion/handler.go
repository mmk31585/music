package ingestion

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	appErr "music/internal/common/errors"
	"music/internal/common/pagination"
	"music/internal/common/response"
	"music/internal/modules/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upload(c *gin.Context) {
	const hardLimit = 200 << 20

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, hardLimit)

	if err := c.Request.ParseMultipartForm(hardLimit); err != nil {
		response.Error(c, appErr.BadRequest("file too large or invalid multipart form", err))
		return
	}

	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		response.Error(c, appErr.BadRequest("no audio file provided in 'audio' field", err))
		return
	}
	defer file.Close()

	userID := auth.UserIDFromContext(c)
	if userID == "" {
		response.Error(c, appErr.Unauthorized("user not authenticated", nil))
		return
	}

	res, err := h.service.Upload(c.Request.Context(), file, header, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidFileType):
			response.Error(c, appErr.BadRequest("unsupported audio format (allowed: mp3, flac, ogg, m4a, aac)", err))
		case errors.Is(err, ErrFileTooLarge):
			response.Error(c, appErr.BadRequest("file exceeds maximum size of 200MB", err))
		case errors.Is(err, ErrNoFileProvided):
			response.Error(c, appErr.BadRequest("no file provided", err))
		case errors.Is(err, ErrStorageFailed):
			response.Error(c, appErr.Internal("failed to store uploaded file", err))
		default:
			response.Error(c, appErr.Internal("upload failed", err))
		}
		return
	}

	response.Created(c, "file uploaded and enrichment started", res)
}

func (h *Handler) GetDraft(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.BadRequest("draft ID is required", nil))
		return
	}

	draft, err := h.service.GetDraftByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.NotFound("draft not found", nil))
			return
		}
		response.Error(c, appErr.Internal("failed to get draft", err))
		return
	}

	response.OK(c, "draft retrieved successfully", draft)
}

func (h *Handler) ListDrafts(c *gin.Context) {
	status := c.Query("status")

	p := pagination.FromRequest(c.Request)

	res, err := h.service.ListDrafts(c.Request.Context(), status, p.Page, p.Limit)
	if err != nil {
		if errors.Is(err, ErrInvalidStatus) {
			response.Error(c, appErr.BadRequest("invalid status filter", err))
			return
		}
		response.Error(c, appErr.Internal("failed to list drafts", err))
		return
	}

	meta := pagination.NewMeta(p, res.Total)

	response.SuccessWithMeta(c, http.StatusOK, "drafts retrieved successfully", res.Items, meta)
}

func (h *Handler) EnrichDraft(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.BadRequest("draft ID is required", nil))
		return
	}

	if err := h.service.EnrichDraft(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.NotFound("draft not found", nil))
			return
		}
		response.Error(c, appErr.Internal("failed to start enrichment", err))
		return
	}

	response.OK(c, "enrichment started", gin.H{"draftId": id, "status": "enriching"})
}

func (h *Handler) GetDraftSuggestions(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.BadRequest("draft ID is required", nil))
		return
	}

	suggestions, err := h.service.GetDraftSuggestions(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.NotFound("draft not found", nil))
			return
		}
		response.Error(c, appErr.Internal("failed to get suggestions", err))
		return
	}

	response.OK(c, "suggestions retrieved successfully", suggestions)
}

func (h *Handler) SaveFinalMetadata(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.BadRequest("draft ID is required", nil))
		return
	}

	var req SaveFinalMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", err))
		return
	}

	if err := h.service.SaveFinalMetadata(c.Request.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, ErrDraftNotFound):
			response.Error(c, appErr.NotFound("draft not found", nil))
		case errors.Is(err, ErrInvalidStatus):
			response.Error(c, appErr.BadRequest("draft must be in review status", nil))
		default:
			response.Error(c, appErr.Internal("failed to save final metadata", err))
		}
		return
	}

	response.OK(c, "final metadata saved and draft accepted", gin.H{})
}

func (h *Handler) RejectDraft(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.BadRequest("draft ID is required", nil))
		return
	}

	var req RejectDraftRequest
	req.Reason = c.DefaultQuery("reason", "")
	_ = c.ShouldBindJSON(&req)

	if err := h.service.RejectDraft(c.Request.Context(), id, req.Reason); err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.NotFound("draft not found", nil))
			return
		}
		response.Error(c, appErr.Internal("failed to reject draft", err))
		return
	}

	response.OK(c, "draft rejected", gin.H{})
}

func (h *Handler) SearchArtists(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Error(c, appErr.BadRequest("search query 'q' is required", nil))
		return
	}

	results, err := h.service.SearchArtists(c.Request.Context(), q)
	if err != nil {
		response.Error(c, appErr.Internal("failed to search artists", err))
		return
	}

	response.OK(c, "artists retrieved", results)
}

func (h *Handler) SearchAlbums(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Error(c, appErr.BadRequest("search query 'q' is required", nil))
		return
	}

	results, err := h.service.SearchAlbums(c.Request.Context(), q)
	if err != nil {
		response.Error(c, appErr.Internal("failed to search albums", err))
		return
	}

	response.OK(c, "albums retrieved", results)
}
