package dictionary_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/dictionary"
	"github.com/stretchr/testify/assert"
)

func TestDictionaryFromFile(t *testing.T) {
	entries, err := dictionary.NewDictionaryFrom("testdata/dict.txt", &http.Client{})
	assert.NoError(t, err)

	expectedValue := []string{
		"home",
		"home/index.php",
		"blabla",
	}
	assert.Equal(t, expectedValue, entries)
}

func TestShouldFailToCreateDictionaryFromInvalidPath(t *testing.T) {
	_, err := dictionary.NewDictionaryFrom("http:///home/\n", &http.Client{})
	assert.Error(t, err)
}

func TestDictionaryFromAbsolutePath(t *testing.T) {
	path, err := filepath.Abs("testdata/dict.txt")
	assert.NoError(t, err)

	entries, err := dictionary.NewDictionaryFrom(path, &http.Client{})
	assert.NoError(t, err)

	expectedValue := []string{
		"home",
		"home/index.php",
		"blabla",
	}
	assert.Equal(t, expectedValue, entries)
}

func TestDictionaryWithUnableToReadFolderShouldFail(t *testing.T) {
	newFolderPath := "testdata/" + test.RandStringRunes(10)

	err := os.Mkdir(newFolderPath, 0200)
	assert.NoError(t, err)

	defer removeTestDirectory(t, newFolderPath)

	_, err = dictionary.NewDictionaryFrom(newFolderPath, &http.Client{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
}

func TestDictionaryFromFileWithInvalidPath(t *testing.T) {
	t.Parallel()

	d, err := dictionary.NewDictionaryFrom("testdata/gibberish_nonexisting_file", &http.Client{})
	assert.Error(t, err)
	assert.Nil(t, d)

	assert.Contains(t, err.Error(), "unable to open")
}

func TestNewDictionaryFromRemoteFile(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dict := `/home
/about
/contacts
something
potato
`
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(dict)) //nolint:errcheck
		}),
	)
	defer srv.Close()

	entries, err := dictionary.NewDictionaryFrom(srv.URL, &http.Client{})
	assert.NoError(t, err)

	expectedValue := []string{
		"/home",
		"/about",
		"/contacts",
		"something",
		"potato",
	}
	assert.Equal(t, expectedValue, entries)
}

func TestNewDictionaryFromRemoteFileWillReturnErrorWhenRequestTimeout(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond) // out of paranoia - we dont want unstable tests

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("/home")) //nolint:errcheck
		}),
	)
	defer srv.Close()

	entries, err := dictionary.NewDictionaryFrom(
		srv.URL,
		&http.Client{
			Timeout: time.Microsecond,
		},
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get")
	assert.Contains(t, err.Error(), "Timeout")

	assert.Nil(t, entries)
}

func TestNewDictionaryFromRemoteShouldFailWhenRemoteReturnNon200Status(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}),
	)
	defer srv.Close()

	entries, err := dictionary.NewDictionaryFrom(srv.URL, &http.Client{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), srv.URL)
	assert.Contains(t, err.Error(), "status code 403")

	assert.Nil(t, entries)
}

func removeTestDirectory(t *testing.T, path string) {
	if !strings.Contains(path, "testdata") {
		t.Fatalf("cannot delete `%s`, it is not in a `testdata` folder", path)

		return
	}

	stats, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to read `%s` properties", path)
	}

	if !stats.IsDir() {
		t.Fatalf("cannot delete `%s`, it is not a directory", path)
	}

	err = os.Remove(path)
	if err != nil {
		t.Fatalf("failed to remove `%s`: %s", path, err.Error())
	}
}
