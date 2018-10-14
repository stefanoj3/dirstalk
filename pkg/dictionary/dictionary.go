package dictionary

import (
	"bufio"
	"os"

	"github.com/pkg/errors"
)

func NewDictionaryFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open: %s", path)
	}
	defer file.Close()

	entries := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}

	return entries, nil
}
