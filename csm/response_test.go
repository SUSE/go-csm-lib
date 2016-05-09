package csm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDetails struct {
	One string
	Two int
}

func TestCSMResponse(t *testing.T) {
	assert := assert.New(t)
	details := testDetails{One: "test", Two: 1}
	response := CreateCSMResponse(details)

	assert.Equal(details, response.Details)
}
