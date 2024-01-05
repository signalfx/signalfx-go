package slo

import (
	"encoding/json"
	"fmt"
)

const (
	RequestBased = "RequestBased"
	WindowsBased = "WindowsBased"
)

type BaseSlo struct {
	Creator       string      `json:"creator,omitempty"`
	LastUpdatedBy string      `json:"lastUpdatedBy,omitempty"`
	Created       int64       `json:"created,omitempty"`
	LastUpdated   int64       `json:"lastUpdated,omitempty"`
	Id            string      `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	Description   string      `json:"description,omitempty"`
	Targets       []SloTarget `json:"targets,omitempty"`
	Type          string      `json:"type,omitempty"`
	Metadata      []string    `json:"metadata,omitempty"`
}

type SloObject struct {
	BaseSlo
	*RequestBasedSlo
	*WindowBasedSlo
}

type RequestBasedSlo struct {
	Inputs *RequestBasedSloInput `json:"inputs,omitempty"`
}

type WindowBasedSlo struct {
	Inputs *WindowBasedSloInput `json:"inputs,omitempty"`
}

type RequestBasedSloInput struct {
	ProgramText      string `json:"programText,omitempty"`
	GoodEventsLabel  string `json:"goodEventsLabel,omitempty"`
	TotalEventsLabel string `json:"totalEventsLabel,omitempty"`
}

type WindowBasedSloInput struct {
	ProgramText string `json:"programText,omitempty"`
}

func (slo *SloObject) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &slo.BaseSlo); err != nil {
		return err
	}
	switch slo.Type {
	case RequestBased:
		slo.RequestBasedSlo = &RequestBasedSlo{}
		return json.Unmarshal(data, slo.RequestBasedSlo)
	case WindowsBased:
		slo.WindowBasedSlo = &WindowBasedSlo{}
		return json.Unmarshal(data, slo.WindowBasedSlo)
	default:
		return fmt.Errorf("unrecognized SLO type %s", slo.Type)
	}
}

func (slo SloObject) MarshalJSON() ([]byte, error) {
	switch slo.Type {
	case RequestBased:
		return json.Marshal(struct {
			BaseSlo
			*RequestBasedSlo
		}{slo.BaseSlo, slo.RequestBasedSlo})
	case WindowsBased:
		return json.Marshal(struct {
			BaseSlo
			*WindowBasedSlo
		}{slo.BaseSlo, slo.WindowBasedSlo})
	default:
		return nil, fmt.Errorf("unrecognized SLO type %s", slo.Type)
	}
}
