package scan

import (
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool"
)

const (
	jobQueueSize = 100
)

// Target represents the target to scan
type Target struct {
	Path   string
	Method string
	Depth  int
}

// Result represents the result of the scan of a single URL
type Result struct {
	Target   Target
	URL      *url.URL
	Response *http.Response
}

func NewScanner(
	httpClient Doer,
	eventEmitter *emission.Emitter,
	logger *logrus.Logger,
) *Scanner {
	return &Scanner{
		httpClient:   httpClient,
		eventEmitter: eventEmitter,
		logger:       logger,
		jobQueue:     make(chan Target, jobQueueSize),
		isReleased:   abool.New(),
	}
}

type Scanner struct {
	httpClient   Doer
	eventEmitter *emission.Emitter
	logger       *logrus.Logger
	jobQueue     chan Target
	isReleased   *abool.AtomicBool
}

func (s *Scanner) AddTarget(target Target) {
	s.jobQueue <- target
}

func (s *Scanner) Scan(baseUrl *url.URL, workers int) {
	u := normalizeBaseURL(*baseUrl)

	wg := sync.WaitGroup{}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			s.work(u)
			wg.Done()
		}()

	}

	wg.Wait()
}

func (s *Scanner) work(baseURL url.URL) {
	attempts := 3

	for {
		select {
		case target := <-s.jobQueue:
			s.processTarget(baseURL, target)
			continue
		case <-time.After(400 * time.Millisecond):
		}

		if s.isReleased.IsSet() {
			attempts--
		}

		if attempts == 0 {
			break
		}
	}
}

func (s *Scanner) Release() {
	s.isReleased.Set()
}

func (s *Scanner) processTarget(baseURL url.URL, target Target) {
	l := s.logger.WithFields(logrus.Fields{
		"method": target.Method,
		"depth":  target.Depth,
		"path":   target.Path,
	})

	l.Debug("Working")

	u := buildURL(baseURL, target)
	req, err := http.NewRequest(target.Method, u.String(), nil)
	if err != nil {
		l.WithError(err).Error("failed to build request")
		return
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		l.WithError(err).Warn("failed to perform request")
		return
	}

	if err := res.Body.Close(); err != nil {
		l.WithError(err).Warn("failed to close response body")
	}

	result := &Result{
		Target:   target,
		Response: res,
	}

	s.eventEmitter.Emit(EventResultFound, result)
}

func normalizeBaseURL(baseURL url.URL) url.URL {
	if strings.HasSuffix(baseURL.Path, "/") {
		return baseURL
	}

	baseURL.Path += "/"

	return baseURL
}

func buildURL(baseURL url.URL, target Target) url.URL {
	baseURL.Path = path.Join(baseURL.Path, target.Path)
	return baseURL
}
