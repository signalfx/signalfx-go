package signalfx

import (
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateVictorOpsIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", http.StatusOK, nil, "integration/create_victor_ops_success.json"))

	result, err := client.CreateVictorOpsIntegration(&integration.VictorOpsIntegration{
		Type: "VictorOps",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestGetVictorOpsIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", http.StatusOK, nil, "integration/create_victor_ops_success.json"))

	result, err := client.GetVictorOpsIntegration("id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateVictorOpsIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", http.StatusOK, nil, "integration/create_victor_ops_success.json"))

	result, err := client.UpdateVictorOpsIntegration("id", &integration.VictorOpsIntegration{
		Type: "VictorOps",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeleteVictorOpsIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", http.StatusNoContent, nil, ""))

	err := client.DeleteVictorOpsIntegration("id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
