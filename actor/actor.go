package actor

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type invokeMessageEvent struct {
	parcel *Parcel
}

type Receiver interface {
	Receive(env *Environ, p *Parcel)
}

type actor struct {
	engine   *Engine
	id       *ID
	receiver Receiver
	events   *eventStream
	environ  *Environ

	childrenLock sync.Mutex
	children     map[*ID]Actor
}

func newActor(engine *Engine, parent *ID, recv Receiver, name string, tags ...string) *actor {
	id := newID(engine.address, name, tags...)
	a := &actor{
		id:       id,
		engine:   engine,
		receiver: recv,
		events:   newEventStream(engine.capacity),

		childrenLock: sync.Mutex{},
		children:     make(map[*ID]Actor),
	}
	a.environ = newEnviron(engine, parent, a)

	a.events.Start()
	a.events.Subscribe(a.handleEvents)

	return a
}

func (a *actor) ID() *ID {
	return a.id
}

func (a *actor) Invoke(p *Parcel) {
	a.events.Dispatch(
		invokeMessageEvent{parcel: p},
	)
}

func (a *actor) AddChild(actor Actor) {
	a.childrenLock.Lock()
	defer a.childrenLock.Unlock()

	a.children[actor.ID()] = actor
	a.engine.dispatchActor(actor)
}

func (a *actor) DropChild(ctx context.Context, id *ID) error {
	a.childrenLock.Lock()
	defer a.childrenLock.Unlock()

	if _, ok := a.children[id]; !ok {
		return fmt.Errorf("child with '%s' was not found", id.String())
	}

	delete(a.children, id)

	return a.engine.Drop(ctx, id)
}

func (a *actor) handleEvents(events []any) {
	for _, event := range events {
		switch msg := event.(type) {
		case invokeMessageEvent:
			a.receiver.Receive(a.environ, msg.parcel)

		default:
			log.Println("receive unsupported event")
		}
	}
}

func (a *actor) Shutdown(ctx context.Context) error {
	a.childrenLock.Lock()
	defer a.childrenLock.Unlock()

	for id := range a.children {
		delete(a.children, id)
	}

	return a.events.Stop(ctx)
}
