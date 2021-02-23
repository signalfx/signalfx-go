package signalfx

import (
	"context"
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreatePagerDutyIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", true, http.StatusOK, nil, "integration/create_pd_success.json"))

	result, err := client.CreatePagerDutyIntegration(context.Background(), &integration.PagerDutyIntegration{
		Type: "PagerDuty",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestGetPagerDutyIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", true, http.StatusOK, nil, "integration/create_pd_success.json"))

	result, err := client.GetPagerDutyIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestGetPagerDutyIntegrationByName(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "GET", true, http.StatusOK, nil, "integration/get_by_name_pd_success.json"))

	result, err := client.GetPagerDutyIntegrationByName(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name)
}

func TestUpdatePagerDutyIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", true, http.StatusOK, nil, "integration/create_pd_success.json"))

	result, err := client.UpdatePagerDutyIntegration(context.Background(), "id", &integration.PagerDutyIntegration{
		Type: "PagerDuty",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeletePagerDutyIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeletePagerDutyIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
