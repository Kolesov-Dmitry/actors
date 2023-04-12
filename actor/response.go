package actor

import (
	"context"
)

type Response struct {
	value chan any
}

func newResponse() *Response {
	return &Response{
		value: make(chan any, 1),
	}
}

func (r *Response) SetValue(value any) {
	select {
	case r.value <- value:
	default:
	}
}

func (r *Response) Result(ctx context.Context) (any, error) {
	select {
	case result := <-r.value:
		close(r.value)
		return result, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
