package metric_ruleset

import (
	"encoding/json"
	"errors"
)

type MetricRuleset struct {
	// System-defined identifier for the metric ruleset
	Id               string             `json:"id"`
	Version          int64              `json:"version"`
	Name             string             `json:"name"`
	Enabled          bool               `json:"enabled"`
	Description      *string            `json:"description,omitempty"`
	Destination      int32              `json:"destination"`
	MetricMatcher    MetricMatcher      `json:"-"`
	RawMetricMatcher json.RawMessage    `json:"metricMatcher"`
	Aggregators      []MetricAggregator `json:"-"`
	RawAggregators   []json.RawMessage  `json:"aggregators,omitempty"`
}

func (mr *MetricRuleset) MarshalJSON() ([]byte, error) {
	type metricRuleset MetricRuleset

	if mr.Aggregators != nil {
		for _, a := range mr.Aggregators {
			b, err := json.Marshal(a)
			if err != nil {
				return nil, err
			}
			mr.RawAggregators = append(mr.RawAggregators, b)
		}
	}

	b, err := json.Marshal(mr.MetricMatcher)
	if err != nil {
		return nil, err
	}
	mr.RawMetricMatcher = b

	return json.Marshal((*metricRuleset)(mr))
}

func (mr *MetricRuleset) UnmarshalJSON(b []byte) error {
	type metricRuleset MetricRuleset

	err := json.Unmarshal(b, (*metricRuleset)(mr))
	if err != nil {
		return err
	}

	for _, raw := range mr.RawAggregators {

		// first unmarshals to base aggregator to get type
		var a BaseMetricAggregator
		err = json.Unmarshal(raw, &a)
		if err != nil {
			return err
		}

		var i MetricAggregator
		switch a.Type {
		case "rollup":
			i = &RollupAggregator{}
		default:
			return errors.New("unknown aggregator type")
		}

		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}

		mr.Aggregators = append(mr.Aggregators, i)
	}

	// first unmarshals to base matcher to get type
	var m BaseMetricMatcher
	err = json.Unmarshal(mr.RawMetricMatcher, &m)
	if err != nil {
		return err
	}

	var i MetricMatcher
	switch m.Type {
	case "simple":
		i = &SimpleMetricMatcher{}
	default:
		return errors.New("unknown matcher type")
	}

	err = json.Unmarshal(mr.RawMetricMatcher, i)
	if err != nil {
		return err
	}

	mr.MetricMatcher = i

	return nil
}
