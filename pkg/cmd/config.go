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

const failedToReadPropertyError = "failed to read %s"

func scanConfigFromCmd(cmd *cobra.Command) (*scan.Config, error) {
	c := &scan.Config{}

	var err error

	c.DictionaryPath = cmd.Flag(flagScanDictionary).Value.String()

	if c.DictionaryTimeoutInMilliseconds, err = cmd.Flags().GetInt(flagScanDictionaryGetTimeout); err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanDictionaryGetTimeout)
	}

	if c.HTTPMethods, err = cmd.Flags().GetStringSlice(flagScanHTTPMethods); err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanHTTPMethods)
	}

	if c.HTTPStatusesToIgnore, err = cmd.Flags().GetIntSlice(flagScanHTTPStatusesToIgnore); err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanHTTPStatusesToIgnore)
	}

	if c.Threads, err = cmd.Flags().GetInt(flagScanThreads); err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanThreads)
	}

	if c.TimeoutInMilliseconds, err = cmd.Flags().GetInt(flagScanHTTPTimeout); err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanHTTPTimeout)
	}

	if c.CacheRequests, err = cmd.Flags().GetBool(flagScanHTTPCacheRequests); err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanHTTPCacheRequests)
	}

	if c.ScanDepth, err = cmd.Flags().GetInt(flagScanScanDepth); err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanScanDepth)
	}

	socks5Host := cmd.Flag(flagScanSocks5Host).Value.String()
	if len(socks5Host) > 0 {
		if c.Socks5Url, err = url.Parse("socks5://" + socks5Host); err != nil {
			return nil, errors.Wrapf(err, "invalid value for %s", flagScanSocks5Host)
		}
	}

	c.UserAgent = cmd.Flag(flagScanUserAgent).Value.String()

	if c.UseCookieJar, err = cmd.Flags().GetBool(flagScanCookieJar); err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanCookieJar)
	}

	rawCookies, err := cmd.Flags().GetStringArray(flagScanCookie)
	if err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanCookie)
	}

	if c.Cookies, err = rawCookiesToCookies(rawCookies); err != nil {
		return nil, errors.Wrap(err, "failed to convert rawCookies to objects")
	}

	rawHeaders, err := cmd.Flags().GetStringArray(flagScanHeader)
	if err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagScanHeader)
	}

	if c.Headers, err = rawHeadersToHeaders(rawHeaders); err != nil {
		return nil, errors.Wrapf(err, "failed to convert rawHeaders (%v)", rawHeaders)
	}

	c.Out = cmd.Flag(flagScanResultOutput).Value.String()

	c.ShouldSkipSSLCertificatesValidation, err = cmd.Flags().GetBool(flagShouldSkipSSLCertificatesValidation)
	if err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagShouldSkipSSLCertificatesValidation)
	}

	c.IgnoreEmpty20xResponses, err = cmd.Flags().GetBool(flagIgnore20xWithEmptyBody)
	if err != nil {
		return nil, errors.Wrapf(err, failedToReadPropertyError, flagIgnore20xWithEmptyBody)
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
