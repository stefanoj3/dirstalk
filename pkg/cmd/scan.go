package cmd

import (
	"net/http"
	"net/url"
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/dictionary"
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

const (
	flagDictionary      = "dictionary"
	flagDictionaryShort = "d"
	flagHttpMethods     = "http-methods"
	flagHttpTimeout     = "http-timeout"
	flagScanDepth       = "scan-depth"
	flagThreads         = "threads"
	flagThreadsShort    = "t"
)

func newScanCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [url]",
		Short: "Scan the given URL",
		RunE:  buildScanFunction(logger),
	}

	cmd.Flags().StringP(
		flagDictionary,
		flagDictionaryShort,
		"",
		"dictionary to use for the scan",
	)
	cmd.MarkFlagFilename(flagDictionary)
	cmd.MarkFlagRequired(flagDictionary)

	cmd.Flags().StringSlice(
		flagHttpMethods,
		[]string{"GET"},
		"comma separated list of http methods to use; eg: GET,POST,PUT",
	)

	cmd.Flags().IntP(
		flagThreads,
		flagThreadsShort,
		3,
		"amount of threads for concurrent requests",
	)

	cmd.Flags().IntP(
		flagHttpTimeout,
		"",
		5000,
		"timeout in milliseconds",
	)

	cmd.Flags().IntP(
		flagScanDepth,
		"",
		3,
		"scan depth",
	)

	return cmd
}

func buildScanFunction(logger *logrus.Logger) func(cmd *cobra.Command, args []string) error {
	f := func(cmd *cobra.Command, args []string) error {
		u, err := getUrl(args)
		if err != nil {
			return err
		}

		dict, err := dictionary.NewDictionaryFromFile(cmd.Flag(flagDictionary).Value.String())
		if err != nil {
			return errors.Wrap(err, "failed to generate dictionary from file")
		}

		httpMethods, err := cmd.Flags().GetStringSlice(flagHttpMethods)
		if err != nil {
			return errors.Wrap(err, "failed to read http methods flag")
		}

		threads, err := cmd.Flags().GetInt(flagThreads)
		if err != nil {
			return errors.Wrap(err, "failed to read threads flag")
		}

		timeoutInMilliseconds, err := cmd.Flags().GetInt(flagHttpTimeout)
		if err != nil {
			return errors.Wrap(err, "failed to read http-timeout flag")
		}

		scanDepth, err := cmd.Flags().GetInt(flagScanDepth)
		if err != nil {
			return errors.Wrap(err, "failed to read http-timeout flag")
		}

		eventManager := emission.NewEmitter()

		printer := scan.NewResultLogger(logger)
		eventManager.On(scan.EventResultFound, printer.Log)

		r := scan.ReProcessor{}
		eventManager.On(scan.EventResultFound, r.Process)

		s := scan.NewScanner(
			&http.Client{
				Timeout: time.Millisecond * time.Duration(timeoutInMilliseconds),
			},
			eventManager,
			logger,
		)

		go scan.NewTargetProducer(eventManager, httpMethods, dict, scanDepth).Run()

		eventManager.AddListener(scan.EventTargetProduced, s.AddTarget)
		eventManager.AddListener(scan.EventProducerFinished, s.Release)

		logger.WithFields(logrus.Fields{
			"url":               u.String(),
			"threads":           threads,
			"dictionary.length": len(dict),
		}).Info("Starting scan")

		s.Scan(*u, threads)

		logger.Info("Finished scan")

		return nil
	}

	return f
}

func getUrl(args []string) (*url.URL, error) {
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
