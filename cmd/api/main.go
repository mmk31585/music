package main

import (
	"context"
	"music/internal/common/middleware"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"syscall"
	"time"

	"music/internal/app"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	application, err := app.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	log := application.Logger

	if application.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.GinZapLogger(log))
	r.Use(middleware.GinZapRecovery(log))
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Device-Id",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	application.RegisterRoutes(r)
	registerMediaRoutes(r, "./media")

	// IMPORTANT:
	// If your app routes are currently registered inside Bootstrap using chi,
	// you need to migrate them to Gin too.
	//
	// Example:
	// application.RegisterRoutes(r)
	//
	// For now, we attach Gin router to the existing HTTP server.
	application.HTTPServer.Handler = r

	go func() {
		log.Info("starting server",
			zap.String("addr", application.HTTPServer.Addr),
			zap.String("env", application.Config.App.Env),
		)

		if err := application.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), application.Config.App.ShutdownTimeout)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	}

	log.Info("server stopped")
	time.Sleep(100 * time.Millisecond)
}
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

		// Prevent path traversal
		if strings.HasPrefix(cleanPath, "..") {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "forbidden",
			})
			return
		}

		fullPath := filepath.Join(mediaRoot, cleanPath)

		info, err := os.Stat(fullPath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "file not found",
			})
			return
		}

		if info.IsDir() {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "file not found",
			})
			return
		}

		c.File(fullPath)
	})
}
