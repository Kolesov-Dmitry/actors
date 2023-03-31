package actor

import (
	"context"
	"fmt"
	"sync"
)

type Actor interface {
	ID() ID
	Invoke(p *Parcel)
	AddChild(actor Actor)
	DropChild(ctx context.Context, id ID) error
	Shutdown(ctx context.Context) error
}

type dispatcher struct {
	mu     sync.RWMutex
	actors map[string]Actor
}

func newDispatcher() *dispatcher {
	return &dispatcher{
		actors: make(map[string]Actor),
	}
}

func (d *dispatcher) Add(a Actor) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	id := a.ID().String()
	if _, ok := d.actors[id]; ok {
		return fmt.Errorf("actor with '%s' ID already exists", id)
	}

	d.actors[id] = a

	return nil
}

func (d *dispatcher) Remove(ctx context.Context, id ID) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	actor, ok := d.actors[id.String()]
	if !ok {
		return fmt.Errorf("failed to drop actor '%s': %w", id.String(), ErrActorDoesNotExists)
	}

	delete(d.actors, id.String())

	return actor.Shutdown(ctx)
}

func (d *dispatcher) ActorById(id ID) Actor {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if actor, ok := d.actors[id.String()]; ok {
		return actor
	}

	return nil
}

func (d *dispatcher) Shutdown(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, actor := range d.actors {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context timeout")
		default:
		}

		if err := actor.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}
