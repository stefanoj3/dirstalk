package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/cmd"
	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	logger, _ := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c)
	assert.NoError(t, err)
}

func TestVersionCommand(t *testing.T) {
	logger, buf := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "version")
	assert.NoError(t, err)

	// Ensure the command ran and produced some of the expected output
	// it is not in the scope of this test to ensure the correct output
	assert.Contains(t, buf.String(), "Version: ")
}

func executeCommand(root *cobra.Command, args ...string) (err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)

	a := []string{""}
	os.Args = append(a, args...) //nolint

	_, err = root.ExecuteC()

	return err
}

func removeTestFile(path string) {
	if !strings.Contains(path, "testdata") {
		return
	}

	_ = os.Remove(path) //nolint:errcheck
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
