package scan

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	breakingText = "Found something breaking"
	foundText    = "Found"
	notFoundText = "Not found"
)

func NewResultLogger(logger *logrus.Logger) *ResultLogger {
	return &ResultLogger{logger: logger}
}

type ResultLogger struct {
	logger *logrus.Logger
}

func (c *ResultLogger) Log(result Result) {
	statusCode := result.StatusCode

	l := c.logger.WithFields(logrus.Fields{
		"status-code": statusCode,
		"method":      result.Target.Method,
		"url":         result.URL.String(),
	})

	if statusCode == http.StatusNotFound {
		l.Debug(notFoundText)
	} else if statusCode >= http.StatusInternalServerError {
		l.Warn(breakingText)
	} else {
		l.Info(foundText)
	}
}
