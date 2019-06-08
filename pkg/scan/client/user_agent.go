package client

import (
	"errors"
	"net/http"
)

func decorateTransportWithUserAgentDecorator(decorated http.RoundTripper, userAgent string) (*userAgentTransportDecorator, error) {
	if decorated == nil {
		return nil, errors.New("decorated round tripper is nil")
	}

	return &userAgentTransportDecorator{decorated: decorated, userAgent: userAgent}, nil
}

type userAgentTransportDecorator struct {
	decorated http.RoundTripper
	userAgent string
}

func (u *userAgentTransportDecorator) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", u.userAgent)

	return u.decorated.RoundTrip(r)
}
