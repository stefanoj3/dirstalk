package scan_test

import (
	"net/http"
	"testing"

	"github.com/chuckpreslar/emission"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stretchr/testify/assert"
)

func TestTargetProducer(t *testing.T) {
	actualTargets := make([]scan.Target, 0)

	emitter := emission.NewEmitter()
	emitter.On(scan.EventTargetProduced, func(t scan.Target) {
		actualTargets = append(actualTargets, t)
	})

	producer := scan.NewTargetProducer(
		emitter,
		[]string{http.MethodGet, http.MethodPost},
		[]string{"/path", "/mypath"},
		2,
	)

	producer.Run()

	expectedTargets := []scan.Target{
		{Path: "/path", Method: http.MethodGet, Depth: 2},
		{Path: "/path", Method: http.MethodPost, Depth: 2},
		{Path: "/mypath", Method: http.MethodGet, Depth: 2},
		{Path: "/mypath", Method: http.MethodPost, Depth: 2},
	}

	assert.Equal(t, expectedTargets, actualTargets)
}
