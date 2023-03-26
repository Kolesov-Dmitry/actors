package actor

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
