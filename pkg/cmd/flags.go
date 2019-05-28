package cmd

const (
	// Root flags
	flagVerbose      = "verbose"
	flagVerboseShort = "v"

	// Scan flags
	flagDictionary      = "dictionary"
	flagDictionaryShort = "d"
	flagHTTPMethods     = "http-methods"
	flagHTTPTimeout     = "http-timeout"
	flagScanDepth       = "scan-depth"
	flagThreads         = "threads"
	flagThreadsShort    = "t"
	flagSocks5Host      = "socks5"
	flagUserAgent       = "user-agent"
	flagCookieJar       = "use-cookie-jar"
	flagCookies         = "cookies"

	// Generate dictionary flags
	flagOutput           = "out"
	flagOutputShort      = "o"
	flagAbsolutePathOnly = "absolute-only"
)
