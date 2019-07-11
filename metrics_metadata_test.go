package signalfx

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/signalfx/signalfx-go/metrics_metadata"
)

func TestGetDimension(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/dimension/string/string2", verifyRequest(t, "GET", http.StatusOK, nil, "metrics_metadata/get_dimension_success.json"))

	result, err := client.GetDimension("string", "string2")
	assert.NoError(t, err, "Unexpected error getting dimension")
	assert.Equal(t, result.Key, "string", "Key does not match")
}

func TestGetMissingDimension(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/dimension/string/string2", verifyRequest(t, "GET", http.StatusNotFound, nil, ""))

	result, err := client.GetDimension("string", "string2")
	assert.Error(t, err, "Should have gotten an error from a missing dimension")
	assert.Nil(t, result, "Should have gotten a nil result from a missing dimension")
}

func TestSearchDimension(t *testing.T) {
	teardown := setup()
	defer teardown()

	query := "foo:*"
	limit := 10
	offset := 2
	orderBy := "bar"
	params := url.Values{}
	params.Add("orderBy", orderBy)
	params.Add("query", query)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/dimension", verifyRequest(t, "GET", http.StatusOK, params, "metrics_metadata/dimension_search_success.json"))

	results, err := client.SearchDimension(query, orderBy, limit, offset)
	assert.NoError(t, err, "Unexpected error search dimensions")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestUpdateDimension(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/dimension/string/string", verifyRequest(t, "PUT", http.StatusOK, nil, "metrics_metadata/update_dimension_success.json"))

	result, err := client.UpdateDimension("string", "string", &metrics_metadata.Dimension{
		Key: "string",
	})
	assert.NoError(t, err, "Unexpected error updating dimension")
	assert.Equal(t, "string", result.Key, "Key does not match")
}

func TestUpdateMissingDimension(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/dimension/string/string2", verifyRequest(t, "PUT", http.StatusNotFound, nil, ""))

	result, err := client.UpdateDimension("string", "string2", &metrics_metadata.Dimension{
		Key: "string",
	})
	assert.Error(t, err, "Should have gotten an error from an update on a missing dimension")
	assert.Nil(t, result, "Should have gotten a nil result from an update on a missing dimension")
}

func TestGetMetric(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/metric/string", verifyRequest(t, "GET", http.StatusOK, nil, "metrics_metadata/get_metric_success.json"))

	result, err := client.GetMetric("string")
	assert.NoError(t, err, "Unexpected error getting metric")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingMetric(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/metric/string", verifyRequest(t, "GET", http.StatusNotFound, nil, ""))

	result, err := client.GetMetric("string")
	assert.Error(t, err, "Should have gotten an error from a missing metric")
	assert.Nil(t, result, "Should have gotten a nil result from a missing metric")
}

func TestSearchMetric(t *testing.T) {
	teardown := setup()
	defer teardown()

	query := "foo:*"
	limit := 10
	offset := 2
	orderBy := "bar"
	params := url.Values{}
	params.Add("orderBy", orderBy)
	params.Add("query", query)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/metric", verifyRequest(t, "GET", http.StatusOK, params, "metrics_metadata/metric_search_success.json"))

	results, err := client.SearchMetric(query, orderBy, limit, offset)
	assert.NoError(t, err, "Unexpected error search metrics")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestGetMetricTimeSeries(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/metrictimeseries/string", verifyRequest(t, "GET", http.StatusOK, nil, "metrics_metadata/get_metric_time_series_success.json"))

	result, err := client.GetMetricTimeSeries("string")
	assert.NoError(t, err, "Unexpected error getting metric time series")
	assert.Equal(t, result.Metric, "string", "Metric does not match")
}

func TestGetMissingMetricTimeSeries(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/metrictimeseries/string", verifyRequest(t, "GET", http.StatusNotFound, nil, ""))

	result, err := client.GetMetricTimeSeries("string")
	assert.Error(t, err, "Should have gotten an error from a missing metric time series")
	assert.Nil(t, result, "Should have gotten a nil result from a missing metric time series")
}

func TestSearchMetricTimeSeries(t *testing.T) {
	teardown := setup()
	defer teardown()

	query := "foo:*"
	limit := 10
	offset := 2
	orderBy := "bar"
	params := url.Values{}
	params.Add("orderBy", orderBy)
	params.Add("query", query)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/metrictimeseries", verifyRequest(t, "GET", http.StatusOK, params, "metrics_metadata/metric_time_series_search_success.json"))

	results, err := client.SearchMetricTimeSeries(query, orderBy, limit, offset)
	assert.NoError(t, err, "Unexpected error search metric time series")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestSearchTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	query := "foo:*"
	limit := 10
	offset := 2
	orderBy := "bar"
	params := url.Values{}
	params.Add("orderBy", orderBy)
	params.Add("query", query)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/tag", verifyRequest(t, "GET", http.StatusOK, params, "metrics_metadata/tag_search_success.json"))

	results, err := client.SearchTag(query, orderBy, limit, offset)
	assert.NoError(t, err, "Unexpected error search tags")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestGetTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag/string", verifyRequest(t, "GET", http.StatusOK, nil, "metrics_metadata/get_tag_success.json"))

	result, err := client.GetTag("string")
	assert.NoError(t, err, "Unexpected error getting Tag")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag/string", verifyRequest(t, "GET", http.StatusNotFound, nil, ""))

	result, err := client.GetTag("string")
	assert.Error(t, err, "Should have gotten an error from a missing tag")
	assert.Nil(t, result, "Should have gotten a nil result from a missing tag")
}

func TestDeleteTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag/string", verifyRequest(t, "DELETE", http.StatusNoContent, nil, ""))

	err := client.DeleteTag("string")
	assert.NoError(t, err, "Unexpected error deleting tag")
}

func TestDeleteMissingTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag", verifyRequest(t, "POST", http.StatusNotFound, nil, ""))

	err := client.DeleteTag("example")
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}

func TestUpdateCreateTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag/string", verifyRequest(t, "PUT", http.StatusOK, nil, "metrics_metadata/create_update_tag_success.json"))

	result, err := client.CreateUpdateTag("string", &metrics_metadata.CreateUpdateTagRequest{
		Description: "string",
	})
	assert.NoError(t, err, "Unexpected error updating tag")
	assert.Equal(t, "string", result.Name, "Key does not match")
}
