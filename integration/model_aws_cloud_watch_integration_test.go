package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAwsCloudWatchIntegration_MarshallingToAndFromJSON(t *testing.T) {
	cwIntegration := &AwsCloudWatchIntegration{
		Name:     "abc",
		Services: []AwsService{"AWS/Foo"},
		NamespaceSyncRules: []*AwsNameSpaceSyncRule{
			{Namespace: "AWS/Bar"},
		},
	}
	jsonBytes, err := json.Marshal(cwIntegration)
	assert.NoError(t, err, "Error marshalling CW integration")

	unmarshalledCwIntegration := &AwsCloudWatchIntegration{}
	err = json.Unmarshal(jsonBytes, &unmarshalledCwIntegration)
	assert.NoError(t, err, "Error unmarshalling CW integration")
	assert.Equal(t, cwIntegration, unmarshalledCwIntegration, "Marshalling and unmarshalling mismatch")
	assert.Equal(t, AwsService("AWS/Foo"), unmarshalledCwIntegration.Services[0], "Service does not match")
}

func TestDerivedTypeComparison(t *testing.T) {
	assert.NotEqual(t, AwsService("AWS/Foo"), "AWS/Foo", "Same value, different types")

	assert.True(t, AwsService("AWS/Foo") == "AWS/Foo", "Golang revolution!")
	assert.True(t, AwsService("AWS/Bar") == "AWS/Bar", "Golang revolution!")
	assert.True(t, AwsService("") == "", "Golang revolution!")
}

func TestMarshalAwsCloudWatchIntegrationWithInactiveMetricsPollRate(t *testing.T) {
	cwIntegration := AwsCloudWatchIntegration{
		InactiveMetricsPollRate: 60000,
	}
	payload, err := json.Marshal(&cwIntegration)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","inactiveMetricsPollRate":60000,"metricStreamsManagedExternally":false}`, string(payload), "payload does not match")
	assert.Equal(t, int64(60000), cwIntegration.InactiveMetricsPollRate, "InactiveMetricsPollRate does not match")
}

func TestUnmarshalAwsCloudWatchIntegrationWithInactiveMetricsPollRate(t *testing.T) {
	expectedValue := int64(60000)

	cwIntegration := AwsCloudWatchIntegration{}
	err := json.Unmarshal([]byte(`{"inactiveMetricsPollRate":60000}`), &cwIntegration)

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.Equal(t, expectedValue, cwIntegration.InactiveMetricsPollRate, "InactiveMetricsPollRate does not match")
}
