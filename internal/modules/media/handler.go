package media

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	appErr "music/internal/contracts/errors"
	"music/internal/shared/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
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
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/media/upload [post]
func (h *Handler) UploadAdminMedia(c *gin.Context) {
	const hardLimit = 60 << 20 // 60 MB

	// Set max bytes reader on the request
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, hardLimit)

	// Parse multipart form
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
		err := file.Close()
		if err != nil {

		}
	}(file)

	res, err := h.service.Upload(c.Request.Context(), category, file, header)
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

	// Use GinSuccess with status 201 Created
	response.Success(c, http.StatusCreated, "file uploaded successfully", res)
}

// detectUploadField unchanged because it works on *http.Request
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

var _ multipart.File
