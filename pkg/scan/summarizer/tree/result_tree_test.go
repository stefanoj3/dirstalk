package tree_test

import (
	"net/http"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/summarizer/tree"
	"github.com/stretchr/testify/assert"
)

func TestNewResultTreePrinter(t *testing.T) {
	results := []scan.Result{
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/home",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/home"),
				},
			},
		),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/home/123",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/home/123"),
				},
			},
		),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/about"),
				},
			},
		),
	}

	actual := tree.NewResultTreeProducer().String(results)

	expected := `/
├── about
└── home
    └── 123
`

	assert.Equal(t, expected, actual)
}
