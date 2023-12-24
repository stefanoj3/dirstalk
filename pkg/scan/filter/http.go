package filter

import (
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"regexp"
)

func NewHTTPStatusResultFilter(httpStatusesToIgnore []int, ignoreEmptyBody bool, assumeStatusRegex map[int]regexp.Regexp) (*HTTPStatusResultFilter, error) {
	httpStatusesToIgnoreMap := make(map[int]struct{}, len(httpStatusesToIgnore))
	for _, statusToIgnore := range httpStatusesToIgnore {
		httpStatusesToIgnoreMap[statusToIgnore] = struct{}{}
	}

	return &HTTPStatusResultFilter{httpStatusesToIgnoreMap: httpStatusesToIgnoreMap, ignoreEmptyBody: ignoreEmptyBody, assumeStatusRegex: assumeStatusRegex}, nil
}

type HTTPStatusResultFilter struct {
	httpStatusesToIgnoreMap map[int]struct{}
	ignoreEmptyBody         bool
	assumeStatusRegex       map[int]regexp.Regexp
}

func (f HTTPStatusResultFilter) ShouldIgnore(result scan.Result) bool {
	if f.ignoreEmptyBody && result.StatusCode/100 == 2 && result.ContentLength == 0 {
		return true
	}

	if f.assumeStatusRegex != nil {
		for code, regex := range f.assumeStatusRegex {
			if regex.Match(result.Body) {
				_, found := f.httpStatusesToIgnoreMap[code]
				return found
			}
		}
	}
	_, found := f.httpStatusesToIgnoreMap[result.StatusCode]

	return found
}

func (f HTTPStatusResultFilter) ShouldReadBody() bool {
	return len(f.assumeStatusRegex) > 0
}
