package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"music/internal/common/response"
)

const maintenanceRedisKey = "app:maintenance"

type Handler struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewHandler(db *sqlx.DB, rdb *redis.Client) *Handler {
	return &Handler{db: db, rdb: rdb}
}

type StatsResponse struct {
	TrackCount  int `json:"track_count"`
	ArtistCount int `json:"artist_count"`
	AlbumCount  int `json:"album_count"`
	GenreCount  int `json:"genre_count"`
}

func (h *Handler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	var trackCount, artistCount, albumCount, genreCount int

	if err := h.db.GetContext(ctx, &trackCount, "SELECT COUNT(*) FROM tracks"); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.db.GetContext(ctx, &artistCount, "SELECT COUNT(*) FROM artists"); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.db.GetContext(ctx, &albumCount, "SELECT COUNT(*) FROM albums"); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.db.GetContext(ctx, &genreCount, "SELECT COUNT(*) FROM genres"); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "ok", StatsResponse{
		TrackCount:  trackCount,
		ArtistCount: artistCount,
		AlbumCount:  albumCount,
		GenreCount:  genreCount,
	})
}

// GetMaintenanceStatus returns the current maintenance mode state (public).
func (h *Handler) GetMaintenanceStatus(c *gin.Context) {
	ctx := c.Request.Context()
	enabled := false

	if h.rdb != nil {
		val, err := h.rdb.Exists(ctx, maintenanceRedisKey).Result()
		if err == nil && val > 0 {
			enabled = true
		}
	}

	response.Success(c, http.StatusOK, "ok", gin.H{
		"maintenance": enabled,
	})
}

type MaintenanceToggleRequest struct {
	Enabled bool `json:"enabled"`
}

// ToggleMaintenance enables or disables maintenance mode (admin only).
func (h *Handler) ToggleMaintenance(c *gin.Context) {
	ctx := c.Request.Context()

	var req MaintenanceToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if h.rdb != nil {
		if req.Enabled {
			err := h.rdb.Set(ctx, maintenanceRedisKey, "1", 0).Err()
			if err != nil {
				response.Error(c, err)
				return
			}
		} else {
			err := h.rdb.Del(ctx, maintenanceRedisKey).Err()
			if err != nil {
				response.Error(c, err)
				return
			}
		}
	}

	response.Success(c, http.StatusOK, "maintenance mode updated", gin.H{
		"maintenance": req.Enabled,
	})
}
