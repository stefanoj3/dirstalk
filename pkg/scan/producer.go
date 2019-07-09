package scan

type Producer interface {
	Produce() <-chan Target
}

type ReProducer interface {
	Reproduce() func(r Result) <-chan Target
}
