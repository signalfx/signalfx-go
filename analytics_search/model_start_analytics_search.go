package analytics_search

// StartAnalyticsSearchJob starts an analytics search job
type StartAnalyticsSearchGraphQLResponseData struct {
	Data struct {
		StartAnalyticsSearch struct {
			JobID string `json:"jobId"`
		} `json:"startAnalyticsSearch"`
	} `json:"data"`
}

type StartAnalyticsSearchGraphQLSharedParameters struct {
	TimeRangeMillis TimeRangeMillis `json:"timeRangeMillis"`
	Filters         []Filter        `json:"filters"`
}

type TimeRangeMillis struct {
	Gte int64 `json:"gte"`
	Lte int64 `json:"lte"`
}

type Filter struct {
	TraceFilter *TraceFilter `json:"traceFilter,omitempty"`
	SpanFilters []SpanFilter `json:"spanFilters,omitempty"`
	FilterType  string       `json:"filterType"`
}

type TraceFilter struct {
	Tags []Tag `json:"tags"`
}

type SpanFilter struct {
	Tags []Tag `json:"tags"`
}

type Tag struct {
	Tag       string   `json:"tag"`
	Operation string   `json:"operation"`
	Values    []string `json:"values"`
}

type StartAnalyticsSearchVariables struct {
	Parameters Parameters `json:"parameters,omitempty"`
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
