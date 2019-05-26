package scan

import "net/url"

// Config represents the configuration needed to perform a scan
type Config struct {
	DictionaryPath        string
	HTTPMethods           []string
	Threads               int
	TimeoutInMilliseconds int
	ScanDepth             int
	Socks5Url             *url.URL
	UserAgent             string
	UseCookieJar          bool
}
