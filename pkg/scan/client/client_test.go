package client_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan/client"
	"github.com/stretchr/testify/assert"
)

func TestWhenRemoteIsTooSlowClientShouldTimeout(t *testing.T) {
	testServer, _ := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * 100)
		}),
	)
	defer testServer.Close()

	c, err := client.NewClientFromConfig(
		10,
		nil,
		"",
		false,
		nil,
		nil,
		true,
		false,
		nil,
	)
	assert.NoError(t, err)

	res, err := c.Get(testServer.URL) //nolint
	assert.Error(t, err)
	assert.Nil(t, res)

	assert.Contains(t, err.Error(), "exceeded")
}

func TestShouldForwardProvidedCookiesWhenUsingJar(t *testing.T) {
	const (
		serverCookieName  = "server_cookie_name"
		serverCookieValue = "server_cookie_value"
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(
				w,
				&http.Cookie{
					Name:    serverCookieName,
					Value:   serverCookieValue,
					Expires: time.Now().AddDate(0, 1, 0),
				},
			)
		}),
	)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	cookies := []*http.Cookie{
		{
			Name:  "a_cookie_name",
			Value: "a_cookie_value",
		},
	}

	c, err := client.NewClientFromConfig(
		100,
		nil,
		"",
		true,
		cookies,
		map[string]string{},
		false,
		false,
		u,
	)
	assert.NoError(t, err)

	res, err := c.Get(testServer.URL)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	defer res.Body.Close() //nolint:errcheck

	assert.Equal(t, 1, serverAssertion.Len())

	serverAssertion.At(0, func(r http.Request) {
		assert.Equal(t, 1, len(r.Cookies()))

		assert.Equal(t, r.Cookies()[0].Name, cookies[0].Name)
		assert.Equal(t, r.Cookies()[0].Value, cookies[0].Value)
		assert.Equal(t, r.Cookies()[0].Expires, cookies[0].Expires)
	})

	res, err = c.Get(testServer.URL)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	defer res.Body.Close() //nolint:errcheck

	assert.Equal(t, 2, serverAssertion.Len())

	serverAssertion.At(1, func(r http.Request) {
		assert.Equal(t, 2, len(r.Cookies()))

		assert.Equal(t, r.Cookies()[0].Name, cookies[0].Name)
		assert.Equal(t, r.Cookies()[0].Value, cookies[0].Value)
		assert.Equal(t, r.Cookies()[0].Expires, cookies[0].Expires)

		assert.Equal(t, r.Cookies()[1].Name, serverCookieName)
		assert.Equal(t, r.Cookies()[1].Value, serverCookieValue)
	})
}

func TestShouldForwardCookiesWhenJarIsDisabled(t *testing.T) {
	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	cookies := []*http.Cookie{
		{
			Name:  "a_cookie_name",
			Value: "a_cookie_value",
		},
	}

	c, err := client.NewClientFromConfig(
		100,
		nil,
		"",
		false,
		cookies,
		map[string]string{},
		true,
		false,
		u,
	)
	assert.NoError(t, err)

	res, err := c.Get(testServer.URL)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	defer res.Body.Close() //nolint:errcheck

	assert.Equal(t, 1, serverAssertion.Len())

	serverAssertion.At(0, func(r http.Request) {
		assert.Equal(t, 1, len(r.Cookies()))

		assert.Equal(t, r.Cookies()[0].Name, cookies[0].Name)
		assert.Equal(t, r.Cookies()[0].Value, cookies[0].Value)
		assert.Equal(t, r.Cookies()[0].Expires, cookies[0].Expires)
	})
}

func TestShouldForwardProvidedHeader(t *testing.T) {
	const (
		headerName  = "my_header_name"
		headerValue = "my_header_value_123"
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	c, err := client.NewClientFromConfig(
		100,
		nil,
		"",
		false,
		nil,
		map[string]string{headerName: headerValue},
		true,
		false,
		u,
	)
	assert.NoError(t, err)

	res, err := c.Get(testServer.URL)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	defer res.Body.Close() //nolint:errcheck

	assert.Equal(t, 1, serverAssertion.Len())

	serverAssertion.At(0, func(r http.Request) {
		assert.Equal(t, headerValue, r.Header.Get(headerName))
	})
}

func TestShouldFailToCreateAClientWithInvalidSocks5Url(t *testing.T) {
	u := url.URL{Scheme: "potatoscheme"}

	c, err := client.NewClientFromConfig(
		100,
		&u,
		"",
		false,
		nil,
		map[string]string{},
		true,
		false,
		nil,
	)
	assert.Nil(t, c)
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "unknown scheme")
}

func TestShouldNotRepeatTheSameRequestTwice(t *testing.T) {
	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	c, err := client.NewClientFromConfig(
		100,
		nil,
		"",
		false,
		nil,
		nil,
		true,
		false,
		u,
	)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	assert.NoError(t, err)

	res, err := c.Do(req)
	assert.NoError(t, err)

	res.Body.Close() //nolint:errcheck,gosec

	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = c.Do(req) //nolint
	assert.Contains(t, err.Error(), client.ErrRequestRedundant.Error())
	assert.Nil(t, res)

	assert.Equal(t, 1, serverAssertion.Len())
}

func TestShouldFailToCommunicateWithServerHavingInvalidSSLCertificates(t *testing.T) {
	testServer, serverAssertion := test.NewTSLServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	c, err := client.NewClientFromConfig(
		1500,
		nil,
		"",
		false,
		nil,
		nil,
		true,
		false,
		u,
	)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	assert.NoError(t, err)

	res, err := c.Do(req) //nolint:bodyclose
	assert.Error(t, err)
	assert.Nil(t, res)

	assert.Contains(t, err.Error(), "certificate")

	// the request should NOT hit the handler
	assert.Equal(t, 0, serverAssertion.Len())
}

func TestShouldBeAbleToSkipSSLCertificatesCheck(t *testing.T) {
	testServer, serverAssertion := test.NewTSLServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	assert.NoError(t, err)

	c, err := client.NewClientFromConfig(
		1500,
		nil,
		"",
		false,
		nil,
		nil,
		true,
		true,
		u,
	)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	assert.NoError(t, err)

	res, err := c.Do(req)
	assert.NoError(t, err)

	res.Body.Close() //nolint:errcheck,gosec

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	// the request should hit the handler
	assert.Equal(t, 1, serverAssertion.Len())
}
