package signalfx

import (
	"context"
	"github.com/signalfx/signalfx-go/metric_ruleset"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestCreateArchivedMetricRuleset(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(MetricRulesetApiURL, verifyRequest(t, http.MethodPost, true, http.StatusOK, nil, "metric_ruleset/create_archived_ruleset_success.json"))

	dest := "Archived"
	restorationID := "GWBTAQwAAAA"
	metricName := "container_cpu_utilization"
	ruleName := "TestRule"
	rulesetDescription := "Metric ruleset for container_cpu_utilization"
	exceptionRuleDescription := "exception rule 1"
	filterNot := false
	filterPropertyValue := "container_id"
	startTime := (time.Now().Unix() - 900) * 1000
	stopTime := (time.Now().Unix() - 200) * 1000
	result, err := client.CreateMetricRuleset(context.Background(), &metric_ruleset.CreateMetricRulesetRequest{
		MetricName:  metricName,
		Description: &rulesetDescription,
		Version:     1,
		ExceptionRules: []metric_ruleset.ExceptionRule{
			{
				Name:        ruleName,
				Description: &exceptionRuleDescription,
				Enabled:     true,
				Matcher: metric_ruleset.DimensionMatcher{
					Type: "dimension",
					Filters: []metric_ruleset.PropertyFilter{
						{
							NOT:           &filterNot,
							Property:      &filterPropertyValue,
							PropertyValue: []string{"cont_a", "cont_b"},
						},
					},
				},
				//add restoration
				Restoration: &metric_ruleset.ExceptionRuleRestorationFields{
					RestorationId: (*string)(&restorationID),
					StartTime:     &startTime,
					StopTime:      &stopTime,
				},
			},
		},
		RoutingRule: metric_ruleset.RoutingRule{
			Destination: (*string)(&dest),
		},
	})

	assert.NoError(t, err, "Unexpected error creating metric ruleset")
	assert.NotNil(t, restorationID, "Restoration ID is null")
	assert.Equal(t, restorationID, *result.ExceptionRules[0].Restoration.RestorationId, "Restoration ID does not match")
	assert.Equal(t, int64(1724793174572), *result.ExceptionRules[0].Restoration.StartTime, "StartTime does not match")
	assert.Equal(t, int64(1724796774661), *result.ExceptionRules[0].Restoration.StopTime, "StopTime does not match")
	assert.Equal(t, metricName, *result.MetricName, "MetricName does not match")
	assert.Equal(t, rulesetDescription, *result.Description, "Description does not match")
	assert.Equal(t, 1, len(result.ExceptionRules), "Unexpected length of exception rules array")
	assert.Equal(t, ruleName, result.ExceptionRules[0].Name, "Exception rule name does not match")
	assert.Equal(t, &exceptionRuleDescription, result.ExceptionRules[0].Description, "Exception rule description does not match")
	assert.Equal(t, 1, len(result.ExceptionRules[0].Matcher.Filters), "Unexpected length of exception rule filter array")
	assert.Equal(t, 2, len(result.ExceptionRules[0].Matcher.Filters[0].PropertyValue), "Unexpected length of exception rule filter property values array")
	assert.Equal(t, dest, *result.RoutingRule.Destination, "RoutingRule destination does not match expected")
}

func TestCreateMetricRuleset(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(MetricRulesetApiURL, verifyRequest(t, http.MethodPost, true, http.StatusOK, nil, "metric_ruleset/create_ruleset_success.json"))

	dest := "RealTime"
	dropDimensions := false
	ruleName := "TestRule"
	result, err := client.CreateMetricRuleset(context.Background(), &metric_ruleset.CreateMetricRulesetRequest{
		MetricName: "container_cpu_utilization",
		Version:    1,
		AggregationRules: []metric_ruleset.AggregationRule{
			{
				Name:    &ruleName,
				Enabled: true,
				Matcher: metric_ruleset.MetricMatcher{
					DimensionMatcher: &metric_ruleset.DimensionMatcher{
						Type:    "dimension",
						Filters: []metric_ruleset.PropertyFilter{},
					},
				},
				Aggregator: metric_ruleset.MetricAggregator{
					RollupAggregator: &metric_ruleset.RollupAggregator{
						Type:           "rollup",
						OutputName:     "container_cpu_utilization.by.sfx_realm.sfx_service.agg",
						Dimensions:     []string{"sfx_realm", "sfx_service"},
						DropDimensions: &dropDimensions,
					},
				},
			},
		},
		RoutingRule: metric_ruleset.RoutingRule{
			Destination: (*string)(&dest),
		},
	})

	assert.NoError(t, err, "Unexpected error creating metric ruleset")
	assert.Equal(t, "container_cpu_utilization", *result.MetricName, "MetricName does not match")
	assert.Equal(t, 1, len(result.AggregationRules), "Unexpected length of aggregation rules array")
	assert.Equal(t, "rollup", result.AggregationRules[0].Aggregator.RollupAggregator.Type, "Aggregation Rule type does not match expected")
	assert.Equal(t, dest, *result.RoutingRule.Destination, "RoutingRule destination does not match expected")
}

func TestGetMetricRuleset(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(MetricRulesetApiURL+"/TestId_MemoryUtilization", verifyRequest(t, http.MethodGet, true, http.StatusOK, nil, "metric_ruleset/get_ruleset_success.json"))

	result, err := client.GetMetricRuleset(context.Background(), "TestId_MemoryUtilization")
	assert.NoError(t, err, "Unexpected error getting metric ruleset")
	assert.Equal(t, "memory.utilization", *result.MetricName, "MetricName does not match")
}

func TestUpdateMetricRuleset(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(MetricRulesetApiURL+"/TestId", verifyRequest(t, http.MethodPut, true, http.StatusOK, nil, "metric_ruleset/update_ruleset_success.json"))

	metricName := "container_cpu_utilization"
	dest := "Drop"
	version := int64(2)
	dropDimensions := false
	ruleName := "UpdatedName"
	result, err := client.UpdateMetricRuleset(context.Background(), "TestId", &metric_ruleset.UpdateMetricRulesetRequest{
		MetricName: &metricName,
		Version:    &version,
		AggregationRules: []metric_ruleset.AggregationRule{
			{
				Name:    &ruleName,
				Enabled: false,
				Matcher: metric_ruleset.MetricMatcher{
					DimensionMatcher: &metric_ruleset.DimensionMatcher{
						Filters: []metric_ruleset.PropertyFilter{},
					},
				},
				Aggregator: metric_ruleset.MetricAggregator{
					RollupAggregator: &metric_ruleset.RollupAggregator{
						Type:           "rollup",
						OutputName:     "container_cpu_utilization.by.sfx_realm.sfx_service.agg",
						Dimensions:     []string{"sfx_realm", "sfx_service"},
						DropDimensions: &dropDimensions,
					},
				},
			},
		},
		RoutingRule: &metric_ruleset.RoutingRule{
			Destination: (*string)(&dest),
		},
	})

	assert.NoError(t, err, "Unexpected error updating metric ruleset")
	assert.Equal(t, metricName, *result.MetricName, "Name does not match")
	assert.Equal(t, 1, len(result.AggregationRules), "Unexpected length of aggregation rules array")
	assert.Equal(t, false, result.AggregationRules[0].Enabled)
	assert.Equal(t, "Drop", *result.RoutingRule.Destination)
}

func TestDeleteMetricRuleset(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(MetricRulesetApiURL+"/metric_ruleset_id", verifyRequest(t, http.MethodDelete, true, http.StatusNoContent, nil, ""))

	err := client.DeleteMetricRuleset(context.Background(), "metric_ruleset_id")
	assert.NoError(t, err, "Unexpected error deleting metric ruleset")
}

func TestGenerateAggregationMetricName(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(MetricRulesetApiURL+"/generateAggregationMetricName", verifyRequest(t, http.MethodPost, true, http.StatusOK, nil, "metric_ruleset/generate_aggregation_metric_name_success.txt"))

	dropDimensions := false
	result, err := client.GenerateAggregationMetricName(context.Background(), metric_ruleset.GenerateAggregationNameRequest{
		MetricName:     "cpu.utilization",
		Dimensions:     []string{"sfx_realm"},
		DropDimensions: &dropDimensions,
	})

	assert.NoError(t, err, "Unexpected error generating aggregation metric name")
	assert.Equal(t, "cpu.utilization.by.sfx_realm.agg", result)
}
