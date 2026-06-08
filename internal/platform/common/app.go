package app

import (
	"music/internal/platform/events"
	"net/http"

	"music/internal/common/validator"
	"music/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type App struct {
	Config     *config.Config
	Logger     *zap.Logger
	DB         *pgxpool.Pool
	Redis      *redis.Client
	Validator  *validator.Validator
	Router     *gin.Engine
	Events     *events.Bus
	HTTPServer *http.Server
}
