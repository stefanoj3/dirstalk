package output

import "github.com/stefanoj3/dirstalk/pkg/scan"

func NewNullSaver() NullSaver {
	return NullSaver{}
}

type NullSaver struct{}

func (n NullSaver) Save(r scan.Result) error {
	return nil
}

func (n NullSaver) Close() error {
	return nil
}
