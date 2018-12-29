package dictionary

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

const commentPrefix = "#"

func NewDictionaryFrom(path string, doer Doer) ([]string, error) {
	_, err := url.ParseRequestURI(path)
	if err != nil {
		return newDictionaryFromLocalFile(path)
	}

	return newDictionaryFromRemoteFile(path, doer)
}

func newDictionaryFromLocalFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "dictionary: unable to open: %s", path)
	}
	defer file.Close()

	return dictionaryFromReader(file), nil
}

func dictionaryFromReader(reader io.Reader) []string {
	entries := make([]string, 0)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if isAComment(line) {
			continue
		}

		entries = append(entries, line)
	}
	return entries
}

func newDictionaryFromRemoteFile(path string, doer Doer) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "dictionary: failed to build request for `%s`", path)
	}

	res, err := doer.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "dictionary: failed to get `%s`", path)
	}
	defer res.Body.Close()

	statusCode := res.StatusCode
	if statusCode > 299 || statusCode < 200 {
		return nil, errors.Errorf(
			"dictionary: failed to retrieve from `%s`, status code %d",
			path,
			statusCode,
		)
	}

	return dictionaryFromReader(res.Body), nil
}

func isAComment(line string) bool {
	return line[0:1] == commentPrefix
}
