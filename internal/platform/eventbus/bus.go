package eventbus

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
	logger      *zap.Logger
	async       bool
}

func NewBus(logger *zap.Logger) *Bus {
	if logger == nil {
		logger = zap.NewNop() // or zap.L() for a global logger
	}

	return &Bus{
		subscribers: make(map[string][]Handler),
		logger:      logger,
		async:       false,
	}
}

// SetAsync controls whether handlers run synchronously or in goroutines.
// For early systems, sync is often better because it's deterministic.
// Async can be enabled later when desired.
func (b *Bus) SetAsync(async bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.async = async
}

func (b *Bus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.subscribers[event.EventName()]...)
	async := b.async
	b.mu.RUnlock()

	if len(handlers) == 0 {
		return nil
	}

	if async {
		for _, handler := range handlers {
			h := handler
			go func() {
				if err := h(ctx, event); err != nil {
					b.logger.Error("event handler failed",
						zap.String("event", event.EventName()),
						zap.Error(err),
					)
				}
			}()
		}
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			b.logger.Error("event handler failed",
				zap.String("event", event.EventName()),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (b *Bus) Subscribe(eventName string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[eventName] = append(b.subscribers[eventName], handler)
}
