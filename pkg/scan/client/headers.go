package client

import (
	"errors"
	"net/http"
)

func decorateTransportWithHeadersDecorator(decorated http.RoundTripper, headers map[string]string) (*headersTransportDecorator, error) {
	if decorated == nil {
		return nil, errors.New("decorated round tripper is nil")
	}

	if headers == nil {
		return nil, errors.New("headers is nil")
	}

	return &headersTransportDecorator{decorated: decorated, headers: headers}, nil
}

type headersTransportDecorator struct {
	decorated http.RoundTripper
	headers   map[string]string
}

func (h *headersTransportDecorator) RoundTrip(r *http.Request) (*http.Response, error) {
	for key, value := range h.headers {
		r.Header.Set(key, value)
	}

	return h.decorated.RoundTrip(r)
}
