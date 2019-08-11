package producer

import (
	"sync"

	"github.com/stefanoj3/dirstalk/pkg/common/urlpath"
	"github.com/stefanoj3/dirstalk/pkg/pathutil"
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

const defaultChannelBuffer = 25

func NewReProducer(
	producer scan.Producer,
) *ReProducer {
	return &ReProducer{producer: producer}
}

type ReProducer struct {
	producer scan.Producer
}

// Reproduce will check if it is possible to go deeper on the result provided, if so will
func (r *ReProducer) Reproduce() func(r scan.Result) <-chan scan.Target {
	return r.buildReproducer()
}

func (r *ReProducer) buildReproducer() func(result scan.Result) <-chan scan.Target {
	resultRegistry := sync.Map{}

	return func(result scan.Result) <-chan scan.Target {
		resultChannel := make(chan scan.Target, defaultChannelBuffer)

		go func() {
			defer close(resultChannel)

			if result.Target.Depth <= 0 {
				return
			}

			// no point in appending to a filename
			if pathutil.HasExtension(result.Target.Path) {
				return
			}

			_, inRegistry := resultRegistry.Load(result.Target.Path)
			if inRegistry {
				return
			}
			resultRegistry.Store(result.Target.Path, false)

			for target := range r.producer.Produce() {
				newTarget := result.Target
				newTarget.Depth--
				newTarget.Path = urlpath.Join(newTarget.Path, target.Path)
				newTarget.Method = target.Method

				resultChannel <- newTarget
			}

		}()

		return resultChannel
	}
}
