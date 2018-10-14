package cmd_test

import (
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/cmd"
	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	logger, _ := test.NewLogger()

	c, err := cmd.NewRootCommand(logger)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}
