package cmd_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

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

func TestScanCommand(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer testServer.Close()

	_, _, err = executeCommand(c, "scan", testServer.URL, "--dictionary", "testdata/dict.txt", "-v")
	assert.NoError(t, err)

	assert.Equal(t, 3, serverAssertion.Len())
}

func TestScanWithRemoteDictionary(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	dictionaryServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dict := `home
home/index.php
blabla
`
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(dict))
		}),
	)
	defer dictionaryServer.Close()

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer testServer.Close()

	_, _, err = executeCommand(c, "scan", testServer.URL, "--dictionary", dictionaryServer.URL)
	assert.NoError(t, err)

	assert.Equal(t, 3, serverAssertion.Len())
}

func TestScanWithUserAgentFlag(t *testing.T) {
	const testUserAgent = "my_test_user_agent"

	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer testServer.Close()

	_, _, err = executeCommand(
		c,
		"scan",
		testServer.URL,
		"--user-agent",
		testUserAgent,
		"--dictionary",
		"testdata/dict.txt",
	)
	assert.NoError(t, err)

	assert.Equal(t, 3, serverAssertion.Len())
	serverAssertion.Range(func(_ int, r http.Request) {
		assert.Equal(t, testUserAgent, r.Header.Get("User-Agent"))
	})
}

func TestScanWithCookies(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	defer testServer.Close()

	_, _, err = executeCommand(
		c,
		"scan",
		testServer.URL,
		"--cookie",
		"name1=val1",
		"--cookie",
		"name2=val2",
		"--dictionary",
		"testdata/dict.txt",
	)
	assert.NoError(t, err)

	serverAssertion.Range(func(_ int, r http.Request) {
		assert.Equal(t, 2, len(r.Cookies()))

		assert.Equal(t, r.Cookies()[0].Name, "name1")
		assert.Equal(t, r.Cookies()[0].Value, "val1")

		assert.Equal(t, r.Cookies()[1].Name, "name2")
		assert.Equal(t, r.Cookies()[1].Value, "val2")
	})
}

func TestWhenProvidingCookiesInWrongFormatShouldErr(t *testing.T) {
	const malformedCookie = "gibberish"

	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer testServer.Close()

	_, _, err = executeCommand(
		c,
		"scan",
		testServer.URL,
		"--cookie",
		malformedCookie,
		"--dictionary",
		"testdata/dict.txt",
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cookie format is invalid")
	assert.Contains(t, err.Error(), malformedCookie)

	assert.Equal(t, 0, serverAssertion.Len())
}

func TestScanWithCookieJar(t *testing.T) {
	const (
		serverCookieName  = "server_cookie_name"
		serverCookieValue = "server_cookie_value"
	)

	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	once := sync.Once{}
	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			once.Do(func() {
				http.SetCookie(
					w,
					&http.Cookie{
						Name:    serverCookieName,
						Value:   serverCookieValue,
						Expires: time.Now().AddDate(0, 1, 0),
					},
				)
			})
		}),
	)
	defer testServer.Close()

	_, _, err = executeCommand(
		c,
		"scan",
		testServer.URL,
		"--use-cookie-jar",
		"--dictionary",
		"testdata/dict.txt",
		"-t",
		"1",
	)
	assert.NoError(t, err)

	serverAssertion.Range(func(index int, r http.Request) {
		if index == 0 { // first request should have no cookies
			assert.Equal(t, 0, len(r.Cookies()))
			return
		}

		assert.Equal(t, 1, len(r.Cookies()))
		assert.Equal(t, r.Cookies()[0].Name, serverCookieName)
		assert.Equal(t, r.Cookies()[0].Value, serverCookieValue)
	})
}

func TestScanWithUnknownHeaderShouldErr(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	defer testServer.Close()

	_, _, err = executeCommand(
		c,
		"scan",
		testServer.URL,
		"--gibberishflag",
		"--dictionary",
		"testdata/dict.txt",
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown flag")

	assert.Equal(t, 0, serverAssertion.Len())
}

func TestDictionaryGenerateCommand(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	testFilePath := "testdata/" + test.RandStringRunes(10)
	defer removeTestFile(testFilePath)
	_, _, err = executeCommand(c, "dictionary.generate", ".", "-o", testFilePath)
	assert.NoError(t, err)

	content, err := ioutil.ReadFile(testFilePath)
	assert.NoError(t, err)

	// Ensure the command ran and produced some of the expected output
	// it is not in the scope of this test to ensure the correct output
	assert.Contains(t, string(content), "root_integration_test.go")
}

func TestGenerateDictionaryWithoutOutputPath(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	_, _, err = executeCommand(c, "dictionary.generate", ".")
	assert.NoError(t, err)
}

func TestGenerateDictionaryWithInvalidDirectory(t *testing.T) {
	logger, _ := test.NewLogger()

	fakePath := "./" + test.RandStringRunes(10)
	c, err := createCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	_, _, err = executeCommand(c, "dictionary.generate", fakePath)
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "unable to use the provided path")
	assert.Contains(t, err.Error(), fakePath)
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
	dirStalkCmd.AddCommand(cmd.NewGenerateDictionaryCommand())
	dirStalkCmd.AddCommand(cmd.NewVersionCommand(logger.Out))

	return dirStalkCmd, nil
}
