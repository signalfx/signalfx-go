package signalfx

import (
	"context"
	"net/http"
	"testing"

	automated_archival "github.com/signalfx/signalfx-go/automated-archival"
	"github.com/stretchr/testify/assert"
)

func TestCreateSettings(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(AutomatedArchivalApiURL+"/settings", verifyRequest(t, http.MethodPost, true, http.StatusOK, nil, "automated_archival/create_settings_success.json"))
	result, err := client.CreateSettings(context.Background(), &automated_archival.AutomatedArchivalSettings{
		Version:        1,
		Enabled:        true,
		LookbackPeriod: "P60D",
		GracePeriod:    "P30D",
	})

	assert.NoError(t, err, "Unexpected error creating Automated Archival settings")
	assert.Equal(t, int64(1), result.Version, "Version doesn't match")
	assert.Equal(t, true, result.Enabled, "Enabled doesn't match")
	assert.Equal(t, "P60D", result.LookbackPeriod, "Lookback period doesn't match")
	assert.Equal(t, "P30D", result.GracePeriod, "Grace period doesn't match")
}

func TestGetSettings(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(AutomatedArchivalApiURL+"/settings", verifyRequest(t, http.MethodGet, true, http.StatusOK, nil, "automated_archival/get_settings_success.json"))

	result, err := client.GetSettings(context.Background())
	assert.NoError(t, err, "Unexpected error getting Automated Archival settings")
	assert.Equal(t, int64(1), result.Version, "Version doesn't match")
	assert.Equal(t, true, result.Enabled, "Enabled doesn't match")
	assert.Equal(t, "P60D", result.LookbackPeriod, "Lookback period doesn't match")
	assert.Equal(t, "P30D", result.GracePeriod, "Grace period doesn't match")
	assert.Equal(t, int64(1674598662022), *result.Created, "'Created' timestamp doesn't match")
	assert.Equal(t, int64(1674598662022), *result.LastUpdated, "'Last Updated' timestamp doesn't match")
	assert.Equal(t, "TestCreatorId", *result.Creator, "Creator doesn't match")
	assert.Equal(t, "TestUpdatedId", *result.LastUpdatedBy, "'Last Updated By' user doesn't match")
}

func TestUpdateSettings(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(AutomatedArchivalApiURL+"/settings", verifyRequest(t, http.MethodPut, true, http.StatusOK, nil, "automated_archival/update_settings_success.json"))
	result, err := client.UpdateSettings(context.Background(), &automated_archival.AutomatedArchivalSettings{
		Version:        2,
		Enabled:        true,
		LookbackPeriod: "P45D",
		GracePeriod:    "P30D",
	})

	assert.NoError(t, err, "Unexpected error updating Automated Archival settings")
	assert.Equal(t, int64(2), result.Version, "Version doesn't match")
	assert.Equal(t, true, result.Enabled, "Enabled doesn't match")
	assert.Equal(t, "P45D", result.LookbackPeriod, "Lookback period doesn't match")
	assert.Equal(t, "P30D", result.GracePeriod, "Grace period doesn't match")
}

func TestDeleteSettings(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(AutomatedArchivalApiURL+"/settings", verifyRequest(t, http.MethodDelete, true, http.StatusNoContent, nil, ""))
	var version = int64(1)
	err := client.DeleteSettings(context.Background(), &automated_archival.AutomatedArchivalSettingsDeleteRequest{
		Version: &version,
	})
	assert.NoError(t, err, "Unexpected error deleting Automated Archival settings")
}

func TestCreateExemptMetrics(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(AutomatedArchivalApiURL+"/exempt-metrics", verifyRequest(t, http.MethodPost, true, http.StatusOK, nil, "automated_archival/create_exempt_metrics_success.json"))
	Creator := "AAAAAAAAAAA"
	Created := int64(1741968478768)
	Metric1ID := "GmAb7YHAIBs"
	Metric2ID := "GmAb7YHAIB0"
	Metric3ID := "GmAb7YHAIB8"
	Name1 := "bugbash.automated.archival1"
	Name2 := "bugbash.automated.archival.132"
	Name3 := "bugbash.automated.archival.10"
	result, err := client.CreateExemptMetrics(context.Background(), &[]automated_archival.ExemptMetric{
		{
			Creator:       &Creator,
			LastUpdatedBy: nil,
			Created:       &Created,
			LastUpdated:   &Created,
			Id:            &Metric1ID,
			Name:          Name1,
		},
		{
			Creator:       &Creator,
			LastUpdatedBy: nil,
			Created:       &Created,
			LastUpdated:   &Created,
			Id:            &Metric2ID,
			Name:          Name2,
		},
		{
			Creator:       &Creator,
			LastUpdatedBy: nil,
			Created:       &Created,
			LastUpdated:   &Created,
			Id:            &Metric3ID,
			Name:          Name3,
		},
	})

	assert.NoError(t, err, "Unexpected error creating Automated Archival exempt metrics")
	assert.Equal(t, 3, len(*result), "Unexpected exempt metrics array length")
	exemptMetric1 := (*result)[0]
	assert.Equal(t, Metric1ID, *exemptMetric1.Id, "Exempt metric id doesn't match")
	assert.Equal(t, Name1, exemptMetric1.Name, "Exempt metric name doesn't match")
	exemptMetric2 := (*result)[1]
	assert.Equal(t, Metric2ID, *exemptMetric2.Id, "Exempt metric id doesn't match")
	assert.Equal(t, Name2, exemptMetric2.Name, "Exempt metric name doesn't match")
	exemptMetric3 := (*result)[2]
	assert.Equal(t, Metric3ID, *exemptMetric3.Id, "Exempt metric id doesn't match")
	assert.Equal(t, Name3, exemptMetric3.Name, "Exempt metric name doesn't match")
}

func TestGetExemptMetrics(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(AutomatedArchivalApiURL+"/exempt-metrics", verifyRequest(t, http.MethodGet, true, http.StatusOK, nil, "automated_archival/get_exempt_metrics_success.json"))
	Metric1ID := "GmAb7YHAIBs"
	Metric2ID := "GmAb7YHAIB0"
	Metric3ID := "GmAb7YHAIB8"
	Name1 := "bugbash.automated.archival1"
	Name2 := "bugbash.automated.archival.132"
	Name3 := "bugbash.automated.archival.10"
	result, err := client.GetExemptMetrics(context.Background())

	assert.NoError(t, err, "Unexpected error creating Automated Archival exempt metrics")
	assert.Equal(t, 3, len(*result), "Unexpected exempt metrics array length")
	exemptMetric1 := (*result)[0]
	assert.Equal(t, Metric1ID, *exemptMetric1.Id, "Exempt metric id doesn't match")
	assert.Equal(t, Name1, exemptMetric1.Name, "Exempt metric name doesn't match")
	exemptMetric2 := (*result)[1]
	assert.Equal(t, Metric2ID, *exemptMetric2.Id, "Exempt metric id doesn't match")
	assert.Equal(t, Name2, exemptMetric2.Name, "Exempt metric name doesn't match")
	exemptMetric3 := (*result)[2]
	assert.Equal(t, Metric3ID, *exemptMetric3.Id, "Exempt metric id doesn't match")
	assert.Equal(t, Name3, exemptMetric3.Name, "Exempt metric name doesn't match")
}

func TestDeleteExemptMetrics(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(AutomatedArchivalApiURL+"/exempt-metrics", verifyRequest(t, http.MethodDelete, true, http.StatusNoContent, nil, ""))
	err := client.DeleteExemptMetrics(context.Background(), &automated_archival.ExemptMetricDeleteRequest{
		Ids: []string{"GmAb7YHAIBs", "GmAb7YHAIB0", "GmAb7YHAIB8"},
	})
	assert.NoError(t, err, "Unexpected error deleting Automated Archival exempt metrics")
}
