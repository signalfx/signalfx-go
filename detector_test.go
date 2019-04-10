package signalfx

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/signalfx/signalfx-go/detector"
	"github.com/stretchr/testify/assert"
)

func TestCreateDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector", verifyRequest(t, "POST", http.StatusOK, nil, "detector/create_success.json"))

	result, err := client.CreateDetector(&detector.CreateUpdateDetectorRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error creating detector")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeleteDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "DELETE", http.StatusOK, nil, ""))

	err := client.DeleteDetector("string")
	assert.NoError(t, err, "Unexpected error getting detector")
}

func TestDisableDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/disable", verifyRequest(t, "PUT", http.StatusOK, nil, ""))

	err := client.DisableDetector("string")
	assert.NoError(t, err, "Unexpected error enabling detector")
}

func TestEnableDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/enable", verifyRequest(t, "PUT", http.StatusOK, nil, ""))

	err := client.EnableDetector("string")
	assert.NoError(t, err, "Unexpected error disabling detector")
}

func TestGetDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "GET", http.StatusOK, nil, "detector/get_success.json"))

	result, err := client.GetDetector("string")
	assert.NoError(t, err, "Unexpected error getting detector")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestSearchDetector(t *testing.T) {
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

	mux.HandleFunc("/v2/detector", verifyRequest(t, "GET", http.StatusOK, params, "detector/search_success.json"))

	results, err := client.SearchDetectors(limit, name, offset, tags)
	assert.NoError(t, err, "Unexpected error search detector")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestUpdateDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "PUT", http.StatusOK, nil, "detector/update_success.json"))

	result, err := client.UpdateDetector("string", &detector.CreateUpdateDetectorRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error updating detector")
	assert.Equal(t, "string", result.Name, "Name does not match")
}
