package cmd

import (
	"net/url"

	"github.com/chuckpreslar/emission"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

const (
	flagDictionary      = "dictionary"
	flagDictionaryShort = "d"
	flagHTTPMethods     = "http-methods"
	flagHTTPTimeout     = "http-timeout"
	flagScanDepth       = "scan-depth"
	flagThreads         = "threads"
	flagThreadsShort    = "t"
	flagSocks5Host      = "socks5"
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
		flagHTTPMethods,
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
		flagHTTPTimeout,
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

	cmd.Flags().StringP(
		flagSocks5Host,
		"",
		"",
		"socks5 host to use",
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

		eventManager := emission.NewEmitter()
		printer := scan.NewResultLogger(logger)
		eventManager.On(scan.EventResultFound, printer.Log)

		return scan.StartScan(logger, eventManager, cnf, u)
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
