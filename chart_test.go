package signalfx

import (
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

	mux.HandleFunc("/v2/chart", verifyRequest(t, "POST", http.StatusOK, nil, "chart/create_success.json"))

	result, err := client.CreateChart(&chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error creating chart")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestBadCreateChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart", verifyRequest(t, "POST", http.StatusBadRequest, nil, ""))

	result, err := client.CreateChart(&chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.Error(t, err, "Expected error creating bad chart")
	assert.Nil(t, result, "Exepcted nil result creating bad chart")
}

func TestDeleteChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "DELETE", http.StatusOK, nil, ""))

	err := client.DeleteChart("string")
	assert.NoError(t, err, "Unexpected error deleting chart")
}

func TestDeleteMissingChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "DELETE", http.StatusNotFound, nil, ""))

	err := client.DeleteChart("string")
	assert.Error(t, err, "Expected error deleting missing chart")
}

func TestGetChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "GET", http.StatusOK, nil, "chart/get_success.json"))

	result, err := client.GetChart("string")
	assert.NoError(t, err, "Unexpected error getting chart")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "GET", http.StatusNotFound, nil, ""))

	result, err := client.GetChart("string")
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

	mux.HandleFunc("/v2/chart", verifyRequest(t, "GET", http.StatusOK, params, "chart/search_success.json"))

	results, err := client.SearchCharts(limit, name, offset, tags)
	assert.NoError(t, err, "Unexpected error search chart")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestUpdateChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "PUT", http.StatusOK, nil, "chart/update_success.json"))

	result, err := client.UpdateChart("string", &chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error updating chart")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateMissingChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "PUT", http.StatusNotFound, nil, ""))

	result, err := client.UpdateChart("string", &chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.Error(t, err, "Expected error updating chart")
	assert.Nil(t, result, "Expected nil result updating chart")
}
