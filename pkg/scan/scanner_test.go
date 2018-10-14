package scan_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stretchr/testify/assert"
)

func TestScannerShouldEmitEventWhenScanningATarget(t *testing.T) {
	result := &scan.Result{}

	emitter := emission.NewEmitter()
	emitter.On(scan.EventResultFound, func(r *scan.Result) {
		result = r
	})

	doer := &scan.DoerMock{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			assert.Equal(t, "GET", request.Method)
			assert.Equal(t, "/my/url/to/something", request.URL.Path)

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			}, nil
		},
	}

	logger, loggerBuffer := test.NewLogger()

	sut := scan.NewScanner(doer, emitter, logger)
	sut.AddTarget(
		scan.Target{
			Path:   "my/url/to/something",
			Method: "GET",
		},
	)

	u, err := url.ParseRequestURI("http://127.0.0.1")
	assert.NoError(t, err)

	sut.Release()
	sut.Scan(u, 3)

	time.Sleep(1 * time.Second)

	assert.NotContains(t, loggerBuffer.String(), "error")
	assert.NotContains(t, loggerBuffer.String(), "warn")

	assert.Equal(t, 200, result.Response.StatusCode)
}
