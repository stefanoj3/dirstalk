package scan

import (
	"path"
	"sync"

	"github.com/chuckpreslar/emission"
	"github.com/stefanoj3/dirstalk/pkg/pathutil"
)

var statusCodesToSkip = map[int]bool{
	404: false,
}

func NewReProcessor(eventEmitter *emission.Emitter, httpMethods []string, dictionary []string) *ReProcessor {
	return &ReProcessor{eventEmitter: eventEmitter, httpMethods: httpMethods, dictionary: dictionary}
}

type ReProcessor struct {
	eventEmitter   *emission.Emitter
	httpMethods    []string
	dictionary     []string
	resultRegistry sync.Map
}

func (r *ReProcessor) ReProcess(result *Result) {
	if _, ok := statusCodesToSkip[result.Response.StatusCode]; ok {
		return
	}

	if result.Target.Depth <= 0 {
		return
	}

	// no point in appending to a filename
	if pathutil.HasExtension(result.Target.Path) {
		return
	}

	_, inRegistry := r.resultRegistry.Load(result.Target.Path)
	if inRegistry {
		return
	}
	r.resultRegistry.Store(result.Target.Path, false)

	for _, entry := range r.dictionary {
		for _, httpMethod := range r.httpMethods {
			t := result.Target
			t.Depth--
			t.Path = path.Join(t.Path, entry)
			t.Method = httpMethod

			r.eventEmitter.Emit(EventTargetProduced, t)
		}

	}
}
