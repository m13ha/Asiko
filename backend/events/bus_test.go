package events

import (
	"context"
	"errors"
	"testing"
)

func TestSyncEventBusPublishesToSubscribers(t *testing.T) {
	bus := NewSyncEventBus()
	var received []string

	bus.Subscribe(func(ctx context.Context, event Event) error {
		value, _ := event.Data.(string)
		received = append(received, "a:"+value)
		return nil
	})
	bus.Subscribe(func(ctx context.Context, event Event) error {
		value, _ := event.Data.(string)
		received = append(received, "b:"+value)
		return nil
	})

	if err := bus.Publish(context.Background(), Event{Name: "topic.test", Data: "hello"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(received) != 2 {
		t.Fatalf("expected 2 handlers to receive payload, got %d", len(received))
	}
}

func TestSyncEventBusReturnsErrors(t *testing.T) {
	bus := NewSyncEventBus()
	bus.Subscribe(func(ctx context.Context, event Event) error {
		return errors.New("boom")
	})

	if err := bus.Publish(context.Background(), Event{Name: "topic.fail", Data: "payload"}); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
