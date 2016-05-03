package csm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSMRequest(t *testing.T) {
	assert := assert.New(t)
	workspaceID := "905ab3af-b149-46b6-bc6e-75073c703e19"
	connectionID := "c0232d8b-2cbe-4c5c-b7ee-9024c6ee0006"
	outFile := "/home/user/myfile"

	args := []string{"exe", outFile, workspaceID, connectionID}

	request, err := GetCSMRequest(args)

	assert.Nil(err)
	assert.Equal(connectionID, request.ConnectionID)
	assert.Equal(workspaceID, request.WorkspaceID)
	assert.Equal(outFile, request.OutputPath)
}

func TestCSMRequestNoConnection(t *testing.T) {
	assert := assert.New(t)
	workspaceID := "905ab3af-b149-46b6-bc6e-75073c703e19"
	outFile := "/home/user/myfile"

	args := []string{"exe", outFile, workspaceID}

	request, err := GetCSMRequest(args)

	assert.Nil(err)
	assert.Equal("", request.ConnectionID)
	assert.Equal(workspaceID, request.WorkspaceID)
	assert.Equal(outFile, request.OutputPath)
}

func TestCSMRequestFail(t *testing.T) {
	assert := assert.New(t)
	workspaceID := "905ab3af-b149-46b6-bc6e-75073c703e19"
	outFile := "/home/user/myfile"

	args := []string{"exe", outFile, workspaceID, "extrastring", "extrastring"}

	_, err := GetCSMRequest(args)

	assert.NotNil(err)
}
