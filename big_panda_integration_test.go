package signalfx

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/signalfx/signalfx-go/integration"
	"github.com/stretchr/testify/assert"
)

const (
	bigPandaId                     = "FFF-yyyyZZ"
	bigPandaName                   = "BigPanda integ"
	bigPandaCreated                = int64(1647426008469)
	bigPandaCreatedBy              = "EEEzzzzYYY"
	bigPandaLastUpdated            = int64(1647426008953)
	bigPandaLastUpdatedBy          = "YY_zzzzyyy"
	bigPandaAppKey                 = "app-key"
	bigPandaToken                  = "bearer-token"
	bigPandaAlertTriggeredTemplate = "{\"status\":\"critical\",\"summary\":\"{{{messageTitle}}}\"}"
	bigPandaAlertResolvedTemplate  = "{\"status\":\"ok\",\"summary\":\"{{{messageTitle}}}\"}"
)

var bigPandaRequestData = integration.BigPandaIntegration{
	Enabled:                       true,
	Type:                          integration.BIG_PANDA,
	Name:                          bigPandaName,
	AppKey:                        bigPandaAppKey,
	Token:                         bigPandaToken,
	AlertTriggeredPayloadTemplate: bigPandaAlertTriggeredTemplate,
	AlertResolvedPayloadTemplate:  bigPandaAlertResolvedTemplate,
}

var bigPandaResponseData = integration.BigPandaIntegration{
	Created:                       bigPandaCreated,
	Creator:                       bigPandaCreatedBy,
	Enabled:                       true,
	Id:                            bigPandaId,
	LastUpdated:                   bigPandaLastUpdated,
	LastUpdatedBy:                 bigPandaLastUpdatedBy,
	Name:                          bigPandaName,
	Type:                          integration.BIG_PANDA,
	AlertTriggeredPayloadTemplate: bigPandaAlertTriggeredTemplate,
	AlertResolvedPayloadTemplate:  bigPandaAlertResolvedTemplate,
}

func verifyGetBigPandaIntegrationRequest(t *testing.T, r *http.Request) {
	verifyHeaders(t, r, true)
	verifyParams(t, r, url.Values{})

	assert.Equal(t, "GET", r.Method, "Incorrect HTTP method")

	body, err := ioutil.ReadAll(r.Body)

	assert.NoError(t, err, "Unexpected error getting request body")
	assert.Empty(t, body, "Unexpected request body")
}

func TestGetBigPandaIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL+"/"+bigPandaId,
		createResponse(t, http.StatusOK, "integration/create_big_panda_success.json", verifyGetBigPandaIntegrationRequest))

	result, err := client.GetBigPandaIntegration(context.Background(), bigPandaId)

	assert.NoError(t, err, "Unexpected error getting the integration")
	assert.Equal(t, &bigPandaResponseData, result, "Integration does not match")
}

func verifyBigPandaIntegrationWriteRequest(t *testing.T, r *http.Request, method string) {
	verifyHeaders(t, r, true)
	verifyParams(t, r, url.Values{})

	assert.Equal(t, method, r.Method, "Incorrect HTTP method")

	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)

	if assert.NoError(t, err, "Unexpected error getting request body") {
		assert.Equal(t, string(integration.BIG_PANDA), body["type"], "Sent integration type does not match")
		assert.Equal(t, true, body["enabled"], "Sent integration not enabled")
		assert.Empty(t, body["id"], "Sent unexpected ID")
		assert.Equal(t, bigPandaName, body["name"], "Sent name does not match")
		assert.Empty(t, body["created"], "Sent unexpected create time")
		assert.Empty(t, body["creator"], "Sent unexpected creator")
		assert.Empty(t, body["lastUpdated"], "Sent unexpected update time")
		assert.Empty(t, body["lastUpdatedBy"], "Sent unexpected updater")
		assert.Equal(t, bigPandaAppKey, body["appKey"], "Sent app key does not match")
		assert.Equal(t, bigPandaToken, body["token"], "Sent token does not match")
		assert.Equal(t, bigPandaAlertTriggeredTemplate, body["alertTriggeredPayloadTemplate"], "Sent alert triggered template does not match")
		assert.Equal(t, bigPandaAlertResolvedTemplate, body["alertResolvedPayloadTemplate"], "Sent alert resolved template does not match")
	}
}

func verifyCreateBigPandaIntegrationRequest(t *testing.T, r *http.Request) {
	verifyBigPandaIntegrationWriteRequest(t, r, "POST")
}

func TestCreateBigPandaIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL,
		createResponse(t, http.StatusOK, "integration/create_big_panda_success.json", verifyCreateBigPandaIntegrationRequest))

	result, err := client.CreateBigPandaIntegration(context.Background(), &bigPandaRequestData)

	if assert.NoError(t, err, "Unexpected error creating an integration") {
		assert.Equal(t, integration.BIG_PANDA, result.Type, "Integration type does not match")
		assert.Equal(t, true, result.Enabled, "Integration not enabled")
		assert.Equal(t, bigPandaId, result.Id, "ID does not match")
		assert.Equal(t, bigPandaName, result.Name, "Name does not match")
		assert.Equal(t, bigPandaAlertTriggeredTemplate, result.AlertTriggeredPayloadTemplate, "Alert triggered template does not match")
		assert.Equal(t, bigPandaAlertResolvedTemplate, result.AlertResolvedPayloadTemplate, "Alert resolved template does not match")
		assert.Empty(t, result.AppKey, "App key is not empty")
		assert.Empty(t, result.Token, "Token is not empty")
	}
}

func verifyCreateBigPandaIntegrationWithoutTemplatesRequest(t *testing.T, r *http.Request) {
	verifyHeaders(t, r, true)
	verifyParams(t, r, url.Values{})

	assert.Equal(t, "POST", r.Method, "Incorrect HTTP method")

	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)

	if assert.NoError(t, err, "Unexpected error getting request body") {
		assert.Equal(t, string(integration.BIG_PANDA), body["type"], "Sent integration type does not match")
		assert.Equal(t, true, body["enabled"], "Sent integration not enabled")
		assert.Equal(t, bigPandaName, body["name"], "Sent name does not match")
		assert.Equal(t, bigPandaAppKey, body["appKey"], "Sent app key does not match")
		assert.Equal(t, bigPandaToken, body["token"], "Sent token does not match")
		assert.NotContains(t, body, "alertTriggeredPayloadTemplate", "Triggered template should be omitted when not configured")
		assert.NotContains(t, body, "alertResolvedPayloadTemplate", "Resolved template should be omitted when not configured")
	}
}

func TestCreateBigPandaIntegrationWithoutTemplates(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL,
		createResponse(t, http.StatusOK, "integration/create_big_panda_success.json", verifyCreateBigPandaIntegrationWithoutTemplatesRequest))

	result, err := client.CreateBigPandaIntegration(context.Background(), &integration.BigPandaIntegration{
		Enabled: true,
		Type:    integration.BIG_PANDA,
		Name:    bigPandaName,
		AppKey:  bigPandaAppKey,
		Token:   bigPandaToken,
	})

	assert.NoError(t, err, "Unexpected error creating an integration without templates")
	assert.NotNil(t, result, "Result expected")
}

func verifyUpdateBigPandaIntegrationRequest(t *testing.T, r *http.Request) {
	verifyBigPandaIntegrationWriteRequest(t, r, "PUT")
}

func TestUpdateBigPandaIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL+"/"+bigPandaId,
		createResponse(t, http.StatusOK, "integration/create_big_panda_success.json", verifyUpdateBigPandaIntegrationRequest))

	result, err := client.UpdateBigPandaIntegration(context.Background(), bigPandaId, &bigPandaRequestData)

	if assert.NoError(t, err, "Unexpected error updating the integration") {
		assert.Equal(t, integration.BIG_PANDA, result.Type, "Integration type does not match")
		assert.Equal(t, bigPandaId, result.Id, "ID does not match")
		assert.Equal(t, bigPandaAlertTriggeredTemplate, result.AlertTriggeredPayloadTemplate, "Alert triggered template does not match")
		assert.Equal(t, bigPandaAlertResolvedTemplate, result.AlertResolvedPayloadTemplate, "Alert resolved template does not match")
	}
}

func TestDeleteBigPandaIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL+"/"+bigPandaId,
		createResponse(t, http.StatusNoContent, "", verifyDeleteBigPandaIntegrationRequest))

	err := client.DeleteBigPandaIntegration(context.Background(), bigPandaId)

	assert.NoError(t, err, "Unexpected error deleting the integration")
}

func verifyDeleteBigPandaIntegrationRequest(t *testing.T, r *http.Request) {
	verifyHeaders(t, r, true)
	verifyParams(t, r, url.Values{})

	assert.Equal(t, "DELETE", r.Method, "Incorrect HTTP method")

	body, err := ioutil.ReadAll(r.Body)

	assert.NoError(t, err, "Unexpected error getting request body")
	assert.Empty(t, body, "Unexpected request body")
}
