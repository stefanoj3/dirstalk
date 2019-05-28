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

func (s StatelessJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
}

func (s StatelessJar) Cookies(u *url.URL) []*http.Cookie {
	return s.cookies
}
