package result_test

import (
	"net/url"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/result"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stretchr/testify/assert"
)

func TestLoadResultsFromFile(t *testing.T) {
	results, err := result.LoadResultsFromFile("testdata/out.txt")
	assert.NoError(t, err)

	expectedResults := []scan.Result{
		{
			Target:     scan.Target{Path: "partners", Method: "GET", Depth: 3},
			StatusCode: 200,
			URL: url.URL{
				Scheme: "https",
				User:   (*url.Userinfo)(nil),
				Host:   "www.brucewillisdiesinarmageddon.co.de",
				Path:   "/partners",
			},
		},
		{
			Target:     scan.Target{Path: "s", Method: "GET", Depth: 3},
			StatusCode: 400,
			URL: url.URL{
				Scheme: "https",
				User:   (*url.Userinfo)(nil),
				Host:   "www.brucewillisdiesinarmageddon.co.de",
				Path:   "/s",
			},
		},
		{
			Target:     scan.Target{Path: "adview", Method: "GET", Depth: 3},
			StatusCode: 204,
			URL: url.URL{
				Scheme: "https",
				User:   (*url.Userinfo)(nil),
				Host:   "www.brucewillisdiesinarmageddon.co.de",
				Path:   "/adview",
			},
		},
		{
			Target:     scan.Target{Path: "partners/terms", Method: "GET", Depth: 2},
			StatusCode: 200,
			URL: url.URL{
				Scheme: "https",
				User:   (*url.Userinfo)(nil),
				Host:   "www.brucewillisdiesinarmageddon.co.de",
				Path:   "/partners/terms",
			},
		},
	}

	assert.Equal(t, expectedResults, results)
}

func TestLoadResultsFromFileShouldErrForDirectories(t *testing.T) {
	_, err := result.LoadResultsFromFile("testdata/")
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "is a directory")
}

func TestLoadResultsFromFileShouldErrForInvalidFileFormat(t *testing.T) {
	_, err := result.LoadResultsFromFile("testdata/invalidout.txt")
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "unable to read line")
}
