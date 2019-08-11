package test

import (
	"net/url"
	"testing"
)

func MustParseURL(t *testing.T, rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("failed to parse url: %s", rawurl)
	}

	return u
}
