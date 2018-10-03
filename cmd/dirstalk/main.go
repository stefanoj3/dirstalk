package main

import (
	"github.com/sirupsen/logrus"
	"github.com/stefanoj3/dirstalk/pkg/cmd"
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true}

	dirStalkCmd, err := cmd.NewRootCommand(logger)
	if err != nil {
		logger.WithField("err", err).Fatal("Failed to initialize application")
	}

	if err := dirStalkCmd.Execute(); err != nil {
		logger.WithField("err", err).Fatal("Execution error")
	}
}
