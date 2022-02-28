package signalfx

import (
	"context"
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateAWSCloudWatchIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", true, http.StatusOK, nil, "integration/create_aws_success.json"))

	result, err := client.CreateAWSCloudWatchIntegration(context.Background(), &integration.AwsCloudWatchIntegration{
		Type: "AWSCloudWatch",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestGetAWSCloudWatchIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", true, http.StatusOK, nil, "integration/create_aws_success.json"))

	result, err := client.GetAWSCloudWatchIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateAWSCloudWatchIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", true, http.StatusOK, nil, "integration/create_aws_success.json"))

	result, err := client.UpdateAWSCloudWatchIntegration(context.Background(), "id", &integration.AwsCloudWatchIntegration{
		Type: "AWSCloudWatch",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateAWSCloudWatchIntegrationMetricStatsToSync(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", true, http.StatusOK, nil, "integration/create_aws_metric_stats_to_sync_success.json"))

	metricStatsToSync := map[string]map[string][]string{
		"AWS/EC2": {
			"NetworkPacketsIn": []string{"p95", "p99"},
		},
		"AWS/ECS": {
			"CPUReservation": []string{"mean"},
			"CPUUtilization": []string{"upper", "mean", "p95", "p99", "p99.5"},
		},
	}

	result, err := client.UpdateAWSCloudWatchIntegration(context.Background(), "id", &integration.AwsCloudWatchIntegration{
		Type:              "AWSCloudWatch",
		MetricStatsToSync: metricStatsToSync,
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, metricStatsToSync, result.MetricStatsToSync, "MetricStatsToSync does not match")
}

func TestUpdateAWSCloudWatchIntegrationMetricStreams(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", true, http.StatusOK, nil, "integration/create_aws_metric_streams_success.json"))

	result, err := client.UpdateAWSCloudWatchIntegration(context.Background(), "id", &integration.AwsCloudWatchIntegration{
		Type:                   "AWSCloudWatch",
		MetricStreamsSyncState: "ENABLED",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "ENABLED", result.MetricStreamsSyncState, "MetricStreamsSyncState does not match")
}

func TestUpdateAWSCloudWatchIntegrationLogsSyncState(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", true, http.StatusOK, nil, "integration/create_aws_logs_sync_state_success.json"))

	result, err := client.UpdateAWSCloudWatchIntegration(context.Background(), "id", &integration.AwsCloudWatchIntegration{
		Type:          "AWSCloudWatch",
		LogsSyncState: "ENABLED",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "ENABLED", result.LogsSyncState, "LogsSyncState does not match")
}

func TestDeleteAWSCloudWatchIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteAWSCloudWatchIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
