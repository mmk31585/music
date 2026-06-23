package app

import (
	"context"
	"music/internal/platform/events"
	"net/http"
	"sync"

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

	// lifecycle management
	ctx          context.Context
	cancel       context.CancelFunc
	backgroundWg sync.WaitGroup
}

// AddBackground increments the WaitGroup and returns a derived context
// that is cancelled when the app shuts down. Call Done() when the goroutine exits.
func (a *App) AddBackground() context.Context {
	a.backgroundWg.Add(1)
	return a.ctx
}

// BackgroundDone signals that a background goroutine has finished.
func (a *App) BackgroundDone() {
	a.backgroundWg.Done()
}
