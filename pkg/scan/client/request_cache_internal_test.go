package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestCacheTransportDecorator(t *testing.T) {
	transport, err := decorateTransportWithRequestCacheDecorator(nil)
	assert.Nil(t, transport)
	assert.Error(t, err)
}
