package producer_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/producer"
	"github.com/stretchr/testify/assert"
)

func TestDictionaryProducerShouldProduce(t *testing.T) {
	t.Parallel()

	const depth = 4

	sut := producer.NewDictionaryProducer(
		[]string{http.MethodGet, http.MethodPost},
		[]string{"/home", "/about"},
		depth,
	)

	results := make([]scan.Target, 0, 4)

	producerChannel := sut.Produce(context.Background())
	for r := range producerChannel {
		results = append(results, r)
	}

	assert.Len(t, results, 4)

	expectedResults := []scan.Target{
		{
			Depth:  depth,
			Path:   "/home",
			Method: http.MethodGet,
		},
		{
			Depth:  depth,
			Path:   "/home",
			Method: http.MethodPost,
		},
		{
			Depth:  depth,
			Path:   "/about",
			Method: http.MethodGet,
		},
		{
			Depth:  depth,
			Path:   "/about",
			Method: http.MethodPost,
		},
	}

	assert.Equal(t, expectedResults, results)
}

func TestDictionaryProducerCanBeCanceled(t *testing.T) {
	t.Parallel()

	const depth = 4

	sut := producer.NewDictionaryProducer(
		[]string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		[]string{"/home", "/about", "/index", "/search", "/tomato"},
		depth,
	)

	ctx, cancelFunc := context.WithCancel(context.Background())

	producerChannel := sut.Produce(ctx)

	cancelFunc()

	resultsCount := 0

	for range producerChannel {
		resultsCount++
	}

	// 11 is the size of the producer buffer
	assert.True(t, resultsCount <= 11)
}
