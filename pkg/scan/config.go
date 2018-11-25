package scan

import "net/url"

type Config struct {
	Dictionary            []string
	HTTPMethods           []string
	Threads               int
	TimeoutInMilliseconds int
	ScanDepth             int
	Socks5Url             *url.URL
}
