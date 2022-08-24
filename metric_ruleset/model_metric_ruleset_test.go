package metric_ruleset

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarshalGetMetricRulesetResponseWithAggregators(t *testing.T) {
	dest := FULL_FIDELITY
	metricRuleset := GetMetricRulesetResponse{
		Id:            stringToPointer("metric_ruleset_id"),
		Creator:       stringToPointer("user"),
		Created:       int64ToPointer(time.Now().UnixMilli()),
		Version:       int64ToPointer(11),
		Name:          stringToPointer("name"),
		Destination:   &dest,
		Enabled:       boolToPointer(true),
		LastUpdatedBy: stringToPointer("updater"),
		LastUpdated:   int64ToPointer(time.Now().UnixMilli()),
		MetricMatcher: &MetricMatcher{
			SimpleMetricMatcher: &SimpleMetricMatcher{
				Type:       "simple",
				MetricName: "metricName",
				Filters: []PropertyFilter{
					{
						Property:      stringToPointer("dim1"),
						NOT:           boolToPointer(true),
						PropertyValue: []string{"val1", "val2"},
					},
				},
			},
		},
		Aggregators: []RollupAggregator{
			{
				Type:       "rollup",
				OutputName: "newMetricName1",
				DimensionsToKeep: []string{
					"dim1",
					"dim2",
				},
			},
			{
				Type:       "rollup",
				OutputName: "newMetricName2",
				DimensionsToKeep: []string{
					"dim3",
					"dim4",
				},
			},
		},
	}
	payload, err := json.Marshal(&metricRuleset)

	var expectedPayload = fmt.Sprintf(
		`{
			"id":"metric_ruleset_id",
			"version":11,
			"name":"name",
			"enabled":true,
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
			],
			"creator":"user",
			"created":%d,
			"lastUpdatedBy":"updater",
			"lastUpdated":%d,
			"destination":"FullFidelity"
		}`, *metricRuleset.Created, *metricRuleset.LastUpdated)
	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.JSONEq(t, expectedPayload, string(payload), "payload does not match")
}

func TestMarshalMetricRulesetNoAggregators(t *testing.T) {
	dest := REALTIME_13_MO
	metricRuleset := GetMetricRulesetResponse{
		Id:          stringToPointer("metric_ruleset_id"),
		Version:     int64ToPointer(2),
		Name:        stringToPointer("name"),
		Description: stringToPointer("description"),
		Destination: &dest,
		MetricMatcher: &MetricMatcher{
			SimpleMetricMatcher: &SimpleMetricMatcher{
				Type:       "simple",
				MetricName: "metricName",
			},
		},
	}
	payload, err := json.Marshal(&metricRuleset)

	var expectedPayload = `{
			"id":"metric_ruleset_id",
			"version":2,
			"name":"name",
			"description":"description",
			"destination":"Realtime_13MO",
			"metricMatcher":{
				"type":"simple",
				"metricName":"metricName"
			}
		}`
	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.JSONEq(t, expectedPayload, string(payload), "payload does not match")
}

func TestUnmarshalMetricRulesetWithAggregators(t *testing.T) {
	metricRuleset := GetMetricRulesetResponse{}
	var jsonString = `{
			"id":"metric_ruleset_id",
			"version":2,
			"name":"name",
			"enabled":true,
			"description":"description",
			"destination":"Drop",
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

	dest := DROP
	expectedMetricRuleset := GetMetricRulesetResponse{
		Id:          stringToPointer("metric_ruleset_id"),
		Version:     int64ToPointer(2),
		Name:        stringToPointer("name"),
		Enabled:     boolToPointer(true),
		Description: stringToPointer("description"),
		Destination: &dest,
		MetricMatcher: &MetricMatcher{
			SimpleMetricMatcher: &SimpleMetricMatcher{
				Type:       "simple",
				MetricName: "metricName",
			},
		},
		Aggregators: []RollupAggregator{
			{
				Type:       "rollup",
				OutputName: "newMetricName1",
				DimensionsToKeep: []string{
					"dim1",
					"dim2",
				},
			},
			{
				Type:       "rollup",
				OutputName: "newMetricName2",
				DimensionsToKeep: []string{
					"dim3",
					"dim4",
				},
			},
		},
	}

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.EqualValues(t, expectedMetricRuleset, metricRuleset)
}

func TestUnmarshalMetricRulesetNoAggregators(t *testing.T) {
	metricRuleset := GetMetricRulesetResponse{}
	var jsonString = `{
			"id":"metric_ruleset_id",
			"version":2,
			"name":"name",
			"destination":"FullFidelity",
			"enabled": false,
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
	dest := FULL_FIDELITY
	expectedMetricRuleset := GetMetricRulesetResponse{
		Id:          stringToPointer("metric_ruleset_id"),
		Version:     int64ToPointer(2),
		Name:        stringToPointer("name"),
		Enabled:     boolToPointer(false),
		Destination: &dest,
		MetricMatcher: &MetricMatcher{
			SimpleMetricMatcher: &SimpleMetricMatcher{
				Type:       "simple",
				MetricName: "metricName",
				Filters: []PropertyFilter{
					{
						Property:      stringToPointer("dim1"),
						NOT:           boolToPointer(true),
						PropertyValue: []string{"val1", "val2"},
					},
				},
			},
		},
	}
	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.EqualValues(t, expectedMetricRuleset, metricRuleset)
}

func stringToPointer(s string) *string {
	return &s
}

func int64ToPointer(i int64) *int64 {
	return &i
}

func boolToPointer(b bool) *bool {
	return &b
}
