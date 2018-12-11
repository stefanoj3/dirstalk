package scan_test

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stefanoj3/dirstalk/pkg/scan"
)

func TestResultSummarizer(t *testing.T) {
	b := &bytes.Buffer{}

	summarizer := scan.NewResultSummarizer(b)

	summarizer.Add(
		&scan.Result{
			Response: &http.Response{
				StatusCode: 201,
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    MustParseUrl(t, "http://mysite/home"),
				},
			},
		},
	)

	summarizer.Add(
		&scan.Result{
			Response: &http.Response{
				StatusCode: 201,
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    MustParseUrl(t, "http://mysite/home/hidden"),
				},
			},
		},
	)

	summarizer.Add(
		&scan.Result{
			Response: &http.Response{
				StatusCode: 200,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    MustParseUrl(t, "http://mysite/home/about"),
				},
			},
		},
	)

	summarizer.Add(
		&scan.Result{
			Response: &http.Response{
				StatusCode: 200,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    MustParseUrl(t, "http://mysite/home/about/me"),
				},
			},
		},
	)

	summarizer.Add(
		&scan.Result{
			Response: &http.Response{
				StatusCode: 200,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    MustParseUrl(t, "http://mysite/contacts"),
				},
			},
		},
	)

	summarizer.Add(
		&scan.Result{
			Response: &http.Response{
				StatusCode: 404,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    MustParseUrl(t, "http://mysite/gibberish"),
				},
			},
		},
	)

	summarizer.Summarize()

	expectedResult := `6 requests made, 5 results found
/
├── home
│   ├── hidden
│   └── about
│       └── me
└── contacts

http://mysite/home [201] [POST]
http://mysite/home/hidden [201] [POST]
http://mysite/home/about [200] [GET]
http://mysite/home/about/me [200] [GET]
http://mysite/contacts [200] [GET]
`

	assert.Equal(t, expectedResult, b.String())
}

func MustParseUrl(t *testing.T, rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("failed to parse url: %s", rawurl)
	}

	return u
}
