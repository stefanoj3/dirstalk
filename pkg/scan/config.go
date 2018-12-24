package scan

import "net/url"

type Config struct {
	DictionaryPath        string
	HTTPMethods           []string
	Threads               int
	TimeoutInMilliseconds int
	ScanDepth             int
	Socks5Url             *url.URL
}
