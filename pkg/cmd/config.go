package cmd

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

func scanConfigFromCmd(cmd *cobra.Command) (*scan.Config, error) {
	c := &scan.Config{}

	var err error

	c.DictionaryPath = cmd.Flag(flagDictionary).Value.String()

	c.HTTPMethods, err = cmd.Flags().GetStringSlice(flagHTTPMethods)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http methods flag")
	}

	c.HTTPStatusesToIgnore, err = cmd.Flags().GetIntSlice(flagHTTPStatusesToIgnore)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http methods flag")
	}

	c.Threads, err = cmd.Flags().GetInt(flagThreads)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read threads flag")
	}

	c.TimeoutInMilliseconds, err = cmd.Flags().GetInt(flagHTTPTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http-timeout flag")
	}

	c.CacheRequests, err = cmd.Flags().GetBool(flagHTTPCacheRequests)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http-cache-requests flag")
	}

	c.ScanDepth, err = cmd.Flags().GetInt(flagScanDepth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http-timeout flag")
	}

	socks5Host := cmd.Flag(flagSocks5Host).Value.String()
	if len(socks5Host) > 0 {
		c.Socks5Url, err = url.Parse("socks5://" + socks5Host)
		if err != nil {
			return nil, errors.Wrap(err, "invalid value for "+flagSocks5Host)
		}
	}

	c.UserAgent = cmd.Flag(flagUserAgent).Value.String()

	c.UseCookieJar, err = cmd.Flags().GetBool(flagCookieJar)
	if err != nil {
		return nil, errors.Wrap(err, "cookie jar flag is invalid")
	}

	rawCookies, err := cmd.Flags().GetStringArray(flagCookie)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read cookies flag")
	}

	c.Cookies, err = rawCookiesToCookies(rawCookies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert rawCookies to objects")
	}

	rawHeaders, err := cmd.Flags().GetStringArray(flagHeader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read cookies flag")
	}

	c.Headers, err = rawHeadersToHeaders(rawHeaders)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert rawHeaders")
	}

	return c, nil
}

func rawHeadersToHeaders(rawHeaders []string) (map[string]string, error) {
	headers := make(map[string]string, len(rawHeaders)*2)

	for _, rawHeader := range rawHeaders {
		parts := strings.Split(rawHeader, ":")
		if len(parts) != 2 {
			return nil, errors.Errorf("header is in invalid format: %s", rawHeader)
		}

		headers[parts[0]] = parts[1]
	}

	return headers, nil
}

func rawCookiesToCookies(rawCookies []string) ([]*http.Cookie, error) {
	cookies := make([]*http.Cookie, 0, len(rawCookies))

	for _, rawCookie := range rawCookies {
		parts := strings.Split(rawCookie, "=")
		if len(parts) != 2 {
			return nil, errors.Errorf("cookie format is invalid: %s", rawCookie)
		}

		cookies = append(
			cookies,
			&http.Cookie{
				Name:    parts[0],
				Value:   parts[1],
				Expires: time.Now().AddDate(0, 0, 2),
			},
		)
	}

	return cookies, nil
}
