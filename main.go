package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Song struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	ContentType string    `json:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

const (
	serverAddr = ":8080"

	dataDir   = "data"
	dataFile  = "data/songs.json"
	uploadDir = "uploads"

	maxUploadSize = 50 << 20 // 50 MB
)

var (
	songs      = make(map[string]Song)
	songsMutex sync.RWMutex
)

var allowedOrigins = map[string]bool{
	"http://localhost:3000": true,
	"http://127.0.0.1:3000": true,
	"http://localhost:5173": true,
	"http://127.0.0.1:5173": true,
}

func main() {
	if err := ensureDirectories(); err != nil {
		panic(err)
	}

	if err := loadSongs(); err != nil {
		fmt.Println("failed to load songs:", err)
	}

	r := gin.Default()

	r.MaxMultipartMemory = maxUploadSize

	r.Use(corsMiddleware())

	r.GET("/health", healthCheck)

	r.POST("/upload", uploadSong)
	r.GET("/songs", listSongs)
	r.GET("/songs/:id", getSong)
	r.GET("/stream/:id", streamSong)

	fmt.Println("Server running on http://localhost:8080")

	if err := r.Run(serverAddr); err != nil {
		panic(err)
	}
}

func ensureDirectories() error {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	return nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, x-device-id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")
		c.Writer.Header().Set("Access-Control-Max-Age", "43200")

		if c.Request.Method == http.MethodOptions {
			if origin != "" && !allowedOrigins[origin] {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		if origin != "" && !allowedOrigins[origin] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"type":    "error",
				"message": "origin is not allowed by CORS",
			})
			return
		}

		c.Next()
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "server is running",
		"status":  "ok",
	})
}

func uploadSong(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		jsonError(c, http.StatusBadRequest, "title is required")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			jsonError(c, http.StatusBadRequest, "file is required")
			return
		}

		jsonError(c, http.StatusBadRequest, "invalid file upload")
		return
	}
	defer file.Close()

	originalFileName := filepath.Base(header.Filename)
	ext := strings.ToLower(filepath.Ext(originalFileName))

	if !isAllowedAudioExt(ext) {
		jsonError(c, http.StatusBadRequest, "only mp3, wav, ogg files are allowed")
		return
	}

	id := generateID()
	storedFileName := id + ext
	filePath := filepath.Join(uploadDir, storedFileName)

	dst, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, "failed to create uploaded file")
		return
	}

	written, copyErr := io.Copy(dst, file)
	closeErr := dst.Close()

	if copyErr != nil {
		_ = os.Remove(filePath)
		jsonError(c, http.StatusInternalServerError, "failed to save uploaded file")
		return
	}

	if closeErr != nil {
		_ = os.Remove(filePath)
		jsonError(c, http.StatusInternalServerError, "failed to close uploaded file")
		return
	}

	if written == 0 {
		_ = os.Remove(filePath)
		jsonError(c, http.StatusBadRequest, "uploaded file is empty")
		return
	}

	contentType := detectContentType(header.Header.Get("Content-Type"), ext)

	song := Song{
		ID:          id,
		Title:       title,
		FileName:    originalFileName,
		FilePath:    filePath,
		ContentType: contentType,
		UploadedAt:  time.Now().UTC(),
	}

	songsMutex.Lock()
	songs[id] = song
	songsMutex.Unlock()

	if err := saveSongs(); err != nil {
		_ = os.Remove(filePath)

		songsMutex.Lock()
		delete(songs, id)
		songsMutex.Unlock()

		jsonError(c, http.StatusInternalServerError, "failed to save song metadata")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "song uploaded successfully",
		"song":    song,
	})
}

func listSongs(c *gin.Context) {
	songsMutex.RLock()
	result := make([]Song, 0, len(songs))

	for _, song := range songs {
		result = append(result, song)
	}
	songsMutex.RUnlock()

	sort.Slice(result, func(i, j int) bool {
		return result[i].UploadedAt.After(result[j].UploadedAt)
	})

	c.JSON(http.StatusOK, result)
}

func getSong(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	song, exists := findSong(id)
	if !exists {
		jsonError(c, http.StatusNotFound, "song not found")
		return
	}

	c.JSON(http.StatusOK, song)
}

func streamSong(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	song, exists := findSong(id)
	if !exists {
		jsonError(c, http.StatusNotFound, "song not found")
		return
	}

	file, err := os.Open(song.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			jsonError(c, http.StatusNotFound, "audio file not found")
			return
		}

		jsonError(c, http.StatusInternalServerError, "cannot open audio file")
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		jsonError(c, http.StatusInternalServerError, "cannot read audio file")
		return
	}

	if stat.IsDir() {
		jsonError(c, http.StatusInternalServerError, "invalid audio file")
		return
	}

	c.Header("Content-Type", song.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, sanitizeHeaderFileName(song.FileName)))
	c.Header("Accept-Ranges", "bytes")

	http.ServeContent(c.Writer, c.Request, song.FileName, stat.ModTime(), file)
}

func findSong(id string) (Song, bool) {
	songsMutex.RLock()
	defer songsMutex.RUnlock()

	song, exists := songs[id]
	return song, exists
}

func saveSongs() error {
	songsMutex.RLock()
	list := make([]Song, 0, len(songs))

	for _, song := range songs {
		list = append(list, song)
	}
	songsMutex.RUnlock()

	sort.Slice(list, func(i, j int) bool {
		return list[i].UploadedAt.After(list[j].UploadedAt)
	})

	tempFile := dataFile + ".tmp"

	file, err := os.Create(tempFile)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(list); err != nil {
		_ = file.Close()
		_ = os.Remove(tempFile)
		return err
	}

	if err := file.Close(); err != nil {
		_ = os.Remove(tempFile)
		return err
	}

	return os.Rename(tempFile, dataFile)
}

func loadSongs() error {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}
	defer file.Close()

	var list []Song

	if err := json.NewDecoder(file).Decode(&list); err != nil {
		return err
	}

	songsMutex.Lock()
	defer songsMutex.Unlock()

	for _, song := range list {
		if song.ID == "" {
			continue
		}

		if song.FilePath == "" {
			continue
		}

		songs[song.ID] = song
	}

	return nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func isAllowedAudioExt(ext string) bool {
	switch ext {
	case ".mp3", ".wav", ".ogg":
		return true
	default:
		return false
	}
}

func detectContentType(uploadedContentType string, ext string) string {
	uploadedContentType = strings.TrimSpace(uploadedContentType)

	if strings.HasPrefix(uploadedContentType, "audio/") {
		return uploadedContentType
	}

	if byExt := mime.TypeByExtension(ext); byExt != "" {
		return byExt
	}

	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	default:
		return "application/octet-stream"
	}
}

func sanitizeHeaderFileName(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, `"`, "")
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\r", "")

	if name == "" {
		return "audio"
	}

	return name
}

func jsonError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"type":    "error",
		"message": message,
	})
}
