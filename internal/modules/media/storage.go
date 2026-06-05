package media

import (
	"mime"
	"regexp"
	"strings"
)

func sanitizeCategory(category string) string {
	category = strings.ToLower(strings.TrimSpace(category))

	re := regexp.MustCompile(`[^a-z0-9\-_]+`)
	category = re.ReplaceAllString(category, "-")

	category = strings.Trim(category, "-_")

	return category
}

func sanitizeExtension(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))

	if ext == "" {
		return ""
	}

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	re := regexp.MustCompile(`^\.[a-z0-9]+$`)
	if !re.MatchString(ext) {
		return ""
	}

	return ext
}

func extensionForMime(mimeType string) string {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))

	switch mimeType {
	case "audio/mpeg":
		return ".mp3"
	case "audio/ogg":
		return ".ogg"
	case "audio/flac":
		return ".flac"
	case "audio/wav", "audio/x-wav", "audio/wave":
		return ".wav"
	case "audio/mp4", "audio/aac":
		return ".m4a"
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	}

	if exts, _ := mime.ExtensionsByType(mimeType); len(exts) > 0 {
		return sanitizeExtension(exts[0])
	}

	return ""
}
