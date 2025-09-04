package integration

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
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
