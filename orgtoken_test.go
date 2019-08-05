package signalfx

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/signalfx/signalfx-go/orgtoken"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token", verifyRequest(t, "POST", http.StatusOK, nil, "orgtoken/create_success.json"))

	quota := int32(1234)
	quotaThreshold := int32(1235)
	result, err := client.CreateOrgToken(&orgtoken.CreateUpdateTokenRequest{
		Name: "string",
		Limits: &orgtoken.Limit{
			DpmQuota:                 &quota,
			DpmNotificationThreshold: &quotaThreshold,
		},
	})
	assert.NoError(t, err, "Unexpected error creating orgtoken")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestCreateBadOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token", verifyRequest(t, "POST", http.StatusBadRequest, nil, ""))

	result, err := client.CreateOrgToken(&orgtoken.CreateUpdateTokenRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should have gotten an error from a bad create")
	assert.Nil(t, result, "Should have a null token on bad create")
}

func TestDeleteOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string", verifyRequest(t, "DELETE", http.StatusNoContent, nil, ""))

	err := client.DeleteOrgToken("string")
	assert.NoError(t, err, "Unexpected error deleting token")
}

func TestDeleteMissingOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token", verifyRequest(t, "POST", http.StatusNotFound, nil, ""))

	err := client.DeleteOrgToken("example")
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}

func TestGetOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string", verifyRequest(t, "GET", http.StatusOK, nil, "orgtoken/get_success.json"))

	result, err := client.GetOrgToken("string")
	assert.NoError(t, err, "Unexpected error getting token")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string", verifyRequest(t, "GET", http.StatusNotFound, nil, ""))

	result, err := client.GetOrgToken("string")
	assert.Error(t, err, "Should have gotten an error from a missing token")
	assert.Nil(t, result, "Should have gotten a nil result from a missing token")
}

func TestSearchOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	name := "foo"
	offset := 2
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/token", verifyRequest(t, "GET", http.StatusOK, params, "orgtoken/search_success.json"))

	results, err := client.SearchOrgTokens(limit, name, offset)
	assert.NoError(t, err, "Unexpected error search token")
	assert.Equal(t, int32(2), results.Count, "Incorrect number of results")
}

func TestUpdateOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string", verifyRequest(t, "PUT", http.StatusOK, nil, "orgtoken/update_success.json"))

	result, err := client.UpdateOrgToken("string", &orgtoken.CreateUpdateTokenRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error updating token")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateMissingOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string", verifyRequest(t, "PUT", http.StatusNotFound, nil, ""))

	result, err := client.UpdateOrgToken("string", &orgtoken.CreateUpdateTokenRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should have gotten an error from an update on a missing token")
	assert.Nil(t, result, "Should have gotten a nil result from an update on a missing token")
}
