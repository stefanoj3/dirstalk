package scan_test

import (
	"net/http"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stretchr/testify/assert"
)

func TestResultLogger(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario    string
		statusCode  int
		expectedMsg string
	}{
		{
			scenario:    "status code 200, something found",
			statusCode:  200,
			expectedMsg: "msg=Found",
		},
		{
			scenario:    "status code 201, something found",
			statusCode:  201,
			expectedMsg: "msg=Found",
		},
		{
			scenario:    "status code 300, a redirect, something found",
			statusCode:  300,
			expectedMsg: "msg=Found",
		},
		{
			scenario:    "status code 400, bad request",
			statusCode:  400,
			expectedMsg: "msg=Found",
		},
		{
			scenario:    "status code 401, unauthorized, something found",
			statusCode:  401,
			expectedMsg: "msg=Found",
		},
		{
			scenario:    "status code 403, forbidden, something found",
			statusCode:  403,
			expectedMsg: "msg=Found",
		},
		{
			scenario:    "status code 405, method not allowed, something found",
			statusCode:  405,
			expectedMsg: "msg=Found",
		},
		{
			scenario:    "status code 404, not found, nothing found",
			statusCode:  404,
			expectedMsg: `msg="Not found"`,
		},
		{
			scenario:    "status code 500, internal server error, something breaking found",
			statusCode:  500,
			expectedMsg: `msg="Found something breaking"`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			logger, loggerBuffer := test.NewLogger()

			resultLogger := scan.NewResultLogger(logger)
			resultLogger.Log(
				&scan.Result{
					Response: &http.Response{
						StatusCode: tc.statusCode,
					},
				},
			)

			assert.Contains(t, loggerBuffer.String(), tc.expectedMsg)
		})
	}

}
