package summarizer

import (
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

type ResultTree interface {
	String(results []scan.Result) string
}
