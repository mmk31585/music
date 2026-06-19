package dashboard

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"music/internal/common/response"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

type StatsResponse struct {
	TrackCount  int `json:"track_count"`
	ArtistCount int `json:"artist_count"`
	AlbumCount  int `json:"album_count"`
	GenreCount  int `json:"genre_count"`
}

func (h *Handler) GetStats(c *gin.Context) {
	ctx := context.Background()

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
