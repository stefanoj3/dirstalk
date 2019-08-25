package urlpath

import (
	"path"
)

func HasExtension(p string) bool {
	return len(path.Ext(p)) > 0
}
