package actor

const (
	capacityDefault = 1
)

type engineOptions struct {
	capacity   int
	middleware []Middleware
}

type EngineOption func(opts *engineOptions)

func engineOpts(opts ...EngineOption) *engineOptions {
	options := &engineOptions{
		capacity: capacityDefault,
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithCapacity(capacity int) EngineOption {
	return func(opts *engineOptions) {
		opts.capacity = capacity
	}
}

func WithMiddleware(ms ...Middleware) EngineOption {
	return func(opts *engineOptions) {
		opts.middleware = ms
	}
}
