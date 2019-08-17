package cmd

const (
	// Root flags
	flagVerbose      = "verbose"
	flagVerboseShort = "v"

	// Scan flags
	flagDictionary           = "dictionary"
	flagDictionaryShort      = "d"
	flagHTTPMethods          = "http-methods"
	flagHTTPStatusesToIgnore = "http-statuses-to-ignore"
	flagHTTPTimeout          = "http-timeout"
	flagHTTPCacheRequests    = "http-cache-requests"
	flagScanDepth            = "scan-depth"
	flagThreads              = "threads"
	flagThreadsShort         = "t"
	flagSocks5Host           = "socks5"
	flagUserAgent            = "user-agent"
	flagCookieJar            = "use-cookie-jar"
	flagCookie               = "cookie"
	flagHeader               = "header"
	flagResultOutput         = "out"

	// Generate dictionary flags
	flagOutput           = "out"
	flagOutputShort      = "o"
	flagAbsolutePathOnly = "absolute-only"

	// Result view flags
	flagResultFile      = "result-file"
	flagResultFileShort = "r"
)
