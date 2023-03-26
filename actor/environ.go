package actor

import "context"

type Environ struct {
	engine *Engine
	actor  Actor
}

func newEnviron(e *Engine, a Actor) *Environ {
	return &Environ{
		engine: e,
		actor:  a,
	}
}

func (e *Environ) SpawnChild(receiver Receiver, name string) *ID {
	if receiver == nil || name == "" {
		return nil
	}

	child := newActor(e.engine, receiver, name)
	e.actor.AddChild(child)

	return child.id
}

func (e *Environ) DropChild(ctx context.Context, id *ID) error {
	return e.actor.DropChild(ctx, id)
}

func (e *Environ) Send(id *ID, msg any) bool {
	parcel := &Parcel{
		Sender:  e.actor.ID(),
		Message: msg,
	}

	return e.engine.send(id, parcel)
}
