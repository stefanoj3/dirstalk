package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecorateTransportHeaderShouldFailWithNilDecorated(t *testing.T) {
	transport, err := decorateTransportWithHeadersDecorator(nil, map[string]string{})
	assert.Nil(t, transport)
	assert.Error(t, err)
}

func TestDecorateTransportHeaderShouldFailWithNilHeaderMap(t *testing.T) {
	transport, err := decorateTransportWithHeadersDecorator(http.DefaultTransport, nil)
	assert.Nil(t, transport)
	assert.Error(t, err)
}
