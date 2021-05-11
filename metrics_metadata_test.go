package signalfx

import (
	"context"
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

	mux.HandleFunc("/v2/dimension/string/string2", verifyRequest(t, "GET", true, http.StatusOK, nil, "metrics_metadata/get_dimension_success.json"))

	result, err := client.GetDimension(context.Background(), "string", "string2")
	assert.NoError(t, err, "Unexpected error getting dimension")
	assert.Equal(t, result.Key, "string", "Key does not match")
}

func TestGetMissingDimension(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/dimension/string/string2", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetDimension(context.Background(), "string", "string2")
	assert.Error(t, err, "Should have gotten an error from a missing dimension")
	assert.Nil(t, result, "Should have gotten a nil result from a missing dimension")
}

func TestSearchDimension(t *testing.T) {
	teardown := setup()
	defer teardown()

	query := "foo:*"
	limit := 10
	offset := 2
	params := url.Values{}
	params.Add("query", query)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/dimension", verifyRequest(t, "GET", true, http.StatusOK, params, "metrics_metadata/dimension_search_success.json"))

	results, err := client.SearchDimension(context.Background(), query, "", limit, offset)
	assert.NoError(t, err, "Unexpected error search dimensions")
	assert.Equal(t, int32(1), results.Count, "Incorrect results count")
	assert.Equal(t, 1, len(results.Results), "Incorrect number of results")
}

func TestSearchDimensionWithOrderBy(t *testing.T) {
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

	mux.HandleFunc("/v2/dimension", verifyRequest(t, "GET", true, http.StatusOK, params, "metrics_metadata/dimension_search_success.json"))

	results, err := client.SearchDimension(context.Background(), query, orderBy, limit, offset)
	assert.NoError(t, err, "Unexpected error search dimensions")
	assert.Equal(t, int32(1), results.Count, "Incorrect results count")
	assert.Equal(t, 1, len(results.Results), "Incorrect number of results")
}

func TestSearchDimensionBad(t *testing.T) {
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

	mux.HandleFunc("/v2/dimension", verifyRequest(t, "GET", true, http.StatusBadRequest, params, ""))

	_, err := client.SearchDimension(context.Background(), query, orderBy, limit, offset)
	assert.Error(t, err, "Didn't receive expected error for search dimensions")
}

func TestUpdateDimension(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/dimension/string/string", verifyRequest(t, "PUT", true, http.StatusOK, nil, "metrics_metadata/update_dimension_success.json"))

	result, err := client.UpdateDimension(context.Background(), "string", "string", &metrics_metadata.Dimension{
		Key: "string",
	})
	assert.NoError(t, err, "Unexpected error updating dimension")
	assert.Equal(t, "string", result.Key, "Key does not match")
}

func TestUpdateMissingDimension(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/dimension/string/string2", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	result, err := client.UpdateDimension(context.Background(), "string", "string2", &metrics_metadata.Dimension{
		Key: "string",
	})
	assert.Error(t, err, "Should have gotten an error from an update on a missing dimension")
	assert.Nil(t, result, "Should have gotten a nil result from an update on a missing dimension")
}

func TestGetMetric(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/metric/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "metrics_metadata/get_metric_success.json"))

	result, err := client.GetMetric(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting metric")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingMetric(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/metric/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetMetric(context.Background(), "string")
	assert.Error(t, err, "Should have gotten an error from a missing metric")
	assert.Nil(t, result, "Should have gotten a nil result from a missing metric")
}

func TestSearchMetric(t *testing.T) {
	teardown := setup()
	defer teardown()

	query := "foo:*"
	limit := 10
	offset := 2
	params := url.Values{}
	params.Add("query", query)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/metric", verifyRequest(t, "GET", true, http.StatusOK, params, "metrics_metadata/metric_search_success.json"))

	results, err := client.SearchMetric(context.Background(), query, "", limit, offset)
	assert.NoError(t, err, "Unexpected error search metrics")
	assert.Equal(t, int32(1), results.Count, "Incorrect results count")
	assert.Equal(t, 1, len(results.Results), "Incorrect number of results")
}

func TestSearchMetricWithOrderBy(t *testing.T) {
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

	mux.HandleFunc("/v2/metric", verifyRequest(t, "GET", true, http.StatusOK, params, "metrics_metadata/metric_search_success.json"))

	results, err := client.SearchMetric(context.Background(), query, orderBy, limit, offset)
	assert.NoError(t, err, "Unexpected error search metrics")
	assert.Equal(t, int32(1), results.Count, "Incorrect results count")
	assert.Equal(t, 1, len(results.Results), "Incorrect number of results")
}

func TestSearchMetricBad(t *testing.T) {
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

	mux.HandleFunc("/v2/metric", verifyRequest(t, "GET", true, http.StatusBadRequest, params, ""))

	_, err := client.SearchMetric(context.Background(), query, orderBy, limit, offset)
	assert.Error(t, err, "Unexpected error search metrics")
}

func TestGetMetricTimeSeries(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/metrictimeseries/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "metrics_metadata/get_metric_time_series_success.json"))

	result, err := client.GetMetricTimeSeries(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting metric time series")
	assert.Equal(t, result.Metric, "string", "Metric does not match")
}

func TestGetMissingMetricTimeSeries(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/metrictimeseries/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetMetricTimeSeries(context.Background(), "string")
	assert.Error(t, err, "Should have gotten an error from a missing metric time series")
	assert.Nil(t, result, "Should have gotten a nil result from a missing metric time series")
}

func TestSearchMetricTimeSeries(t *testing.T) {
	teardown := setup()
	defer teardown()

	query := "foo:*"
	limit := 10
	offset := 2
	params := url.Values{}
	params.Add("query", query)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/metrictimeseries", verifyRequest(t, "GET", true, http.StatusOK, params, "metrics_metadata/metric_time_series_search_success.json"))

	results, err := client.SearchMetricTimeSeries(context.Background(), query, "", limit, offset)
	assert.NoError(t, err, "Unexpected error search metric time series")
	assert.Equal(t, int32(1), results.Count, "Incorrect results count")
	assert.Equal(t, 1, len(results.Results), "Incorrect number of results")
}

func TestSearchMetricTimeSeriesWithOrderBy(t *testing.T) {
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

	mux.HandleFunc("/v2/metrictimeseries", verifyRequest(t, "GET", true, http.StatusOK, params, "metrics_metadata/metric_time_series_search_success.json"))

	results, err := client.SearchMetricTimeSeries(context.Background(), query, orderBy, limit, offset)
	assert.NoError(t, err, "Unexpected error search metric time series")
	assert.Equal(t, int32(1), results.Count, "Incorrect results count")
	assert.Equal(t, 1, len(results.Results), "Incorrect number of results")
}

func TestSearchMetricTimeSeriesBad(t *testing.T) {
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

	mux.HandleFunc("/v2/metrictimeseries", verifyRequest(t, "GET", true, http.StatusBadRequest, params, ""))

	_, err := client.SearchMetricTimeSeries(context.Background(), query, orderBy, limit, offset)
	assert.Error(t, err, "Unexpected error search metric time series")
}

func TestSearchTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	query := "foo:*"
	limit := 10
	offset := 2
	params := url.Values{}
	params.Add("query", query)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/tag", verifyRequest(t, "GET", true, http.StatusOK, params, "metrics_metadata/tag_search_success.json"))

	results, err := client.SearchTag(context.Background(), query, "", limit, offset)
	assert.NoError(t, err, "Unexpected error search tags")
	assert.Equal(t, int32(1), results.Count, "Incorrect results count")
	assert.Equal(t, 1, len(results.Results), "Incorrect number of results")
}

func TestSearchTagWithOrderBy(t *testing.T) {
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

	mux.HandleFunc("/v2/tag", verifyRequest(t, "GET", true, http.StatusOK, params, "metrics_metadata/tag_search_success.json"))

	results, err := client.SearchTag(context.Background(), query, orderBy, limit, offset)
	assert.NoError(t, err, "Unexpected error search tags")
	assert.Equal(t, int32(1), results.Count, "Incorrect results count")
	assert.Equal(t, 1, len(results.Results), "Incorrect number of results")
}

func TestSearchTagBad(t *testing.T) {
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

	mux.HandleFunc("/v2/tag", verifyRequest(t, "GET", true, http.StatusBadRequest, params, ""))

	_, err := client.SearchTag(context.Background(), query, orderBy, limit, offset)
	assert.Error(t, err, "Unexpected error search tags")
}

func TestGetTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "metrics_metadata/get_tag_success.json"))

	result, err := client.GetTag(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting Tag")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetTag(context.Background(), "string")
	assert.Error(t, err, "Should have gotten an error from a missing tag")
	assert.Nil(t, result, "Should have gotten a nil result from a missing tag")
}

func TestDeleteTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag/string", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteTag(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error deleting tag")
}

func TestDeleteMissingTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag", verifyRequest(t, "POST", true, http.StatusNotFound, nil, ""))

	err := client.DeleteTag(context.Background(), "example")
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}

func TestUpdateCreateTag(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/tag/string", verifyRequest(t, "PUT", true, http.StatusOK, nil, "metrics_metadata/create_update_tag_success.json"))

	result, err := client.CreateUpdateTag(context.Background(), "string", &metrics_metadata.CreateUpdateTagRequest{
		Description: "string",
	})
	assert.NoError(t, err, "Unexpected error updating tag")
	assert.Equal(t, "string", result.Name, "Key does not match")
}

func TestUpdateCreateMetric(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/metric/string", verifyRequest(t, "PUT", true, http.StatusOK, nil, "metrics_metadata/create_update_metric_success.json"))

	result, err := client.CreateUpdateMetric(context.Background(), "string", &metrics_metadata.CreateUpdateMetricRequest{
		Description: "string",
	})
	assert.NoError(t, err, "Unexpected error updating tag")
	assert.Equal(t, "string", result.Name, "Key does not match")
}
