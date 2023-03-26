package actor

import (
	"context"
	"sync"
)

type Environ struct {
	engine  *Engine
	actorId *ID

	childrenLock sync.Mutex
	children     map[*ID]Actor
}

func newEnviron(e *Engine, id *ID) *Environ {
	return &Environ{
		engine:       e,
		actorId:      id,
		childrenLock: sync.Mutex{},
		children:     make(map[*ID]Actor),
	}
}

func (e *Environ) SpawnChild(receiver Receiver, name string) *ID {
	if receiver == nil || name == "" {
		return nil
	}

	actor := newActor(e.engine, receiver, name)
	e.engine.dispatchActor(actor)

	e.childrenLock.Lock()
	defer e.childrenLock.Unlock()

	e.children[actor.id] = actor

	return actor.id
}

func (e *Environ) shutdown(ctx context.Context) error {
	e.childrenLock.Lock()
	defer e.childrenLock.Unlock()

	for id, child := range e.children {
		if err := child.Shutdown(ctx); err != nil {
			return err
		}

		delete(e.children, id)
	}

	return nil
}
