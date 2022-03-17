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
	serviceNowId                     = "FFF-yyyyZZ"
	serviceNowName                   = "SNOW integ"
	serviceNowInstanceName           = "inst.service-now.com"
	serviceNowIssueTypeIncident      = "Incident"
	serviceNowCreated                = int64(1647426008469)
	serviceNowCreatedBy              = "EEEzzzzYYY"
	serviceNowLastUpdated            = int64(1647426008953)
	serviceNowLastUpdatedBy          = "YY_zzzzyyy"
	serviceNowAlertTriggeredTemplate = "{\"short_description\": \"{{{messageTitle}}} (customized)\"}"
	serviceNowAlertResolvedTemplate  = "{\"close_notes\": \"{{{messageTitle}}} (customized close msg)\"}"
	serviceNowUsername               = "i.am"
	serviceNowPassword               = "my#pass"
)

var serviceNowRequestData = integration.ServiceNowIntegration{
	Enabled:                       true,
	Type:                          integration.SERVICE_NOW,
	Name:                          serviceNowName,
	InstanceName:                  serviceNowInstanceName,
	IssueType:                     serviceNowIssueTypeIncident,
	AlertTriggeredPayloadTemplate: serviceNowAlertTriggeredTemplate,
	AlertResolvedPayloadTemplate:  serviceNowAlertResolvedTemplate,
	Username:                      serviceNowUsername,
	Password:                      serviceNowPassword,
}

func verifyGetServiceNowIntegrationRequest(t *testing.T, r *http.Request) {
	verifyHeaders(t, r, true)
	verifyParams(t, r, url.Values{})

	assert.Equal(t, "GET", r.Method, "Incorrect HTTP method")

	body, err := ioutil.ReadAll(r.Body)

	assert.NoError(t, err, "Unexpected error getting request body")
	assert.Empty(t, body, "Unexpected request body")
}

func TestGetServiceNowIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL+"/"+serviceNowId,
		createResponse(t, http.StatusOK, "integration/create_service_now_success.json", verifyGetServiceNowIntegrationRequest))

	result, err := client.GetServiceNowIntegration(context.Background(), serviceNowId)

	if assert.NoError(t, err, "Unexpected error getting the integration") {
		assert.Equal(t, integration.SERVICE_NOW, result.Type, "Integration type does not match")
		assert.Equal(t, true, result.Enabled, "Integration not enabled")
		assert.Equal(t, serviceNowId, result.Id, "ID does not match")
		assert.Equal(t, serviceNowName, result.Name, "Name does not match")
		assert.Equal(t, serviceNowInstanceName, result.InstanceName, "Instance name does not match")
		assert.Equal(t, serviceNowIssueTypeIncident, result.IssueType, "Issue type does not match")
		assert.Equal(t, serviceNowCreated, result.Created, "Create time does not match")
		assert.Equal(t, serviceNowCreatedBy, result.Creator, "Creator does not match")
		assert.Equal(t, serviceNowLastUpdated, result.LastUpdated, "Update time does not match")
		assert.Equal(t, serviceNowLastUpdatedBy, result.LastUpdatedBy, "Updater does not match")
		assert.Equal(t, serviceNowAlertTriggeredTemplate, result.AlertTriggeredPayloadTemplate, "Alert triggered template does not match")
		assert.Equal(t, serviceNowAlertResolvedTemplate, result.AlertResolvedPayloadTemplate, "Alert resolved template does not match")
		assert.Empty(t, result.Username, "User name is not empty")
		assert.Empty(t, result.Password, "Password is not empty")
	}
}

func TestFailGetServiceNowIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL+"/"+serviceNowId,
		createResponse(t, http.StatusNotFound, "", verifyGetServiceNowIntegrationRequest))

	result, err := client.GetServiceNowIntegration(context.Background(), serviceNowId)

	assert.Error(t, err, "Error expected getting the integration")
	assert.Nil(t, result, "No result expected")
}

func verifyCreateServiceNowIntegrationRequest(t *testing.T, r *http.Request) {
	verifyHeaders(t, r, true)
	verifyParams(t, r, url.Values{})

	assert.Equal(t, "POST", r.Method, "Incorrect HTTP method")

	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)

	if assert.NoError(t, err, "Unexpected error getting request body") {
		//goland:noinspection GoRedundantConversion
		assert.Equal(t, string(integration.SERVICE_NOW), body["type"], "Sent integration type does not match")
		assert.Equal(t, true, body["enabled"], "Sent integration not enabled")
		assert.Empty(t, body["id"], "Sent unexpected ID")
		assert.Equal(t, serviceNowName, body["name"], "Sent name does not match")
		assert.Equal(t, serviceNowInstanceName, body["instanceName"], "Sent instance name does not match")
		assert.Equal(t, serviceNowIssueTypeIncident, body["issueType"], "Sent issue type does not match")
		assert.Empty(t, body["created"], "Sent unexpected create time")
		assert.Empty(t, body["creator"], "Sent unexpected creator")
		assert.Empty(t, body["lastUpdated"], "Sent unexpected update time")
		assert.Empty(t, body["lastUpdatedBy"], "Sent unexpected updater")
		assert.Equal(t, serviceNowAlertTriggeredTemplate, body["alertTriggeredPayloadTemplate"], "Sent alert triggered template does not match")
		assert.Equal(t, serviceNowAlertResolvedTemplate, body["alertResolvedPayloadTemplate"], "Sent alert resolved template does not match")
		assert.Equal(t, serviceNowUsername, body["username"], "Sent user name does not match")
		assert.Equal(t, serviceNowPassword, body["password"], "Sent password does not match")
	}
}

func TestCreateServiceNowIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL,
		createResponse(t, http.StatusOK, "integration/create_service_now_success.json", verifyCreateServiceNowIntegrationRequest))

	result, err := client.CreateServiceNowIntegration(context.Background(), &serviceNowRequestData)

	if assert.NoError(t, err, "Unexpected error creating an integration") {
		assert.Equal(t, integration.SERVICE_NOW, result.Type, "Integration type does not match")
		assert.Equal(t, true, result.Enabled, "Integration not enabled")
		assert.Equal(t, serviceNowId, result.Id, "ID does not match")
		assert.Equal(t, serviceNowName, result.Name, "Name does not match")
		assert.Equal(t, serviceNowInstanceName, result.InstanceName, "Instance name does not match")
		assert.Equal(t, serviceNowIssueTypeIncident, result.IssueType, "Issue type does not match")
		assert.Equal(t, serviceNowCreated, result.Created, "Create time does not match")
		assert.Equal(t, serviceNowCreatedBy, result.Creator, "Creator does not match")
		assert.Equal(t, serviceNowLastUpdated, result.LastUpdated, "Update time does not match")
		assert.Equal(t, serviceNowLastUpdatedBy, result.LastUpdatedBy, "Updater does not match")
		assert.Equal(t, serviceNowAlertTriggeredTemplate, result.AlertTriggeredPayloadTemplate, "Alert triggered template does not match")
		assert.Equal(t, serviceNowAlertResolvedTemplate, result.AlertResolvedPayloadTemplate, "Alert resolved template does not match")
		assert.Empty(t, result.Username, "User name is not empty")
		assert.Empty(t, result.Password, "Password is not empty")
	}
}

func TestFailCreateServiceNowIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL,
		createResponse(t, http.StatusForbidden, "", verifyCreateServiceNowIntegrationRequest))

	result, err := client.CreateServiceNowIntegration(context.Background(), &serviceNowRequestData)

	assert.Error(t, err, "Error expected creating an integration")
	assert.Nil(t, result, "No result expected")
}

func verifyUpdateServiceNowIntegrationRequest(t *testing.T, r *http.Request) {
	verifyHeaders(t, r, true)
	verifyParams(t, r, url.Values{})

	assert.Equal(t, "PUT", r.Method, "Incorrect HTTP method")

	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)

	if assert.NoError(t, err, "Unexpected error getting request body") {
		//goland:noinspection GoRedundantConversion
		assert.Equal(t, string(integration.SERVICE_NOW), body["type"], "Sent integration type does not match")
		assert.Equal(t, true, body["enabled"], "Sent integration not enabled")
		assert.Empty(t, body["id"], "Sent unexpected ID")
		assert.Equal(t, serviceNowName, body["name"], "Sent name does not match")
		assert.Equal(t, serviceNowInstanceName, body["instanceName"], "Sent instance name does not match")
		assert.Equal(t, serviceNowIssueTypeIncident, body["issueType"], "Sent issue type does not match")
		assert.Empty(t, body["created"], "Sent unexpected create time")
		assert.Empty(t, body["creator"], "Sent unexpected creator")
		assert.Empty(t, body["lastUpdated"], "Sent unexpected update time")
		assert.Empty(t, body["lastUpdatedBy"], "Sent unexpected updater")
		assert.Equal(t, serviceNowAlertTriggeredTemplate, body["alertTriggeredPayloadTemplate"], "Sent alert triggered template does not match")
		assert.Equal(t, serviceNowAlertResolvedTemplate, body["alertResolvedPayloadTemplate"], "Sent alert resolved template does not match")
		assert.Equal(t, serviceNowUsername, body["username"], "Sent user name does not match")
		assert.Equal(t, serviceNowPassword, body["password"], "Sent password does not match")
	}
}

func TestUpdateServiceNowIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL+"/"+serviceNowId,
		createResponse(t, http.StatusOK, "integration/create_service_now_success.json", verifyUpdateServiceNowIntegrationRequest))

	result, err := client.UpdateServiceNowIntegration(context.Background(), serviceNowId, &serviceNowRequestData)

	if assert.NoError(t, err, "Unexpected error updating the integration") {
		assert.Equal(t, integration.SERVICE_NOW, result.Type, "Integration type does not match")
		assert.Equal(t, true, result.Enabled, "Integration not enabled")
		assert.Equal(t, serviceNowId, result.Id, "ID does not match")
		assert.Equal(t, serviceNowName, result.Name, "Name does not match")
		assert.Equal(t, serviceNowInstanceName, result.InstanceName, "Instance name does not match")
		assert.Equal(t, serviceNowIssueTypeIncident, result.IssueType, "Issue type does not match")
		assert.Equal(t, serviceNowCreated, result.Created, "Create time does not match")
		assert.Equal(t, serviceNowCreatedBy, result.Creator, "Creator does not match")
		assert.Equal(t, serviceNowLastUpdated, result.LastUpdated, "Update time does not match")
		assert.Equal(t, serviceNowLastUpdatedBy, result.LastUpdatedBy, "Updater does not match")
		assert.Equal(t, serviceNowAlertTriggeredTemplate, result.AlertTriggeredPayloadTemplate, "Alert triggered template does not match")
		assert.Equal(t, serviceNowAlertResolvedTemplate, result.AlertResolvedPayloadTemplate, "Alert resolved template does not match")
		assert.Empty(t, result.Username, "User name is not empty")
		assert.Empty(t, result.Password, "Password is not empty")
	}
}

func TestFailUpdateServiceNowIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL+"/"+serviceNowId,
		createResponse(t, http.StatusForbidden, "", verifyUpdateServiceNowIntegrationRequest))

	result, err := client.UpdateServiceNowIntegration(context.Background(), serviceNowId, &serviceNowRequestData)

	assert.Error(t, err, "Error expected updating the integration")
	assert.Nil(t, result, "No result expected")
}

func verifyDeleteServiceNowIntegrationRequest(t *testing.T, r *http.Request) {
	verifyHeaders(t, r, true)
	verifyParams(t, r, url.Values{})

	assert.Equal(t, "DELETE", r.Method, "Incorrect HTTP method")

	body, err := ioutil.ReadAll(r.Body)

	assert.NoError(t, err, "Unexpected error getting request body")
	assert.Empty(t, body, "Unexpected request body")
}

func TestDeleteServiceNowIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL+"/"+serviceNowId,
		createResponse(t, http.StatusNoContent, "", verifyDeleteServiceNowIntegrationRequest))

	err := client.DeleteServiceNowIntegration(context.Background(), serviceNowId)

	assert.NoError(t, err, "Unexpected error deleting the integration")
}

func TestFailDeleteServiceNowIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(IntegrationAPIURL+"/"+serviceNowId,
		createResponse(t, http.StatusNotFound, "", verifyDeleteServiceNowIntegrationRequest))

	err := client.DeleteServiceNowIntegration(context.Background(), serviceNowId)

	assert.Error(t, err, "Error expected deleting the integration")
}
