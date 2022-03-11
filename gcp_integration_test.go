package signalfx

import (
	"context"
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateGCPIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", true, http.StatusOK, nil, "integration/create_gcp_success.json"))

	result, err := client.CreateGCPIntegration(context.Background(), &integration.GCPIntegration{
		Type: "GCP",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
	assert.Equal(t, integration.FiveMinutely, *result.PollRate, "PollRate does not match")
	assert.Equal(t, int64(300000), result.PollRateMs, "PollRateMs does not match")
}

func TestGetGCPIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", true, http.StatusOK, nil, "integration/create_gcp_success.json"))

	result, err := client.GetGCPIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateGCPIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", true, http.StatusOK, nil, "integration/create_gcp_success.json"))

	result, err := client.UpdateGCPIntegration(context.Background(), "id", &integration.GCPIntegration{
		Type: "GCP",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeleteGCPIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteGCPIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
