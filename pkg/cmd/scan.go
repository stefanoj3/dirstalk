package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/common"
	"github.com/stefanoj3/dirstalk/pkg/dictionary"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/client"
	"github.com/stefanoj3/dirstalk/pkg/scan/filter"
	"github.com/stefanoj3/dirstalk/pkg/scan/output"
	"github.com/stefanoj3/dirstalk/pkg/scan/producer"
	"github.com/stefanoj3/dirstalk/pkg/scan/summarizer"
	"github.com/stefanoj3/dirstalk/pkg/scan/summarizer/tree"
)

func NewScanCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [url]",
		Short: "Scan the given URL",
		RunE:  buildScanFunction(logger),
	}

	cmd.Flags().StringP(
		flagDictionary,
		flagDictionaryShort,
		"",
		"dictionary to use for the scan (path to local file or remote url)",
	)
	common.Must(cmd.MarkFlagFilename(flagDictionary))
	common.Must(cmd.MarkFlagRequired(flagDictionary))

	cmd.Flags().StringSlice(
		flagHTTPMethods,
		[]string{"GET"},
		"comma separated list of http methods to use; eg: GET,POST,PUT",
	)

	cmd.Flags().IntSlice(
		flagHTTPStatusesToIgnore,
		[]int{http.StatusNotFound},
		"comma separated list of http statuses to ignore when showing and processing results; eg: 404,301",
	)

	cmd.Flags().IntP(
		flagThreads,
		flagThreadsShort,
		3,
		"amount of threads for concurrent requests",
	)

	cmd.Flags().IntP(
		flagHTTPTimeout,
		"",
		5000,
		"timeout in milliseconds",
	)

	cmd.Flags().BoolP(
		flagHTTPCacheRequests,
		"",
		true,
		"cache requests to avoid performing the same request multiple times within the same scan (EG if the "+
			"server reply with the same redirect location multiple times, dirstalk will follow it only once)",
	)

	cmd.Flags().IntP(
		flagScanDepth,
		"",
		3,
		"scan depth",
	)

	cmd.Flags().StringP(
		flagSocks5Host,
		"",
		"",
		"socks5 host to use",
	)

	cmd.Flags().StringP(
		flagUserAgent,
		"",
		"",
		"user agent to use for http requests",
	)

	cmd.Flags().BoolP(
		flagCookieJar,
		"",
		false,
		"enables the use of a cookie jar: it will retain any cookie sent "+
			"from the server and send them for the following requests",
	)

	cmd.Flags().StringArray(
		flagCookie,
		[]string{},
		"cookie to add to each request; eg name=value (can be specified multiple times)",
	)

	cmd.Flags().StringArray(
		flagHeader,
		[]string{},
		"header to add to each request; eg name=value (can be specified multiple times)",
	)

	cmd.Flags().String(
		flagResultOutput,
		"",
		"path where to store result output",
	)

	return cmd
}

func buildScanFunction(logger *logrus.Logger) func(cmd *cobra.Command, args []string) error {
	f := func(cmd *cobra.Command, args []string) error {
		u, err := getURL(args)
		if err != nil {
			return err
		}

		cnf, err := scanConfigFromCmd(cmd)
		if err != nil {
			return errors.Wrap(err, "failed to build config")
		}

		return startScan(logger, cnf, u)
	}

	return f
}

func getURL(args []string) (*url.URL, error) {
	if len(args) == 0 {
		return nil, errors.New("no URL provided")
	}

	arg := args[0]

	u, err := url.ParseRequestURI(arg)
	if err != nil {
		return nil, errors.Wrap(err, "the first argument must be a valid url")
	}

	return u, nil
}

// startScan is a convenience method that wires together all the dependencies needed to start a scan
func startScan(logger *logrus.Logger, cnf *scan.Config, u *url.URL) error {
	c, err := client.NewClientFromConfig(
		cnf.TimeoutInMilliseconds,
		cnf.Socks5Url,
		cnf.UserAgent,
		cnf.UseCookieJar,
		cnf.Cookies,
		cnf.Headers,
		cnf.CacheRequests,
		u,
	)
	if err != nil {
		return errors.Wrap(err, "failed to build client")
	}

	dict, err := dictionary.NewDictionaryFrom(cnf.DictionaryPath, c)
	if err != nil {
		return errors.Wrap(err, "failed to build dictionary")
	}

	targetProducer := producer.NewDictionaryProducer(cnf.HTTPMethods, dict, cnf.ScanDepth)
	reproducer := producer.NewReProducer(targetProducer)

	resultFilter := filter.NewHTTPStatusResultFilter(cnf.HTTPStatusesToIgnore)

	s := scan.NewScanner(
		c,
		targetProducer,
		reproducer,
		resultFilter,
		logger,
	)

	logger.WithFields(logrus.Fields{
		"url":               u.String(),
		"threads":           cnf.Threads,
		"dictionary-length": len(dict),
		"scan-depth":        cnf.ScanDepth,
		"timeout":           cnf.TimeoutInMilliseconds,
		"socks5":            cnf.Socks5Url,
		"cookies":           stringifyCookies(cnf.Cookies),
		"cookie-jar":        cnf.UseCookieJar,
		"headers":           stringifyHeaders(cnf.Headers),
		"user-agent":        cnf.UserAgent,
	}).Info("Starting scan")

	resultSummarizer := summarizer.NewResultSummarizer(tree.NewResultTreeProducer(), logger)

	osSigint := make(chan os.Signal, 1)
	signal.Notify(osSigint, os.Interrupt)

	outputSaver, err := newOutputSaver(cnf.Out)
	if err != nil {
		return errors.Wrap(err, "failed to create output saver")
	}

	finishFunc := func() {
		resultSummarizer.Summarize()
		err := outputSaver.Close()
		if err != nil {
			logger.WithError(err).Error("failed to close output file")
		}
		logger.Info("Finished scan")
	}

	for result := range s.Scan(u, cnf.Threads) {
		select {
		case <-osSigint:
			logger.Info("Received sigint, terminating...")
			finishFunc()
			return nil
		default:
			resultSummarizer.Add(result)
			err := outputSaver.Save(result)
			if err != nil {
				return errors.Wrap(err, "failed to add output to file")
			}
		}
	}

	finishFunc()

	return nil
}

func newOutputSaver(path string) (OutputSaver, error) {
	if path == "" {
		return output.NewNullSaver(), nil
	}

	return output.NewFileSaver(path)
}

func stringifyCookies(cookies []*http.Cookie) string {
	result := ""

	for _, cookie := range cookies {
		result += fmt.Sprintf("{%s=%s}", cookie.Name, cookie.Value)
	}

	return result
}

func stringifyHeaders(headers map[string]string) string {
	result := ""

	for name, value := range headers {
		result += fmt.Sprintf("{%s:%s}", name, value)
	}

	return result
}
