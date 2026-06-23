package social

import (
	"go.uber.org/zap"
)

type Handler struct {
	service      *Service
	queueEngine  *QueueEngine
	stageManager *StageManager
	logger       *zap.Logger
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service, logger: zap.L()}
}

func NewHandlerWithEngine(service *Service, queueEngine *QueueEngine) *Handler {
	return &Handler{service: service, queueEngine: queueEngine, logger: zap.L()}
}

func NewHandlerFull(service *Service, queueEngine *QueueEngine, stageManager *StageManager) *Handler {
	return &Handler{service: service, queueEngine: queueEngine, stageManager: stageManager, logger: zap.L()}
}

func NewHandlerWithLogger(service *Service, queueEngine *QueueEngine, stageManager *StageManager, logger *zap.Logger) *Handler {
	return &Handler{service: service, queueEngine: queueEngine, stageManager: stageManager, logger: logger}
}
