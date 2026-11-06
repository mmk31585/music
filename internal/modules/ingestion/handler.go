package ingestion

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appErr "music/internal/common/errors"
	"music/internal/common/pagination"
	"music/internal/common/response"
	"music/internal/modules/auth"
	"music/internal/modules/ingestion/enrichment"
	"music/internal/modules/ingestion/finalization"
	"music/internal/modules/lyrics"
	platformstorage "music/internal/platform/storage"
)

// uploadStep labels each phase of the upload flow for structured logs.
const (
	uploadStepParse      = "parse_multipart_form"
	uploadStepFileRead   = "read_uploaded_file"
	uploadStepMetadata   = "extract_metadata"
	uploadStepS3Put      = "s3_put_object"
	uploadStepS3GetURL   = "s3_get_url"
	uploadStepDBCreate   = "db_create_draft"
	uploadStepPresign    = "s3_presign_upload_url"
)

type Handler struct {
	service       *Service
	finalizer     *finalization.Service
	cleanup       *CleanupService
	mlClient      *lyrics.MLServiceClient
	storage       platformstorage.Storage
	storageDriver string
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// SetMLClient configures the optional ML service client for AI lyrics
// generation. Set to nil to disable.
func (h *Handler) SetMLClient(m *lyrics.MLServiceClient) {
	h.mlClient = m
}

// SetStorage configures the storage backend for determining audio URLs.
func (h *Handler) SetStorage(s platformstorage.Storage, driver string) {
	h.storage = s
	h.storageDriver = driver
}

func (h *Handler) SetCleanup(c *CleanupService) {
	h.cleanup = c
}

func (h *Handler) Upload(c *gin.Context) {
	const hardLimit = 200 << 20

	log := zap.L().With(zap.String("step", uploadStepParse))

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, hardLimit)

	if err := c.Request.ParseMultipartForm(hardLimit); err != nil {
		log.Error("parse multipart form failed", zap.Error(err))
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "file too large or invalid multipart form", err))
		return
	}

	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		log.Error("read audio file field failed", zap.Error(err))
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "no audio file provided in 'audio' field", err))
		return
	}
	defer file.Close()

	userID := auth.UserIDFromContext(c)
	if userID == "" {
		log.Error("user not authenticated", zap.String("step", uploadStepFileRead))
		response.Error(c, appErr.New(http.StatusUnauthorized, appErr.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	zap.L().Info("parsed multipart upload",
		zap.String("step", uploadStepParse),
		zap.String("filename", header.Filename),
		zap.Int64("size", header.Size),
		zap.String("user_id", userID),
	)

	res, err := h.service.Upload(c.Request.Context(), file, header, userID)
	if err != nil {
		zap.L().Error("service upload failed", zap.String("step", uploadStepS3Put), zap.Error(err))
		switch {
		case errors.Is(err, ErrInvalidFileType):
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "unsupported audio format (allowed: mp3, flac, ogg, m4a, aac)", err))
		case errors.Is(err, ErrFileTooLarge):
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "file exceeds maximum size of 200MB", err))
		case errors.Is(err, ErrNoFileProvided):
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "no file provided", err))
		case errors.Is(err, ErrStorageFailed):
			response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to store uploaded file", err))
		case errors.Is(err, ErrDuplicateFile):
			response.Error(c, appErr.New(http.StatusConflict, appErr.CodeConflict, "this file already exists in the catalog", err))
		default:
			response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "upload failed", err))
		}
		return
	}

	zap.L().Info("upload complete",
		zap.String("step", uploadStepDBCreate),
		zap.String("draft_id", res.DraftID),
		zap.String("user_id", userID),
	)

	response.Success(c, http.StatusCreated, "file uploaded and enrichment started", res)
}

// UploadURL returns a presigned PUT URL for direct client-to-S3 upload.
// Used for files larger than ~4MB to avoid buffering the full body in serverless.
func (h *Handler) UploadURL(c *gin.Context) {
	log := zap.L().With(zap.String("step", uploadStepPresign))

	var req struct {
		Filename    string `json:"filename" validate:"required"`
		ContentType string `json:"contentType" validate:"required"`
		Size        int64  `json:"size" validate:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("invalid request body", zap.Error(err))
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "invalid request body", err))
		return
	}

	userID := auth.UserIDFromContext(c)
	if userID == "" {
		log.Error("user not authenticated")
		response.Error(c, appErr.New(http.StatusUnauthorized, appErr.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	resp, err := h.service.PrepareUpload(c.Request.Context(), req.Filename, req.ContentType, req.Size, userID)
	if err != nil {
		log.Error("prepare upload failed", zap.Error(err))
		switch {
		case errors.Is(err, ErrInvalidFileType):
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "unsupported audio format (allowed: mp3, flac, ogg, m4a, aac)", err))
		case errors.Is(err, ErrFileTooLarge):
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "file exceeds maximum size of 200MB", err))
		case errors.Is(err, ErrNoFileProvided):
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "no file provided", err))
		default:
			response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "upload url failed", err))
		}
		return
	}

	log.Info("presigned upload url issued",
		zap.String("draft_id", resp.DraftID),
		zap.String("key", resp.Key),
		zap.String("user_id", userID),
	)

	response.Success(c, http.StatusOK, "upload url ready", resp)
}

// UploadComplete is called by the client after a direct S3 PUT. It verifies
// the object, extracts metadata, creates the draft, and returns the response.
func (h *Handler) UploadComplete(c *gin.Context) {
	log := zap.L().With(zap.String("step", uploadStepDBCreate))

	var req UploadCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("invalid request body", zap.Error(err))
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "invalid request body", err))
		return
	}

	userID := auth.UserIDFromContext(c)
	if userID == "" {
		log.Error("user not authenticated")
		response.Error(c, appErr.New(http.StatusUnauthorized, appErr.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	resp, err := h.service.CompleteUpload(c.Request.Context(), req.DraftID, req.Key, userID)
	if err != nil {
		log.Error("complete upload failed", zap.String("draft_id", req.DraftID), zap.Error(err))
		switch {
		case errors.Is(err, ErrInvalidFileType):
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "unsupported audio format (allowed: mp3, flac, ogg, m4a, aac)", err))
		case errors.Is(err, ErrNoFileProvided):
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "no file provided", err))
		case errors.Is(err, ErrDuplicateFile):
			response.Error(c, appErr.New(http.StatusConflict, appErr.CodeConflict, "this file already exists in the catalog", err))
		case errors.Is(err, ErrStorageFailed):
			response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to verify stored file", err))
		default:
			response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "upload complete failed", err))
		}
		return
	}

	log.Info("direct S3 upload complete",
		zap.String("draft_id", resp.DraftID),
		zap.String("user_id", userID),
	)

	response.Success(c, http.StatusCreated, "upload complete", resp)
}

func (h *Handler) GetDraft(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft ID is required", nil))
		return
	}

	draft, err := h.service.GetDraftByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.New(http.StatusNotFound, appErr.CodeNotFound, "draft not found", nil))
			return
		}
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to get draft", err))
		return
	}

	response.Success(c, http.StatusOK, "draft retrieved successfully", draft)
}

func (h *Handler) ListDrafts(c *gin.Context) {
	status := c.Query("status")

	p := pagination.FromRequest(c.Request)

	res, err := h.service.ListDrafts(c.Request.Context(), status, p.Page, p.Limit)
	if err != nil {
		if errors.Is(err, ErrInvalidStatus) {
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "invalid status filter", err))
			return
		}
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to list drafts", err))
		return
	}

	if res.Items == nil {
		res.Items = []DraftListItem{}
	}

	meta := pagination.NewMeta(p, res.Total)

	response.Success(c, http.StatusOK, "drafts retrieved successfully", res.Items, meta)
}

func (h *Handler) EnrichDraft(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft ID is required", nil))
		return
	}

	// Parse optional scope from body (targeted re-enrichment).
	// Body is optional — treat empty/missing body as full enrichment.
	var req EnrichDraftRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			req = EnrichDraftRequest{}
		}
	}

	var scope enrichment.EnrichScope
	if len(req.Scope) > 0 {
		for _, s := range req.Scope {
			scope = append(scope, enrichment.Source(s))
		}
	}

	if err := h.service.EnrichDraft(c.Request.Context(), id, scope); err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.New(http.StatusNotFound, appErr.CodeNotFound, "draft not found", nil))
			return
		}
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to start enrichment", err))
		return
	}

	status := "enriching"
	if len(scope) > 0 {
		status = "refetching"
	}
	response.Success(c, http.StatusOK, "enrichment started", gin.H{"draftId": id, "status": status})
}

func (h *Handler) GetDraftSuggestions(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft ID is required", nil))
		return
	}

	suggestions, err := h.service.GetDraftSuggestions(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.New(http.StatusNotFound, appErr.CodeNotFound, "draft not found", nil))
			return
		}
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to get suggestions", err))
		return
	}

	response.Success(c, http.StatusOK, "suggestions retrieved successfully", suggestions)
}

// UpdateDraftMetadata updates a draft's extracted metadata fields (title, artist, album)
// without changing its status. This allows saving in-progress edits before finalization.
func (h *Handler) UpdateDraftMetadata(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft ID is required", nil))
		return
	}

	var req UpdateDraftMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "invalid request body", err))
		return
	}

	if err := h.service.UpdateExtractedMetadata(c.Request.Context(), id, req.Title, req.Artist, req.Album, req.CoverURL); err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.New(http.StatusNotFound, appErr.CodeNotFound, "draft not found", nil))
			return
		}
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to update draft metadata", err))
		return
	}

	response.Success(c, http.StatusOK, "draft metadata updated", gin.H{})
}

func (h *Handler) SaveFinalMetadata(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft ID is required", nil))
		return
	}

	var req SaveFinalMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "invalid request body", err))
		return
	}

	if err := h.service.SaveFinalMetadata(c.Request.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, ErrDraftNotFound):
			response.Error(c, appErr.New(http.StatusNotFound, appErr.CodeNotFound, "draft not found", nil))
		case errors.Is(err, ErrInvalidStatus):
			response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft must be in review status", nil))
		default:
			response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to save final metadata", err))
		}
		return
	}

	response.Success(c, http.StatusOK, "final metadata saved and draft accepted", gin.H{})
}

func (h *Handler) DeleteDraft(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft ID is required", nil))
		return
	}

	userID := auth.UserIDFromContext(c)
	if userID == "" {
		response.Error(c, appErr.New(http.StatusUnauthorized, appErr.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	if err := h.service.DeleteDraft(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.New(http.StatusNotFound, appErr.CodeNotFound, "draft not found", nil))
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			response.Error(c, appErr.New(http.StatusForbidden, appErr.CodeForbidden, "you do not have permission to delete this draft", nil))
			return
		}
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to delete draft", err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) RejectDraft(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft ID is required", nil))
		return
	}

	var req RejectDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Reason = c.DefaultQuery("reason", "")
	}

	if err := h.service.RejectDraft(c.Request.Context(), id, req.Reason); err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.New(http.StatusNotFound, appErr.CodeNotFound, "draft not found", nil))
			return
		}
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to reject draft", err))
		return
	}

	response.Success(c, http.StatusOK, "draft rejected", gin.H{})
}

func (h *Handler) SearchArtists(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "search query 'q' is required", nil))
		return
	}

	results, err := h.service.SearchArtists(c.Request.Context(), q)
	if err != nil {
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to search artists", err))
		return
	}

	response.Success(c, http.StatusOK, "artists retrieved", results)
}

func (h *Handler) SetFinalizer(f *finalization.Service) {
	h.finalizer = f
}

func (h *Handler) FinalizeDraft(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft ID is required", nil))
		return
	}

	if h.finalizer == nil {
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "finalization not configured", nil))
		return
	}

	draft, err := h.service.GetDraftRaw(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.New(http.StatusNotFound, appErr.CodeNotFound, "draft not found", nil))
			return
		}
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to get draft", err))
		return
	}

	var finalMeta string
	if draft.FinalMetadata != nil {
		finalMeta = *draft.FinalMetadata
	}

	result, err := h.finalizer.Finalize(c.Request.Context(), draft.ID, string(draft.Status), draft.FilePath, draft.Format, finalMeta, draft.ExtractedMetadata)
	if err != nil {
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to finalize draft: "+err.Error(), err))
		return
	}

	// Fire-and-forget AI lyrics generation if ML client is configured
	if h.mlClient != nil && result.TrackID != "" {
		h.enqueueAILyricsGeneration(c.Request.Context(), draft, result)
	}

	response.Success(c, http.StatusOK, "draft published successfully", result)
}

// enqueueAILyricsGeneration fires a non-blocking request to the ML service
// to generate AI lyrics for tracks that don't already have them.
func (h *Handler) enqueueAILyricsGeneration(ctx context.Context, draft *IngestionDraft, result *finalization.FinalizeResult) {
	go func() {
		// Parse final metadata for track info
		title := ""
		artist := ""
		if draft.FinalMetadata != nil && *draft.FinalMetadata != "" {
			var meta finalization.DraftMetadata
			if err := json.Unmarshal([]byte(*draft.FinalMetadata), &meta); err == nil {
				title = meta.Track.Title
				artist = meta.Artist.Name
			}
		}

		// Build the request to the ML service.
		// For local storage: pass the relative storage key as AudioFilePath.
		// For S3 storage: resolve a (pre)signed URL as AudioDownloadURL.
		var audioFilePath, audioDownloadURL string
		if h.storageDriver == "local" {
			audioFilePath = draft.FilePath
		} else if h.storage != nil {
			if url, err := h.storage.GetURL(ctx, draft.FilePath); err == nil {
				audioDownloadURL = url
			}
		}

		req := lyrics.CreateLyricsJobRequest{
			TrackID:          result.TrackID,
			AudioFilePath:    audioFilePath,
			AudioDownloadURL: audioDownloadURL,
			TrackTitle:       title,
			TrackArtist:      artist,
		}

		jobID, err := h.mlClient.EnqueueLyricsJob(ctx, req)
		if err != nil {
			zap.L().Error("failed to enqueue lyrics job", zap.String("track_id", result.TrackID), zap.Error(err))
			return
		}
		zap.L().Info("enqueued lyrics job", zap.String("job_id", jobID), zap.String("track_id", result.TrackID))
	}()
}

func (h *Handler) SearchAlbums(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "search query 'q' is required", nil))
		return
	}

	results, err := h.service.SearchAlbums(c.Request.Context(), q)
	if err != nil {
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to search albums", err))
		return
	}

	response.Success(c, http.StatusOK, "albums retrieved", results)
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetIngestionStats(c.Request.Context())
	if err != nil {
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to get ingestion stats", err))
		return
	}
	response.Success(c, http.StatusOK, "ingestion stats retrieved", stats)
}

func (h *Handler) TriggerCleanup(c *gin.Context) {
	if h.cleanup == nil {
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "cleanup service not configured", nil))
		return
	}
	flagged, err := h.cleanup.FlagStaleDrafts(c.Request.Context(), 24*time.Hour)
	if err != nil {
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "cleanup failed", err))
		return
	}
	response.Success(c, http.StatusOK, "cleanup completed", gin.H{"flagged": flagged})
}

func (h *Handler) UploadDraftImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "draft ID is required", nil))
		return
	}

	entityType := c.Param("entity")
	if entityType != "artist" && entityType != "album" && entityType != "track" && entityType != "cover" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "entity must be 'artist', 'album', 'track', or 'cover'", nil))
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "no image file provided in 'image' field", err))
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		response.Error(c, appErr.New(http.StatusBadRequest, appErr.CodeBadRequest, "unsupported image format (allowed: jpg, png, webp)", nil))
		return
	}

	imageURL, err := h.service.UploadDraftImage(c.Request.Context(), id, entityType, file, ext)
	if err != nil {
		if errors.Is(err, ErrDraftNotFound) {
			response.Error(c, appErr.New(http.StatusNotFound, appErr.CodeNotFound, "draft not found", nil))
			return
		}
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to upload image", err))
		return
	}

	response.Success(c, http.StatusOK, "image uploaded", gin.H{"url": imageURL})
}

func (h *Handler) GetConfig(c *gin.Context) {
	response.Success(c, http.StatusOK, "ingestion config", IngestionConfigResponse{
		MaxUploadSize:     200 << 20,
		EnrichmentEnabled: true,
	})
}
