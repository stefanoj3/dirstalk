package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/cmd"
	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	_, out, err := executeCommand(c)
	assert.NoError(t, err)

	// ensure the summary is printed
	assert.Contains(t, out, "dirstalk is a tool that attempts")
	assert.Contains(t, out, "Usage")
	assert.Contains(t, out, "dictionary.generate")
	assert.Contains(t, out, "scan")
}

func TestVersionCommand(t *testing.T) {
	logger, buf := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	_, _, err = executeCommand(c, "version")
	assert.NoError(t, err)

	// Ensure the command ran and produced some of the expected output
	// it is not in the scope of this test to ensure the correct output
	assert.Contains(t, buf.String(), "Version: ")
}

func executeCommand(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)

	a := []string{""}
	os.Args = append(a, args...) //nolint

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func removeTestFile(path string) {
	if !strings.Contains(path, "testdata") {
		return
	}

	_ = os.Remove(path)
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
