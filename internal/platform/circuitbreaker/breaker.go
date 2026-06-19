package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var ErrBreakerOpen = errors.New("circuit breaker is open")

type Breaker struct {
	mu           sync.RWMutex
	state        State
	failureCount int
	lastFailure  time.Time

	threshold     int
	recoveryTime  time.Duration
	halfOpenMax   int
	halfOpenCount int
}

type Config struct {
	Threshold    int
	RecoveryTime time.Duration
	HalfOpenMax  int
}

func New(cfg Config) *Breaker {
	if cfg.Threshold <= 0 {
		cfg.Threshold = 3
	}
	if cfg.RecoveryTime <= 0 {
		cfg.RecoveryTime = 30 * time.Second
	}
	if cfg.HalfOpenMax <= 0 {
		cfg.HalfOpenMax = 1
	}
	return &Breaker{
		state:        StateClosed,
		threshold:    cfg.Threshold,
		recoveryTime: cfg.RecoveryTime,
		halfOpenMax:  cfg.HalfOpenMax,
	}
}

func (b *Breaker) State() State {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.state
}

func (b *Breaker) Call(ctx context.Context, fn func(context.Context) error) error {
	state := b.beforeCall()
	switch state {
	case StateOpen:
		return ErrBreakerOpen
	case StateHalfOpen:
		if !b.tryAcquireHalfOpen() {
			return ErrBreakerOpen
		}
	}

	err := fn(ctx)
	b.afterCall(err)
	return err
}

func (b *Breaker) beforeCall() State {
	b.mu.RLock()
	state := b.state
	b.mu.RUnlock()

	if state != StateOpen {
		return state
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateOpen && time.Since(b.lastFailure) >= b.recoveryTime {
		b.state = StateHalfOpen
		b.halfOpenCount = 0
		return StateHalfOpen
	}
	return StateOpen
}

func (b *Breaker) tryAcquireHalfOpen() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.halfOpenCount < b.halfOpenMax {
		b.halfOpenCount++
		return true
	}
	return false
}

func (b *Breaker) afterCall(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err == nil {
		b.state = StateClosed
		b.failureCount = 0
		b.halfOpenCount = 0
		return
	}

	b.failureCount++
	b.lastFailure = time.Now()

	if b.state == StateHalfOpen || b.failureCount >= b.threshold {
		b.state = StateOpen
	}
}

func (b *Breaker) FailureCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.failureCount
}
