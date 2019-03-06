package signalfx

import (
	"net/http"
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

func TestDeleteChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "DELETE", http.StatusOK, nil, ""))

	err := client.DeleteChart("string")
	assert.NoError(t, err, "Unexpected error getting chart")
}

func TestGetChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "GET", http.StatusOK, nil, "chart/get_success.json"))

	result, err := client.GetChart("string")
	assert.NoError(t, err, "Unexpected error getting chart")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

// func TestSearchChart(t *testing.T) {
// 	teardown := setup()
// 	defer teardown()
//
// 	limit := 10
// 	name := "foo"
// 	offset := 2
// 	tags := "bar"
// 	params := url.Values{}
// 	params.Add("limit", strconv.Itoa(limit))
// 	params.Add("name", name)
// 	params.Add("offset", strconv.Itoa(offset))
// 	params.Add("tags", tags)
//
// 	mux.HandleFunc("/v2/chart", verifyRequest(t, "GET", http.StatusOK, params, "chart/get_success.json"))
//
// 	results, err := client.SearchChart(limit, name, offset, tags)
// 	assert.NoError(t, err, "Unexpected error search chart")
// 	assert.Equal(t, int64(0), results.Count, "Incorrect number of results")
// }
//
// func TestUpdateChart(t *testing.T) {
// 	teardown := setup()
// 	defer teardown()
//
// 	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "PUT", http.StatusOK, nil, "chart/update_success.json"))
//
// 	result, err := client.UpdateChart("string", &Chart{
// 		Name: "string",
// 	})
// 	assert.NoError(t, err, "Unexpected error updating chart")
// 	assert.Equal(t, "string", result.Name, "Name does not match")
// }
