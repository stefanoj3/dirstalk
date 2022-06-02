package cmd

const (
	// Root flags.
	flagRootVerbose      = "verbose"
	flagRootVerboseShort = "v"

	// Scan flags.
	flagScanDictionary                      = "dictionary"
	flagScanDictionaryShort                 = "d"
	flagScanDictionaryGetTimeout            = "dictionary-get-timeout"
	flagScanHTTPMethods                     = "http-methods"
	flagScanHTTPStatusesToIgnore            = "http-statuses-to-ignore"
	flagScanHTTPTimeout                     = "http-timeout"
	flagScanHTTPCacheRequests               = "http-cache-requests"
	flagScanScanDepth                       = "scan-depth"
	flagScanThreads                         = "threads"
	flagScanThreadsShort                    = "t"
	flagScanSocks5Host                      = "socks5"
	flagScanUserAgent                       = "user-agent"
	flagScanCookieJar                       = "use-cookie-jar"
	flagScanCookie                          = "cookie"
	flagScanHeader                          = "header"
	flagScanResultOutput                    = "out"
	flagShouldSkipSSLCertificatesValidation = "no-check-certificate"

	flagIgnore20xWithEmptyBody = "ignore-empty-body"

	// Generate dictionary flags.
	flagDictionaryGenerateOutput           = "out"
	flagDictionaryGenerateOutputShort      = "o"
	flagDictionaryGenerateAbsolutePathOnly = "absolute-only"

	// Result view flags.
	flagResultViewResultFile      = "result-file"
	flagResultViewResultFileShort = "r"

	// Result diff flags.
	flagResultDiffFirstFile       = "first"
	flagResultDiffFirstFileShort  = "f"
	flagResultDiffSecondFile      = "second"
	flagResultDiffSecondFileShort = "s"
)
