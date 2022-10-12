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

func TestMarshalGCPIntegrationWithWhitelist(t *testing.T) {
	whitelist := []string{"key"}
	gcpInt := GCPIntegration{
		Whitelist: whitelist,
	}
	payload, err := json.Marshal(&gcpInt)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","includeList":["key"]}`, string(payload), "payload does not match")
	assert.Nil(t, gcpInt.IncludeList, "IncludeList has been changed")
}

func TestMarshalGCPIntegrationWithIncludeList(t *testing.T) {
	includeList := []string{"key"}
	payload, err := json.Marshal(GCPIntegration{
		IncludeList: includeList,
	})

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","includeList":["key"]}`, string(payload), "payload does not match")
}

func TestUnmarshalGCPIntegrationWithIncludeList(t *testing.T) {
	expectedValue := []string{"key"}

	GCP := GCPIntegration{}
	err := json.Unmarshal([]byte(`{"includeList":["key"]}`), &GCP)

	assert.NoError(t, err, "Unexpected error unmarshalling integration")
	assert.Equal(t, expectedValue, GCP.IncludeList, "IncludeList does not match")
	assert.Equal(t, expectedValue, GCP.Whitelist, "Whitelist does not match")
}

func TestMarshalGCPIntegrationWithUseMetricSourceProjectForQuotaEnabled(t *testing.T) {
	gcpInt := GCPIntegration{
		UseMetricSourceProjectForQuota: true,
	}
	payload, err := json.Marshal(&gcpInt)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","useMetricSourceProjectForQuota":true}`, string(payload), "payload does not match")
}

func TestMarshalGCPIntegrationWithImportGCPMetricsEnabled(t *testing.T) {
	gcpInt := GCPIntegration{
		ImportGCPMetrics: newBoolPtr(true),
	}

	payload, err := json.Marshal(&gcpInt)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","importGCPMetrics":true}`, string(payload), "payload does not match")

}

func TestMarshalGCPIntegrationWithImportGCPMetricsDisabled(t *testing.T) {
	gcpInt := GCPIntegration{
		ImportGCPMetrics: newBoolPtr(false),
	}

	payload, err := json.Marshal(&gcpInt)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","importGCPMetrics":false}`, string(payload), "payload does not match")

}

func TestMarshalGCPIntegrationWithImportGCPMetricsEmpty(t *testing.T) {
	GCP := GCPIntegration{}

	payload, err := json.Marshal(GCP)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":""}`, string(payload), "payload does not match")

}

func TestUnmarshalGCPIntegrationWithImportGCPMetricsEnabled(t *testing.T) {
	GCP := GCPIntegration{}

	err := json.Unmarshal([]byte(`{"importGCPMetrics":true}`), &GCP)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, true, *GCP.ImportGCPMetrics, "ImportGCPMetrics does not match")

}

func TestUnmarshalGCPIntegrationWithImportGCPMetricsDisabled(t *testing.T) {
	GCP := GCPIntegration{}

	err := json.Unmarshal([]byte(`{"importGCPMetrics":false}`), &GCP)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, false, *GCP.ImportGCPMetrics, "ImportGCPMetrics does not match")

}

func TestUnmarshalGCPIntegrationWithImportGCPMetricsEmpty(t *testing.T) {
	GCP := GCPIntegration{}

	err := json.Unmarshal([]byte(`{}`), &GCP)

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, (*bool)(nil), GCP.ImportGCPMetrics, "ImportGCPMetrics does not match")

}
