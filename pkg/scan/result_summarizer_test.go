package scan_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stretchr/testify/assert"
)

func TestResultSummarizer(t *testing.T) {
	b := &bytes.Buffer{}

	summarizer := scan.NewResultSummarizer(b)

	summarizer.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/home",
			},
			&http.Response{
				StatusCode: 201,
				Request: &http.Request{
					URL: test.MustParseUrl(t, "http://mysite/home"),
				},
			},
		),
	)

	summarizer.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/home/hidden",
			},
			&http.Response{
				StatusCode: 201,
				Request: &http.Request{
					URL: test.MustParseUrl(t, "http://mysite/home/hidden"),
				},
			},
		),
	)

	summarizer.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/home/about",
			},
			&http.Response{
				StatusCode: 200,
				Request: &http.Request{
					URL: test.MustParseUrl(t, "http://mysite/home/about"),
				},
			},
		),
	)

	summarizer.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/home/about/me",
			},
			&http.Response{
				StatusCode: 200,
				Request: &http.Request{
					URL: test.MustParseUrl(t, "http://mysite/home/about/me"),
				},
			},
		),
	)

	summarizer.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/home/home",
			},
			&http.Response{
				StatusCode: 200,
				Request: &http.Request{
					URL: test.MustParseUrl(t, "http://mysite/home/home"),
				},
			},
		),
	)

	summarizer.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/contacts",
			},
			&http.Response{
				StatusCode: 200,
				Request: &http.Request{
					URL: test.MustParseUrl(t, "http://mysite/contacts"),
				},
			},
		),
	)

	summarizer.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/gibberish",
			},
			&http.Response{
				StatusCode: 404,
				Request: &http.Request{
					URL: test.MustParseUrl(t, "http://mysite/gibberish"),
				},
			},
		),
	)

	summarizer.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/path/to/my/files",
			},
			&http.Response{
				StatusCode: 200,
				Request: &http.Request{
					URL: test.MustParseUrl(t, "http://mysite/path/to/my/files"),
				},
			},
		),
	)

	// Adding twice the same result should not change the outcome
	summarizer.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/path/to/my/files",
			},
			&http.Response{
				StatusCode: 200,
				Request: &http.Request{
					URL: test.MustParseUrl(t, "http://mysite/path/to/my/files"),
				},
			},
		),
	)

	summarizer.Summarize()

	expectedResult := `9 requests made, 7 results found
/
├── contacts
├── home
│   ├── about
│   │   └── me
│   ├── hidden
│   └── home
└── path
    └── to
        └── my
            └── files

http://mysite/contacts [200] [GET]
http://mysite/home [201] [POST]
http://mysite/home/about [200] [GET]
http://mysite/home/about/me [200] [GET]
http://mysite/home/hidden [201] [POST]
http://mysite/home/home [200] [GET]
http://mysite/path/to/my/files [200] [GET]
`
	assert.Equal(t, expectedResult, b.String())
}
