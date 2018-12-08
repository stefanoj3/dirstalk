package cmd_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/cmd"
	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := cmd.NewRootCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	_, out, err := executeCommandC(c)
	assert.NoError(t, err)

	// ensure the summary is printed
	assert.Contains(t, out, "dirstalk is a tool that attempts")
	assert.Contains(t, out, "Usage")
	assert.Contains(t, out, "dictionary.generate")
	assert.Contains(t, out, "scan")
}

func TestScanCommand(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := cmd.NewRootCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	var calls int32
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&calls, 1)
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer srv.Close()

	_, _, err = executeCommandC(c, "scan", srv.URL, "--dictionary", "testdata/dict.txt")
	assert.NoError(t, err)

	assert.Equal(t, int32(3), calls)
}

func TestDictionaryGenerateCommand(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := cmd.NewRootCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	testFilePath := "testdata/" + test.RandStringRunes(10)
	defer removeTestFile(testFilePath)
	_, _, err = executeCommandC(c, "dictionary.generate", ".", "-o", testFilePath)
	assert.NoError(t, err)

	content, err := ioutil.ReadFile(testFilePath)
	assert.NoError(t, err)

	// Ensure the command ran and produced some of the expected output
	// it is not in the scope of this test to ensure the correct output
	assert.Contains(t, string(content), "root_integration_test.go")
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)

	a := []string{""}
	os.Args = append(a, args...)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func removeTestFile(path string) {
	if !strings.Contains(path, "testdata") {
		return
	}

	_ = os.Remove(path)
}
