package csm

import (
	"encoding/json"
	"os"

	"github.com/hpcloud/sidecar-extensions/go/csm/status"
)

type csmFileConnection struct {
	filePath string
}

func NewCSMFileConnection(filePath string) CSMConnection {
	return &csmFileConnection{filePath: filePath}
}

func (c *csmFileConnection) Write(response CSMResponse) error {
	f, err := os.OpenFile(c.filePath, os.O_RDWR|os.O_APPEND, 0660)
	defer f.Close()
	if err != nil {
		return err
	}

	b, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func (c *csmFileConnection) WriteError(input error) error {
	response := NewCSMResponse(500, input.Error(), status.Failed)

	err := c.Write(response)
	return err
}
