package navigator

import "github.com/signalfx/signalfx-go/util"

type JobFilter struct {
	Property      string             `json:"property,omitempty"`
	PropertyValue util.StringOrSlice `json:"propertyValue,omitempty"`
	Not           bool               `json:"not,omitempty"`
	Type          string             `json:"type,omitempty"`
}
