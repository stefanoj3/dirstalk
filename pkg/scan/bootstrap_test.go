package scan_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/scan"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stretchr/testify/assert"
)

func TestStartScan(t *testing.T) {
	logger, _ := test.NewLogger()

	requestMap := &sync.Map{}
	testServer := buildTestServer(requestMap)

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	config := &scan.Config{
		Threads:               3,
		Dictionary:            []string{"/home", "/about"},
		HttpMethods:           []string{http.MethodGet},
		TimeoutInMilliseconds: 2000,
		ScanDepth:             3,
	}

	err = scan.StartScan(logger, config, u)
	assert.NoError(t, err)

	assertRequest(t, true, "GET", "/home", requestMap)
	assertRequest(t, true, "GET", "/about", requestMap)
	assertRequest(t, true, "GET", "/home/home", requestMap)
	assertRequest(t, true, "GET", "/home/about", requestMap)

	assertRequest(t, false, "DELETE", "/home", requestMap)
	assertRequest(t, false, "DELETE", "/about", requestMap)

	assertRequest(t, false, "PATCH", "/home", requestMap)
	assertRequest(t, false, "PATCH", "/about", requestMap)

	assertRequest(t, false, "POST", "/home", requestMap)
	assertRequest(t, false, "POST", "/about", requestMap)

	assertRequest(t, false, "PUT", "/home", requestMap)
	assertRequest(t, false, "PUT", "/about", requestMap)

	assertRequest(t, false, "GET", "/about/home", requestMap)
	assertRequest(t, false, "GET", "/about/about", requestMap)
}

func assertRequest(t *testing.T, expected bool, method, path string, requestMap *sync.Map) {
	if expected {
		assert.True(
			t,
			hasRequest(requestMap, method, path),
			"expected request for `%s %s`, none received",
			method,
			path,
		)
	} else {
		assert.False(
			t,
			hasRequest(requestMap, method, path),
			"no request was expected for `%s %s`",
			method,
			path,
		)
	}
}

func buildTestServer(requestMap *sync.Map) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			storeRequest(requestMap, r)

			if r.Method == http.MethodGet && r.URL.Path == "/home" {
				w.WriteHeader(200)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}),
	)
}

func storeRequest(requestMap *sync.Map, r *http.Request) {
	requestMap.Store(methodAndPathToString(r.Method, r.URL.Path), true)
}

func hasRequest(requestMap *sync.Map, method, path string) bool {
	_, ok := requestMap.Load(methodAndPathToString(method, path))
	return ok
}

func methodAndPathToString(method, path string) string {
	return fmt.Sprintf("%s_%s", method, path)
}
