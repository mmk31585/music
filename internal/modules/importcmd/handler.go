package importcmd

import (
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"

	appErr "music/internal/common/errors"
	"music/internal/common/response"
	ingestion "music/internal/modules/ingestion"
)

type Handler struct {
	service     *Service
	ingestion   *ingestion.Service
	downloadDir string
}

func NewHandler(service *Service, ingestionSvc *ingestion.Service, downloadDir string) *Handler {
	return &Handler{service: service, ingestion: ingestionSvc, downloadDir: downloadDir}
}

func (h *Handler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Error(c, appErr.BadRequest("search query 'q' is required", nil))
		return
	}

	results, err := h.service.Search(c.Request.Context(), q)
	if err != nil {
		response.Error(c, appErr.Internal("search failed: "+err.Error(), err))
		return
	}

	response.OK(c, "search results retrieved", results)
}

func (h *Handler) Import(c *gin.Context) {
	if err := os.MkdirAll(h.downloadDir, 0755); err != nil {
		response.Error(c, appErr.Internal("failed to create temp directory", err))
		return
	}

	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", err))
		return
	}

	entry, localPath, err := h.service.Download(c.Request.Context(), req.URL, h.downloadDir)
	if err != nil {
		response.Error(c, appErr.Internal("download failed: "+err.Error(), err))
		return
	}

	uploadedBy := getUserID(c)
	if uploadedBy == "" {
		os.Remove(localPath)
		response.Error(c, appErr.Unauthorized("user not authenticated", nil))
		return
	}

	f, err := os.Open(localPath)
	if err != nil {
		os.Remove(localPath)
		response.Error(c, appErr.Internal("failed to open downloaded file", err))
		return
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		os.Remove(localPath)
		response.Error(c, appErr.Internal("failed to stat file", err))
		return
	}

	header := &multipart.FileHeader{
		Filename: entry.Title + ".mp3",
		Size:     stat.Size(),
	}

	result, err := h.ingestion.Upload(c.Request.Context(), f, header, uploadedBy)
	f.Close()
	os.Remove(localPath)

	if err != nil {
		response.Error(c, appErr.Internal("failed to create draft: "+err.Error(), err))
		return
	}

	resp := ImportResponse{
		DraftID:  result.DraftID,
		Title:    entry.Title,
		Artist:   entry.Uploader,
		Duration: int(entry.Duration),
		Message:  "track added to ingestion. review and publish from the ingestion page.",
	}

	response.Created(c, "track imported successfully", resp)
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
