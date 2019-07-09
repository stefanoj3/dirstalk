package scan_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/producer"
	"github.com/stretchr/testify/assert"
)

func TestScanningWithEmptyProducerWillProduceNoResults(t *testing.T) {
	logger, _ := test.NewLogger()

	prod := producer.NewDictionaryProducer([]string{}, []string{}, 1)
	client := &http.Client{Timeout: time.Microsecond}
	sut := scan.NewScanner(
		client,
		prod,
		producer.NewReProducer(prod),
		logger,
	)

	results := sut.Scan(test.MustParseUrl(t, "http://localhost/"), 10)

	for r := range results {
		t.Fatalf("No results expected, got %s", r.Target.Path)
	}
}

func TestScannerWillLogAnErrorWithInvalidDictionary(t *testing.T) {
	logger, loggerBuffer := test.NewLogger()

	prod := producer.NewDictionaryProducer(
		[]string{"\n"},
		[]string{"/home"},
		1,
	)

	client := &http.Client{Timeout: time.Microsecond}
	sut := scan.NewScanner(
		client,
		prod,
		producer.NewReProducer(prod),
		logger,
	)

	testServer, serverAssertion := test.NewServerWithAssertion(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	defer testServer.Close()

	results := sut.Scan(test.MustParseUrl(t, testServer.URL), 10)

	for r := range results {
		t.Fatalf("No results expected, got %s", r.Target.Path)
	}

	assert.Contains(t, loggerBuffer.String(), "failed to build request")
	assert.Contains(t, loggerBuffer.String(), "invalid method")
	assert.Equal(t, 0, serverAssertion.Len())
}
