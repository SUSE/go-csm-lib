package csm

import "github.com/hpcloud/sidecar-extensions/go/csm/status"

type CSMResponse struct {
	HttpCode       int         `json:"http_code"`
	Payload        interface{} `json:"details"`
	Status         string      `json:"status"`
	ProcessingType string      `json:"processing_type"`
}

func NewCSMResponse(httpCode int, payload interface{}, status status.Status) CSMResponse {
	return CSMResponse{
		HttpCode: httpCode, Payload: payload, Status: string(status), ProcessingType: "Extension",
	}
}
