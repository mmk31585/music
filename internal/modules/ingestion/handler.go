package ingestion

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	appErr "music/internal/common/errors"
	"music/internal/common/pagination"
	"music/internal/common/response"
	"music/internal/modules/auth"
	"music/internal/modules/ingestion/finalization"
)

type Handler struct {
	service   *Service
	finalizer *finalization.Service
	cleanup   *CleanupService
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SetCleanup(c *CleanupService) {
	h.cleanup = c
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
		case errors.Is(err, ErrDuplicateFile):
			response.Error(c, appErr.Conflict("this file already exists in the catalog", err))
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

	if res.Items == nil {
		res.Items = []DraftListItem{}
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
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Reason = c.DefaultQuery("reason", "")
	}

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

func (h *Handler) SetFinalizer(f *finalization.Service) {
	h.finalizer = f
}

func (h *Handler) FinalizeDraft(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.BadRequest("draft ID is required", nil))
		return
	}

	if h.finalizer == nil {
		response.Error(c, appErr.Internal("finalization not configured", nil))
		return
	}

	draft, err := h.service.GetDraftRaw(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.NotFound("draft not found", nil))
			return
		}
		response.Error(c, appErr.Internal("failed to get draft", err))
		return
	}

	var finalMeta string
	if draft.FinalMetadata != nil {
		finalMeta = *draft.FinalMetadata
	}

	result, err := h.finalizer.Finalize(c.Request.Context(), draft.ID, string(draft.Status), draft.FilePath, draft.Format, finalMeta)
	if err != nil {
		response.Error(c, appErr.Internal("failed to finalize draft: "+err.Error(), err))
		return
	}

	response.OK(c, "draft published successfully", result)
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

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetIngestionStats(c.Request.Context())
	if err != nil {
		response.Error(c, appErr.Internal("failed to get ingestion stats", err))
		return
	}
	response.OK(c, "ingestion stats retrieved", stats)
}

func (h *Handler) TriggerCleanup(c *gin.Context) {
	if h.cleanup == nil {
		response.Error(c, appErr.Internal("cleanup service not configured", nil))
		return
	}
	flagged, err := h.cleanup.FlagStaleDrafts(c.Request.Context(), 24*time.Hour)
	if err != nil {
		response.Error(c, appErr.Internal("cleanup failed", err))
		return
	}
	response.OK(c, "cleanup completed", gin.H{"flagged": flagged})
}

func (h *Handler) UploadDraftImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.BadRequest("draft ID is required", nil))
		return
	}

	entityType := c.Param("entity")
	if entityType != "artist" && entityType != "album" && entityType != "track" && entityType != "cover" {
		response.Error(c, appErr.BadRequest("entity must be 'artist', 'album', 'track', or 'cover'", nil))
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		response.Error(c, appErr.BadRequest("no image file provided in 'image' field", err))
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		response.Error(c, appErr.BadRequest("unsupported image format (allowed: jpg, png, webp)", nil))
		return
	}

	imageURL, err := h.service.UploadDraftImage(c.Request.Context(), id, entityType, file, ext)
	if err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.NotFound("draft not found", nil))
			return
		}
		response.Error(c, appErr.Internal("failed to upload image", err))
		return
	}

	response.OK(c, "image uploaded", gin.H{"url": imageURL})
}

func (h *Handler) GetConfig(c *gin.Context) {
	response.OK(c, "ingestion config", IngestionConfigResponse{
		MaxUploadSize:     200 << 20,
		EnrichmentEnabled: true,
	})
}
