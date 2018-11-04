package scan_test

import (
	"net/http"
	"testing"

	"github.com/chuckpreslar/emission"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stretchr/testify/assert"
)

func TestReProcessor(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		scenario        string
		target          scan.Target
		response        *http.Response
		expectedTargets []scan.Target
	}{
		{
			scenario: "target is not a file, response is 200, depth is > 0",
			target: scan.Target{
				Path:   "mybasepath/",
				Method: "GET",
				Depth:  3,
			},
			response: &http.Response{StatusCode: 200},
			expectedTargets: []scan.Target{
				{
					Path:   "mybasepath/path",
					Method: "GET",
					Depth:  2,
				},
				{
					Path:   "mybasepath/path",
					Method: "POST",
					Depth:  2,
				},
				{
					Path:   "mybasepath/mypath",
					Method: "GET",
					Depth:  2,
				},
				{
					Path:   "mybasepath/mypath",
					Method: "POST",
					Depth:  2,
				},
				{
					Path:   "mybasepath/myfolder/path",
					Method: "GET",
					Depth:  2,
				},
				{
					Path:   "mybasepath/myfolder/path",
					Method: "POST",
					Depth:  2,
				},
				{
					Path:   "mybasepath/myfolder2/path",
					Method: "GET",
					Depth:  2,
				},
				{
					Path:   "mybasepath/myfolder2/path",
					Method: "POST",
					Depth:  2,
				},
			},
		},
		{
			scenario: "target is not a file, response is 404, depth is > 0",
			target: scan.Target{
				Path:   "mybasepath/",
				Method: "GET",
				Depth:  3,
			},
			response:        &http.Response{StatusCode: 404},
			expectedTargets: []scan.Target{},
		},
		{
			scenario: "target is not a file, response is 200, depth is 0",
			target: scan.Target{
				Path:   "mybasepath/",
				Method: "GET",
				Depth:  0,
			},
			response:        &http.Response{StatusCode: 200},
			expectedTargets: []scan.Target{},
		},
		{
			scenario: "target is a file, response is 200, depth is 3",
			target: scan.Target{
				Path:   "image.jpg",
				Method: "GET",
				Depth:  3,
			},
			response:        &http.Response{StatusCode: 200},
			expectedTargets: []scan.Target{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actualTargets := []scan.Target{}

			emitter := emission.NewEmitter()
			emitter.On(scan.EventTargetProduced, func(t scan.Target) {
				actualTargets = append(actualTargets, t)
			})

			reprocessor := scan.NewReProcessor(
				emitter,
				[]string{"GET", "POST"},
				[]string{"/path", "mypath/", "/myfolder/path", "myfolder2/path/"},
			)

			reprocessor.ReProcess(
				&scan.Result{
					Target:   tc.target,
					Response: tc.response,
				},
			)

			assert.Equal(t, tc.expectedTargets, actualTargets)
		})
	}
}

func TestReProcessorShouldNotProcessTwiceTheSamePath(t *testing.T) {
	t.Parallel()

	actualTargets := []scan.Target{}

	emitter := emission.NewEmitter()
	emitter.On(scan.EventTargetProduced, func(t scan.Target) {
		actualTargets = append(actualTargets, t)
	})

	reprocessor := scan.NewReProcessor(
		emitter,
		[]string{"GET"},
		[]string{"/image.jpg"},
	)

	expectedTargets := []scan.Target{
		{
			Path:   "mybasepath/image.jpg",
			Method: "GET",
			Depth:  1,
		},
	}

	reprocessor.ReProcess(
		&scan.Result{
			Target: scan.Target{
				Path:   "mybasepath",
				Method: "GET",
				Depth:  2,
			},
			Response: &http.Response{StatusCode: 200},
		},
	)

	assert.Equal(t, expectedTargets, actualTargets)

	reprocessor.ReProcess(
		&scan.Result{
			Target: scan.Target{
				Path:   "mybasepath",
				Method: "POST",
				Depth:  2,
			},
			Response: &http.Response{StatusCode: 200},
		},
	)

	assert.Equal(t, expectedTargets, actualTargets)
}
