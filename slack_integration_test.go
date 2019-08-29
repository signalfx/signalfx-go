package signalfx

import (
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateSlackIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", http.StatusOK, nil, "integration/create_slack_success.json"))

	result, err := client.CreateSlackIntegration(&integration.SlackIntegration{
		Type: "Slack",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestGetSlackIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", http.StatusOK, nil, "integration/create_slack_success.json"))

	result, err := client.GetSlackIntegration("id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateSlackIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", http.StatusOK, nil, "integration/create_slack_success.json"))

	result, err := client.UpdateSlackIntegration("id", &integration.SlackIntegration{
		Type: "Slack",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeleteSlackIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", http.StatusNoContent, nil, ""))

	err := client.DeleteSlackIntegration("id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
