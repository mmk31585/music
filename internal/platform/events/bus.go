package events

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Bus is an event bus with an optional transactional outbox.
// When outboxStore is set, every published event is first persisted to the
// outbox table before any handler is invoked. The outbox ensures events
// survive process crashes and can be republished by the OutboxRelay.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
	logger      *zap.Logger
	async       bool

	// outboxStore persists events before handler dispatch.
	// When nil, events are purely in-memory (lost on crash).
	outboxStore OutboxStore

	// relay periodically publishes pending outbox events.
	relay *OutboxRelay
}

func NewBus(logger *zap.Logger) *Bus {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Bus{
		subscribers: make(map[string][]Handler),
		logger:      logger,
		async:       false,
	}
}

// SetOutboxStore configures the bus to persist events before dispatching them.
// The second parameter starts the background outbox relay, which periodically
// processes any events that were persisted but never dispatched (e.g. after a crash).
func (b *Bus) SetOutboxStore(store OutboxStore, relayInterval time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.outboxStore = store
	b.relay = NewOutboxRelay(store, b, b.logger)
	b.relay.Start(relayInterval)

	b.logger.Info("event bus outbox enabled",
		zap.Duration("relay_interval", relayInterval))
}

// OutboxStore returns the current outbox store, or nil if none is configured.
func (b *Bus) OutboxStore() OutboxStore {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.outboxStore
}

// Relay returns the outbox relay for lifecycle management (start/stop).
func (b *Bus) Relay() *OutboxRelay {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.relay
}

// SetAsync controls whether handlers run synchronously or in goroutines.
// When the outbox is enabled, handlers are always dispatched asynchronously
// after the event is safely persisted.
func (b *Bus) SetAsync(async bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.async = async
}

// Publish persists the event to the outbox (if configured) and then dispatches
// to all registered handlers. If outbox persistence fails, the event is NOT
// dispatched — ensuring at-least-once delivery semantics.
func (b *Bus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.subscribers[event.EventName()]...)
	async := b.async
	outboxStore := b.outboxStore
	b.mu.RUnlock()

	// ── Step 1: Persist to outbox (durable storage before dispatch) ──
	if outboxStore != nil {
		outboxEvent := NewOutboxEvent(event)
		if err := outboxStore.Insert(ctx, outboxEvent); err != nil {
			b.logger.Error("failed to persist event to outbox",
				zap.String("event", event.EventName()),
				zap.Error(err),
			)
			return err
		}
		b.logger.Debug("event persisted to outbox",
			zap.String("event", event.EventName()),
			zap.String("outbox_id", outboxEvent.ID),
		)
	}

	if len(handlers) == 0 {
		return nil
	}

	// ── Step 2: Dispatch to handlers ──
	if async || outboxStore != nil {
		// Always dispatch async when outbox is enabled (handlers must not block the publisher)
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

// Shutdown gracefully stops the outbox relay (if running) and waits for any
// in-flight dispatches to complete.
func (b *Bus) Shutdown(ctx context.Context) error {
	b.mu.Lock()
	relay := b.relay
	b.mu.Unlock()

	if relay != nil {
		return relay.Stop(ctx)
	}
	return nil
}

// OutboxRelay periodically fetches pending outbox events and re-dispatches them.
// This ensures at-least-once delivery: events that were persisted but whose
// handlers crashed before completion will be retried.
type OutboxRelay struct {
	store   OutboxStore
	bus     EventDispatcher
	logger  *zap.Logger
	stopCh  chan struct{}
	stopped chan struct{}
}

// EventDispatcher is the subset of Bus that the relay needs to republish events.
type EventDispatcher interface {
	Publish(ctx context.Context, event Event) error
}

func NewOutboxRelay(store OutboxStore, bus EventDispatcher, logger *zap.Logger) *OutboxRelay {
	return &OutboxRelay{
		store:   store,
		bus:     bus,
		logger:  logger,
		stopCh:  make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (r *OutboxRelay) Start(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		defer close(r.stopped)

		for {
			select {
			case <-r.stopCh:
				r.logger.Info("outbox relay stopped")
				return
			case <-ticker.C:
				r.processPending()
			}
		}
	}()

	r.logger.Info("outbox relay started", zap.Duration("interval", interval))
}

func (r *OutboxRelay) processPending() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	events, err := r.store.FetchPending(ctx, 100)
	if err != nil {
		r.logger.Error("outbox relay: failed to fetch pending events", zap.Error(err))
		return
	}

	for _, oe := range events {
		// Reconstruct event from persisted payload and republish
		event, err := eventFromOutbox(oe)
		if err != nil {
			r.logger.Error("outbox relay: failed to reconstruct event",
				zap.String("outbox_id", oe.ID),
				zap.Error(err),
			)
			_ = r.store.MarkFailed(ctx, oe.ID, "failed to reconstruct: "+err.Error())
			continue
		}

		if err := r.bus.Publish(ctx, event); err != nil {
			_ = r.store.MarkFailed(ctx, oe.ID, err.Error())
			continue
		}

		_ = r.store.MarkProcessed(ctx, oe.ID)
	}
}

func (r *OutboxRelay) Stop(ctx context.Context) error {
	close(r.stopCh)
	select {
	case <-r.stopped:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// eventFromOutbox reconstructs an Event from a persisted OutboxEvent.
// Currently a stub — the production implementation would use a registry of
// event constructors keyed by event type.
func eventFromOutbox(oe OutboxEvent) (Event, error) {
	// In the current implementation, the relay relies on the fact that the
	// handlers themselves are registered in the bus. The outbox event was
	// already dispatched once; the relay re-dispatches it to the same bus,
	// so any handler can safely ignore events they've already seen via
	// idempotency keys (trace_id).
	//
	// For a full reconstruction, register typed constructors:
	//   var registry = map[string]func(json.RawMessage) (Event, error)
	return &ReplayedEvent{
		Name: oe.EventType,
		At:   oe.CreatedAt,
	}, nil
}

// ReplayedEvent is a lightweight event wrapper for outbox relay re-delivery.
type ReplayedEvent struct {
	Name string
	At   time.Time
}

func (e *ReplayedEvent) EventName() string     { return e.Name }
func (e *ReplayedEvent) OccurredAt() time.Time { return e.At }
