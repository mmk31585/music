package media

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appErr "music/internal/common/errors"
	"music/internal/common/response"
)

// ServiceInterface defines the service methods needed by the HTTP handler.
type ServiceInterface interface {
	Upload(ctx context.Context, category UploadCategory, file multipart.File, header *multipart.FileHeader, createdBy *uuid.UUID) (*UploadResponse, error)
	ListMedia(ctx context.Context) ([]Media, error)
	DeleteMedia(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// UploadAdminMedia godoc
// @Summary Upload media file
// @Description Uploads a single media file for artist images, album covers, track covers, or track audio.
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param artistImage formData file false "Artist image file"
// @Param albumCover formData file false "Album cover file"
// @Param trackCover formData file false "Track cover file"
// @Param trackAudio formData file false "Track audio file"
// @Param playlistCover formData file false "Playlist cover file"
// @Param video formData file false "Video file"
// @Param videoAudio formData file false "Video audio file"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/media/upload [post]
func (h *Handler) UploadAdminMedia(c *gin.Context) {
	const hardLimit = 60 << 20 // 60 MB

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, hardLimit)

	if err := c.Request.ParseMultipartForm(hardLimit); err != nil {
		response.Error(c, appErr.BadRequest("invalid multipart form or file too large", err))
		return
	}

	category, fieldName, err := detectUploadField(c.Request)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoFileProvided):
			response.Error(c, appErr.BadRequest("no valid file field provided", err))
		case errors.Is(err, ErrMultipleFiles):
			response.Error(c, appErr.BadRequest("only one upload field is allowed", err))
		default:
			response.Error(c, appErr.BadRequest("invalid upload field", err))
		}
		return
	}

	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		response.Error(c, appErr.BadRequest("failed to read uploaded file", err))
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	createdBy := getUserIDFromContext(c)

	res, err := h.service.Upload(c.Request.Context(), category, file, header, createdBy)

	if err != nil {
		switch {
		case errors.Is(err, ErrFileTooLarge):
			response.Error(c, appErr.BadRequest("file too large", err))
		case errors.Is(err, ErrInvalidMimeType):
			response.Error(c, appErr.BadRequest("invalid file type", err))
		case errors.Is(err, ErrNoFileProvided):
			response.Error(c, appErr.BadRequest("no file provided", err))
		case errors.Is(err, ErrEmptyFile):
			response.Error(c, appErr.BadRequest("uploaded file is empty", err))
		case errors.Is(err, ErrInvalidFieldName):
			response.Error(c, appErr.BadRequest("invalid upload field name", err))
		case errors.Is(err, ErrStorageFailed):
			response.Error(c, appErr.Internal("failed to store file", err))
		default:
			response.Error(c, appErr.Internal("failed to upload file", err))
		}
		return
	}

	response.Success(c, http.StatusCreated, "file uploaded successfully", res)
}

// ListAdminMedia godoc
// @Summary List all media
// @Description Returns all media items ordered by creation date descending
// @Tags media
// @Produce json
// @Security Bearer
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/media [get]
func (h *Handler) ListAdminMedia(c *gin.Context) {
	items, err := h.service.ListMedia(c.Request.Context())
	if err != nil {
		response.Error(c, appErr.Internal("failed to list media", err))
		return
	}

	response.Success(c, http.StatusOK, "media list retrieved successfully", items)
}

// DeleteAdminMedia godoc
// @Summary Delete a media item
// @Description Deletes a media item by ID (removes file from storage and DB record)
// @Tags media
// @Produce json
// @Security Bearer
// @Param id path string true "Media ID (UUID)"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/media/{id} [delete]
func (h *Handler) DeleteAdminMedia(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid media ID", err))
		return
	}

	if err := h.service.DeleteMedia(c.Request.Context(), id); err != nil {
		response.Error(c, appErr.Internal("failed to delete media", err))
		return
	}

	response.SuccessNoContent(c)
}

func detectUploadField(r *http.Request) (UploadCategory, string, error) {
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return "", "", ErrNoFileProvided
	}

	fields := []struct {
		Name     string
		Category UploadCategory
	}{
		{
			Name:     "artistImage",
			Category: UploadCategoryArtistImage,
		},
		{
			Name:     "albumCover",
			Category: UploadCategoryAlbumCover,
		},
		{
			Name:     "trackCover",
			Category: UploadCategoryTrackCover,
		},
		{
			Name:     "trackAudio",
			Category: UploadCategoryTrackAudio,
		},
		{
			Name:     "playlistCover",
			Category: UploadCategoryPlaylistCover,
		},
		{
			Name:     "video",
			Category: UploadCategoryVideo,
		},
		{
			Name:     "videoAudio",
			Category: UploadCategoryVideoAudio,
		},
	}

	var selectedName string
	var selectedCategory UploadCategory
	count := 0

	for _, field := range fields {
		files := r.MultipartForm.File[field.Name]
		if len(files) > 0 {
			count++
			selectedName = field.Name
			selectedCategory = field.Category
		}
	}

	if count == 0 {
		return "", "", ErrNoFileProvided
	}

	if count > 1 {
		return "", "", ErrMultipleFiles
	}

	return selectedCategory, selectedName, nil
}
func getUserIDFromContext(c *gin.Context) *uuid.UUID {
	keys := []string{"auth_user_id", "user_id", "userID", "userId", "sub"}

	for _, key := range keys {
		raw, exists := c.Get(key)
		if !exists || raw == nil {
			continue
		}

		switch v := raw.(type) {
		case uuid.UUID:
			return &v

		case *uuid.UUID:
			return v

		case string:
			if v == "" {
				continue
			}

			parsed, err := uuid.Parse(v)
			if err == nil {
				return &parsed
			}
		}
	}

	return nil
}
