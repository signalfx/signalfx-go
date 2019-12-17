package signalfx

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/signalfx/signalfx-go/datalink"
	"github.com/stretchr/testify/assert"
)

func TestCreateDataLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/crosslink", verifyRequest(t, "POST", true, http.StatusOK, nil, "datalink/create_success.json"))

	result, err := client.CreateDataLink(&datalink.CreateUpdateDataLinkRequest{
		PropertyName:  "string",
		PropertyValue: "string",
	})
	assert.NoError(t, err, "Unexpected error creating data link")
	assert.Equal(t, "string", result.PropertyName, "PropertyName does not match")
}

func TestBadCreateDataLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/crosslink", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	result, err := client.CreateDataLink(&datalink.CreateUpdateDataLinkRequest{
		PropertyName: "string",
	})
	assert.Error(t, err, "Should have gotten an error from a bad create")
	assert.Nil(t, result, "Should have a null data link on bad create")
}

func TestDeleteDataLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/crosslink/string", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteDataLink("string")
	assert.NoError(t, err, "Unexpected error deleting data link")
}

func TestDeleteMissingDataLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/crosslink", verifyRequest(t, "POST", true, http.StatusNotFound, nil, ""))

	err := client.DeleteDataLink("example")
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}

func TestGetDataLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/crosslink/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "datalink/get_success.json"))

	result, err := client.GetDataLink("string")
	assert.NoError(t, err, "Unexpected error getting data link")
	assert.Equal(t, result.PropertyName, "string", "Name does not match")
}

func TestGetMissingDataLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/crosslink/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetDataLink("string")
	assert.Error(t, err, "Should have gotten an error from a missing data link")
	assert.Nil(t, result, "Should have gotten a nil result from a missing data link")
}

func TestSearchDataLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	context := "foo"
	offset := 2
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("context", context)
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/crosslink", verifyRequest(t, "GET", true, http.StatusOK, params, "datalink/search_success.json"))

	results, err := client.SearchDataLinks(limit, context, offset)
	assert.NoError(t, err, "Unexpected error search data link")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestUpdateDataLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/crosslink/string", verifyRequest(t, "PUT", true, http.StatusOK, nil, "datalink/update_success.json"))

	result, err := client.UpdateDataLink("string", &datalink.CreateUpdateDataLinkRequest{
		PropertyName: "string",
	})
	assert.NoError(t, err, "Unexpected error updating data link")
	assert.Equal(t, "string", result.PropertyName, "PropertyName does not match")
}

func TestUpdateMissingDataLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/crosslink/string", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	result, err := client.UpdateDataLink("string", &datalink.CreateUpdateDataLinkRequest{
		PropertyName: "string",
	})
	assert.Error(t, err, "Should have gotten an error from an update on a missing data link")
	assert.Nil(t, result, "Should have gotten a nil result from an update on a missing data link")
}
