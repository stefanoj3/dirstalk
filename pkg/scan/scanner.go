package scan

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/stefanoj3/dirstalk/pkg/common/urlpath"
)

// Target represents the target to scan
type Target struct {
	Path   string
	Method string
	Depth  int
}

// Result represents the result of the scan of a single URL
type Result struct {
	Target Target

	StatusCode int
	URL        url.URL

	Response *http.Response
}

// NewResult creates a new instance of the Result entity based on the Target and Response
func NewResult(target Target, response *http.Response) Result {
	return Result{
		Target:     target,
		StatusCode: response.StatusCode,
		URL:        *response.Request.URL,
	}
}

func NewScanner(
	httpClient Doer,
	producer Producer,
	reproducer ReProducer,
	logger *logrus.Logger,
) *Scanner {
	return &Scanner{
		httpClient: httpClient,
		producer:   producer,
		reproducer: reproducer,
		logger:     logger,
	}
}

type Scanner struct {
	httpClient Doer
	producer   Producer
	reproducer ReProducer
	logger     *logrus.Logger
}

func (s *Scanner) Scan(baseUrl *url.URL, workers int) <-chan Result {
	resultChannel := make(chan Result, workers)

	u := normalizeBaseURL(*baseUrl)

	wg := sync.WaitGroup{}

	producerChannel := s.producer.Produce()
	reproducer := s.reproducer.Reproduce()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for target := range producerChannel {
				s.processTarget(u, target, reproducer, resultChannel)
			}
		}()

	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	return resultChannel
}

func (s *Scanner) processTarget(
	baseURL url.URL,
	target Target,
	reproducer func(r Result) <-chan Target,
	results chan<- Result,
) {
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
		l.WithError(err).Error("failed to perform request")
		return
	}

	if err := res.Body.Close(); err != nil {
		l.WithError(err).Warn("failed to close response body")
	}

	result := NewResult(target, res)
	results <- result

	for newTarget := range reproducer(result) {
		s.processTarget(baseURL, newTarget, reproducer, results)
	}
}

func normalizeBaseURL(baseURL url.URL) url.URL {
	if strings.HasSuffix(baseURL.Path, "/") {
		return baseURL
	}

	baseURL.Path += "/"

	return baseURL
}

func buildURL(baseURL url.URL, target Target) url.URL {
	baseURL.Path = urlpath.Join(baseURL.Path, target.Path)

	return baseURL
}
