package client

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/stefanoj3/dirstalk/pkg/scan/client/cookie"
	"golang.org/x/net/proxy"
)

func NewClientFromConfig(
	timeoutInMilliseconds int,
	socks5Url *url.URL,
	userAgent string,
	useCookieJar bool,
	cookies []*http.Cookie,
	headers map[string]string,
	shouldCacheRequests bool,
	shouldSkipSSLCertificatesValidation bool,
	u *url.URL,
) (*http.Client, error) {
	transport := buildTransport(shouldSkipSSLCertificatesValidation)

	c := &http.Client{
		Timeout:   time.Millisecond * time.Duration(timeoutInMilliseconds),
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if useCookieJar {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, errors.Wrap(err, "NewClientFromConfig: failed to create cookie jar")
		}

		c.Jar = jar
	}

	if c.Jar != nil {
		c.Jar.SetCookies(u, cookies)
	}

	if len(cookies) > 0 && c.Jar == nil {
		c.Jar = cookie.NewStatelessJar(cookies)
	}

	if socks5Url != nil {
		tbDialer, err := proxy.FromURL(socks5Url, proxy.Direct)
		if err != nil {
			return nil, errors.Wrap(err, "NewClientFromConfig: failed to create socks5 proxy")
		}

		transport.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			return tbDialer.Dial(network, addr)
		}
	}

	var err error

	c.Transport, err = decorateTransportWithUserAgentDecorator(c.Transport, userAgent)
	if err != nil {
		return nil, errors.Wrap(err, "NewClientFromConfig: failed to decorate transport")
	}

	if len(headers) > 0 {
		c.Transport, err = decorateTransportWithHeadersDecorator(c.Transport, headers)
		if err != nil {
			return nil, errors.Wrap(err, "NewClientFromConfig: failed to decorate transport")
		}
	}

	if shouldCacheRequests {
		c.Transport, err = decorateTransportWithRequestCacheDecorator(c.Transport)
		if err != nil {
			return nil, errors.Wrap(err, "NewClientFromConfig: failed to decorate transport")
		}
	}

	return c, nil
}

func buildTransport(shouldSkipSSLCertificatesValidation bool) *http.Transport {
	transport := http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if shouldSkipSSLCertificatesValidation {
		//nolint:gosec
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &transport
}
