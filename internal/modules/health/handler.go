package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"
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

// Health godoc
// @Summary Health check
// @Description Returns basic service health status.
// @Tags health
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	response.Success(c, http.StatusOK, "service is healthy", gin.H{
		"status": "ok",
	})
}

// Live godoc
// @Summary Liveness check
// @Description Returns service liveness status.
// @Tags health
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /health/live [get]
func (h *Handler) Live(c *gin.Context) {
	response.Success(c, http.StatusOK, "service is live", gin.H{
		"status": "live",
	})
}

// Ready godoc
// @Summary Readiness check
// @Description Returns service readiness status after checking dependencies such as database and Redis.
// @Tags health
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /health/ready [get]
func (h *Handler) Ready(c *gin.Context) {
	ctx := c.Request.Context()

	if h.db == nil {
		response.Error(c, apperrors.New(
			http.StatusServiceUnavailable,
			"SERVICE_UNAVAILABLE",
			"database not configured",
			"database connection is nil",
		))
		return
	}

	if err := h.db.Ping(ctx); err != nil {
		response.Error(c, apperrors.New(
			http.StatusServiceUnavailable,
			"SERVICE_UNAVAILABLE",
			"database not ready",
			err.Error(),
		))
		return
	}

	if h.redis == nil {
		response.Error(c, apperrors.New(
			http.StatusServiceUnavailable,
			"SERVICE_UNAVAILABLE",
			"redis not configured",
			"redis client is nil",
		))
		return
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		response.Error(c, apperrors.New(
			http.StatusServiceUnavailable,
			"SERVICE_UNAVAILABLE",
			"redis not ready",
			err.Error(),
		))
		return
	}

	response.Success(c, http.StatusOK, "service is ready", gin.H{
		"status": "ready",
	})
}
