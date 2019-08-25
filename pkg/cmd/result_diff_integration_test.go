package cmd_test

import (
	"fmt"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stretchr/testify/assert"
)

func TestNewResultDiff(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "result.diff", "-f", "testdata/out.txt", "-s", "testdata/out2.txt")
	assert.NoError(t, err)

	// to keep compatibility with other systems open, the language should take care to use the correct newline symbol
	newlineSymbol := fmt.Sprintln()

	expected := "/" + newlineSymbol +
		"├── adview" + newlineSymbol +
		"├── partners" + newlineSymbol +
		"│   └── \x1b[31mterms\x1b[0m\x1b[32m123\x1b[0m" + newlineSymbol +
		"└── s"

	assert.Contains(t, loggerBuffer.String(), expected)
}

func TestNewResultDiffShouldErrWithInvalidFirstFile(t *testing.T) {
	logger, _ := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "result.diff", "-f", "/root/123/bla", "-s", "testdata/out2.txt")
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "/root/123/bla")
}

func TestNewResultDiffShouldErrWithInvalidSecondFile(t *testing.T) {
	logger, _ := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "result.diff", "-f", "testdata/out2.txt", "-s", "/root/123/bla")
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "/root/123/bla")
}

func TestDiffForSameFileShouldErr(t *testing.T) {
	logger, _ := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "result.diff", "-f", "testdata/out.txt", "-s", "testdata/out.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no diffs found")
}
