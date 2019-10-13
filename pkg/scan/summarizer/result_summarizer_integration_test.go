package summarizer_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/summarizer"
	"github.com/stefanoj3/dirstalk/pkg/scan/summarizer/tree"
	"github.com/stretchr/testify/assert"
)

func TestResultSummarizerShouldSummarizeResults(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()
	logger.SetLevel(logrus.FatalLevel)

	sut := summarizer.NewResultSummarizer(tree.NewResultTreeProducer(), logger)

	sut.Add(
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
	)

	sut.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodPost,
				Path:   "/home/hidden",
			},
			&http.Response{
				StatusCode: http.StatusCreated,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/home/hidden"),
				},
			},
		),
	)

	sut.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/home/about",
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/home/about"),
				},
			},
		),
	)

	sut.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/home/about/me",
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/home/about/me"),
				},
			},
		),
	)

	sut.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/home/home",
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/home/home"),
				},
			},
		),
	)

	sut.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/contacts",
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/contacts"),
				},
			},
		),
	)

	sut.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/gibberish",
			},
			&http.Response{
				StatusCode: http.StatusNotFound,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/gibberish"),
				},
			},
		),
	)

	sut.Add(
		scan.NewResult(
			scan.Target{
				Method: http.MethodGet,
				Path:   "/path/to/my/files",
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					URL: test.MustParseURL(t, "http://mysite/path/to/my/files"),
				},
			},
		),
	)

	// Adding multiple times the same result should not change the outcome
	wg := &sync.WaitGroup{}

	const workers = 10

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			sut.Add(
				scan.NewResult(
					scan.Target{
						Method: http.MethodGet,
						Path:   "/path/to/my/files",
					},
					&http.Response{
						StatusCode: http.StatusOK,
						Request: &http.Request{
							URL: test.MustParseURL(t, "http://mysite/path/to/my/files"),
						},
					},
				),
			)
		}()
	}

	wg.Wait()

	sut.Summarize()

	expectedResult := `8 results found
/
├── contacts
├── gibberish
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
http://mysite/gibberish [404] [GET]
http://mysite/home [201] [POST]
http://mysite/home/about [200] [GET]
http://mysite/home/about/me [200] [GET]
http://mysite/home/hidden [201] [POST]
http://mysite/home/home [200] [GET]
http://mysite/path/to/my/files [200] [GET]
`
	assert.Equal(t, expectedResult, loggerBuffer.String())
}

func TestResultSummarizerShouldLogResults(t *testing.T) {
	testCases := []struct {
		result            scan.Result
		expectedToContain []string
	}{
		{
			result: scan.NewResult(
				scan.Target{
					Method: http.MethodPost,
					Path:   "/home",
				},
				&http.Response{
					StatusCode: http.StatusOK,
					Request: &http.Request{
						URL: test.MustParseURL(t, "http://mysite/home"),
					},
				},
			),
			expectedToContain: []string{
				"Found",
				"method=POST",
				"status-code=200",
				`url="http://mysite/home"`,
			},
		},
		{
			result: scan.NewResult(
				scan.Target{
					Method: http.MethodGet,
					Path:   "/index",
				},
				&http.Response{
					StatusCode: http.StatusBadGateway,
					Request: &http.Request{
						URL: test.MustParseURL(t, "http://mysite/index"),
					},
				},
			),
			expectedToContain: []string{
				"Found something breaking",
				"method=GET",
				"status-code=502",
				`url="http://mysite/index"`,
			},
		},
		{
			result: scan.NewResult(
				scan.Target{
					Method: http.MethodGet,
					Path:   "/gibberish",
				},
				&http.Response{
					StatusCode: http.StatusNotFound,
					Request: &http.Request{
						URL: test.MustParseURL(t, "http://mysite/gibberish"),
					},
				},
			),
			expectedToContain: []string{
				"Found",
				"method=GET",
				"status-code=404",
				`url="http://mysite/gibberish"`,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint
		t.Run(tc.result.Target.Path, func(t *testing.T) {
			t.Parallel()
			logger, loggerBuffer := test.NewLogger()
			sut := summarizer.NewResultSummarizer(tree.NewResultTreeProducer(), logger)

			sut.Add(tc.result)

			bufferAsString := loggerBuffer.String()
			for _, expectedToContain := range tc.expectedToContain {
				assert.Contains(t, bufferAsString, expectedToContain)
			}
		})
	}
}
