package signalfx

import (
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestDeleteIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "DELETE", http.StatusNoContent, nil, ""))

	err := client.DeleteIntegration("string")
	assert.NoError(t, err, "Unexpected error deleting integration")
}

func TestDeleteMissingIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "DELETE", http.StatusNotFound, nil, ""))

	err := client.DeleteIntegration("string")
	assert.Error(t, err, "Should get error error deleting missing integration")
}

func TestGetIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "GET", http.StatusOK, nil, "integration/get_success.json"))

	result, err := client.GetIntegration("string")
	assert.NoError(t, err, "Unexpected error getting integration")
	id := result["id"].(string)
	assert.Equal(t, id, "string", "Missing ID")
}

func TestGetMissingIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "GET", http.StatusNotFound, nil, ""))

	result, err := client.GetIntegration("string")
	assert.Error(t, err, "Should get an error getting missing integration")
	assert.Nil(t, result, "Should get a nil result from a missing integration")
}

func TestCreateAWSCloudWatchIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", http.StatusOK, nil, "integration/create_aws_success.json"))

	result, err := client.CreateAWSCloudWatchIntegration(&integration.AwsCloudWatchIntegration{
		Type: "AWSCloudWatch",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestGetAWSCloudWatchIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", http.StatusOK, nil, "integration/create_aws_success.json"))

	result, err := client.GetAWSCloudWatchIntegration("id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateAWSCloudWatchIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", http.StatusOK, nil, "integration/create_aws_success.json"))

	result, err := client.UpdateAWSCloudWatchIntegration("id", &integration.AwsCloudWatchIntegration{
		Type: "AWSCloudWatch",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeleteAWSCloudWatchIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", http.StatusOK, nil, ""))

	err := client.DeleteAWSCloudWatchIntegration("id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
