package signalfx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/signalfx/signalfx-go/chart"
	"github.com/stretchr/testify/assert"
)

func TestCreateChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart", verifyRequest(t, "POST", true, http.StatusOK, nil, "chart/create_success.json"))

	result, err := client.CreateChart(context.Background(), &chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error creating chart")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestBadCreateChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	result, err := client.CreateChart(context.Background(), &chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.Error(t, err, "Expected error creating bad chart")
	assert.Nil(t, result, "Expected nil result creating bad chart")
}

func TestCreateSloChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/createSloChart", verifyRequest(t, "POST", true, http.StatusOK, nil, "chart/create_success.json"))

	result, err := client.CreateSloChart(context.Background(), &chart.CreateUpdateSloChartRequest{
		SloId: "slo-id",
	})
	assert.NoError(t, err, "Unexpected error creating chart")
	assert.Equal(t, "slo-id", result.SloId, "SLO ID does not match")
}

func TestBadCreateSloChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/createSloChart", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	result, err := client.CreateSloChart(context.Background(), &chart.CreateUpdateSloChartRequest{
		SloId: "slo-id",
	})
	assert.Error(t, err, "Expected error creating bad chart")
	assert.Nil(t, result, "Expected nil result creating bad chart")
}

func TestDeleteChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "DELETE", true, http.StatusOK, nil, ""))

	err := client.DeleteChart(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error deleting chart")
}

func TestDeleteMissingChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "DELETE", true, http.StatusNotFound, nil, ""))

	err := client.DeleteChart(context.Background(), "string")
	assert.Error(t, err, "Expected error deleting missing chart")
}

func TestGetChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "chart/get_success.json"))

	result, err := client.GetChart(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting chart")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetChart(context.Background(), "string")
	assert.Error(t, err, "Expected error getting missing chart")
	assert.Nil(t, result, "Expected nil result getting chart")
}

func TestSearchChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	name := "foo"
	offset := 2
	tags := "bar"
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	mux.HandleFunc("/v2/chart", verifyRequest(t, "GET", true, http.StatusOK, params, "chart/search_success.json"))

	results, err := client.SearchCharts(context.Background(), limit, name, offset, tags)
	assert.NoError(t, err, "Unexpected error search chart")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestSearchChartBad(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	name := "foo"
	offset := 2
	tags := "bar"
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	mux.HandleFunc("/v2/chart", verifyRequest(t, "GET", true, http.StatusBadRequest, params, ""))

	_, err := client.SearchCharts(context.Background(), limit, name, offset, tags)
	assert.Error(t, err, "Unexpected error search chart")
}

func TestUpdateChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "PUT", true, http.StatusOK, nil, "chart/update_success.json"))

	result, err := client.UpdateChart(context.Background(), "string", &chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error updating chart")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateMissingChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	result, err := client.UpdateChart(context.Background(), "string", &chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.Error(t, err, "Expected error updating chart")
	assert.Nil(t, result, "Expected nil result updating chart")
}

func TestUpdateSloChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/updateSloChart/slo-id", verifyRequest(t, "PUT", true, http.StatusOK, nil, "chart/update_success.json"))

	result, err := client.UpdateSloChart(context.Background(), "slo-id", &chart.CreateUpdateSloChartRequest{
		SloId: "slo-id",
	})
	assert.NoError(t, err, "Unexpected error updating chart")
	assert.Equal(t, "slo-id", result.SloId, "SLO ID does not match")
}

func TestUpdateMissingSloChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/updateSloChart/slo-id", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	result, err := client.UpdateSloChart(context.Background(), "slo-id", &chart.CreateUpdateSloChartRequest{
		SloId: "slo-id",
	})
	assert.Error(t, err, "Expected error updating chart")
	assert.Nil(t, result, "Expected nil result updating chart")
}

func TestValidateChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/validate", verifyRequest(t, "POST", true, http.StatusNoContent, nil, ""))

	err := client.ValidateChart(context.Background(), &chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error creating chart")
}

func TestBadValidateChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/validate", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	err := client.ValidateChart(context.Background(), &chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.Error(t, err, "Expected error creating bad chart")
}
