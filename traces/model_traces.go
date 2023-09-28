package traces

import "time"

type Span struct {
	DurationMicros int64                  `json:"durationMicros"`
	Logs           []Log                  `json:"logs"`
	ObjectType     string                 `json:"objectType"`
	OperationName  string                 `json:"operationName"`
	ParentId       string                 `json:"parentId"`
	ProcessTags    map[string]interface{} `json:"processTags"`
	ServiceName    string                 `json:"serviceName"`
	SpanId         string                 `json:"spanId"`
	Splunk         map[string]interface{} `json:"splunk"`
	StartTime      time.Time              `json:"startTime"`
	Tags           map[string]interface{} `json:"tags"`
	TraceId        string                 `json:"traceId"`
}

type Log struct {
	Fields    map[string]string `json:"fields"`
	Timestamp string            `json:"timestamp"`
}

type Trace []Span
