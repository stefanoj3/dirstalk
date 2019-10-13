package cmd_test

import (
	"io/ioutil"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stretchr/testify/assert"
)

func TestDictionaryGenerateCommand(t *testing.T) {
	logger, _ := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	testFilePath := "testdata/" + test.RandStringRunes(10)
	defer removeTestFile(testFilePath)
	err := executeCommand(c, "dictionary.generate", ".", "-o", testFilePath)
	assert.NoError(t, err)

	//nolint:gosec
	content, err := ioutil.ReadFile(testFilePath)
	assert.NoError(t, err)

	// Ensure the command ran and produced some of the expected output
	// it is not in the scope of this test to ensure the correct output
	assert.Contains(t, string(content), "root_integration_test.go")
}

func TestDictionaryGenerateCommandShouldErrWhenNoTargetIsProvided(t *testing.T) {
	logger, _ := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "dictionary.generate")
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "no path provided")
}

func TestDictionaryGenerateShouldFailWhenAFilePathIsProvidedInsteadOfADirectory(t *testing.T) {
	logger, _ := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	testFilePath := "testdata/" + test.RandStringRunes(10)
	defer removeTestFile(testFilePath)

	err := executeCommand(c, "dictionary.generate", "./root_integration_test.go")
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "the path should be a directory")
}

func TestGenerateDictionaryWithoutOutputPath(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "dictionary.generate", ".")
	assert.NoError(t, err)

	assert.Contains(t, loggerBuffer.String(), "root_integration_test.go")
}

func TestGenerateDictionaryWithInvalidDirectory(t *testing.T) {
	logger, _ := test.NewLogger()

	fakePath := "./" + test.RandStringRunes(10)
	c := createCommand(logger)
	assert.NotNil(t, c)

	err := executeCommand(c, "dictionary.generate", fakePath)
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "unable to use the provided path")
	assert.Contains(t, err.Error(), fakePath)
}
