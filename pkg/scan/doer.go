package scan

import "net/http"

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func newUserAgentDoerDecorator(doer Doer, userAgent string) *userAgentDoerDecorator {
	return &userAgentDoerDecorator{doer: doer, userAgent: userAgent}
}

type userAgentDoerDecorator struct {
	doer      Doer
	userAgent string
}

func (u *userAgentDoerDecorator) Do(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", u.userAgent)

	return u.doer.Do(r)
}
