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

	jsonBody := `{
		"type": "Webhook",
		"name": "webhoooooook",
		"enabled": true,
		"url": "https://webhook.site/<key>",
		"method": "POST",
		"payloadTemplate": "{\"incidentId\": \"{{{incidentId}}}\"}"
	}`
	mux.HandleFunc("/v2/integration", verifyRequestWithJsonBody(t, "POST", true, http.StatusOK, nil, jsonBody, "integration/create_webhook_success.json"))

	result, err := client.CreateWebhookIntegration(context.Background(), &integration.WebhookIntegration{
		Type:            "Webhook",
		Name:            "webhoooooook",
		Enabled:         true,
		Url:             "https://webhook.site/<key>",
		Method:          "POST",
		PayloadTemplate: "{\"incidentId\": \"{{{incidentId}}}\"}",
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

	jsonBody := `{
		"type": "Webhook",
		"name": "webhoooooook",
		"enabled": true,
		"url": "https://webhook.site/<key>",
		"method": "POST",
		"payloadTemplate": "{\"incidentId\": \"{{{incidentId}}}\"}"
	}`
	mux.HandleFunc("/v2/integration/id", verifyRequestWithJsonBody(t, "PUT", true, http.StatusOK, nil, jsonBody, "integration/create_webhook_success.json"))

	result, err := client.UpdateWebhookIntegration(context.Background(), "id", &integration.WebhookIntegration{
		Type:            "Webhook",
		Name:            "webhoooooook",
		Enabled:         true,
		Url:             "https://webhook.site/<key>",
		Method:          "POST",
		PayloadTemplate: "{\"incidentId\": \"{{{incidentId}}}\"}",
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
