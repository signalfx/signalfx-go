package signalfx

import (
	"context"
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateAzureIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", true, http.StatusOK, nil, "integration/create_azure_success.json"))

	result, err := client.CreateAzureIntegration(context.Background(), &integration.AzureIntegration{
		Type: "Azure",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestGetAzureIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", true, http.StatusOK, nil, "integration/create_azure_success.json"))

	result, err := client.GetAzureIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateAzureIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", true, http.StatusOK, nil, "integration/create_azure_success.json"))

	result, err := client.UpdateAzureIntegration(context.Background(), "id", &integration.AzureIntegration{
		Type: "Azure",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeleteAzureIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteAzureIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
