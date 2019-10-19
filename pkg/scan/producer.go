package scan

import "context"

type Producer interface {
	Produce(ctx context.Context) <-chan Target
}

type ReProducer interface {
	Reproduce(ctx context.Context) func(r Result) <-chan Target
}
