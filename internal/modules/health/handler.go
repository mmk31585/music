package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewHandler(db *pgxpool.Pool, redis *redis.Client) *Handler {
	return &Handler{
		db:    db,
		redis: redis,
	}
}

func (h *Handler) Health(c *gin.Context) {
	response.GinSuccess(c, http.StatusOK, "service is healthy", map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) Live(c *gin.Context) {
	response.GinSuccess(c, http.StatusOK, "service is live", map[string]interface{}{
		"status": "live",
	})
}

func (h *Handler) Ready(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.db.Ping(ctx); err != nil {
		response.GinError(c, apperrors.New(
			http.StatusServiceUnavailable,
			"SERVICE_UNAVAILABLE",
			"database not ready",
			err.Error(),
		))
		return
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		response.GinError(c, apperrors.New(
			http.StatusServiceUnavailable,
			"SERVICE_UNAVAILABLE",
			"redis not ready",
			err.Error(),
		))
		return
	}

	response.GinSuccess(c, http.StatusOK, "service is ready", map[string]interface{}{
		"status": "ready",
	})
}
