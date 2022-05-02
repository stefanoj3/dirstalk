package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/cmd/termination"
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
		flagScanDictionary,
		flagScanDictionaryShort,
		"",
		"dictionary to use for the scan (path to local file or remote url)",
	)
	common.Must(cmd.MarkFlagFilename(flagScanDictionary))
	common.Must(cmd.MarkFlagRequired(flagScanDictionary))

	cmd.Flags().IntP(
		flagScanDictionaryGetTimeout,
		"",
		50000,
		"timeout in milliseconds (used when fetching remote dictionary)",
	)

	cmd.Flags().StringSlice(
		flagScanHTTPMethods,
		[]string{"GET"},
		"comma separated list of http methods to use; eg: GET,POST,PUT",
	)

	cmd.Flags().IntSlice(
		flagScanHTTPStatusesToIgnore,
		[]int{http.StatusNotFound},
		"comma separated list of http statuses to ignore when showing and processing results; eg: 404,301",
	)

	cmd.Flags().IntP(
		flagScanThreads,
		flagScanThreadsShort,
		3,
		"amount of threads for concurrent requests",
	)

	cmd.Flags().IntP(
		flagScanHTTPTimeout,
		"",
		5000,
		"timeout in milliseconds",
	)

	cmd.Flags().BoolP(
		flagScanHTTPCacheRequests,
		"",
		true,
		"cache requests to avoid performing the same request multiple times within the same scan (EG if the "+
			"server reply with the same redirect location multiple times, dirstalk will follow it only once)",
	)

	cmd.Flags().IntP(
		flagScanScanDepth,
		"",
		3,
		"scan depth",
	)

	cmd.Flags().StringP(
		flagScanSocks5Host,
		"",
		"",
		"socks5 host to use",
	)

	cmd.Flags().StringP(
		flagScanUserAgent,
		"",
		"",
		"user agent to use for http requests",
	)

	cmd.Flags().BoolP(
		flagScanCookieJar,
		"",
		false,
		"enables the use of a cookie jar: it will retain any cookie sent "+
			"from the server and send them for the following requests",
	)

	cmd.Flags().StringArray(
		flagScanCookie,
		[]string{},
		"cookie to add to each request; eg name=value (can be specified multiple times)",
	)

	cmd.Flags().StringArray(
		flagScanHeader,
		[]string{},
		"header to add to each request; eg name=value (can be specified multiple times)",
	)

	cmd.Flags().String(
		flagScanResultOutput,
		"",
		"path where to store result output",
	)

	cmd.Flags().Bool(
		flagShouldSkipSSLCertificatesValidation,
		false,
		"to skip checking the validity of SSL certificates",
	)

	cmd.Flags().Bool(
		flagIgnore20xWithEmptyBody,
		false,
		"ignore HTTP 20x responses with empty body",
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

// startScan is a convenience method that wires together all the dependencies needed to start a scan.
func startScan(logger *logrus.Logger, cnf *scan.Config, u *url.URL) error {
	dict, err := buildDictionary(cnf, u)
	if err != nil {
		return err
	}

	s, err := buildScanner(cnf, dict, u, logger)
	if err != nil {
		return err
	}

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

	defer func() {
		resultSummarizer.Summarize()

		err := outputSaver.Close()
		if err != nil {
			logger.WithError(err).Error("failed to close output file")
		}

		logger.Info("Finished scan")
	}()

	ctx, cancellationFunc := context.WithCancel(context.Background())
	defer cancellationFunc()

	resultsChannel := s.Scan(ctx, u, cnf.Threads)

	terminationHandler := termination.NewTerminationHandler(2)

	for {
		select {
		case <-osSigint:
			terminationHandler.SignalTermination()
			cancellationFunc()

			if terminationHandler.ShouldTerminate() {
				logger.Info("Received sigint, terminating...")

				return nil
			}

			logger.Info(
				"Received sigint, trying to shutdown gracefully, another SIGNINT will terminate the application",
			)
		case result, ok := <-resultsChannel:
			if !ok {
				logger.Debug("result channel is being closed, scan should be complete")

				return nil
			}

			resultSummarizer.Add(result)

			if err := outputSaver.Save(result); err != nil {
				return errors.Wrap(err, "failed to add output to file")
			}
		}
	}
}

func buildScanner(cnf *scan.Config, dict []string, u *url.URL, logger *logrus.Logger) (*scan.Scanner, error) {
	targetProducer := producer.NewDictionaryProducer(cnf.HTTPMethods, dict, cnf.ScanDepth)
	reproducer := producer.NewReProducer(targetProducer)

	resultFilter := filter.NewHTTPStatusResultFilter(cnf.HTTPStatusesToIgnore, cnf.IgnoreEmpty20xResponses)

	scannerClient, err := buildScannerClient(cnf, u)
	if err != nil {
		return nil, err
	}

	s := scan.NewScanner(
		scannerClient,
		targetProducer,
		reproducer,
		resultFilter,
		logger,
	)

	return s, nil
}

func buildDictionary(cnf *scan.Config, u *url.URL) ([]string, error) {
	c, err := buildDictionaryClient(cnf, u)
	if err != nil {
		return nil, err
	}

	dict, err := dictionary.NewDictionaryFrom(cnf.DictionaryPath, c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build dictionary")
	}

	return dict, nil
}

func buildScannerClient(cnf *scan.Config, u *url.URL) (*http.Client, error) {
	c, err := client.NewClientFromConfig(
		cnf.TimeoutInMilliseconds,
		cnf.Socks5Url,
		cnf.UserAgent,
		cnf.UseCookieJar,
		cnf.Cookies,
		cnf.Headers,
		cnf.CacheRequests,
		cnf.ShouldSkipSSLCertificatesValidation,
		u,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build scanner client")
	}

	return c, nil
}

func buildDictionaryClient(cnf *scan.Config, u *url.URL) (*http.Client, error) {
	c, err := client.NewClientFromConfig(
		cnf.DictionaryTimeoutInMilliseconds,
		cnf.Socks5Url,
		cnf.UserAgent,
		cnf.UseCookieJar,
		cnf.Cookies,
		cnf.Headers,
		cnf.CacheRequests,
		cnf.ShouldSkipSSLCertificatesValidation,
		u,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build dictionary client")
	}

	return c, nil
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
