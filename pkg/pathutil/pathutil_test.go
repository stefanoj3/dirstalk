package pathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasExtension(t *testing.T) {
	testCases := []struct {
		path           string
		expectedResult bool
	}{
		{
			path:           "images/image.jpg",
			expectedResult: true,
		},
		{
			path:           "file.pdf",
			expectedResult: true,
		},
		{
			path:           "home/page.php",
			expectedResult: true,
		},
		{
			path:           "src/code.cpp",
			expectedResult: true,
		},
		{
			path:           "src/code.h",
			expectedResult: true,
		},
		{
			path:           "folder/script.sh",
			expectedResult: true,
		},
		{
			path:           "myfile",
			expectedResult: false,
		},
		{
			path:           "myfolder/myfile",
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expectedResult, HasExtension(tc.path))
		})
	}
}
