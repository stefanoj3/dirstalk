package filter

import (
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

func NewHTTPStatusResultFilter(httpStatusesToIgnore []int) HTTPStatusResultFilter {
	httpStatusesToIgnoreMap := make(map[int]struct{}, len(httpStatusesToIgnore))
	for _, statusToIgnore := range httpStatusesToIgnore {
		httpStatusesToIgnoreMap[statusToIgnore] = struct{}{}
	}

	return HTTPStatusResultFilter{httpStatusesToIgnoreMap: httpStatusesToIgnoreMap}
}

type HTTPStatusResultFilter struct {
	httpStatusesToIgnoreMap map[int]struct{}
}

func (f HTTPStatusResultFilter) ShouldIgnore(result scan.Result) bool {
	_, found := f.httpStatusesToIgnoreMap[result.StatusCode]

	return found
}
