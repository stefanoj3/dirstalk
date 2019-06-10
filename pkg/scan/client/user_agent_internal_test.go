package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecorateTransportUserAgent(t *testing.T) {
	transport, err := decorateTransportWithUserAgentDecorator(nil, "")
	assert.Nil(t, transport)
	assert.Error(t, err)
}
