package test

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

func NewLogger() (*logrus.Logger, *bytes.Buffer) {
	b := &bytes.Buffer{}

	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetOutput(b)

	return l, b
}
