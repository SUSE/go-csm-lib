package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeGuid(t *testing.T) {
	assert := assert.New(t)
	guid := "8b490a70-a892-4eff-a495-81e905f3960f"

	value := NormalizeGuid(guid)

	assert.Equal(value, "d8b490a70a8924effa49581e905f3960f")
}

func TestGetMD5Hash(t *testing.T) {
	assert := assert.New(t)

	someText := "teststring"

	hash, err := GetMD5Hash(someText, 5)

	assert.Nil(err)

	assert.Equal(5, len(hash))

	sec, err := GetMD5Hash(someText, 5)
	assert.Nil(err)

	assert.Equal(5, len(sec))

	assert.Equal(hash, sec)

}

func TestGetSecureString(t *testing.T) {
	assert := assert.New(t)

	ss, err := SecureRandomString(32)

	assert.Nil(err)
	assert.NotNil(ss)
}
