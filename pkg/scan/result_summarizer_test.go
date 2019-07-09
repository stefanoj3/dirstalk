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
		scan.Result{
			Response: &http.Response{
				StatusCode: 201,
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    test.MustParseUrl(t, "http://mysite/home"),
				},
			},
		},
	)

	summarizer.Add(
		scan.Result{
			Response: &http.Response{
				StatusCode: 201,
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    test.MustParseUrl(t, "http://mysite/home/hidden"),
				},
			},
		},
	)

	summarizer.Add(
		scan.Result{
			Response: &http.Response{
				StatusCode: 200,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    test.MustParseUrl(t, "http://mysite/home/about"),
				},
			},
		},
	)

	summarizer.Add(
		scan.Result{
			Response: &http.Response{
				StatusCode: 200,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    test.MustParseUrl(t, "http://mysite/home/about/me"),
				},
			},
		},
	)

	summarizer.Add(
		scan.Result{
			Response: &http.Response{
				StatusCode: 200,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    test.MustParseUrl(t, "http://mysite/home/home"),
				},
			},
		},
	)

	summarizer.Add(
		scan.Result{
			Response: &http.Response{
				StatusCode: 200,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    test.MustParseUrl(t, "http://mysite/contacts"),
				},
			},
		},
	)

	summarizer.Add(
		scan.Result{
			Response: &http.Response{
				StatusCode: 404,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    test.MustParseUrl(t, "http://mysite/gibberish"),
				},
			},
		},
	)

	summarizer.Add(
		scan.Result{
			Response: &http.Response{
				StatusCode: 200,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    test.MustParseUrl(t, "http://mysite/path/to/my/files"),
				},
			},
		},
	)

	summarizer.Summarize()

	expectedResult := `8 requests made, 7 results found
/
├── home
│   ├── hidden
│   ├── about
│   │   └── me
│   └── home
├── contacts
└── path
    └── to
        └── my
            └── files

http://mysite/home [201] [POST]
http://mysite/home/hidden [201] [POST]
http://mysite/home/about [200] [GET]
http://mysite/home/about/me [200] [GET]
http://mysite/home/home [200] [GET]
http://mysite/contacts [200] [GET]
http://mysite/path/to/my/files [200] [GET]
`

	assert.Equal(t, expectedResult, b.String())
}
