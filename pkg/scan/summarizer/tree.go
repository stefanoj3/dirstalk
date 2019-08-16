package summarizer

import (
	"io"

	"github.com/stefanoj3/dirstalk/pkg/scan"
)

type ResultTreePrinter interface {
	Print(results []scan.Result, out io.Writer)
}
