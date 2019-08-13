package signalfx

import (
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreatePagerDutyIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", http.StatusOK, nil, "integration/create_pd_success.json"))

	result, err := client.CreatePagerDutyIntegration(&integration.PagerDutyIntegration{
		Type: "PagerDuty",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestGetPagerDutyIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", http.StatusOK, nil, "integration/create_pd_success.json"))

	result, err := client.GetPagerDutyIntegration("id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdatePagerDutyIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", http.StatusOK, nil, "integration/create_pd_success.json"))

	result, err := client.UpdatePagerDutyIntegration("id", &integration.PagerDutyIntegration{
		Type: "PagerDuty",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeletePagerDutyIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", http.StatusNoContent, nil, ""))

	err := client.DeletePagerDutyIntegration("id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
