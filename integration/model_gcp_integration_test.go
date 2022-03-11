package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalGCPIntegrationWithPollRate(t *testing.T) {
	pollRate := OneMinutely
	gcpInt := GCPIntegration{
		PollRate: &pollRate,
	}
	payload, err := json.Marshal(&gcpInt)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","pollRate":60000}`, string(payload), "payload does not match")
	assert.Equal(t, int64(0), gcpInt.PollRateMs, "PollRateMs has been changed")
}

func TestMarshalGCPIntegrationWithPollRateMs(t *testing.T) {
	payload, err := json.Marshal(GCPIntegration{
		PollRateMs: 90000,
	})

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","pollRate":90000}`, string(payload), "payload does not match")
}

func TestUnmarshalGCPIntegrationWithPollRate(t *testing.T) {
	GCP := GCPIntegration{}
	err := json.Unmarshal([]byte(`{"pollRate":60000}`), &GCP)

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.Equal(t, OneMinutely, *GCP.PollRate, "PollRate does not match")
	assert.Equal(t, int64(60000), GCP.PollRateMs, "PollRateMs does not match")
}

func TestUnmarshalGCPIntegrationWithPollRateMs(t *testing.T) {
	GCP := GCPIntegration{}
	err := json.Unmarshal([]byte(`{"pollRate":90000}`), &GCP)

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.Nil(t, GCP.PollRate, "PollRate does not match")
	assert.Equal(t, int64(90000), GCP.PollRateMs, "PollRateMs does not match")
}
