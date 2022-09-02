package navigator

type Metric struct {
	Id              string          `json:"id,omitempty"`
	Type            string          `json:"type,omitempty"`
	DisplayName     string          `json:"displayName,omitempty"`
	ValueLabel      string          `json:"valueLabel,omitempty"`
	ValueFormat     string          `json:"valueFormat,omitempty"`
	MetricSelectors []string        `json:"metricSelectors,omitempty"`
	Description     string          `json:"description,omitempty"`
	Job             *Job            `json:"job,omitempty"`
	ColoringScheme  *ColoringScheme `json:"coloringScheme,omitempty"`
}
