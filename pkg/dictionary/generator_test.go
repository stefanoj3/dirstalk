package dictionary_test

import (
	"bytes"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/dictionary"
	"github.com/stretchr/testify/assert"
)

func TestAbsolutePathsGenerator(t *testing.T) {
	t.Parallel()

	b := &bytes.Buffer{}

	dictionaryGenerator := dictionary.NewGenerator(b)

	err := dictionaryGenerator.GenerateDictionaryFrom(
		"./testdata/directory_to_generate_dictionary",
		true,
	)
	assert.NoError(t, err)

	expectedOutput := `testdata/directory_to_generate_dictionary/myfile.php
testdata/directory_to_generate_dictionary/subfolder/image.jpg
testdata/directory_to_generate_dictionary/subfolder/image2.gif
testdata/directory_to_generate_dictionary/subfolder/subsubfolder/myfile.php
testdata/directory_to_generate_dictionary/subfolder/subsubfolder/myfile2.php
`

	assert.Equal(t, expectedOutput, b.String())
}

func TestFilenamePathsGenerator(t *testing.T) {
	t.Parallel()

	b := &bytes.Buffer{}

	dictionaryGenerator := dictionary.NewGenerator(b)

	err := dictionaryGenerator.GenerateDictionaryFrom(
		"testdata/directory_to_generate_dictionary",
		false,
	)
	assert.NoError(t, err)

	expectedOutput := `directory_to_generate_dictionary
myfile.php
subfolder
image.jpg
image2.gif
subsubfolder
myfile2.php
`

	assert.Equal(t, expectedOutput, b.String())
}

func BenchmarkGenerateDictionaryFrom(b *testing.B) {
	buf := &bytes.Buffer{}

	dictionaryGenerator := dictionary.NewGenerator(buf)

	for i := 0; i < b.N; i++ {
		//nolint:errcheck
		_ = dictionaryGenerator.GenerateDictionaryFrom(
			"testdata/directory_to_generate_dictionary",
			false,
		)
	}
}
