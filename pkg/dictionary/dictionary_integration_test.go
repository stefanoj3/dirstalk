package dictionary_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

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

func TestDictionaryFromFileWithInvalidPath(t *testing.T) {
	t.Parallel()

	d, err := dictionary.NewDictionaryFrom("testdata/gibberish_nonexisting_file", &http.Client{})
	assert.Error(t, err)
	assert.Nil(t, d)
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
			_, _ = w.Write([]byte(dict))
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
			_, _ = w.Write([]byte("/home"))
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
