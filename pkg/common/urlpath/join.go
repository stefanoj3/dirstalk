package urlpath

import (
	"path"
	"strings"
)

// Join joins any number of path elements into a single path, adding a
// separating slash if necessary. The result is Cleaned; in particular,
// all empty strings are ignored.
// If the last element end in a slash it will preserve it.
func Join(elem ...string) string {
	joined := path.Join(elem...)

	last := elem[len(elem)-1]

	if strings.HasSuffix(last, "/") && !strings.HasSuffix(joined, "/") {
		joined += "/"
	}

	return joined
}
