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
