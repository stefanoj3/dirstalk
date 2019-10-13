package test

import (
	"net/url"
)

type TestingT interface {
	Fatalf(format string, args ...interface{})
}

func MustParseURL(t TestingT, rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("failed to parse url: %s", rawurl)
	}

	return u
}
