package result

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

func LoadResultsFromFile(resultFilePath string) ([]scan.Result, error) {
	file, err := os.Open(resultFilePath) // #nosec
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %s", resultFilePath)
	}

	defer file.Close() //nolint

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read properties of %s", resultFilePath)
	}

	if fileInfo.IsDir() {
		return nil, errors.Errorf("`%s` is a directory, you need to specify a valid result file", resultFilePath)
	}

	fileScanner := bufio.NewScanner(file)

	lineCounter := 0
	results := make([]scan.Result, 0, 10)

	for fileScanner.Scan() {
		lineCounter++

		r := scan.Result{}

		if err := json.Unmarshal(fileScanner.Bytes(), &r); err != nil {
			return nil, errors.Wrapf(err, "unable to read line %d", lineCounter)
		}

		results = append(results, r)
	}

	if err := fileScanner.Err(); err != nil {
		return nil, errors.Wrap(err, "an error occurred while reading the result file")
	}

	return results, nil
}
