package upload

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"
	"music/internal/modules/auth"
	platformstorage "music/internal/platform/storage"
)

type Handler struct {
	svc     ServiceInterface
	storage platformstorage.Storage
}

func NewHandler(svc ServiceInterface, storage platformstorage.Storage) *Handler {
	return &Handler{svc: svc, storage: storage}
}

func (h *Handler) CreateDraft(c *gin.Context) error {
	userID := auth.UserIDFromContext(c)
	if userID == "" {
		return apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "authentication required", nil)
	}

	var req CreateUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	draft, err := h.svc.CreateDraft(c.Request.Context(), userID, req)
	if err != nil {
		return err
	}

	response.Success(c, http.StatusCreated, "upload draft created", draft)
	return nil
}

func (h *Handler) UploadDraftFile(c *gin.Context) error {
	const hardLimit = 200 << 20

	userID := auth.UserIDFromContext(c)
	if userID == "" {
		return apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "authentication required", nil)
	}

	draftID := c.Param("draftId")

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, hardLimit)
	if err := c.Request.ParseMultipartForm(hardLimit); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeBadRequest, "file too large or invalid multipart form", err)
	}

	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeBadRequest, "no audio file provided in 'audio' field", err)
	}
	defer file.Close()

	// Verify content type is audio
	contentType := header.Header.Get("Content-Type")
	if !isAudioContentType(contentType) {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeBadRequest, "invalid file type; audio files only", nil)
	}

	ext := filepathExt(header.Filename)
	key := fmt.Sprintf("uploads/listener/%s/%d%s", userID, time.Now().UnixMilli(), ext)

	if err := h.storage.Upload(c.Request.Context(), key, file, header.Size, contentType); err != nil {
		return apperrors.New(http.StatusInternalServerError, apperrors.CodeInternal, "failed to store file", err)
	}

	url, err := h.storage.GetURL(c.Request.Context(), key)
	if err != nil {
		return apperrors.New(http.StatusInternalServerError, apperrors.CodeInternal, "failed to get file URL", err)
	}

	if err := h.svc.AttachFile(c.Request.Context(), draftID, userID, key, url, header.Size); err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "file uploaded", map[string]any{
		"draftId":  draftID,
		"fileUrl":  url,
		"fileSize": header.Size,
	})
	return nil
}

func isAudioContentType(ct string) bool {
	if strings.HasPrefix(ct, "audio/") {
		return true
	}
	return false
}

func filepathExt(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ".bin"
}

func (h *Handler) GetDraft(c *gin.Context) error {
	draftID := c.Param("draftId")
	draft, err := h.svc.GetDraft(c.Request.Context(), draftID)
	if err != nil {
		return err
	}
	response.Success(c, http.StatusOK, "draft retrieved", draft)
	return nil
}

func (h *Handler) ListPendingDrafts(c *gin.Context) error {
	source := c.Query("source")
	var sourcePtr *string
	if source != "" {
		sourcePtr = &source
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	drafts, total, err := h.svc.ListPendingDrafts(c.Request.Context(), sourcePtr, limit, offset)
	if err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "pending drafts retrieved", drafts, map[string]any{
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
	return nil
}

func (h *Handler) ListMyDrafts(c *gin.Context) error {
	userID := auth.UserIDFromContext(c)
	if userID == "" {
		return apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "authentication required", nil)
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	drafts, total, err := h.svc.ListUserDrafts(c.Request.Context(), userID, limit, offset)
	if err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "your drafts retrieved", drafts, map[string]any{
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
	return nil
}

func (h *Handler) ReviewDraft(c *gin.Context) error {
	reviewerID := auth.UserIDFromContext(c)
	if reviewerID == "" {
		return apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "authentication required", nil)
	}

	draftID := c.Param("draftId")
	var req ReviewDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	if err := h.svc.ReviewDraft(c.Request.Context(), draftID, reviewerID, req); err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "draft reviewed", struct{}{})
	return nil
}

func (h *Handler) GetUploadSlots(c *gin.Context) error {
	userID := auth.UserIDFromContext(c)
	if userID == "" {
		return apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "authentication required", nil)
	}

	slots, err := h.svc.GetUploadSlots(c.Request.Context(), userID)
	if err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "upload slots retrieved", slots)
	return nil
}

func (h *Handler) GetCoUploaders(c *gin.Context) error {
	trackIDStr := c.Param("trackId")
	trackID, err := strconv.ParseInt(trackIDStr, 10, 64)
	if err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeBadRequest, "invalid track ID", nil)
	}

	cos, err := h.svc.GetCoUploaders(c.Request.Context(), trackID)
	if err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "co-uploaders retrieved", cos)
	return nil
}

func (h *Handler) AddCoUploader(c *gin.Context) error {
	trackIDStr := c.Param("trackId")
	trackID, err := strconv.ParseInt(trackIDStr, 10, 64)
	if err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeBadRequest, "invalid track ID", nil)
	}

	var req CoUploaderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	co, err := h.svc.AddCoUploader(c.Request.Context(), trackID, req)
	if err != nil {
		return err
	}

	response.Success(c, http.StatusCreated, "co-uploader added", co)
	return nil
}

func (h *Handler) RemoveCoUploader(c *gin.Context) error {
	trackIDStr := c.Param("trackId")
	trackID, err := strconv.ParseInt(trackIDStr, 10, 64)
	if err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeBadRequest, "invalid track ID", nil)
	}

	userID := c.Param("userId")
	if err := h.svc.RemoveCoUploader(c.Request.Context(), trackID, userID); err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "co-uploader removed", struct{}{})
	return nil
}
