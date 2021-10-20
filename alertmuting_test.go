package signalfx

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/signalfx/signalfx-go/alertmuting"
	"github.com/stretchr/testify/assert"
)

func TestCreateAlertMutingRule(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alertmuting", verifyRequest(t, "POST", true, http.StatusCreated, nil, "alertmuting/create_success.json"))

	result, err := client.CreateAlertMutingRule(context.Background(), &alertmuting.CreateUpdateAlertMutingRuleRequest{
		Description: "string",
	})
	assert.NoError(t, err, "Unexpected error creating alert muting rule")
	assert.Equal(t, "string", result.Description, "Description does not match")
	assert.Equal(t, "server5", result.Filters[0].PropertyValue.Values[0], "Property Value does not match")
}

func TestCreateBadAlertMutingRule(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alertmuting", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	result, err := client.CreateAlertMutingRule(context.Background(), &alertmuting.CreateUpdateAlertMutingRuleRequest{
		Description: "string",
	})
	assert.Error(t, err, "Should have gotten an error from a bad create")
	assert.Nil(t, result, "Should have a null alert muting rule on bad create")
}

func TestDeleteAlertMutingRule(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alertmuting/string", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteAlertMutingRule(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error deleting alert muting rule")
}

func TestDeleteMissingAlertMutingRule(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alertmuting", verifyRequest(t, "POST", true, http.StatusNotFound, nil, ""))

	err := client.DeleteAlertMutingRule(context.Background(), "example")
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}

func TestGetAlertMutingRule(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alertmuting/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "alertmuting/get_success.json"))

	result, err := client.GetAlertMutingRule(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting alert mutnig rule")
	assert.Equal(t, result.Description, "string", "Name does not match")
	assert.Equal(t, "server6", result.Filters[0].PropertyValue.Values[1], "Property Value does not match")
}

func TestGetMissingAlertMutingRule(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alertmuting/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetAlertMutingRule(context.Background(), "string")
	assert.Error(t, err, "Should have gotten an error from a missing alert muting rule")
	assert.Nil(t, result, "Should have gotten a nil result from a missing alert muting rule")
}

func TestSearchAlertMutingRule(t *testing.T) {
	teardown := setup()
	defer teardown()

	include := "all"
	limit := 10
	query := "creator:AAXYAAAAAZ3"
	offset := 2
	params := url.Values{}
	params.Add("include", include)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("query", query)
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/alertmuting", verifyRequest(t, "GET", true, http.StatusOK, params, "alertmuting/search_success.json"))

	results, err := client.SearchAlertMutingRules(context.Background(), include, limit, query, offset)
	assert.NoError(t, err, "Unexpected error search alert muting rule")
	assert.Equal(t, int32(2), results.Count, "Incorrect number of results")
	assert.Equal(t, 2, len(results.Results), "Incorrect number of results in results")
	assert.Equal(t, "server6", results.Results[0].Filters[0].PropertyValue.Values[1], "Property Value does not match")
	assert.Equal(t, "server5", results.Results[1].Filters[0].PropertyValue.Values[0], "Property Value does not match")
}

func TestUpdateAlertMutingRule(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alertmuting/string", verifyRequest(t, "PUT", true, http.StatusOK, nil, "alertmuting/update_success.json"))

	result, err := client.UpdateAlertMutingRule(context.Background(), "string", &alertmuting.CreateUpdateAlertMutingRuleRequest{
		Description: "string",
	})
	assert.NoError(t, err, "Unexpected error updating alert muting rule")
	assert.Equal(t, "string", result.Description, "Description does not match")
}

func TestUpdateMissingAlertMutingRule(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alertmuting/string", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	result, err := client.UpdateAlertMutingRule(context.Background(), "string", &alertmuting.CreateUpdateAlertMutingRuleRequest{
		Description: "string",
	})
	assert.Error(t, err, "Should have gotten an error from an update on a missing alert muting rule")
	assert.Nil(t, result, "Should have gotten a nil result from an update on a missing alert muting rule")
}

func TestAlertMutingCustomJSONMarshaling(t *testing.T) {

	file, _ := ioutil.ReadFile("testdata/fixtures/alertmuting/update_success.json")
	var muting alertmuting.AlertMutingRule
	err := json.Unmarshal(file, &muting)
	assert.Nil(t, err, "Unexpected error unmarshaling muting rules")
	assert.Equal(t, "server5", muting.Filters[0].PropertyValue.Values[0], "Wrong propertyValue when unmarshalling")

	// marshall it back
	data, err := json.Marshal(muting)
	assert.Nil(t, err, "Unexpected error marshaling muting rules")
	// unmarshal again to check it
	err = json.Unmarshal(data, &muting)
	assert.Nil(t, err, "Unexpected error re-unmarshaling muting rule")
	assert.Equal(t, "server5", muting.Filters[0].PropertyValue.Values[0], "Wrong propertyValue when unmarshalling the second time")
}
