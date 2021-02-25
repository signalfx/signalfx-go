package signalfx

import (
	"context"
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
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
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
