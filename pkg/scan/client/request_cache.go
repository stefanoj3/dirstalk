package client

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

var (
	// ErrRequestRedundant this error is returned when trying to perform the
	// same request (method, host, path) more than one time.
	ErrRequestRedundant = errors.New("this request has been made already")
)

func decorateTransportWithRequestCacheDecorator(decorated http.RoundTripper) (*requestCacheTransportDecorator, error) {
	if decorated == nil {
		return nil, errors.New("decorated round tripper is nil")
	}

	return &requestCacheTransportDecorator{decorated: decorated}, nil
}

type requestCacheTransportDecorator struct {
	decorated  http.RoundTripper
	requestMap sync.Map
}

func (u *requestCacheTransportDecorator) RoundTrip(r *http.Request) (*http.Response, error) {
	key := u.keyForRequest(r)

	_, found := u.requestMap.Load(key)
	if found {
		return nil, ErrRequestRedundant
	}

	u.requestMap.Store(key, struct{}{})

	return u.decorated.RoundTrip(r)
}

func (u *requestCacheTransportDecorator) keyForRequest(r *http.Request) string {
	return fmt.Sprintf("%s~%s~%s", r.Method, r.Host, r.URL.Path)
}
