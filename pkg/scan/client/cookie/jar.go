package cookie

import (
	"net/http"
	"net/url"
)

func NewStatelessJar(cookies []*http.Cookie) StatelessJar {
	return StatelessJar{cookies: cookies}
}

type StatelessJar struct {
	cookies []*http.Cookie
}

func (s StatelessJar) SetCookies(_ *url.URL, _ []*http.Cookie) {
}

func (s StatelessJar) Cookies(_ *url.URL) []*http.Cookie {
	return s.cookies
}
