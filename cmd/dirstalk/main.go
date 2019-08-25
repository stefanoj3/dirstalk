package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/cmd"
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true}

	dirStalkCmd := createCommand(logger)

	if err := dirStalkCmd.Execute(); err != nil {
		logger.WithField("err", err).Fatal("Execution error")
	}
}

func createCommand(logger *logrus.Logger) *cobra.Command {
	dirStalkCmd := cmd.NewRootCommand(logger)

	dirStalkCmd.AddCommand(cmd.NewScanCommand(logger))
	dirStalkCmd.AddCommand(cmd.NewResultViewCommand(logger.Out))
	dirStalkCmd.AddCommand(cmd.NewResultDiffCommand(logger.Out))
	dirStalkCmd.AddCommand(cmd.NewGenerateDictionaryCommand(logger.Out))
	dirStalkCmd.AddCommand(cmd.NewVersionCommand(logger.Out))

	return dirStalkCmd
}
