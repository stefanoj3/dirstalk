package filter

import (
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

func NewHTTPStatusResultFilter(httpStatusesToIgnore []int, ignoreEmptyBody bool) HTTPStatusResultFilter {
	httpStatusesToIgnoreMap := make(map[int]struct{}, len(httpStatusesToIgnore))
	for _, statusToIgnore := range httpStatusesToIgnore {
		httpStatusesToIgnoreMap[statusToIgnore] = struct{}{}
	}

	return HTTPStatusResultFilter{httpStatusesToIgnoreMap: httpStatusesToIgnoreMap, ignoreEmptyBody: ignoreEmptyBody}
}

type HTTPStatusResultFilter struct {
	httpStatusesToIgnoreMap map[int]struct{}
	ignoreEmptyBody         bool
}

func (f HTTPStatusResultFilter) ShouldIgnore(result scan.Result) bool {
	if f.ignoreEmptyBody && result.StatusCode/100 == 2 && result.ContentLength == 0 {
		return true
	}

	_, found := f.httpStatusesToIgnoreMap[result.StatusCode]

	return found
}
