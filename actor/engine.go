package actor

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrActorDoesNotExists = errors.New("actor with provided ID does not exists")
)

type Parcel struct {
	Sender   ID
	Response *Response
	Message  any
}

type Engine struct {
	engineOptions

	disp *dispatcher
}

func NewEngine(opts ...EngineOption) *Engine {
	options := engineOpts(opts...)

	engine := &Engine{
		engineOptions: *options,
		disp:          newDispatcher(),
	}

	return engine
}

func (e *Engine) Spawn(receiver Receiver, name string, tags ...string) (ID, error) {
	if receiver == nil {
		return ID{}, fmt.Errorf("Receiver was not provided")
	}

	if name == "" {
		return ID{}, fmt.Errorf("actor name was not provided")
	}

	actor := newActor(e, &actorConfig{
		receiver:   receiver,
		name:       name,
		tags:       tags,
		parent:     ID{},
		middleware: e.middleware,
	})

	e.dispatchActor(actor)

	return actor.id, nil
}

func (e *Engine) Drop(ctx context.Context, id ID) error {
	e.Send(id, AboutToStopEvent{})

	return e.disp.Remove(ctx, id)
}

func (e *Engine) Send(id ID, msg any) bool {
	parcel := &Parcel{
		Message: msg,
	}

	return e.send(id, parcel)
}

func (e *Engine) SendWithResponse(id ID, msg any) *Response {
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

func (e *Engine) Broadcast(group BroadcastGroup, msg any) bool {
	parcel := &Parcel{
		Message: msg,
	}

	result := true
	for _, id := range group.ids {
		if !e.send(id, parcel) {
			result = false
		}
	}

	return result
}

func (e *Engine) Shutdown(ctx context.Context) error {
	if err := e.disp.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (e *Engine) send(id ID, parcel *Parcel) bool {
	actor := e.disp.ActorById(id)
	if actor == nil {
		return false
	}

	actor.Invoke(parcel)
	return true
}

func (e *Engine) dispatchActor(a Actor) error {
	if err := e.disp.Add(a); err != nil {
		return err
	}

	e.Send(a.ID(), StartedEvent{})

	return nil
}
