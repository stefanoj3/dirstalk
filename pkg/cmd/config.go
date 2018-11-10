package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/dictionary"
	"net/url"
)

type Config struct {
	Dictionary []string
	HttpMethods []string
	Threads int
	TimeoutInMilliseconds int
	ScanDepth int
	Socks5Host *url.URL
}

func configFromCmd(cmd *cobra.Command) (*Config, error) {
	c := &Config{}

	var err error

	c.Dictionary, err = dictionary.NewDictionaryFromFile(cmd.Flag(flagDictionary).Value.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate dictionary from file")
	}

	c.HttpMethods, err = cmd.Flags().GetStringSlice(flagHttpMethods)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http methods flag")
	}

	c.Threads, err = cmd.Flags().GetInt(flagThreads)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read threads flag")
	}

	c.TimeoutInMilliseconds, err = cmd.Flags().GetInt(flagHttpTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http-timeout flag")
	}

	c.ScanDepth, err = cmd.Flags().GetInt(flagScanDepth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http-timeout flag")
	}

	socks5Host := cmd.Flag(flagSocks5Host).Value.String()
	if len(socks5Host) > 0 {
		c.Socks5Host, err = url.Parse("socks5://" + socks5Host)
		if err != nil {
			return nil, errors.Wrap(err, "invalid value for " + flagSocks5Host)
		}
	}

	return c, nil
}