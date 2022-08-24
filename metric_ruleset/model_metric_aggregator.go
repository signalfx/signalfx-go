package metric_ruleset

type BaseMetricAggregator struct {
	Type string `json:"type"`
}

type RollupAggregator struct {
	BaseMetricAggregator
	OutputName       string   `json:"outputName"`
	DimensionsToKeep []string `json:"dimensionsToKeep"`
}

type MetricAggregator interface {
	GetAggregatorType() string
}

func (ba BaseMetricAggregator) GetAggregatorType() string {
	return ba.Type
}
