package navigator

type Job struct {
	Resolution int32        `json:"resolution,omitempty"`
	Template   string       `json:"template,omitempty"`
	VarName    string       `json:"varName,omitempty"`
	Filters    []*JobFilter `json:"filters,omitempty"`
}
