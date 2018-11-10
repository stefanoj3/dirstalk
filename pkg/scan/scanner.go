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

type Target struct {
	Path   string
	Method string
	Depth  int
}

// Result represent the result of the scan of a single URL
type Result struct {
	Target   Target
	URL      *url.URL
	Response *http.Response
}

type Scanner struct {
	httpClient      Doer
	eventEmitter    *emission.Emitter
	logger          *logrus.Logger
	jobQueue        chan Target
	workerWaitGroup sync.WaitGroup
	isReleased      *abool.AtomicBool
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

func (s *Scanner) AddTarget(target Target) {
	s.jobQueue <- target
}

func (s *Scanner) Scan(baseUrl *url.URL, workers int) {
	u := normalizeBaseUrl(*baseUrl)

	for i := 0; i < workers; i++ {
		go s.work(u)
		s.workerWaitGroup.Add(1)
	}

	s.workerWaitGroup.Wait()
}

func (s *Scanner) work(baseUrl url.URL) {
	attempts := 3

	for {
		select {
		case target := <-s.jobQueue:
			s.processTarget(baseUrl, target)
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

	s.workerWaitGroup.Done()
}

func (s *Scanner) Release() {
	s.isReleased.Set()
}

func (s *Scanner) processTarget(baseUrl url.URL, target Target) {
	s.logger.WithFields(logrus.Fields{
		"method": target.Method,
		"depth":  target.Depth,
		"path":   target.Path,
	}).Debug("Working")

	u := buildUrl(baseUrl, target)
	req, err := http.NewRequest(target.Method, u.String(), nil)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"url":    target.Path,
			"method": target.Method,
			"depth":  target.Depth,
			"error":  err.Error(),
		}).Error(
			"failed to build request",
		)
		return
	}

	res, err := s.httpClient.Do(req)
	if err != nil {

		s.logger.WithFields(logrus.Fields{
			"url":    target.Path,
			"method": target.Method,
			"depth":  target.Depth,
			"error":  err.Error(),
		}).Warn(
			"failed to perform request",
		)
		return
	}

	res.Body.Close()

	result := &Result{
		Target:   target,
		URL:      &u,
		Response: res,
	}

	s.eventEmitter.Emit(EventResultFound, result)
}

func normalizeBaseUrl(baseUrl url.URL) url.URL {
	if strings.HasSuffix(baseUrl.Path, "/") {
		return baseUrl
	}

	baseUrl.Path += "/"

	return baseUrl
}

func buildUrl(baseUrl url.URL, target Target) url.URL {
	baseUrl.Path = path.Join(baseUrl.Path, target.Path)
	return baseUrl
}
