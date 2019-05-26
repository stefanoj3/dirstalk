package scan_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/armon/go-socks5"
	"github.com/chuckpreslar/emission"
	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stretchr/testify/assert"
)

const socks5TestServerHost = "127.0.0.1:8899"

func TestStartScan(t *testing.T) {
	logger, _ := test.NewLogger()

	requestMap := &sync.Map{}
	testServer := buildTestServer(requestMap)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	config := &scan.Config{
		Threads:               3,
		DictionaryPath:        "testdata/dictionary1.txt",
		HTTPMethods:           []string{http.MethodGet},
		TimeoutInMilliseconds: 50,
		ScanDepth:             3,
	}

	actualResults := make([]scan.Target, 0, 2)
	mx := sync.Mutex{}
	eventManager := emission.NewEmitter()
	eventManager.On(scan.EventResultFound, func(r *scan.Result) {
		if r.Response.StatusCode != 200 {
			return
		}

		mx.Lock()
		defer mx.Unlock()
		actualResults = append(actualResults, r.Target)
	})

	err = scan.StartScan(logger, eventManager, config, u)
	assert.NoError(t, err)

	// Asserting which requests are made to the remote service
	assertRequest(t, true, http.MethodGet, "/home", requestMap)
	assertRequest(t, true, http.MethodGet, "/about", requestMap)
	assertRequest(t, true, http.MethodGet, "/home/home", requestMap)
	assertRequest(t, true, http.MethodGet, "/home/about", requestMap)
	assertRequest(t, true, http.MethodGet, "/home/about/home", requestMap)
	assertRequest(t, true, http.MethodGet, "/home/about/about", requestMap)

	assertRequest(t, false, http.MethodDelete, "/home", requestMap)
	assertRequest(t, false, http.MethodDelete, "/about", requestMap)

	assertRequest(t, false, http.MethodPatch, "/home", requestMap)
	assertRequest(t, false, http.MethodPatch, "/about", requestMap)

	assertRequest(t, false, http.MethodPost, "/home", requestMap)
	assertRequest(t, false, http.MethodPost, "/about", requestMap)

	assertRequest(t, false, http.MethodPut, "/home", requestMap)
	assertRequest(t, false, http.MethodPut, "/about", requestMap)

	assertRequest(t, false, http.MethodGet, "/about/home", requestMap)
	assertRequest(t, false, http.MethodGet, "/about/about", requestMap)
	// -----------------------------------------------------------

	// Asserting on the actual results found - considering only 200s for this test
	assert.Len(t, actualResults, 2)
	assertTargetsContains(t, scan.Target{Depth: 3, Path: "/home", Method: http.MethodGet}, actualResults)
	assertTargetsContains(t, scan.Target{Depth: 2, Path: "/home/about", Method: http.MethodGet}, actualResults)
}

func TestStartScanWithSocks5ShouldFindResultsWhenAServerIsAvailable(t *testing.T) {
	logger, _ := test.NewLogger()

	testServer := buildTestServer(&sync.Map{})
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	socks5URL, err := url.Parse("socks5://" + socks5TestServerHost)
	assert.NoError(t, err)

	config := &scan.Config{
		Threads:               3,
		DictionaryPath:        "testdata/dictionary1.txt",
		HTTPMethods:           []string{http.MethodGet},
		TimeoutInMilliseconds: 50,
		ScanDepth:             3,
		Socks5Url:             socks5URL,
	}

	listener := startSocks5TestServer(t)
	defer listener.Close()

	actualResults := make([]scan.Target, 0, 2)
	mx := sync.Mutex{}
	eventManager := emission.NewEmitter()
	eventManager.On(scan.EventResultFound, func(r *scan.Result) {
		if r.Response.StatusCode != 200 {
			return
		}

		mx.Lock()
		defer mx.Unlock()
		actualResults = append(actualResults, r.Target)
	})

	err = scan.StartScan(logger, eventManager, config, u)
	assert.NoError(t, err)

	assert.Len(t, actualResults, 2)
	assertTargetsContains(t, scan.Target{Depth: 3, Path: "/home", Method: http.MethodGet}, actualResults)
	assertTargetsContains(t, scan.Target{Depth: 2, Path: "/home/about", Method: http.MethodGet}, actualResults)
}

func TestShouldUseTheSpecifiedUserAgent(t *testing.T) {
	const testUserAgent = "my_test_user_agent"

	logger, _ := test.NewLogger()

	var request *http.Request
	doneChannel := make(chan bool)

	testServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			request = r
			doneChannel <- true
		}),
	)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	config := &scan.Config{
		Threads:               3,
		DictionaryPath:        "testdata/one_element_dictionary.txt",
		HTTPMethods:           []string{http.MethodGet},
		TimeoutInMilliseconds: 50,
		ScanDepth:             3,
		UserAgent:             testUserAgent,
	}

	eventManager := emission.NewEmitter()
	err = scan.StartScan(logger, eventManager, config, u)
	assert.NoError(t, err)

	select {
	case <-doneChannel:
		assert.Equal(t, testUserAgent, request.Header.Get("User-Agent"))
	case <-time.After(time.Second * 1):
		t.Fatal("failed to receive request")
	}
}

func TestShouldFailToScanWithAnUnreachableSocks5Server(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	requestMap := &sync.Map{}
	testServer := buildTestServer(requestMap)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	socks5URL, err := url.Parse("socks5://127.0.0.1:12345")
	assert.NoError(t, err)

	config := &scan.Config{
		Threads:               3,
		DictionaryPath:        "testdata/dictionary2.txt",
		HTTPMethods:           []string{http.MethodGet},
		TimeoutInMilliseconds: 50,
		ScanDepth:             3,
		Socks5Url:             socks5URL,
	}

	listener := startSocks5TestServer(t)
	defer listener.Close()

	actualResults := make([]scan.Target, 0, 2)
	mx := sync.Mutex{}
	eventManager := emission.NewEmitter()
	eventManager.On(scan.EventResultFound, func(r *scan.Result) {
		if r.Response.StatusCode != 200 {
			return
		}

		mx.Lock()
		defer mx.Unlock()
		actualResults = append(actualResults, r.Target)
	})

	err = scan.StartScan(logger, eventManager, config, u)
	assert.NoError(t, err)

	assert.Len(t, actualResults, 0)
	assert.Contains(t, loggerBuffer.String(), "connection refused")

	requestMap.Range(func(key, value interface{}) bool {
		t.Fatal("no request was supposed to be recorded: socks5 is down, the server should remain unreachable")
		return true
	})
}

func TestShouldRetainCookiesSetByTheServerWhenCookieJarIsEnabled(t *testing.T) {
	const (
		cookieName  = "my_cookies_name"
		cookieValue = "my_cookie_value_123"
	)
	logger, _ := test.NewLogger()

	once := sync.Once{}
	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			once.Do(func() {
				http.SetCookie(
					w,
					&http.Cookie{
						Name:    cookieName,
						Value:   cookieValue,
						Expires: time.Now().AddDate(0, 1, 0),
					},
				)
			})
		}),
	)

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	config := &scan.Config{
		Threads:               1,
		DictionaryPath:        "testdata/dictionary1.txt",
		HTTPMethods:           []string{http.MethodGet},
		TimeoutInMilliseconds: 400,
		ScanDepth:             2,
		UseCookieJar:          true,
	}
	eventManager := emission.NewEmitter()

	err = scan.StartScan(logger, eventManager, config, u)
	assert.NoError(t, err)

	assert.Equal(t, 14, serverAssertion.Len())
	serverAssertion.Range(func(index int, r http.Request) {
		if index == 0 { // the first request should have no cookies
			assert.Equal(t, 0, len(r.Cookies()))
			return
		}

		assert.Equal(t, 1, len(r.Cookies()), "Only one cookie expected, got: %v", r.Cookies())
		assert.Equal(t, cookieName, r.Cookies()[0].Name)
		assert.Equal(t, cookieValue, r.Cookies()[0].Value)
	})
}

func TestShouldNotSendAnyCookieIfServerSetNoneWhenUsingCookieJar(t *testing.T) {
	logger, _ := test.NewLogger()

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	config := &scan.Config{
		Threads:               1,
		DictionaryPath:        "testdata/dictionary1.txt",
		HTTPMethods:           []string{http.MethodGet},
		TimeoutInMilliseconds: 400,
		ScanDepth:             2,
		UseCookieJar:          true,
	}
	eventManager := emission.NewEmitter()

	err = scan.StartScan(logger, eventManager, config, u)
	assert.NoError(t, err)

	assert.Equal(t, 14, serverAssertion.Len())
	serverAssertion.Range(func(index int, r http.Request) {
		assert.Equal(t, 0, len(r.Cookies()), "No cookies expected, got: %v", r.Cookies())
	})
}

func assertTargetsContains(t *testing.T, target scan.Target, results []scan.Target) {
	for _, actualResult := range results {
		if target == actualResult {
			return
		}
	}

	t.Errorf("Target %v not found in %v", target, results)
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
			storeRequest(requestMap, r.Method, r.URL.Path)

			if r.Method == http.MethodGet && r.URL.Path == "/home" {
				w.WriteHeader(200)
				return
			}

			if r.Method == http.MethodGet && r.URL.Path == "/home/about" {
				w.WriteHeader(200)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}),
	)
}

func startSocks5TestServer(t *testing.T) net.Listener {
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		t.Fatalf("failed to create socks5: %s", err.Error())
	}

	listener, err := net.Listen("tcp", socks5TestServerHost)
	if err != nil {
		t.Fatalf("failed to create listener: %s", err.Error())
	}

	go func() {
		// Create SOCKS5 proxy on localhost port 8000
		if err := server.Serve(listener); err != nil {
			t.Logf("socks5 stopped serving: %s", err.Error())
		}
	}()

	return listener
}

func storeRequest(requestMap *sync.Map, method, path string) {
	requestMap.Store(methodAndPathToString(method, path), true)
}

func hasRequest(requestMap *sync.Map, method, path string) bool {
	_, ok := requestMap.Load(methodAndPathToString(method, path))
	return ok
}

func methodAndPathToString(method, path string) string {
	return fmt.Sprintf("%s_%s", method, path)
}
