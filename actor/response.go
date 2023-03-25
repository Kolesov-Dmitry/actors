package actor

import (
	"context"

	"github.com/google/uuid"
)

type Response struct {
	engine *Engine
	id     *ID

	value chan any
}

func newResponse(e *Engine) *Response {
	id := newID(e.address, "response", uuid.NewString())
	return &Response{
		id:     id,
		engine: e,
		value:  make(chan any),
	}
}

func (r *Response) ID() *ID {
	return r.id
}

func (r *Response) Invoke(_ *ID, msg any) {
	r.value <- msg
}

func (r *Response) Result(ctx context.Context) (any, error) {
	defer func() {
		r.engine.Drop(context.Background(), r.id)
	}()

	select {
	case result := <-r.value:
		return result, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *Response) Shutdown(_ context.Context) error {
	return nil
}
