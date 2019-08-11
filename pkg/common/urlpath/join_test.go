package urlpath_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/urlpath"
	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	testCases := []struct {
		input          []string
		expectedOutput string
	}{
		{
			input:          []string{"/home"},
			expectedOutput: "/home",
		},
		{
			input:          []string{"/home/"},
			expectedOutput: "/home/",
		},
		{
			input:          []string{"/home/", "test"},
			expectedOutput: "/home/test",
		},
		{
			input:          []string{"/home/", "/test"},
			expectedOutput: "/home/test",
		},
		{
			input:          []string{"/home/", "/test/"},
			expectedOutput: "/home/test/",
		},
		{
			input:          []string{"/home", "/test/"},
			expectedOutput: "/home/test/",
		},
		{
			input:          []string{"/home", "test/"},
			expectedOutput: "/home/test/",
		},
		{
			input:          []string{"/home", "test"},
			expectedOutput: "/home/test",
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint

		scenario := fmt.Sprintf(
			"Input: `%s`, Expected output: `%s`",
			strings.Join(tc.input, ","),
			tc.expectedOutput,
		)

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			output := urlpath.Join(tc.input...)

			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}
