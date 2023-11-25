package filter

import (
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"regexp"
)

func NewHTTPStatusResultFilter(httpStatusesToIgnore []int, ignoreEmptyBody bool, assume404Substring string) (*HTTPStatusResultFilter, error) {
	httpStatusesToIgnoreMap := make(map[int]struct{}, len(httpStatusesToIgnore))
	for _, statusToIgnore := range httpStatusesToIgnore {
		httpStatusesToIgnoreMap[statusToIgnore] = struct{}{}
	}
	var assume404Regex *regexp.Regexp
	var err error
	if assume404Substring != "" {
		assume404Regex, err = regexp.Compile(assume404Substring)
		if err != nil {
			return nil, err
		}
	}

	return &HTTPStatusResultFilter{httpStatusesToIgnoreMap: httpStatusesToIgnoreMap, ignoreEmptyBody: ignoreEmptyBody, assume404Regex: assume404Regex}, nil
}

type HTTPStatusResultFilter struct {
	httpStatusesToIgnoreMap map[int]struct{}
	ignoreEmptyBody         bool
	assume404Regex          *regexp.Regexp
}

func (f HTTPStatusResultFilter) ShouldIgnore(result scan.Result) bool {
	if f.ignoreEmptyBody && result.StatusCode/100 == 2 && result.ContentLength == 0 {
		return true
	}

	if f.assume404Regex != nil {
		if f.assume404Regex.Match(result.Body) {
			_, found := f.httpStatusesToIgnoreMap[404]
			return found
		}
	}
	_, found := f.httpStatusesToIgnoreMap[result.StatusCode]

	return found
}

func (f HTTPStatusResultFilter) ShouldReadBody() bool {
	return f.assume404Regex != nil
}
