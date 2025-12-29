package events

import (
	"context"
	"fmt"
	"sync"
)

// EventHandler is a function that handles an event.
type EventHandler func(ctx context.Context, event Event) error

// EventBus defines the interface for an event bus.
type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(handler EventHandler)
}

// SyncEventBus is a simple synchronous event bus implementation.
type SyncEventBus struct {
	handlers []EventHandler
	mu       sync.RWMutex
}

// NewSyncEventBus creates a new SyncEventBus.
func NewSyncEventBus() *SyncEventBus {
	return &SyncEventBus{
		handlers: make([]EventHandler, 0),
	}
}

// Subscribe registers a handler for all events.
func (b *SyncEventBus) Subscribe(handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, handler)
}

// Publish publishes an event to all subscribers.
// In this synchronous implementation, it blocks until all handlers return.
func (b *SyncEventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers := append([]EventHandler(nil), b.handlers...)
	b.mu.RUnlock()

	var errs []error
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("event bus: %d handlers failed for event %s: %v", len(errs), event.Name, errs)
	}

	return nil
}
