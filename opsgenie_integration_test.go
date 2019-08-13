package signalfx

import (
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateOpsgenieIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", http.StatusOK, nil, "integration/create_opsgenie_success.json"))

	result, err := client.CreateOpsgenieIntegration(&integration.OpsgenieIntegration{
		Type: "Opsgenie",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestGetOpsgenieIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", http.StatusOK, nil, "integration/create_opsgenie_success.json"))

	result, err := client.GetOpsgenieIntegration("id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateOpsgenieIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", http.StatusOK, nil, "integration/create_opsgenie_success.json"))

	result, err := client.UpdateOpsgenieIntegration("id", &integration.OpsgenieIntegration{
		Type: "Opsgenie",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeleteOpsgenieIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", http.StatusNoContent, nil, ""))

	err := client.DeleteOpsgenieIntegration("id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
