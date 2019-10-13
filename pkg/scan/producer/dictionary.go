package producer

import (
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

func NewDictionaryProducer(
	methods []string,
	dictionary []string,
	depth int,
) *DictionaryProducer {
	return &DictionaryProducer{
		methods:    methods,
		dictionary: dictionary,
		depth:      depth,
	}
}

type DictionaryProducer struct {
	methods    []string
	dictionary []string
	depth      int
}

func (p *DictionaryProducer) Produce() <-chan scan.Target {
	targets := make(chan scan.Target, 10)

	go func() {
		defer close(targets)

		for _, entry := range p.dictionary {
			for _, method := range p.methods {
				targets <- scan.Target{
					Path:   entry,
					Method: method,
					Depth:  p.depth,
				}
			}
		}
	}()

	return targets
}
