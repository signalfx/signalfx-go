package signalfx

import (
	"context"
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateWebhookIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration", verifyRequest(t, "POST", true, http.StatusOK, nil, "integration/create_webhook_success.json"))

	result, err := client.CreateWebhookIntegration(context.Background(), &integration.WebhookIntegration{
		Type: "Webhook",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "webhoooooook", result.Name, "Name does not match")
}

func TestGetWebhookIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "GET", true, http.StatusOK, nil, "integration/create_webhook_success.json"))

	result, err := client.GetWebhookIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error getting integration")
	assert.Equal(t, "webhoooooook", result.Name, "Name does not match")
}

func TestUpdateWebhookIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "PUT", true, http.StatusOK, nil, "integration/create_webhook_success.json"))

	result, err := client.UpdateWebhookIntegration(context.Background(), "id", &integration.WebhookIntegration{
		Type: "Webhook",
	})
	assert.NoError(t, err, "Unexpected error creating integration")
	assert.Equal(t, "webhoooooook", result.Name, "Name does not match")
}

func TestDeleteWebhookIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/id", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteWebhookIntegration(context.Background(), "id")
	assert.NoError(t, err, "Unexpected error creating integration")
}
