package cmd

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand(logger *logrus.Logger) (*cobra.Command, error) {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "dirstalk",
		Short: "Stalk the given url trying to enumerate files and folders",
		Long:  `dirstalk is a tool that attempts to enumerate files and folders starting from a given URL`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				logger.SetLevel(logrus.DebugLevel)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(
		&verbose,
		flagVerbose,
		flagVerboseShort,
		false,
		"verbose mode",
	)

	scanCmd, err := newScanCommand(logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scan command")
	}

	cmd.AddCommand(scanCmd)
	cmd.AddCommand(newGenerateDictionaryCommand())
	cmd.AddCommand(newVersionCommand(logger.Out))

	return cmd, nil
}
