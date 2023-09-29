package graphql

// The make up of a graphql request to the SignalFx API
type Request struct {
	OperationName string    `json:"operationName"`
	Variables     Variables `json:"variables"`
	Query         string    `json:"query"`
}

type Variables struct {
	Parameters Parameters `json:"parameters"`
}

type Parameters struct {
	SectionsParameters []SectionParams  `json:"sectionsParameters"`
	SharedParameters   SharedParameters `json:"sharedParameters"`
}

// These variables are dependent on the type of request being made
type SharedParameters interface{}

type SectionParams struct {
	SectionType string `json:"sectionType"`
	Limit       int    `json:"limit,omitempty"`
}

// The make up of a graphql response from the SignalFx API
type Response struct {
	Data interface{} `json:"data"`
}
