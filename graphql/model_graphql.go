package graphql

// The make up of a graphql request to the SignalFx API
type Request struct {
	OperationName string      `json:"operationName"`
	Variables     interface{} `json:"variables"`
	Query         string      `json:"query"`
}

// The make up of a graphql response from the SignalFx API
type Response struct {
	Data interface{} `json:"data"`
}
