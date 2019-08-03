package scan_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/client"
	"github.com/stefanoj3/dirstalk/pkg/scan/producer"
	"github.com/stretchr/testify/assert"
)

func TestScanningWithEmptyProducerWillProduceNoResults(t *testing.T) {
	logger, _ := test.NewLogger()

	prod := producer.NewDictionaryProducer([]string{}, []string{}, 1)
	c := &http.Client{Timeout: time.Microsecond}
	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod, []int{http.StatusNotFound}),
		logger,
	)

	results := sut.Scan(test.MustParseUrl(t, "http://localhost/"), 10)

	for r := range results {
		t.Fatalf("No results expected, got %s", r.Target.Path)
	}
}

func TestScannerWillLogAnErrorWithInvalidDictionary(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	prod := producer.NewDictionaryProducer(
		[]string{"\n"},
		[]string{"/home"},
		1,
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	defer testServer.Close()

	c, err := client.NewClientFromConfig(
		1000,
		nil,
		"",
		false,
		nil,
		nil,
		test.MustParseUrl(t, testServer.URL),
	)
	assert.NoError(t, err)
	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod, []int{http.StatusNotFound}),
		logger,
	)

	results := sut.Scan(test.MustParseUrl(t, testServer.URL), 10)

	for r := range results {
		t.Fatalf("No results expected, got %s", r.Target.Path)
	}

	assert.Contains(t, loggerBuffer.String(), "failed to build request")
	assert.Contains(t, loggerBuffer.String(), "invalid method")
	assert.Equal(t, 0, serverAssertion.Len())
}

func TestScannerWillNotRedirectIfStatusCodeIsInvalid(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	prod := producer.NewDictionaryProducer(
		[]string{http.MethodGet},
		[]string{"/home"},
		3,
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("location", "/potato")
			if r.URL.Path == "/home" {
				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer testServer.Close()

	c, err := client.NewClientFromConfig(
		1000,
		nil,
		"",
		false,
		nil,
		nil,
		test.MustParseUrl(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod, []int{http.StatusNotFound}),
		logger,
	)

	results := make([]scan.Result, 0, 2)
	resultsChannel := sut.Scan(test.MustParseUrl(t, testServer.URL), 10)

	for r := range resultsChannel {
		results = append(results, r)
	}

	expectedsResults := []scan.Result{
		{
			Target:     scan.Target{Path: "/home", Method: http.MethodGet, Depth: 3},
			StatusCode: http.StatusOK,
			URL:        *test.MustParseUrl(t, testServer.URL+"/home"),
		},
		{
			Target:     scan.Target{Path: "/home/home", Method: http.MethodGet, Depth: 2},
			StatusCode: http.StatusNotFound,
			URL:        *test.MustParseUrl(t, testServer.URL+"/home/home"),
		},
	}

	assert.Equal(t, expectedsResults, results)

	assert.Contains(t, loggerBuffer.String(), "/home")
	assert.Contains(t, loggerBuffer.String(), "/home/home")
	assert.Equal(t, 2, serverAssertion.Len())
}

func TestScannerWillChangeMethodForRedirect(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	prod := producer.NewDictionaryProducer(
		[]string{http.MethodPatch},
		[]string{"/home"},
		3,
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/home" {
				http.Redirect(w, r, "/potato", http.StatusMovedPermanently)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer testServer.Close()

	c, err := client.NewClientFromConfig(
		1000,
		nil,
		"",
		false,
		nil,
		nil,
		test.MustParseUrl(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod, []int{http.StatusNotFound}),
		logger,
	)

	results := make([]scan.Result, 0, 4)
	resultsChannel := sut.Scan(test.MustParseUrl(t, testServer.URL), 1)

	for r := range resultsChannel {
		results = append(results, r)
	}

	expectedResults := []scan.Result{
		{
			Target:     scan.Target{Path: "/home", Method: http.MethodPatch, Depth: 3},
			StatusCode: http.StatusMovedPermanently,
			URL:        *test.MustParseUrl(t, testServer.URL+"/home"),
		},
		{
			Target:     scan.Target{Path: "/potato", Method: http.MethodGet, Depth: 2},
			StatusCode: http.StatusNotFound,
			URL:        *test.MustParseUrl(t, testServer.URL+"/potato"),
		},
		{
			Target:     scan.Target{Path: "/home/home", Method: http.MethodPatch, Depth: 2},
			StatusCode: http.StatusNotFound,
			URL:        *test.MustParseUrl(t, testServer.URL+"/home/home"),
		},
	}

	assert.Equal(t, expectedResults, results)

	loggerBufferAsString := loggerBuffer.String()
	assert.Contains(t, loggerBufferAsString, "/home")
	assert.Contains(t, loggerBufferAsString, "/potato")
	assert.Contains(t, loggerBufferAsString, "/home/home")
	assert.NotContains(t, loggerBufferAsString, "error")
	assert.Equal(t, 3, serverAssertion.Len())
}

func TestScannerWhenOutOfDepthWillNotFollowRedirect(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	prod := producer.NewDictionaryProducer(
		[]string{http.MethodPatch},
		[]string{"/home"},
		0,
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/home" {
				http.Redirect(w, r, "/potato", http.StatusMovedPermanently)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer testServer.Close()

	c, err := client.NewClientFromConfig(
		1000,
		nil,
		"",
		false,
		nil,
		nil,
		test.MustParseUrl(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod, []int{http.StatusNotFound}),
		logger,
	)

	results := make([]scan.Result, 0, 1)
	resultsChannel := sut.Scan(test.MustParseUrl(t, testServer.URL), 1)

	for r := range resultsChannel {
		results = append(results, r)
	}

	expectedResults := []scan.Result{
		{
			Target:     scan.Target{Path: "/home", Method: http.MethodPatch, Depth: 0},
			StatusCode: http.StatusMovedPermanently,
			URL:        *test.MustParseUrl(t, testServer.URL+"/home"),
		},
	}

	assert.Equal(t, expectedResults, results)

	loggerBufferAsString := loggerBuffer.String()
	assert.Contains(t, loggerBufferAsString, "/home")
	assert.Contains(t, loggerBufferAsString, "depth is 0, not following any redirect")
	assert.NotContains(t, loggerBufferAsString, "error")
	assert.Equal(t, 1, serverAssertion.Len())
}

func TestScannerWillSkipRedirectWhenLocationHostIsDifferent(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	prod := producer.NewDictionaryProducer(
		[]string{http.MethodPatch},
		[]string{"/home"},
		3,
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/home" {
				http.Redirect(w, r, "http://gibberish/potato", http.StatusMovedPermanently)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer testServer.Close()

	c, err := client.NewClientFromConfig(
		1000,
		nil,
		"",
		false,
		nil,
		nil,
		test.MustParseUrl(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod, []int{http.StatusNotFound}),
		logger,
	)

	results := make([]scan.Result, 0, 4)
	resultsChannel := sut.Scan(test.MustParseUrl(t, testServer.URL), 1)

	for r := range resultsChannel {
		results = append(results, r)
	}

	expectedResults := []scan.Result{
		{
			Target:     scan.Target{Path: "/home", Method: http.MethodPatch, Depth: 3},
			StatusCode: http.StatusMovedPermanently,
			URL:        *test.MustParseUrl(t, testServer.URL+"/home"),
		},
		{
			Target:     scan.Target{Path: "/home/home", Method: http.MethodPatch, Depth: 2},
			StatusCode: http.StatusNotFound,
			URL:        *test.MustParseUrl(t, testServer.URL+"/home/home"),
		},
	}

	assert.Equal(t, expectedResults, results)

	loggerBufferAsString := loggerBuffer.String()
	assert.Contains(t, loggerBufferAsString, "/home")
	assert.Contains(t, loggerBufferAsString, "/home/home")
	assert.Contains(t, loggerBufferAsString, "skipping redirect, pointing to a different host")
	assert.NotContains(t, loggerBufferAsString, "error")
	assert.Equal(t, 2, serverAssertion.Len())
}
