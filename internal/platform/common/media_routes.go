package common

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerMediaRoutes(r *gin.Engine, mediaRoot string) {
	r.GET("/media/*filepath", func(c *gin.Context) {
		requestedPath := c.Param("filepath")

		requestedPath = strings.TrimPrefix(requestedPath, "/")
		if requestedPath == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "file not found",
			})
			return
		}

		cleanPath := filepath.Clean(requestedPath)
		if strings.HasPrefix(cleanPath, "..") {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "forbidden",
			})
			return
		}

		fullPath := filepath.Join(mediaRoot, cleanPath)

		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "file not found",
			})
			return
		}

		c.File(fullPath)
	})
}
