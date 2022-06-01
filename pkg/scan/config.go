package scan

import (
	"net/http"
	"net/url"
)

// Config represents the configuration needed to perform a scan.
type Config struct {
	DictionaryPath                      string
	DictionaryTimeoutInMilliseconds     int
	HTTPMethods                         []string
	HTTPStatusesToIgnore                []int
	Threads                             int
	TimeoutInMilliseconds               int
	CacheRequests                       bool
	ScanDepth                           int
	Socks5Url                           *url.URL
	UserAgent                           string
	UseCookieJar                        bool
	Cookies                             []*http.Cookie
	Headers                             map[string]string
	Out                                 string
	ShouldSkipSSLCertificatesValidation bool
	IgnoreEmpty20xResponses             bool
}
