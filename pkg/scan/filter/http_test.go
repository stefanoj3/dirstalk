package filter_test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/filter"
	"github.com/stretchr/testify/assert"
)

func TestHTTPStatusResultFilter(t *testing.T) {
	testCases := []struct {
		statusCodesToIgnore []int
		result              scan.Result
		expectedResult      bool
	}{
		{
			statusCodesToIgnore: []int{http.StatusCreated, http.StatusNotFound},
			result:              scan.Result{StatusCode: http.StatusOK},
			expectedResult:      false,
		},
		{
			statusCodesToIgnore: []int{http.StatusCreated, http.StatusNotFound},
			result:              scan.Result{StatusCode: http.StatusNotFound},
			expectedResult:      true,
		},
		{
			statusCodesToIgnore: []int{},
			result:              scan.Result{StatusCode: http.StatusNotFound},
			expectedResult:      false,
		},
		{
			statusCodesToIgnore: []int{},
			result:              scan.Result{StatusCode: http.StatusOK},
			expectedResult:      false,
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint

		scenario := fmt.Sprintf("ignored: %v, result: %d", tc.statusCodesToIgnore, tc.result.StatusCode)

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			actual := filter.NewHTTPStatusResultFilter(tc.statusCodesToIgnore, false).ShouldIgnore(tc.result)
			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestHTTPStatusResultFilterShouldWorkConcurrently(_ *testing.T) {
	sut := filter.NewHTTPStatusResultFilter(nil, false)

	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func(i int) {
			sut.ShouldIgnore(scan.Result{StatusCode: i})
			wg.Done()
		}(i)
	}

	wg.Wait()
}
