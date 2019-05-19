package scan

import (
	"github.com/chuckpreslar/emission"
)

func NewTargetProducer(
	eventEmitter *emission.Emitter,
	methods []string,
	dictionary []string,
	depth int,
) *TargetProducer {
	return &TargetProducer{
		eventEmitter: eventEmitter,
		methods:      methods,
		dictionary:   dictionary,
		depth:        depth,
	}

}

type TargetProducer struct {
	eventEmitter *emission.Emitter
	methods      []string
	dictionary   []string
	depth        int
}

func (p *TargetProducer) Run() {
	for _, entry := range p.dictionary {
		for _, method := range p.methods {
			p.eventEmitter.Emit(
				EventTargetProduced,
				Target{
					Path:   entry,
					Method: method,
					Depth:  p.depth,
				},
			)
		}
	}

	p.eventEmitter.Emit(EventProducerFinished)
}
