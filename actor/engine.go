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

type Parcel struct {
	Sender   *ID
	Response *Response
	Message  any
}

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

func (e *Engine) Send(id *ID, msg any) bool {
	parcel := &Parcel{
		Message: msg,
	}

	return e.send(id, parcel)
}

func (e *Engine) SendWithResponse(id *ID, msg any) *Response {
	response := newResponse()
	parcel := &Parcel{
		Response: response,
		Message:  msg,
	}
	if !e.send(id, parcel) {
		return nil
	}

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

func (e *Engine) send(id *ID, parcel *Parcel) bool {
	actor := e.disp.ActorById(id)
	if actor == nil {
		return false
	}

	actor.Invoke(parcel)
	return true
}

func (e *Engine) dispatchActor(a Actor) {
	e.disp.Add(a)
	e.eventStream.Dispatch(StartedEvent{ID: a.ID()})
}
