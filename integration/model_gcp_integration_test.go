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

func TestMarshalGCPIntegrationWithSAKeyConfig(t *testing.T) {
	payload, err := json.Marshal(GCPIntegration{
		ProjectServiceKeys: []*GCPProject{
			{
				ProjectId:  "prj-id-123",
				ProjectKey: "{\"some\":\"key\"}",
			},
		},
	})

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","projectServiceKeys":[{"projectId":"prj-id-123","projectKey":"{\"some\":\"key\"}"}]}`, string(payload), "payload does not match")
}

func TestMarshalGCPIntegrationWithWIFConfig(t *testing.T) {
	payload, err := json.Marshal(GCPIntegration{
		AuthMethod: WORKLOAD_IDENTITY_FEDERATION,
		WifConfigs: []*GCPProjectWIFConfig{
			{
				ProjectId: "prj-id-123",
				WIFConfig: "{\"some\":\"config\"}",
			},
		},
	})

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.Equal(t, `{"enabled":false,"type":"","authMethod":"WORKLOAD_IDENTITY_FEDERATION","workloadIdentityFederationConfigs":[{"projectId":"prj-id-123","wifConfig":"{\"some\":\"config\"}"}]}`, string(payload), "payload does not match")
}

func TestUnMarshalGCPIntegrationWithWIFConfig(t *testing.T) {
	GCP := GCPIntegration{}
	err := json.Unmarshal([]byte(`{"authMethod":"WORKLOAD_IDENTITY_FEDERATION","wifSplunkIdentity":{"account_id": "123", "aws_role_arn": "arn:aws:sts::123:assumed-role/splunk-o11y"},"workloadIdentityFederationConfigs":[{"projectId":"prj-id-123","wifConfig":"{\"some\":\"config\"}"}]}`), &GCP)

	expectedSplunkIdentity := map[string]string{
		"account_id":   "123",
		"aws_role_arn": "arn:aws:sts::123:assumed-role/splunk-o11y",
	}
	expectedConfigs := []*GCPProjectWIFConfig{
		{
			ProjectId: "prj-id-123",
			WIFConfig: "{\"some\":\"config\"}",
		},
	}
	expectedAuthMethod := WORKLOAD_IDENTITY_FEDERATION

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.EqualValues(t, expectedConfigs, GCP.WifConfigs, "WifConfigs do not match")
	assert.EqualValues(t, expectedSplunkIdentity, GCP.WifSplunkIdentity, "WifSplunkIdentity does not match")
	assert.EqualValues(t, expectedAuthMethod, GCP.AuthMethod, "AuthMethod does not match")
}

func TestUnMarshalGCPIntegrationWithSAKeysHidden(t *testing.T) {
	GCP := GCPIntegration{}
	err := json.Unmarshal([]byte(`{"projectServiceKeys":[{"projectId":"prj-id-123"}]}`), &GCP)

	expectedConfigs := []*GCPProject{
		{
			ProjectId: "prj-id-123",
		},
	}

	assert.NoError(t, err, "Unexpected error marshalling integration")
	assert.EqualValues(t, expectedConfigs, GCP.ProjectServiceKeys, "ProjectServiceKeys do not match")
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
