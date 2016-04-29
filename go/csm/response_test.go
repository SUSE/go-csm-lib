package csm

import (
	"testing"

	"github.com/hpcloud/sidecar-extensions/go/csm/status"
	"github.com/stretchr/testify/assert"
)

type testDetails struct {
	One string
	Two int
}

func TestCSMResponse(t *testing.T) {
	assert := assert.New(t)
	details := testDetails{One: "test", Two: 1}
	response := NewCSMResponse(200, details, status.Successful)

	assert.Equal(200, response.HttpCode)
	assert.Equal(details, response.Details)
	assert.Equal("Successful", response.Status)
}
