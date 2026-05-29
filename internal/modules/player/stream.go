package player

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func detectAudioContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg", ".oga":
		return "audio/ogg"
	case ".m4a":
		return "audio/mp4"
	case ".aac":
		return "audio/aac"
	case ".flac":
		return "audio/flac"
	case ".webm":
		return "audio/webm"
	default:
		contentType := mime.TypeByExtension(ext)
		if contentType != "" {
			return contentType
		}

		return "application/octet-stream"
	}
}

func setStreamingHeaders(w http.ResponseWriter, contentType string) {
	header := w.Header()
	header.Set("Content-Type", contentType)
	header.Set("Accept-Ranges", "bytes")

	// Browser/player performance.
	header.Set("Cache-Control", "public, max-age=86400")

	// Useful for frontend/debugging when CORS is enabled.
	header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges, Content-Type")
}
