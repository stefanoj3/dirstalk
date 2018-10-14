package scan

import "github.com/sirupsen/logrus"

const (
	breakingText = "Found something breaking"
	foundText    = "Found"
	notFoundText = "Not Found"
)

type ResultLogger struct {
	logger *logrus.Logger
}

func NewResultLogger(logger *logrus.Logger) *ResultLogger {
	return &ResultLogger{logger: logger}
}

func (c *ResultLogger) Log(result *Result) {
	statusCode := result.Response.StatusCode

	l := c.logger.WithFields(logrus.Fields{
		"status-code": statusCode,
		"method":      result.Target.Method,
		"url":         result.URL,
	})

	if statusCode == 404 {
		l.Debug(notFoundText)
	} else if statusCode >= 500 {
		l.Warn(breakingText)
	} else {
		l.Info(foundText)
	}
}

