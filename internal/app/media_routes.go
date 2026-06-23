package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterMediaRoutes(r *gin.Engine, mediaRoot string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestedPath := req.URL.Path[len("/uploads/"):]

		if requestedPath == "" {
			http.Error(w, `{"success":false,"message":"file not found"}`, http.StatusNotFound)
			return
		}

		cleanPath := filepath.Clean(requestedPath)
		if strings.HasPrefix(cleanPath, "..") {
			http.Error(w, `{"success":false,"message":"forbidden"}`, http.StatusForbidden)
			return
		}

		fullPath := filepath.Join(mediaRoot, cleanPath)

		file, err := os.Open(fullPath)
		if err != nil {
			http.Error(w, `{"success":false,"message":"file not found"}`, http.StatusNotFound)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil || stat.IsDir() {
			http.Error(w, `{"success":false,"message":"file not found"}`, http.StatusNotFound)
			return
		}

		// Compute ETag from file info
		etag := computeETag(fullPath, stat)
		w.Header().Set("ETag", etag)
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")

		// Support If-None-Match
		if match := req.Header.Get("If-None-Match"); match != "" {
			if match == etag || match == `"`+etag+`"` {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		// Support If-Modified-Since
		if modSince := req.Header.Get("If-Modified-Since"); modSince != "" {
			if t, err := time.Parse(http.TimeFormat, modSince); err == nil {
				if !stat.ModTime().After(t) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}

		// Serve with full range support
		http.ServeContent(w, req, stat.Name(), stat.ModTime(), file)
	})

	r.GET("/uploads/*filepath", gin.WrapH(handler))
}

func computeETag(path string, stat os.FileInfo) string {
	h := sha256.New()
	h.Write([]byte(path))
	h.Write([]byte(strconv.FormatInt(stat.Size(), 10)))
	h.Write([]byte(stat.ModTime().UTC().Format(time.RFC3339)))
	return fmt.Sprintf(`"%s"`, hex.EncodeToString(h.Sum(nil)))
}
