package test

import (
	"bytes"
	"sync"

	"github.com/sirupsen/logrus"
)

type ThreadSafeBuffer struct {
	buf *bytes.Buffer
	rw  sync.RWMutex
}

func (t *ThreadSafeBuffer) String() string {
	t.rw.RLock()
	defer t.rw.RUnlock()

	return t.buf.String()
}

func (t *ThreadSafeBuffer) Write(p []byte) (n int, err error) {
	t.rw.Lock()
	defer t.rw.Unlock()

	return t.buf.Write(p)
}

func NewLogger() (*logrus.Logger, *ThreadSafeBuffer) {
	b := &bytes.Buffer{}
	tsb := &ThreadSafeBuffer{buf: b}

	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetOutput(tsb)

	return l, tsb
}
