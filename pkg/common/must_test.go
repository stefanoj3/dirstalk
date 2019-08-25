package common_test

import (
	"errors"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common"
)

func TestMustShouldNotPanicForNoErr(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("no panic expected")
		}
	}()

	common.Must(nil)
}

func TestMustShouldPanicOnErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("panic expected")
		}
	}()

	common.Must(errors.New("my error"))
}
