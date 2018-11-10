package scan

import "net/url"

type Config struct {
	Dictionary            []string
	HttpMethods           []string
	Threads               int
	TimeoutInMilliseconds int
	ScanDepth             int
	Socks5Host            *url.URL
}
