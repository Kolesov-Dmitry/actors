package actor

const (
	localhostAddress = "localhost"
	capacityDefault  = 1
)

type engineOptions struct {
	address  string
	capacity int
}

type EngineOption func(opts *engineOptions)

func engineOpts(opts ...EngineOption) *engineOptions {
	options := &engineOptions{
		address:  localhostAddress,
		capacity: capacityDefault,
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithAddress(addr string) EngineOption {
	return func(opts *engineOptions) {
		opts.address = addr
	}
}

func WithCapacity(capacity int) EngineOption {
	return func(opts *engineOptions) {
		opts.capacity = capacity
	}
}
