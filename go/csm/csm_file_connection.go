package csm

type csmFileConnection struct {
	filePath string
}

type CSMResponse struct {
	httpCode int         `json:"http_code"`
	payload  interface{} `json:"payload"`
}

func NewCSM(filePath string) CSMConnection {
	return &csmFileConnection{filePath: filePath}
}

func (*csmFileConnection) SendResponse() {

}
