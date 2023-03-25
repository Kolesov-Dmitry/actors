package actor

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrActorDoesNotExists = errors.New("actor with provided ID does not exists")
)

const (
	hooksStreamSize = 100
)

type Engine struct {
	engineOptions

	disp        *dispatcher
	eventStream *eventStream
}

func NewEngine(opts ...EngineOption) *Engine {
	options := engineOpts(opts...)

	engine := &Engine{
		engineOptions: *options,
		disp:          newDispatcher(),
		eventStream:   newEventStream(hooksStreamSize),
	}

	engine.eventStream.Start()

	return engine
}

func (e *Engine) Subscribe(handler EventHandlerFunc) uuid.UUID {
	return e.eventStream.Subscribe(handler)
}

func (e *Engine) Spawn(receiver Receiver, name string) *ID {
	if receiver == nil || name == "" {
		return nil
	}

	actor := newActor(e, receiver, name)

	e.dispatchActor(actor)

	return actor.id
}

func (e *Engine) Drop(ctx context.Context, id *ID) error {
	err := e.disp.Remove(ctx, id)
	if !errors.Is(err, ErrActorDoesNotExists) {
		e.eventStream.Dispatch(DroppedEvent{ID: id})
	}

	return err
}

func (e *Engine) Send(id *ID, msg any) {
	actor := e.disp.ActorById(id)
	if actor != nil {
		actor.Invoke(nil, msg)
	}
}

func (e *Engine) SendWithResponse(id *ID, msg any) *Response {
	actor := e.disp.ActorById(id)
	if actor == nil {
		return nil
	}

	response := newResponse(e)
	e.dispatchActor(response)

	actor.Invoke(response.id, msg)

	return response
}

func (e *Engine) Shutdown(ctx context.Context) error {
	if err := e.disp.Shutdown(ctx); err != nil {
		return err
	}

	if err := e.eventStream.Stop(ctx); err != nil {
		return err
	}

	return nil
}

func (e *Engine) dispatchActor(a Actor) {
	e.disp.Add(a)
	e.eventStream.Dispatch(StartedEvent{ID: a.ID()})
}
