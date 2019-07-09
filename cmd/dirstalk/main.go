package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/cmd"
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true}

	dirStalkCmd, err := createCommand(logger)
	if err != nil {
		logger.WithField("err", err).Fatal("Failed to initialize application")
	}

	if err := dirStalkCmd.Execute(); err != nil {
		logger.WithField("err", err).Fatal("Execution error")
	}
}

func createCommand(logger *logrus.Logger) (*cobra.Command, error) {
	dirStalkCmd, err := cmd.NewRootCommand(logger)
	if err != nil {
		return nil, err
	}

	scanCmd, err := cmd.NewScanCommand(logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scan command")
	}

	dirStalkCmd.AddCommand(scanCmd)
	dirStalkCmd.AddCommand(cmd.NewGenerateDictionaryCommand(logger.Out))
	dirStalkCmd.AddCommand(cmd.NewVersionCommand(logger.Out))

	return dirStalkCmd, nil
}
