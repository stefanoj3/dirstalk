package output

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

var (
	errNilWriteCloser = errors.New("Saver: writeCloser is nil")
)

func NewFileSaver(path string) (Saver, error) {
	file, err := os.Create(path)
	if err != nil {
		return Saver{}, errors.Wrapf(err, "failed to create file `%s` for output", path)
	}

	return Saver{writeCloser: file}, nil
}

type Saver struct {
	writeCloser io.WriteCloser
}

func (f Saver) Save(r scan.Result) error {
	if f.writeCloser == nil {
		return errNilWriteCloser
	}

	rawResult, err := convertResultToRawData(r)
	if err != nil {
		return errors.Wrap(err, "Saver: failed to convert result")
	}

	_, err = fmt.Fprintln(f.writeCloser, string(rawResult))

	return errors.Wrapf(err, "Saver: failed to write result: %s", rawResult)
}

func (f Saver) Close() error {
	if f.writeCloser == nil {
		return errNilWriteCloser
	}

	return f.writeCloser.Close()
}
