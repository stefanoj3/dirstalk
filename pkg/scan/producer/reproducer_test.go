package producer_test

import (
	"net/http"
	"sort"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/test"

	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/producer"
	"github.com/stretchr/testify/assert"
)

func TestNewReProducer(t *testing.T) {
	t.Parallel()

	methods := []string{http.MethodGet, http.MethodPost}
	dictionary := []string{"/home", "/about"}

	dictionaryProducer := producer.NewDictionaryProducer(methods, dictionary, 1)

	sut := producer.NewReProducer(dictionaryProducer)

	result := scan.NewResult(
		scan.Target{
			Path:   "/home",
			Method: http.MethodGet,
			Depth:  1,
		},
		&http.Response{
			StatusCode: http.StatusOK,
			Request: &http.Request{
				URL: test.MustParseUrl(t, "http://mysite/contacts"),
			},
		},
	)

	reproducerFunc := sut.Reproduce()
	reproducerChannel := reproducerFunc(result)

	targets := make([]scan.Target, 0, 10)
	for tar := range reproducerChannel {
		targets = append(targets, tar)
	}

	sort.Slice(targets, func(i, j int) bool {
		return targets[i].Path < targets[j].Path && targets[i].Method < targets[j].Method
	})

	assert.Len(t, targets, 4)

	expectedTargets := []scan.Target{
		{
			Path:   "/home/home",
			Method: http.MethodGet,
			Depth:  0,
		},
		{
			Path:   "/home/about",
			Method: http.MethodGet,
			Depth:  0,
		},
		{
			Path:   "/home/home",
			Method: http.MethodPost,
			Depth:  0,
		},
		{
			Path:   "/home/about",
			Method: http.MethodPost,
			Depth:  0,
		},
	}
	assert.Equal(t, expectedTargets, targets)

	// reproducing again on the same result should not yield more targets
	reproducerChannel = reproducerFunc(result)

	targets = make([]scan.Target, 0)
	for tar := range reproducerChannel {
		targets = append(targets, tar)
	}
	assert.Len(t, targets, 0)
}

func TestReProducerShouldProduceNothingForDepthZero(t *testing.T) {
	t.Parallel()

	methods := []string{http.MethodGet, http.MethodPost}
	dictionary := []string{"/home", "/about"}

	dictionaryProducer := producer.NewDictionaryProducer(methods, dictionary, 1)

	sut := producer.NewReProducer(dictionaryProducer)

	result := scan.NewResult(
		scan.Target{
			Path:   "/home",
			Method: http.MethodGet,
			Depth:  0,
		},
		&http.Response{
			StatusCode: http.StatusOK,
			Request: &http.Request{
				URL: test.MustParseUrl(t, "http://mysite/contacts"),
			},
		},
	)

	reproducerFunc := sut.Reproduce()
	reproducerChannel := reproducerFunc(result)

	targets := make([]scan.Target, 0)
	for tar := range reproducerChannel {
		targets = append(targets, tar)
	}

	assert.Len(t, targets, 0)
}
