package path_test

import (
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/path"
	"github.com/stretchr/testify/assert"
)

func TestDictionaryFromFile(t *testing.T) {
	entries, err := path.NewDictionaryFromFile("datafile/dict.txt")
	assert.NoError(t, err)

	expectedValue := []string{
		"home",
		"home/index.php",
		"blabla",
	}
	assert.Equal(t, expectedValue, entries)
}

func TestDictionaryFromFileWithInvalidPath(t *testing.T) {
	d, err := path.NewDictionaryFromFile("datafile/gibberish_nonexisting_file")
	assert.Error(t, err)
	assert.Nil(t, d)
}
