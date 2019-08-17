package result_test

import (
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/result"
	"github.com/stretchr/testify/assert"
)

func TestLoadResultsFromFileShouldErrForInvalidPath(t *testing.T) {
	_, err := result.LoadResultsFromFile("/root/123/abc")
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "permission denied")
}
