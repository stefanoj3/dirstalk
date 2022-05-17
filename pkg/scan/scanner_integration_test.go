package scan_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/client"
	"github.com/stefanoj3/dirstalk/pkg/scan/filter"
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
		producer.NewReProducer(prod),
		filter.NewHTTPStatusResultFilter([]int{http.StatusNotFound}, false),
		logger,
	)

	results := sut.Scan(context.Background(), test.MustParseURL(t, "http://localhost/"), 10)

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
		true,
		false,
		test.MustParseURL(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod),
		filter.NewHTTPStatusResultFilter([]int{http.StatusNotFound}, false),
		logger,
	)

	results := sut.Scan(context.Background(), test.MustParseURL(t, testServer.URL), 10)

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
		true,
		false,
		test.MustParseURL(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod),
		filter.NewHTTPStatusResultFilter([]int{http.StatusNotFound}, false),
		logger,
	)

	results := make([]scan.Result, 0, 2)
	resultsChannel := sut.Scan(context.Background(), test.MustParseURL(t, testServer.URL), 10)

	for r := range resultsChannel {
		results = append(results, r)
	}

	expectedsResults := []scan.Result{
		{
			Target:     scan.Target{Path: "/home", Method: http.MethodGet, Depth: 3},
			StatusCode: http.StatusOK,
			URL:        *test.MustParseURL(t, testServer.URL+"/home"),
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

			if r.URL.Path == "/potato" {
				w.WriteHeader(http.StatusCreated)

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
		true,
		false,
		test.MustParseURL(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod),
		filter.NewHTTPStatusResultFilter([]int{http.StatusNotFound}, false),
		logger,
	)

	results := make([]scan.Result, 0, 3)
	resultsChannel := sut.Scan(context.Background(), test.MustParseURL(t, testServer.URL), 1)

	for r := range resultsChannel {
		results = append(results, r)
	}

	expectedResults := []scan.Result{
		{
			Target:     scan.Target{Path: "/home", Method: http.MethodPatch, Depth: 3},
			StatusCode: http.StatusMovedPermanently,
			URL:        *test.MustParseURL(t, testServer.URL+"/home"),
		},
		{
			Target:     scan.Target{Path: "/potato", Method: http.MethodGet, Depth: 2},
			StatusCode: http.StatusCreated,
			URL:        *test.MustParseURL(t, testServer.URL+"/potato"),
		},
	}

	assert.Equal(t, expectedResults, results)

	assert.NotContains(t, loggerBuffer.String(), "error")
	assert.Equal(t, 4, serverAssertion.Len())
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
		true,
		false,
		test.MustParseURL(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod),
		filter.NewHTTPStatusResultFilter([]int{http.StatusNotFound}, false),
		logger,
	)

	results := make([]scan.Result, 0, 1)
	resultsChannel := sut.Scan(context.Background(), test.MustParseURL(t, testServer.URL), 1)

	for r := range resultsChannel {
		results = append(results, r)
	}

	expectedResults := []scan.Result{
		{
			Target:     scan.Target{Path: "/home", Method: http.MethodPatch, Depth: 0},
			StatusCode: http.StatusMovedPermanently,
			URL:        *test.MustParseURL(t, testServer.URL+"/home"),
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
		true,
		false,
		test.MustParseURL(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod),
		filter.NewHTTPStatusResultFilter([]int{http.StatusNotFound}, false),
		logger,
	)

	results := make([]scan.Result, 0, 2)
	resultsChannel := sut.Scan(context.Background(), test.MustParseURL(t, testServer.URL), 1)

	for r := range resultsChannel {
		results = append(results, r)
	}

	expectedResults := []scan.Result{
		{
			Target:     scan.Target{Path: "/home", Method: http.MethodPatch, Depth: 3},
			StatusCode: http.StatusMovedPermanently,
			URL:        *test.MustParseURL(t, testServer.URL+"/home"),
		},
	}

	assert.Equal(t, expectedResults, results)

	loggerBufferAsString := loggerBuffer.String()
	assert.Contains(t, loggerBufferAsString, "skipping redirect, pointing to a different host")
	assert.NotContains(t, loggerBufferAsString, "error")
	assert.Equal(t, 2, serverAssertion.Len())
}

func TestScannerWillIgnoreRequestRedundantError(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	prod := producer.NewDictionaryProducer(
		[]string{http.MethodGet},
		[]string{"/home", "/home"}, // twice the same entry to trick the client into doing the same request twice
		3,
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		true,
		false,
		test.MustParseURL(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod),
		filter.NewHTTPStatusResultFilter([]int{http.StatusNotFound}, false),
		logger,
	)

	results := make([]scan.Result, 0, 1)
	resultsChannel := sut.Scan(context.Background(), test.MustParseURL(t, testServer.URL), 1)

	for r := range resultsChannel {
		results = append(results, r)
	}

	assert.Equal(t, 0, len(results))

	loggerBufferAsString := loggerBuffer.String()
	assert.Contains(t, loggerBufferAsString, "/home")
	assert.Contains(t, loggerBufferAsString, "this request has been made already")
	assert.Equal(t, 1, serverAssertion.Len())
}

func TestCanCancelScanUsingContext(t *testing.T) {
	logger, _ := test.NewLogger()

	prod := producer.NewDictionaryProducer(
		[]string{http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodPost, http.MethodPut},
		[]string{"/home", "/index", "/about", "/search", "/jobs", "robots.txt", "/subscription", "/orders"},
		200000000,
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer testServer.Close()

	// the depth of the dictionary and the fact that the server returns always a
	// http.StatusOK should keep this test running forever in case the cancellation would not work

	c, err := client.NewClientFromConfig(
		1000,
		nil,
		"",
		false,
		nil,
		nil,
		true,
		false,
		test.MustParseURL(t, testServer.URL),
	)
	assert.NoError(t, err)

	sut := scan.NewScanner(
		c,
		prod,
		producer.NewReProducer(prod),
		filter.NewHTTPStatusResultFilter([]int{http.StatusNotFound}, false),
		logger,
	)

	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	resultsChannel := sut.Scan(ctx, test.MustParseURL(t, testServer.URL), 100)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancelFunc()
	}()

	done := make(chan struct{})

	go func() {
		for range resultsChannel {

		}
		done <- struct{}{}
	}()

	select {
	case <-done:
		t.Log("result channel closed")
	case <-time.After(time.Second * 8):
		t.Fatalf("the scan should have terminated by now, something is wrong with the context cancellation")
	}

	assert.True(t, serverAssertion.Len() > 1)
}
