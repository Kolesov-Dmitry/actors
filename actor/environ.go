package actor

import (
	"context"
	"fmt"
)

type Environ struct {
	engine *Engine
	parent ID
	actor  Actor
}

func newEnviron(e *Engine, p ID, a Actor) *Environ {
	return &Environ{
		engine: e,
		parent: p,
		actor:  a,
	}
}

func (e *Environ) SpawnChild(receiver Receiver, name string) (ID, error) {
	if receiver == nil {
		return ID{}, fmt.Errorf("Receiver was not provided")
	}

	if name == "" {
		return ID{}, fmt.Errorf("actor name was not provided")
	}

	child := newActor(e.engine, &actorConfig{
		receiver:   receiver,
		name:       name,
		parent:     e.actor.ID(),
		middleware: e.engine.middleware,
	})
	e.actor.AddChild(child)

	return child.id, nil
}

func (e *Environ) DropChild(ctx context.Context, id ID) error {
	return e.actor.DropChild(ctx, id)
}

func (e *Environ) Send(id ID, msg any) bool {
	parcel := &Parcel{
		Sender:  e.actor.ID(),
		Message: msg,
	}

	return e.engine.send(id, parcel)
}

func (e *Environ) Self() ID {
	return e.actor.ID()
}

func (e *Environ) Parent() ID {
	return e.parent
}
