# Events Contract

This package provides a synchronous, in-memory event bus with a simple event envelope.

## Event Envelope

All events use the same structure:

```go
type Event struct {
    Name string
    Data interface{}
}
```

## Publishing

Publish uses the same API regardless of event data:

```go
bus.Publish(ctx, events.Event{
    Name: events.EventBookingCreated,
    Data: events.BookingEventData{ /* ... */ },
})
```

## Subscribing

Subscribers receive **all** events and must filter by `event.Name`:

```go
bus.Subscribe(func(ctx context.Context, event events.Event) error {
    if event.Name != events.EventBookingCreated {
        return nil
    }
    // Handle event.Data
    return nil
})
```

## Conventions

- Use `events.Event*` constants for event names.
- Keep `Data` minimal but sufficient for consumers.
- Handlers should return errors to surface failures to the publisher.
