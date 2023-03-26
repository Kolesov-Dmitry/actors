package actor

import (
	"context"
	"log"
)

type invokeMessageEvent struct {
	parcel *Parcel
}

type Receiver interface {
	Receive(p *Parcel)
}

type actor struct {
	engine   *Engine
	id       *ID
	receiver Receiver
	events   *eventStream
}

func newActor(e *Engine, recv Receiver, name string, tags ...string) *actor {
	id := newID(e.address, name, tags...)
	p := &actor{
		id:       id,
		engine:   e,
		receiver: recv,
		events:   newEventStream(e.capacity),
	}

	p.events.Start()
	p.events.Subscribe(p.handleEvents)

	return p
}

func (a *actor) ID() *ID {
	return a.id
}

func (a *actor) Invoke(p *Parcel) {
	a.events.Dispatch(
		invokeMessageEvent{parcel: p},
	)
}

func (a *actor) handleEvents(events []any) {
	for _, event := range events {
		switch msg := event.(type) {
		case invokeMessageEvent:
			a.receiver.Receive(msg.parcel)

		default:
			log.Println("receive unsupported event")
		}
	}
}

func (a *actor) Shutdown(ctx context.Context) error {
	return a.events.Stop(ctx)
}
