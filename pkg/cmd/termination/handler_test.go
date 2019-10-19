package termination_test

import (
	"sync"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/cmd/termination"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	handler := termination.NewTerminationHandler(3)

	handler.SignalTermination()
	assert.False(t, handler.ShouldTerminate())

	handler.SignalTermination()
	assert.False(t, handler.ShouldTerminate())

	handler.SignalTermination()
	assert.True(t, handler.ShouldTerminate())

	handler.SignalTermination()
	assert.True(t, handler.ShouldTerminate())
}

func TestHandlerShouldWorkWithMultipleRoutines(_ *testing.T) {
	handler := termination.NewTerminationHandler(10)

	const workers = 1000

	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			handler.ShouldTerminate()
			handler.SignalTermination()
		}()
	}

	wg.Wait()
}
