package output_test

import (
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/output"
	"github.com/stretchr/testify/assert"
)

func TestNullSaver(t *testing.T) {
	t.Parallel()

	sut := output.NewNullSaver()

	assert.NoError(t, sut.Save(scan.Result{}))
	assert.NoError(t, sut.Close())
	assert.NoError(t, sut.Save(scan.Result{}))
}
