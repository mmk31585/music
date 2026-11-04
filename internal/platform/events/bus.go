package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
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
	outboxStore := b.outboxStore
	b.mu.RUnlock()

	// ── Step 1: Persist to outbox (durable storage before dispatch) ──
	var outboxID string
	if outboxStore != nil {
		outboxEvent := NewOutboxEvent(event)
		if err := outboxStore.Insert(ctx, outboxEvent); err != nil {
			b.logger.Error("failed to persist event to outbox",
				zap.String("event", event.EventName()),
				zap.Error(err),
			)
			return err
		}
		outboxID = outboxEvent.ID
		b.logger.Debug("event persisted to outbox",
			zap.String("event", event.EventName()),
			zap.String("outbox_id", outboxEvent.ID),
		)
	}

	// ── Step 2: Dispatch to handlers ──
	b.dispatch(ctx, handlers, event)

	// ── Step 3: Mark the row processed ──
	// The event has been dispatched (handlers may still be running async). We
	// mark it processed so the relay / DB-polling workers do not re-deliver the
	// same side effect. Rows left pending (e.g. a crash between insert and this
	// update) are recovered by the relay.
	if outboxStore != nil && outboxID != "" {
		if err := outboxStore.MarkProcessed(ctx, outboxID); err != nil {
			b.logger.Error("failed to mark outbox event processed",
				zap.String("event", event.EventName()),
				zap.String("outbox_id", outboxID),
				zap.Error(err),
			)
		}
	}

	return nil
}

// dispatch invokes handlers according to the bus's sync/async mode. It never
// persists the event; callers are responsible for outbox bookkeeping.
func (b *Bus) dispatch(ctx context.Context, handlers []Handler, event Event) {
	b.mu.RLock()
	async := b.async
	hasOutbox := b.outboxStore != nil
	b.mu.RUnlock()

	if async || hasOutbox {
		// Always dispatch async when outbox is enabled (handlers must not block
		// the publisher). Detach from the caller's context: the HTTP request
		// context is cancelled as soon as the response is written, which would
		// abort in-flight handlers (e.g. creating a welcome notification).
		dctx := context.WithoutCancel(ctx)
		for _, handler := range handlers {
			h := handler
			go func() {
				if err := h(dctx, event); err != nil {
					b.logger.Error("event handler failed",
						zap.String("event", event.EventName()),
						zap.Error(err),
					)
				}
			}()
		}
		return
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			b.logger.Error("event handler failed",
				zap.String("event", event.EventName()),
				zap.Error(err),
			)
		}
	}
}

// Dispatch delivers an already-persisted event to in-memory subscribers
// WITHOUT writing it back to the outbox. Used by the outbox relay for
// re-delivery of events that were persisted but never dispatched.
func (b *Bus) Dispatch(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.subscribers[event.EventName()]...)
	b.mu.RUnlock()

	b.dispatch(ctx, handlers, event)
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

// EventDispatcher is the subset of Bus that the relay needs to re-deliver events.
// It must dispatch WITHOUT persisting, otherwise every relay pass would insert a
// new outbox row (unbounded growth) and re-deliver the event forever.
type EventDispatcher interface {
	Dispatch(ctx context.Context, event Event)
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
		// Reconstruct event from persisted payload and re-dispatch.
		event, err := eventFromOutbox(oe)
		if err != nil {
			r.logger.Error("outbox relay: failed to reconstruct event",
				zap.String("outbox_id", oe.ID),
				zap.Error(err),
			)
			_ = r.store.MarkFailed(ctx, oe.ID, "failed to reconstruct: "+err.Error())
			continue
		}

		// Dispatch only — never persist again (that would insert a new row and
		// loop forever).
		r.bus.Dispatch(ctx, event)

		if err := r.store.MarkProcessed(ctx, oe.ID); err != nil {
			r.logger.Error("outbox relay: failed to mark processed",
				zap.String("outbox_id", oe.ID),
				zap.Error(err),
			)
		}
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

// eventFromOutbox reconstructs a typed Event from its persisted outbox payload.
// Subscribers type-assert the concrete event type, so a generic ReplayedEvent
// would be silently ignored and recovery would be a no-op.
func eventFromOutbox(oe OutboxEvent) (Event, error) {
	var p struct {
		OccurredAt     time.Time  `json:"occurred_at"`
		UserID         *uuid.UUID `json:"user_id"`
		TrackID        *uuid.UUID `json:"track_id"`
		ArtistID       *uuid.UUID `json:"artist_id"`
		AlbumID        *uuid.UUID `json:"album_id"`
		PlaylistID     *uuid.UUID `json:"playlist_id"`
		SubscriptionID *uuid.UUID `json:"subscription_id"`
		Duration       int        `json:"duration"`
		Completed      bool       `json:"completed"`
		Source         string     `json:"source"`
		Name           string     `json:"name"`
		Email          string     `json:"email"`
		IsPublic       bool       `json:"is_public"`
		Plan           string     `json:"plan"`
		Amount         int64      `json:"amount"`
		Currency       string     `json:"currency"`
	}
	if err := json.Unmarshal(oe.Payload, &p); err != nil {
		return nil, fmt.Errorf("unmarshal outbox payload: %w", err)
	}

	at := oe.CreatedAt
	if !p.OccurredAt.IsZero() {
		at = p.OccurredAt
	}
	base := BaseEvent{At: at}
	uuidOrZero := func(v *uuid.UUID) uuid.UUID {
		if v == nil {
			return uuid.Nil
		}
		return *v
	}

	switch oe.EventType {
	case EventTrackPlayed:
		return TrackPlayedEvent{
			BaseEvent: base,
			UserID:    uuidOrZero(p.UserID),
			TrackID:   uuidOrZero(p.TrackID),
			ArtistID:  uuidOrZero(p.ArtistID),
			AlbumID:   uuidOrZero(p.AlbumID),
			Duration:  p.Duration,
			Completed: p.Completed,
			Source:    p.Source,
		}, nil
	case EventPlaylistCreated:
		return PlaylistCreatedEvent{
			BaseEvent:  base,
			UserID:     uuidOrZero(p.UserID),
			PlaylistID: uuidOrZero(p.PlaylistID),
			Name:       p.Name,
			IsPublic:   p.IsPublic,
		}, nil
	case EventSubscriptionPurchased:
		return SubscriptionPurchasedEvent{
			BaseEvent:      base,
			UserID:         uuidOrZero(p.UserID),
			SubscriptionID: uuidOrZero(p.SubscriptionID),
			Plan:           p.Plan,
			Amount:         p.Amount,
			Currency:       p.Currency,
		}, nil
	case EventUserRegistered:
		return UserRegisteredEvent{
			BaseEvent: base,
			UserID:    uuidOrZero(p.UserID),
			Email:     p.Email,
			Name:      p.Name,
		}, nil
	}

	return nil, fmt.Errorf("unknown outbox event type %q", oe.EventType)
}
