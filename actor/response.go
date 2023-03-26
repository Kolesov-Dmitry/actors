package actor

import (
	"context"
)

type Response struct {
	value chan any
}

func newResponse() *Response {
	return &Response{
		value: make(chan any),
	}
}

func (r *Response) SetValue(value any) {
	r.value <- value
}

func (r *Response) Result(ctx context.Context) (any, error) {
	select {
	case result := <-r.value:
		return result, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
