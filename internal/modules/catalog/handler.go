package catalog

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appErr "music/internal/common/errors"
	"music/internal/common/pagination"
	"music/internal/common/response"
	"music/internal/common/validator"
	"music/internal/modules/media"
)

type Handler struct {
	service      *Service
	validator    *validator.Validator
	mediaService *media.Service
}

func NewHandler(service *Service, validator *validator.Validator, mediaService ...*media.Service) *Handler {
	var uploader *media.Service
	if len(mediaService) > 0 {
		uploader = mediaService[0]
	}

	return &Handler{
		service:      service,
		validator:    validator,
		mediaService: uploader,
	}
}

// --------------------
// Artists
// --------------------

func (h *Handler) RegisterArtist(c *gin.Context) {
	var req CreateArtistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", nil))
		return
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		response.Error(c, appErr.Validation("validation failed", errs))
		return
	}

	artist, err := h.service.CreateArtist(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicateArtistSlug):
			response.Error(c, appErr.Conflict("artist with similar name already exists",
				map[string]string{"name": req.Name}))
		default:
			response.Error(c, appErr.Internal("failed to create artist", nil))
		}
		return
	}

	response.Success(c, http.StatusCreated, "artist created successfully", artist)
}

func (h *Handler) GetArtist(c *gin.Context) {
	idStr := c.Param("artistID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid artist id",
			map[string]interface{}{"id": idStr}))
		return
	}

	artist, err := h.service.GetArtist(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrArtistNotFound):
			response.Error(c, appErr.NotFound("artist not found",
				map[string]string{"id": id.String()}))
		default:
			response.Error(c, appErr.Internal("failed to fetch artist", nil))
		}
		return
	}

	response.Success(c, http.StatusOK, "artist fetched successfully", artist)
}

func (h *Handler) ListArtists(c *gin.Context) {
	filter := ParseArtistListFilter(c) // changed to accept gin.Context

	artists, total, err := h.service.ListArtists(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, appErr.Internal("failed to list artists", nil))
		return
	}

	meta := pagination.NewMeta(filter.Pagination, total)
	response.SuccessWithMeta(c, http.StatusOK, "artists fetched successfully", artists, meta)
}

func (h *Handler) UpdateArtist(c *gin.Context) {
	idStr := c.Param("artistID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid artist id",
			map[string]interface{}{"id": idStr}))
		return
	}

	var req UpdateArtistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", nil))
		return
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		response.Error(c, appErr.Validation("validation failed", errs))
		return
	}

	artist, err := h.service.UpdateArtist(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrArtistNotFound):
			response.Error(c, appErr.NotFound("artist not found",
				map[string]string{"id": id.String()}))
		case errors.Is(err, ErrDuplicateArtistSlug):
			detail := map[string]string{}
			if req.Name != nil {
				detail["name"] = *req.Name
			}
			response.Error(c, appErr.Conflict("artist with similar name already exists", detail))
		default:
			response.Error(c, appErr.Internal("failed to update artist", nil))
		}
		return
	}

	response.Success(c, http.StatusOK, "artist updated successfully", artist)
}

func (h *Handler) DeleteArtist(c *gin.Context) {
	idStr := c.Param("artistID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid artist id",
			map[string]interface{}{"id": idStr}))
		return
	}

	if err := h.service.DeleteArtist(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, ErrArtistNotFound):
			response.Error(c, appErr.NotFound("artist not found",
				map[string]string{"id": id.String()}))
		default:
			response.Error(c, appErr.Internal("failed to delete artist", nil))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "artist deleted successfully", nil)
}

// --------------------
// Albums
// --------------------

func (h *Handler) RegisterAlbum(c *gin.Context) {
	var req CreateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", nil))
		return
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		response.Error(c, appErr.Validation("validation failed", errs))
		return
	}

	album, err := h.service.CreateAlbum(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrArtistNotFound):
			response.Error(c, appErr.NotFound("artist not found",
				map[string]string{"artistId": req.ArtistID.String()}))
		case errors.Is(err, ErrDuplicateAlbumSlug):
			response.Error(c, appErr.Conflict("album with similar title already exists for this artist",
				map[string]string{"title": req.Title}))
		default:
			response.Error(c, appErr.Internal("failed to create album", nil))
		}
		return
	}

	response.Success(c, http.StatusCreated, "album created successfully", album)
}

func (h *Handler) GetAlbum(c *gin.Context) {
	idStr := c.Param("albumID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid album id",
			map[string]interface{}{"id": idStr}))
		return
	}

	album, err := h.service.GetAlbum(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrAlbumNotFound):
			response.Error(c, appErr.NotFound("album not found",
				map[string]string{"id": id.String()}))
		default:
			response.Error(c, appErr.Internal("failed to fetch album", nil))
		}
		return
	}

	response.Success(c, http.StatusOK, "album fetched successfully", album)
}

func (h *Handler) ListAlbums(c *gin.Context) {
	filter := ParseAlbumListFilter(c)

	albums, total, err := h.service.ListAlbums(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, appErr.Internal("failed to list albums", nil))
		return
	}

	meta := pagination.NewMeta(filter.Pagination, total)
	response.SuccessWithMeta(c, http.StatusOK, "albums fetched successfully", albums, meta)
}

func (h *Handler) UpdateAlbum(c *gin.Context) {
	idStr := c.Param("albumID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid album id",
			map[string]interface{}{"id": idStr}))
		return
	}

	var req UpdateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", nil))
		return
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		response.Error(c, appErr.Validation("validation failed", errs))
		return
	}

	album, err := h.service.UpdateAlbum(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrAlbumNotFound):
			response.Error(c, appErr.NotFound("album not found",
				map[string]string{"id": id.String()}))
		case errors.Is(err, ErrDuplicateAlbumSlug):
			detail := map[string]string{}
			if req.Title != nil {
				detail["title"] = *req.Title
			}
			response.Error(c, appErr.Conflict("album with similar title already exists for this artist", detail))
		default:
			response.Error(c, appErr.Internal("failed to update album", nil))
		}
		return
	}

	response.Success(c, http.StatusOK, "album updated successfully", album)
}

func (h *Handler) DeleteAlbum(c *gin.Context) {
	idStr := c.Param("albumID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid album id",
			map[string]interface{}{"id": idStr}))
		return
	}

	if err := h.service.DeleteAlbum(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, ErrAlbumNotFound):
			response.Error(c, appErr.NotFound("album not found",
				map[string]string{"id": id.String()}))
		default:
			response.Error(c, appErr.Internal("failed to delete album", nil))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "album deleted successfully", nil)
}

// --------------------
// Tracks
// --------------------

func (h *Handler) RegisterTrack(c *gin.Context) {
	var req CreateTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// log the raw body for debugging
		body, _ := c.GetRawData()
		log.Printf("invalid JSON: %v, body: %s", err, string(body))
		response.Error(c, appErr.BadRequest("invalid request body", nil))
		return
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		response.Error(c, appErr.Validation("validation failed", errs))
		return
	}

	track, err := h.service.CreateTrack(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrArtistNotFound):
			response.Error(c, appErr.NotFound("artist not found",
				map[string]string{"artistId": req.ArtistID.String()}))
		case errors.Is(err, ErrAlbumNotFound):
			detail := map[string]string{}
			if req.AlbumID != nil {
				detail["albumId"] = req.AlbumID.String()
			}
			response.Error(c, appErr.NotFound("album not found", detail))
		case errors.Is(err, ErrInvalidGenreIDs):
			response.Error(c, appErr.BadRequest("one or more genre ids are invalid",
				map[string]interface{}{"genreIds": req.GenreIDs}))
		case errors.Is(err, ErrDuplicateTrackSlug):
			response.Error(c, appErr.Conflict("track with similar title already exists for this artist",
				map[string]string{"title": req.Title}))
		default:
			response.Error(c, appErr.Internal("failed to create track", nil))
		}
		return
	}

	response.Success(c, http.StatusCreated, "track created successfully", track)
}

func (h *Handler) RegisterTrackWithAudio(c *gin.Context) {
	if h.mediaService == nil {
		response.Error(c, appErr.Internal("media upload service is not configured", nil))
		return
	}

	const hardLimit = 70 << 20
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, hardLimit)

	if err := c.Request.ParseMultipartForm(hardLimit); err != nil {
		response.Error(c, appErr.BadRequest("invalid multipart form or file too large", err))
		return
	}

	file, header, err := c.Request.FormFile("trackAudio")
	if err != nil {
		response.Error(c, appErr.BadRequest("track audio file is required", err))
		return
	}
	defer file.Close()

	req, err := createTrackRequestFromMultipart(c)
	if err != nil {
		response.Error(c, appErr.BadRequest(err.Error(), nil))
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		req.Title = strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		response.Error(c, appErr.Validation("validation failed", errs))
		return
	}

	upload, err := h.mediaService.Upload(c.Request.Context(), media.UploadCategoryTrackAudio, file, header)
	if err != nil {
		switch {
		case errors.Is(err, media.ErrFileTooLarge):
			response.Error(c, appErr.BadRequest("file too large", err))
		case errors.Is(err, media.ErrInvalidMimeType):
			response.Error(c, appErr.BadRequest("invalid audio file type", err))
		case errors.Is(err, media.ErrEmptyFile):
			response.Error(c, appErr.BadRequest("uploaded file is empty", err))
		default:
			response.Error(c, appErr.Internal("failed to upload track audio", err))
		}
		return
	}

	req.AudioURL = &upload.URL

	track, err := h.service.CreateTrack(c.Request.Context(), req)
	if err != nil {
		_ = h.mediaService.DeleteUploadedFile(c.Request.Context(), upload)
		writeTrackCreateError(c, err, req)
		return
	}

	response.Success(c, http.StatusCreated, "track uploaded and created successfully", gin.H{
		"track": track,
		"media": upload,
	})
}

func createTrackRequestFromMultipart(c *gin.Context) (CreateTrackRequest, error) {
	artistID, err := uuid.Parse(strings.TrimSpace(c.PostForm("artistId")))
	if err != nil {
		return CreateTrackRequest{}, errors.New("valid artistId is required")
	}

	req := CreateTrackRequest{
		ArtistID: artistID,
		Title:    strings.TrimSpace(c.PostForm("title")),
	}

	if value := strings.TrimSpace(c.PostForm("albumId")); value != "" {
		albumID, err := uuid.Parse(value)
		if err != nil {
			return CreateTrackRequest{}, errors.New("albumId is invalid")
		}
		req.AlbumID = &albumID
	}

	if value := strings.TrimSpace(c.PostForm("durationSeconds")); value != "" {
		duration, err := strconv.Atoi(value)
		if err != nil || duration < 0 {
			return CreateTrackRequest{}, errors.New("durationSeconds must be a positive number")
		}
		req.DurationSeconds = duration
	}

	if value := strings.TrimSpace(c.PostForm("trackNumber")); value != "" {
		trackNumber, err := strconv.Atoi(value)
		if err != nil || trackNumber < 1 {
			return CreateTrackRequest{}, errors.New("trackNumber must be a positive number")
		}
		req.TrackNumber = &trackNumber
	}

	if value := strings.TrimSpace(c.PostForm("explicit")); value != "" {
		explicit, err := strconv.ParseBool(value)
		if err != nil {
			return CreateTrackRequest{}, errors.New("explicit must be true or false")
		}
		req.Explicit = &explicit
	}

	if value := strings.TrimSpace(c.PostForm("isPublic")); value != "" {
		isPublic, err := strconv.ParseBool(value)
		if err != nil {
			return CreateTrackRequest{}, errors.New("isPublic must be true or false")
		}
		req.IsPublic = &isPublic
	}

	if value := strings.TrimSpace(c.PostForm("coverUrl")); value != "" {
		req.CoverURL = &value
	}

	genreIDs, err := parseMultipartGenreIDs(c)
	if err != nil {
		return CreateTrackRequest{}, err
	}
	req.GenreIDs = genreIDs

	return req, nil
}

func parseMultipartGenreIDs(c *gin.Context) ([]uuid.UUID, error) {
	values := c.PostFormArray("genreIds")
	if single := strings.TrimSpace(c.PostForm("genreIds")); single != "" {
		values = append(values, strings.Split(single, ",")...)
	}

	ids := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		id, err := uuid.Parse(value)
		if err != nil {
			return nil, errors.New("one or more genreIds are invalid")
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func writeTrackCreateError(c *gin.Context, err error, req CreateTrackRequest) {
	switch {
	case errors.Is(err, ErrArtistNotFound):
		response.Error(c, appErr.NotFound("artist not found",
			map[string]string{"artistId": req.ArtistID.String()}))
	case errors.Is(err, ErrAlbumNotFound):
		detail := map[string]string{}
		if req.AlbumID != nil {
			detail["albumId"] = req.AlbumID.String()
		}
		response.Error(c, appErr.NotFound("album not found", detail))
	case errors.Is(err, ErrInvalidGenreIDs):
		response.Error(c, appErr.BadRequest("one or more genre ids are invalid",
			map[string]interface{}{"genreIds": req.GenreIDs}))
	case errors.Is(err, ErrDuplicateTrackSlug):
		response.Error(c, appErr.Conflict("track with similar title already exists for this artist",
			map[string]string{"title": req.Title}))
	default:
		response.Error(c, appErr.Internal("failed to create track", nil))
	}
}

func (h *Handler) GetTrack(c *gin.Context) {
	idStr := c.Param("trackID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid track id",
			map[string]interface{}{"id": idStr}))
		return
	}

	track, err := h.service.GetTrack(c.Request.Context(), id, false)
	if err != nil {
		switch {
		case errors.Is(err, ErrTrackNotFound):
			response.Error(c, appErr.NotFound("track not found",
				map[string]string{"id": id.String()}))
		default:
			response.Error(c, appErr.Internal("failed to fetch track", nil))
		}
		return
	}

	response.Success(c, http.StatusOK, "track fetched successfully", track)
}

func (h *Handler) ListTracks(c *gin.Context) {
	filter := ParseTrackListFilter(c, false)

	tracks, total, err := h.service.ListTracks(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, appErr.Internal("failed to list tracks", nil))
		return
	}

	meta := pagination.NewMeta(filter.Pagination, total)
	response.SuccessWithMeta(c, http.StatusOK, "tracks fetched successfully", tracks, meta)
}

func (h *Handler) ListAdminTracks(c *gin.Context) {
	filter := ParseTrackListFilter(c, true)

	tracks, total, err := h.service.ListTracks(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, appErr.Internal("failed to list admin tracks", nil))
		return
	}

	meta := pagination.NewMeta(filter.Pagination, total)
	response.SuccessWithMeta(c, http.StatusOK, "admin tracks fetched successfully", tracks, meta)
}

func (h *Handler) UpdateTrack(c *gin.Context) {
	idStr := c.Param("trackID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid track id",
			map[string]interface{}{"id": idStr}))
		return
	}

	var req UpdateTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", nil))
		return
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		response.Error(c, appErr.Validation("validation failed", errs))
		return
	}

	track, err := h.service.UpdateTrack(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTrackNotFound):
			response.Error(c, appErr.NotFound("track not found",
				map[string]string{"id": id.String()}))
		case errors.Is(err, ErrAlbumNotFound):
			detail := map[string]string{}
			if req.AlbumID != nil {
				detail["albumId"] = req.AlbumID.String()
			}
			response.Error(c, appErr.NotFound("album not found", detail))
		case errors.Is(err, ErrInvalidGenreIDs):
			response.Error(c, appErr.BadRequest("one or more genre ids are invalid",
				map[string]interface{}{"genreIds": req.GenreIDs}))
		case errors.Is(err, ErrDuplicateTrackSlug):
			detail := map[string]string{}
			if req.Title != nil {
				detail["title"] = *req.Title
			}
			response.Error(c, appErr.Conflict("track with similar title already exists for this artist", detail))
		default:
			response.Error(c, appErr.Internal("failed to update track", nil))
		}
		return
	}

	response.Success(c, http.StatusOK, "track updated successfully", track)
}

func (h *Handler) DeleteTrack(c *gin.Context) {
	idStr := c.Param("trackID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid track id",
			map[string]interface{}{"id": idStr}))
		return
	}

	if err := h.service.DeleteTrack(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, ErrTrackNotFound):
			response.Error(c, appErr.NotFound("track not found",
				map[string]string{"id": id.String()}))
		default:
			response.Error(c, appErr.Internal("failed to delete track", nil))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "track deleted successfully", nil)
}

// --------------------
// Genres
// --------------------

func (h *Handler) CreateGenre(c *gin.Context) {
	var req CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", nil))
		return
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		response.Error(c, appErr.Validation("validation failed", errs))
		return
	}

	genre, err := h.service.CreateGenre(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicateGenreName):
			response.Error(c, appErr.Conflict("genre name already exists",
				map[string]string{"name": req.Name}))
		case errors.Is(err, ErrDuplicateGenreSlug):
			response.Error(c, appErr.Conflict("genre slug already exists",
				map[string]string{"name": req.Name}))
		default:
			response.Error(c, appErr.Internal("failed to create genre", nil))
		}
		return
	}

	response.Success(c, http.StatusCreated, "genre created successfully", genre)
}

func (h *Handler) GetGenre(c *gin.Context) {
	idStr := c.Param("genreID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid genre id",
			map[string]interface{}{"id": idStr}))
		return
	}

	genre, err := h.service.GetGenre(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrGenreNotFound):
			response.Error(c, appErr.NotFound("genre not found",
				map[string]string{"id": id.String()}))
		default:
			response.Error(c, appErr.Internal("failed to fetch genre", nil))
		}
		return
	}

	response.Success(c, http.StatusOK, "genre fetched successfully", genre)
}

func (h *Handler) ListGenres(c *gin.Context) {
	genres, err := h.service.ListGenres(c.Request.Context())
	if err != nil {
		response.Error(c, appErr.Internal("failed to list genres", nil))
		return
	}

	response.Success(c, http.StatusOK, "genres fetched successfully", genres)
}

func (h *Handler) UpdateGenre(c *gin.Context) {
	idStr := c.Param("genreID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid genre id",
			map[string]interface{}{"id": idStr}))
		return
	}

	var req UpdateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErr.BadRequest("invalid request body", nil))
		return
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		response.Error(c, appErr.Validation("validation failed", errs))
		return
	}

	genre, err := h.service.UpdateGenre(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrGenreNotFound):
			response.Error(c, appErr.NotFound("genre not found",
				map[string]string{"id": id.String()}))
		case errors.Is(err, ErrDuplicateGenreName):
			detail := map[string]string{}
			if req.Name != nil {
				detail["name"] = *req.Name
			}
			response.Error(c, appErr.Conflict("genre name already exists", detail))
		case errors.Is(err, ErrDuplicateGenreSlug):
			detail := map[string]string{}
			if req.Name != nil {
				detail["name"] = *req.Name
			}
			response.Error(c, appErr.Conflict("genre slug already exists", detail))
		default:
			response.Error(c, appErr.Internal("failed to update genre", nil))
		}
		return
	}

	response.Success(c, http.StatusOK, "genre updated successfully", genre)
}

func (h *Handler) DeleteGenre(c *gin.Context) {
	idStr := c.Param("genreID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, appErr.BadRequest("invalid genre id",
			map[string]interface{}{"id": idStr}))
		return
	}

	if err := h.service.DeleteGenre(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, ErrGenreNotFound):
			response.Error(c, appErr.NotFound("genre not found",
				map[string]string{"id": id.String()}))
		default:
			response.Error(c, appErr.Internal("failed to delete genre", nil))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "genre deleted successfully", nil)
}

// --------------------
// Search
// --------------------

func (h *Handler) Search(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))

	result, err := h.service.Search(c.Request.Context(), q)
	if err != nil {
		response.Error(c, appErr.Internal("failed to search catalog", nil))
		return
	}

	response.Success(c, http.StatusOK, "search completed successfully", result)
}
