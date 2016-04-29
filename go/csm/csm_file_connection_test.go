package csm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hpcloud/sidecar-extensions/go/csm/status"
	"github.com/stretchr/testify/assert"
)

func TestCSMFileConnectionWrite(t *testing.T) {
	assert := assert.New(t)

	testFile, err := ioutil.TempFile(os.TempDir(), "prefix")
	defer os.Remove(testFile.Name())

	connection := NewCSMFileConnection(testFile.Name())

	response := NewCSMResponse(200, "", status.Successful)
	err = connection.Write(response)

	assert.Nil(err)

	content, err := ioutil.ReadFile(testFile.Name())
	assert.Nil(err)

	newResponse := CSMResponse{}
	err = json.Unmarshal(content, &newResponse)
	assert.Nil(err)
	assert.Equal(200, newResponse.HttpCode)
	assert.Equal("", newResponse.Details)
	assert.Equal("Extension", newResponse.ProcessingType)
}

func TestCSMFileConnectionWriteError(t *testing.T) {
	assert := assert.New(t)

	testFile, err := ioutil.TempFile(os.TempDir(), "prefix")
	defer os.Remove(testFile.Name())

	connection := NewCSMFileConnection(testFile.Name())

	err = connection.WriteError(errors.New("an error"))

	assert.Nil(err)

	content, err := ioutil.ReadFile(testFile.Name())
	assert.Nil(err)

	newResponse := CSMResponse{}
	err = json.Unmarshal(content, &newResponse)
	assert.Nil(err)
	assert.Equal(500, newResponse.HttpCode)
	assert.Equal("an error", newResponse.Details)
	assert.Equal("Extension", newResponse.ProcessingType)
}

func TestCSMFileConnectionWriteStruct(t *testing.T) {
	assert := assert.New(t)

	testFile, err := ioutil.TempFile(os.TempDir(), "prefix")
	fmt.Println(testFile.Name())
	defer os.Remove(testFile.Name())

	connection := NewCSMFileConnection(testFile.Name())
	details := testDetails{One: "test", Two: 1}
	response := NewCSMResponse(200, details, status.Successful)
	err = connection.Write(response)

	assert.Nil(err)

	content, err := ioutil.ReadFile(testFile.Name())
	assert.Nil(err)

	newResponse := CSMResponse{}
	err = json.Unmarshal(content, &newResponse)
	assert.Nil(err)
	assert.Equal(200, newResponse.HttpCode)

	assert.Equal("test", newResponse.Details.(map[string]interface{})["One"])

	var twoValue float64
	twoValue = 1
	assert.Equal(twoValue, newResponse.Details.(map[string]interface{})["Two"])

	assert.Equal("Extension", newResponse.ProcessingType)
}
