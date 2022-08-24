package metric_ruleset

type BaseMetricMatcher struct {
	Type string `json:"type"`
}

type SimpleMetricMatcher struct {
	BaseMetricMatcher
	MetricName string           `json:"metricName"`
	Filters    []PropertyFilter `json:"filters,omitempty"`
}

type MetricMatcher interface {
	GetMatcherType() string
}

func (bm BaseMetricMatcher) GetMatcherType() string {
	return bm.Type
}

type PropertyFilter struct {
	Property      string   `json:"property"`
	Not           bool     `json:"NOT"`
	PropertyValue []string `json:"propertyValue"`
}
