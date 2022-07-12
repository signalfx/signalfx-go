package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalAzureIntegrationWithPollRate(t *testing.T) {
	pollRate := OneMinutely
	azureInt := AzureIntegration{
		PollRate: &pollRate,
	}
	payload, err := json.Marshal(&azureInt)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","pollRate":60000,"syncGuestOsNamespaces":false}`, string(payload), "payload does not match")
	assert.Equal(t, int64(0), azureInt.PollRateMs, "PollRateMs has been changed")
}

func TestMarshalAzureIntegrationWithPollRateMs(t *testing.T) {
	payload, err := json.Marshal(AzureIntegration{
		PollRateMs: 90000,
	})

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","pollRate":90000,"syncGuestOsNamespaces":false}`, string(payload), "payload does not match")
}

func TestUnmarshalAzureIntegrationWithPollRate(t *testing.T) {
	azure := AzureIntegration{}
	err := json.Unmarshal([]byte(`{"pollRate":60000}`), &azure)

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.Equal(t, OneMinutely, *azure.PollRate, "PollRate does not match")
	assert.Equal(t, int64(60000), azure.PollRateMs, "PollRateMs does not match")
}

func TestUnmarshalAzureIntegrationWithPollRateMs(t *testing.T) {
	azure := AzureIntegration{}
	err := json.Unmarshal([]byte(`{"pollRate":90000}`), &azure)

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.Nil(t, azure.PollRate, "PollRate does not match")
	assert.Equal(t, int64(90000), azure.PollRateMs, "PollRateMs does not match")
}

func TestMarshalAzureIntegrationWithResourceFilterRules(t *testing.T) {
	payload, err := json.Marshal(AzureIntegration{
		ResourceFilterRules: []AzureFilterRule{{Filter: AzureFilterExpression{Source: "foobar"}}},
	})

	assert.NoError(t, err, "Unexpected error marshalling integration")
	expectedPayload := `{"enabled":false,"type":"","resourceFilterRules":[{"filter":{"source":"foobar"}}],"syncGuestOsNamespaces":false}`
	assert.Equal(t, expectedPayload, string(payload), "payload does not match")
}

func TestUnmarshalAzureIntegrationWithResourceFilterRules(t *testing.T) {
	azure := AzureIntegration{}
	err := json.Unmarshal([]byte(`{"resourceFilterRules":[{"filter": {"source": "foobar"}}]}`), &azure)

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	expectedFilterRules := []AzureFilterRule{{Filter: AzureFilterExpression{Source: "foobar"}}}
	assert.Equal(t, expectedFilterRules, azure.ResourceFilterRules, "ResourceFilterRules does not match")
}

func TestMarshalAzureIntegrationWithAdditionalServices(t *testing.T) {
	payload, err := json.Marshal(AzureIntegration{AdditionalServices: []string{"qwe", "abc"}})

	assert.NoError(t, err, "Unexpected error marshalling integration")
	expectedPayload := `{"enabled":false,"type":"","additionalServices":["qwe","abc"],"syncGuestOsNamespaces":false}`
	assert.Equal(t, expectedPayload, string(payload), "payload does not match")
}

func TestUnmarshalAzureIntegrationWithAdditionalServices(t *testing.T) {
	azure := AzureIntegration{}
	err := json.Unmarshal([]byte(`{"additionalServices":["qwe","abc"]}`), &azure)

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.Equal(t, []string{"qwe", "abc"}, azure.AdditionalServices, "AdditionalServices does not match")
}

func TestMarshalAzureIntegrationWithImportAzureMonitor(t *testing.T) {
	payload, err := json.Marshal(AzureIntegration{ImportAzureMonitor: true})

	assert.NoError(t, err, "Unexpected error marshalling integration")
	expectedPayload := `{"enabled":false,"type":"","syncGuestOsNamespaces":false,"importAzureMonitor":true}`
	assert.Equal(t, expectedPayload, string(payload), "payload does not match")
}

func TestUnmarshalAzureIntegrationWithImportAzureMonitor(t *testing.T) {
	azure := AzureIntegration{}
	err := json.Unmarshal([]byte(`{"importAzureMonitor":false}`), &azure)

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.Equal(t, false, azure.ImportAzureMonitor, "ImportAzureMonitor does not match")
}
