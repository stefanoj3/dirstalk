package common

import "github.com/pkg/errors"

// Must accepts an error in input, if the error is not nil it panics
// It helps to simplify code where no error is expected.
func Must(err error) {
	if err != nil {
		panic(errors.WithStack(err))
	}
}
