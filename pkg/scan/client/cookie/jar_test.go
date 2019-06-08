package cookie_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/scan/client/cookie"
	"github.com/stretchr/testify/assert"
)

func TestStatelessJarShouldWorkWithNilCookies(t *testing.T) {
	assert.Nil(t, cookie.NewStatelessJar(nil).Cookies(nil))
}

func TestStatelessJarShouldBeStateless(t *testing.T) {
	cookies := []*http.Cookie{
		{
			Name:  "a_cookie_name",
			Value: "a_cookie_value",
		},
	}

	jar := cookie.NewStatelessJar(cookies)

	u, err := url.Parse("http://github.com/stefanoj3")
	assert.NoError(t, err)

	assert.Equal(t, cookies, jar.Cookies(u))

	jar.SetCookies(
		u,
		[]*http.Cookie{
			{
				Name:  "another_cookie_name",
				Value: "another_cookie_value",
			},
		},
	)

	assert.Equal(t, cookies, jar.Cookies(u))
}
