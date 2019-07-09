package producer_test

import (
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

	producerChannel := sut.Produce()
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
