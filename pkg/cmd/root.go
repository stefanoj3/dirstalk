package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand(logger *logrus.Logger) *cobra.Command {
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
		flagRootVerbose,
		flagRootVerboseShort,
		false,
		"verbose mode",
	)

	return cmd
}
