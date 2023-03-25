package actor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

type EventHandlerFunc func(events []any)

type eventStream struct {
	capacity int
	events   chan any
	subs     map[uuid.UUID]EventHandlerFunc
	mu       sync.RWMutex

	exit chan struct{}
	done atomic.Bool
}

func newEventStream(capacity int) *eventStream {
	return &eventStream{
		capacity: capacity,
		events:   nil,
		subs:     make(map[uuid.UUID]EventHandlerFunc),
		mu:       sync.RWMutex{},
		exit:     make(chan struct{}),
	}
}

func (e *eventStream) Subscribe(handler EventHandlerFunc) uuid.UUID {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.events == nil || handler == nil {
		return uuid.UUID{}
	}

	id := uuid.New()
	e.subs[id] = handler

	return id
}

func (e *eventStream) Unsubscribe(id uuid.UUID) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.subs, id)
}

func (e *eventStream) Start() {
	if e.events != nil {
		return
	}

	e.events = make(chan any, e.capacity)
	e.done.Store(false)

	go func() {
	loop:
		for {
			select {
			case <-e.exit:
				break loop
			default:
			}

			events, ok := e.consumeEvents()
			e.processEvents(events)

			if !ok {
				break loop
			}
		}

		e.done.Store(true)
	}()
}

func (e *eventStream) Stop(ctx context.Context) error {
	if e.events == nil {
		return nil
	}

	close(e.events)

	for {
		select {
		case <-ctx.Done():
			close(e.exit)
			return fmt.Errorf("context timeout")

		default:
		}

		if e.done.Load() {
			close(e.exit)
			break
		}
	}

	return nil
}

func (e *eventStream) Dispatch(event any) {
	if e.events == nil || event == nil {
		return
	}

	e.events <- event
}

func (e *eventStream) consumeEvents() ([]any, bool) {
	size := len(e.events)

	if size == 0 {
		select {
		case event := <-e.events:
			if event == nil {
				return nil, false
			}

			return []any{event}, true

		default:
			return nil, true
		}
	}

	events := make([]any, size)
	for idx := 0; idx < size; idx++ {
		event, ok := <-e.events
		if !ok {
			return events[:idx], false
		}

		events[idx] = event
	}

	return events, true
}

func (e *eventStream) processEvents(events []any) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, handler := range e.subs {
		handler(events)
	}
}
