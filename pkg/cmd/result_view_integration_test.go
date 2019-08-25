package cmd_test

import (
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stretchr/testify/assert"
)

func TestResultViewShouldErrWhenCalledWithoutResultFlag(t *testing.T) {
	logger, _ := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "result.view")
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "result-file")
	assert.Contains(t, err.Error(), "not set")
}

func TestResultViewShouldErrWhenCalledWithInvalidPath(t *testing.T) {
	logger, _ := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "result.view", "-r", "/root/123/abc")
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "failed to load results from")
}

func TestResultView(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "result.view", "-r", "testdata/out.txt")
	assert.NoError(t, err)

	expected := `/
├── adview
├── partners
│   └── terms
└── s
`

	assert.Contains(t, loggerBuffer.String(), expected)
}
