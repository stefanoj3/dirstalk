package dictionary

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const commentPrefix = "#"

func NewDictionaryFrom(path string, doer Doer) ([]string, error) {
	if strings.HasPrefix(path, "http") {
		return newDictionaryFromRemoteFile(path, doer)
	}

	return newDictionaryFromLocalFile(path)
}

func newDictionaryFromLocalFile(path string) ([]string, error) {
	file, err := os.Open(path) // #nosec
	if err != nil {
		return nil, errors.Wrapf(err, "dictionary: unable to open: %s", path)
	}

	defer file.Close() //nolint

	return dictionaryFromReader(file), nil
}

func dictionaryFromReader(reader io.Reader) []string {
	entries := make([]string, 0)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if isAComment(line) {
			continue
		}

		entries = append(entries, line)
	}

	return entries
}

func newDictionaryFromRemoteFile(path string, doer Doer) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil) //nolint
	if err != nil {
		return nil, errors.Wrapf(err, "dictionary: failed to build request for `%s`", path)
	}

	res, err := doer.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "dictionary: failed to get `%s`", path)
	}

	defer res.Body.Close() //nolint:errcheck

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
