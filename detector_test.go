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

func TestCreateBadDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector", verifyRequest(t, "POST", http.StatusBadRequest, nil, ""))

	result, err := client.CreateDetector(&detector.CreateUpdateDetectorRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should have gotten an error from a bad create")
	assert.Nil(t, result, "Should have a null detector on bad create")
}

func TestDeleteDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "DELETE", http.StatusNoContent, nil, ""))

	err := client.DeleteDetector("string")
	assert.NoError(t, err, "Unexpected error deleting detector")
}

func TestDeleteMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector", verifyRequest(t, "POST", http.StatusNotFound, nil, ""))

	err := client.DeleteDetector("example")
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}

func TestDisableDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/disable", verifyRequest(t, "PUT", http.StatusNoContent, nil, ""))

	err := client.DisableDetector("string", []string{"example"})
	assert.NoError(t, err, "Unexpected error disabling detector")
}

func TestDisableMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/disable", verifyRequest(t, "PUT", http.StatusNotFound, nil, ""))

	err := client.DisableDetector("string", []string{"example"})
	assert.Error(t, err, "Should have gotten an error from a missing disable")
}

func TestEnableDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/enable", verifyRequest(t, "PUT", http.StatusNoContent, nil, ""))

	err := client.EnableDetector("string", []string{"example"})
	assert.NoError(t, err, "Unexpected error disabling detector")
}

func TestEnableMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/enable", verifyRequest(t, "PUT", http.StatusNotFound, nil, ""))

	err := client.EnableDetector("string", []string{"example"})
	assert.Error(t, err, "Should have gotten an error from a missing enable")
}

func TestGetDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "GET", http.StatusOK, nil, "detector/get_success.json"))

	result, err := client.GetDetector("string")
	assert.NoError(t, err, "Unexpected error getting detector")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "GET", http.StatusNotFound, nil, ""))

	result, err := client.GetDetector("string")
	assert.Error(t, err, "Should have gotten an error from a missing detector")
	assert.Nil(t, result, "Should have gotten a nil result from a missing detector")
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

func TestUpdateMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "PUT", http.StatusNotFound, nil, ""))

	result, err := client.UpdateDetector("string", &detector.CreateUpdateDetectorRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should have gotten an error from an update on a missing detector")
	assert.Nil(t, result, "Should have gotten a nil result from an update on a missing detector")
}
