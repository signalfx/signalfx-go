package metric_ruleset

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalMetricRulesetWithAggregators(t *testing.T) {
	metricRuleset := MetricRuleset{
		Id: "metric_ruleset_id",
		Version: 2,
		Name: "name",
		Destination: 1,
		MetricMatcher: &SimpleMetricMatcher{
			BaseMetricMatcher: BaseMetricMatcher{ Type: "simple" },
			MetricName: "metricName",
			Filters: []PropertyFilter{
				{
					Property: "dim1",
					Not: true,
					PropertyValue: []string { "val1", "val2" },
				},
			},
		},
		Aggregators: []MetricAggregator{
			&RollupAggregator{
				BaseMetricAggregator: BaseMetricAggregator{Type: "rollup"},
				OutputName: "newMetricName1",
				DimensionsToKeep: []string{
					"dim1",
					"dim2",
				},
			},
			&RollupAggregator{
				BaseMetricAggregator: BaseMetricAggregator{Type: "rollup"},
				OutputName: "newMetricName2",
				DimensionsToKeep: []string{
					"dim3",
					"dim4",
				},
			},
		},
	}
	payload, err := json.Marshal(&metricRuleset)

	var expectedPayload =
		`{
			"id":"metric_ruleset_id",
			"version":2,
			"name":"name",
			"enabled":false,
			"destination":1,
			"metricMatcher":{
				"type":"simple",
				"metricName":"metricName",
				"filters":[
					{
						"property":"dim1",
						"NOT":true,
						"propertyValue":[
							"val1","val2"
						]
					}
				]
			},
			"aggregators":[
				{
					"type":"rollup",
					"outputName":"newMetricName1",
					"dimensionsToKeep":["dim1","dim2"]
				},
				{
					"type":"rollup",
					"outputName":"newMetricName2",
					"dimensionsToKeep":["dim3","dim4"]
				}
			]
		}`
	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.JSONEq(t, expectedPayload, string(payload), "payload does not match")
}

func TestMarshalMetricRulesetNoAggregators(t *testing.T) {
	description := "description"
	metricRuleset := MetricRuleset{
		Id: "metric_ruleset_id",
		Version: 2,
		Name: "name",
		Description: &description,
		Enabled: true,
		Destination: 1,
		MetricMatcher: &SimpleMetricMatcher{
			BaseMetricMatcher: BaseMetricMatcher{ Type: "simple" },
			MetricName: "metricName",
		},
	}
	payload, err := json.Marshal(&metricRuleset)

	var expectedPayload =
		`{
			"id":"metric_ruleset_id",
			"version":2,
			"name":"name",
			"enabled":true,
			"description":"description",
			"destination":1,
			"metricMatcher":{
				"type":"simple",
				"metricName":"metricName"
			}
		}`
	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.JSONEq(t, expectedPayload, string(payload), "payload does not match")
}

func TestUnmarshalMetricRulesetWithAggregators(t *testing.T) {
	metricRuleset := MetricRuleset{}
	var jsonString =
		`{
			"id":"metric_ruleset_id",
			"version":2,
			"name":"name",
			"enabled":true,
			"description":"description",
			"destination":2,
			"metricMatcher":{
				"type":"simple",
				"metricName":"metricName"
			},
			"aggregators":[
				{
					"type":"rollup",
					"outputName":"newMetricName1",
					"dimensionsToKeep":["dim1","dim2"]
				},
				{
					"type":"rollup",
					"outputName":"newMetricName2",
					"dimensionsToKeep":["dim3","dim4"]
				}
			]
		}`
	err := json.Unmarshal([]byte(jsonString), &metricRuleset)

	description := "description"
	expectedMetricRuleset := MetricRuleset{
		Id: "metric_ruleset_id",
		Version: 2,
		Name: "name",
		Enabled: true,
		Description: &description,
		Destination: 2,
		MetricMatcher: &SimpleMetricMatcher{
			BaseMetricMatcher: BaseMetricMatcher{ Type: "simple" },
			MetricName: "metricName",
		},
		Aggregators: []MetricAggregator{
			&RollupAggregator{
				BaseMetricAggregator: BaseMetricAggregator{Type: "rollup"},
				OutputName: "newMetricName1",
				DimensionsToKeep: []string{
					"dim1",
					"dim2",
				},
			},
			&RollupAggregator{
				BaseMetricAggregator: BaseMetricAggregator{Type: "rollup"},
				OutputName: "newMetricName2",
				DimensionsToKeep: []string{
					"dim3",
					"dim4",
				},
			},
		},
	}
	metricRuleset.RawMetricMatcher = nil
	metricRuleset.RawAggregators = nil

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.EqualValues(t, expectedMetricRuleset, metricRuleset)
}

func TestUnmarshalMetricRulesetNoAggregators(t *testing.T) {
	metricRuleset := MetricRuleset{}
	var jsonString =
		`{
			"id":"metric_ruleset_id",
			"version":2,
			"name":"name",
			"destination":1,
			"metricMatcher":{
				"type":"simple",
				"metricName":"metricName",
				"filters":[
					{
						"property":"dim1",
						"NOT":true,
						"propertyValue":[
							"val1","val2"
						]
					}
				]
			}
		}`
	err := json.Unmarshal([]byte(jsonString), &metricRuleset)
	metricRuleset.RawMetricMatcher = nil
	metricRuleset.RawAggregators = nil
	expectedMetricRuleset := MetricRuleset{
		Id: "metric_ruleset_id",
		Version: 2,
		Name: "name",
		Enabled: false,
		Destination: 1,
		MetricMatcher: &SimpleMetricMatcher{
			BaseMetricMatcher: BaseMetricMatcher{ Type: "simple" },
			MetricName: "metricName",
			Filters: []PropertyFilter{
				{
					Property: "dim1",
					Not: true,
					PropertyValue: []string {"val1","val2"},
				},
			},
		},
	}
	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.EqualValues(t, expectedMetricRuleset, metricRuleset)
}
