package handler

import (
	"fmt"
	"music/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	service *service.SongService
}

func NewSongHandler(svc *service.SongService) *SongHandler {
	return &SongHandler{service: svc}
}

func (h *SongHandler) Upload(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	song, err := h.service.Upload(c.Request.Context(), title, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "song uploaded successfully",
		"song":    song,
	})
}

func (h *SongHandler) List(c *gin.Context) {
	songs, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}

func (h *SongHandler) Get(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	song, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}

func (h *SongHandler) Delete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "song deleted"})
}

func (h *SongHandler) Stream(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "song id is required"})
		return
	}

	filePath, contentType, err := h.service.GetAudioFilePath(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "song not found",
			"id":    id,
		})
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "audio file not found on disk",
				"path":  filePath,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot open audio file",
		})
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot read audio file",
		})
		return
	}

	if stat.IsDir() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid audio file",
		})
		return
	}

	fileName := filepath.Base(filePath)

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, fileName))
	c.Header("Accept-Ranges", "bytes")

	http.ServeContent(c.Writer, c.Request, fileName, stat.ModTime(), file)
}
