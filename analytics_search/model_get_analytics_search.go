package analytics_search

type GetAnalyticsSearchGraphQLResponseData struct {
	Data struct {
		GetAnalyticsSearch AnalyticsSearch `json:"getAnalyticsSearch"`
	} `json:"data"`
}

type AnalyticsSearch struct {
	JobID    string    `json:"jobId"`
	Sections []Section `json:"sections"`
}

type Section struct {
	SectionType         string               `json:"sectionType"`
	IsComplete          bool                 `json:"isComplete"`
	ProcessingRange     ProcessingRange      `json:"processingRange"`
	LegacyTraceExamples []LegacyTraceExample `json:"legacyTraceExamples,omitempty"`
}

type ProcessingRange struct {
	StartTimestampMillis int64 `json:"startTimestampMillis"`
	EndTimestampMillis   int64 `json:"endTimestampMillis"`
}

type LegacyTraceExample struct {
	InitiatingService      string             `json:"initiatingService"`
	InitiatingOperation    string             `json:"initiatingOperation"`
	InitiatingHttpMethod   string             `json:"initiatingHttpMethod"`
	InitiatingSpanWasError bool               `json:"initiatingSpanWasError"`
	InitiatorSpanType      *string            `json:"initiatorSpanType"` // Pointer because it can be null
	StartTimeMicros        int64              `json:"startTimeMicros"`
	DurationMicros         int64              `json:"durationMicros"`
	TraceId                string             `json:"traceId"`
	ServiceSpanCounts      []ServiceSpanCount `json:"serviceSpanCounts"`
}

type ServiceSpanCount struct {
	Service   string  `json:"service"`
	SpanCount int     `json:"spanCount"`
	Errors    []Error `json:"errors"`
}

type Error struct {
	SpanID      *string     `json:"spanID"` // Pointer because it can be null
	ErrorDetail ErrorDetail `json:"error"`
	IsRootCause bool        `json:"isRootCause"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type GetAnalyticsSearchVariables struct {
	JobID string `json:"jobId,omitempty"`
}
