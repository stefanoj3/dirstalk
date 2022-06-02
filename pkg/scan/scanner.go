package scan

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/stefanoj3/dirstalk/pkg/common/urlpath"
	"github.com/stefanoj3/dirstalk/pkg/scan/client"
)

// Target represents the target to scan.
type Target struct {
	Path   string
	Method string
	Depth  int
}

// Result represents the result of the scan of a single URL.
type Result struct {
	Target        Target
	StatusCode    int
	URL           url.URL
	ContentLength int64
}

// NewResult creates a new instance of the Result entity based on the Target and Response.
func NewResult(target Target, response *http.Response) Result {
	return Result{
		Target:        target,
		StatusCode:    response.StatusCode,
		URL:           *response.Request.URL,
		ContentLength: response.ContentLength,
	}
}

func NewScanner(
	httpClient Doer,
	producer Producer,
	reproducer ReProducer,
	resultFilter ResultFilter,
	logger *logrus.Logger,
) *Scanner {
	return &Scanner{
		httpClient:   httpClient,
		producer:     producer,
		reproducer:   reproducer,
		resultFilter: resultFilter,
		logger:       logger,
	}
}

type Scanner struct {
	httpClient   Doer
	producer     Producer
	reproducer   ReProducer
	resultFilter ResultFilter
	logger       *logrus.Logger
}

func (s *Scanner) Scan(ctx context.Context, baseURL *url.URL, workers int) <-chan Result {
	resultChannel := make(chan Result, workers)

	u := normalizeBaseURL(*baseURL)

	wg := sync.WaitGroup{}

	producerChannel := s.producer.Produce(ctx)
	reproducer := s.reproducer.Reproduce(ctx)

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					s.logger.Debug("terminating worker: context cancellation")
				case target, ok := <-producerChannel:
					if !ok {
						s.logger.Debug("terminating worker: producer channel closed")

						return
					}

					s.processTarget(ctx, u, target, reproducer, resultChannel)
				}
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
	ctx context.Context,
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

	req, err := http.NewRequestWithContext(ctx, target.Method, u.String(), nil)
	if err != nil {
		l.WithError(err).Error("failed to build request")

		return
	}

	s.processRequest(ctx, l, req, target, results, reproducer, baseURL)
}

func (s *Scanner) processRequest(
	ctx context.Context,
	l *logrus.Entry,
	req *http.Request,
	target Target,
	results chan<- Result,
	reproducer func(r Result) <-chan Target,
	baseURL url.URL,
) {
	res, err := s.httpClient.Do(req)
	if err != nil && strings.Contains(err.Error(), client.ErrRequestRedundant.Error()) {
		l.WithError(err).Debug("skipping, request was already made")

		return
	}

	if err != nil {
		l.WithError(err).Error("failed to perform request")

		return
	}

	if err := res.Body.Close(); err != nil {
		l.WithError(err).Warn("failed to close response body")
	}

	result := NewResult(target, res)

	if s.resultFilter.ShouldIgnore(result) {
		return
	}

	results <- result

	redirectTarget, shouldRedirect := s.shouldRedirect(l, req, res, target.Depth)
	if shouldRedirect {
		s.processTarget(ctx, baseURL, redirectTarget, reproducer, results)
	}

	for newTarget := range reproducer(result) {
		s.processTarget(ctx, baseURL, newTarget, reproducer, results)
	}
}

func (s *Scanner) shouldRedirect(l *logrus.Entry, req *http.Request, res *http.Response, targetDepth int) (Target, bool) {
	if targetDepth == 0 {
		l.Debug("depth is 0, not following any redirect")

		return Target{}, false
	}

	redirectMethod := req.Method
	location := res.Header.Get("Location")

	if location == "" {
		return Target{}, false
	}

	redirectStatusCodes := map[int]bool{
		http.StatusMovedPermanently:  true,
		http.StatusFound:             true,
		http.StatusSeeOther:          true,
		http.StatusTemporaryRedirect: false,
		http.StatusPermanentRedirect: false,
	}

	shouldOverrideRequestMethod, shouldRedirect := redirectStatusCodes[res.StatusCode]
	if !shouldRedirect {
		return Target{}, false
	}

	// RFC 2616 allowed automatic redirection only with GET and
	// HEAD requests. RFC 7231 lifts this restriction, but we still
	// restrict other methods to GET to maintain compatibility.
	// See Issue 18570.
	if shouldOverrideRequestMethod {
		if req.Method != "GET" && req.Method != "HEAD" {
			redirectMethod = "GET"
		}
	}

	u, err := url.Parse(location)
	if err != nil {
		l.WithError(err).
			WithField("location", location).
			Warn("failed to parse location for redirect")

		return Target{}, false
	}

	if u.Host != "" && u.Host != req.Host {
		l.Debug("skipping redirect, pointing to a different host")

		return Target{}, false
	}

	return Target{
		Path:   u.Path,
		Method: redirectMethod,
		Depth:  targetDepth - 1,
	}, true
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
