package dictionary

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func NewGenerator(out io.Writer) *Generator {
	return &Generator{out: out}
}

type Generator struct {
	out io.Writer
}

func (g *Generator) GenerateDictionaryFrom(path string, absoluteOnly bool) error {
	var (
		dictionary []string
		err        error
	)

	if absoluteOnly {
		dictionary, err = findAbsolutePaths(path)
	} else {
		dictionary, err = findFileNames(path)
	}

	if err != nil {
		return errors.Wrap(err, "failed to generate dictionary")
	}

	for _, entry := range dictionary {
		_, err = fmt.Fprintln(g.out, entry)
		if err != nil {
			return errors.Wrap(err, "failed to write to buffer")
		}
	}

	return nil
}

func findAbsolutePaths(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "findAbsolutePaths: failed to walk")
		}

		if !info.IsDir() {
			files = append(files, p)
		}

		return nil
	})

	return files, err
}

func findFileNames(root string) ([]string, error) {
	var files []string

	filesByKey := make(map[string]bool)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "findFileNames: failed to walk")
		}

		if _, ok := filesByKey[info.Name()]; !ok {
			filesByKey[info.Name()] = true
			files = append(files, info.Name())

			return nil
		}

		return nil
	})

	return files, err
}
