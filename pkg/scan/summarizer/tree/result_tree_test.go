package tree_test

import (
	"net/http"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/summarizer/tree"
	"github.com/stretchr/testify/assert"
)

var testResult string

func TestNewResultTreePrinter(t *testing.T) {
	results := []scan.Result{
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/"),
				},
			}, nil),
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
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/home/123/",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/home/123"),
				},
			}, nil),
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
			}, nil),
	}

	actual := tree.NewResultTreeProducer().String(results)

	expected := `/
├── about
└── home
    └── 123
`

	assert.Equal(t, expected, actual)
}

func BenchmarkResultTree(b *testing.B) {
	results := []scan.Result{
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/home",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "http://mysite/home"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/home/123",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "http://mysite/home/123"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "http://mysite/about"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "http://mysite/about"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "http://mysite/about"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/12",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/12"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/123",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/123"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/3",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/3"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/b",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/b"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/b/c/d/e",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/b/c/d/e"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/b/c/f/e",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/b/c/f/e"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/i/l/m/n/o/p/q",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/i/l/m/n/o/p/q"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/1/l/m/n/o/p/q",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/1/l/m/n/o/p/q"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/2/l/m/n/o/p/q",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/2/l/m/n/o/p/q"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/3/l/m/n/o/p/q",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/3/l/m/n/o/p/q"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/4/l/m/n/o/p/q",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/4/l/m/n/o/p/q"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/4/1/m/n/o/p/q",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/4/1/m/n/o/p/q"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/4/2/m/n/o/p/q",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/4/2/m/n/o/p/q"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/4/2/m/n/o/p/u",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/4/2/m/n/o/p/u"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/4/2/m/n/o/p/z",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/4/2/m/n/o/p/z"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/about/1/2/a/c/f/e/g/h/4/2/m/n/o/p/z/1",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/about/1/2/a/c/f/e/g/h/4/2/m/n/o/p/z/1"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/somepage",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/somepage"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/anotherpage",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/anotherpage"),
				},
			}, nil),
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/anotherpage2",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(b, "/anotherpage2"),
				},
			}, nil),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		testResult = tree.NewResultTreeProducer().String(results)
	}

	testResult += "1"
}
