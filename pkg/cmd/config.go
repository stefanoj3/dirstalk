package cmd

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/dictionary"
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

func scanConfigFromCmd(cmd *cobra.Command) (*scan.Config, error) {
	c := &scan.Config{}

	var err error

	c.Dictionary, err = dictionary.NewDictionaryFromFile(cmd.Flag(flagDictionary).Value.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate dictionary from file")
	}

	c.HTTPMethods, err = cmd.Flags().GetStringSlice(flagHTTPMethods)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http methods flag")
	}

	c.Threads, err = cmd.Flags().GetInt(flagThreads)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read threads flag")
	}

	c.TimeoutInMilliseconds, err = cmd.Flags().GetInt(flagHTTPTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http-timeout flag")
	}

	c.ScanDepth, err = cmd.Flags().GetInt(flagScanDepth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http-timeout flag")
	}

	socks5Host := cmd.Flag(flagSocks5Host).Value.String()
	if len(socks5Host) > 0 {
		c.Socks5Url, err = url.Parse("socks5://" + socks5Host)
		if err != nil {
			return nil, errors.Wrap(err, "invalid value for "+flagSocks5Host)
		}
	}

	return c, nil
}
