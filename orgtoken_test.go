package signalfx

import (
	"context"
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

	mux.HandleFunc("/v2/token", verifyRequest(t, "POST", true, http.StatusOK, nil, "orgtoken/create_success.json"))

	quota := int32(1234)
	quotaThreshold := int32(1235)
	result, err := client.CreateOrgToken(context.Background(), &orgtoken.CreateUpdateTokenRequest{
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

	mux.HandleFunc("/v2/token", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	result, err := client.CreateOrgToken(context.Background(), &orgtoken.CreateUpdateTokenRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should have gotten an error from a bad create")
	assert.Nil(t, result, "Should have a null token on bad create")
}

func TestDeleteOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string%2Ffart", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteOrgToken(context.Background(), "string/fart")
	assert.NoError(t, err, "Unexpected error deleting token")
}

func TestDeleteMissingOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token", verifyRequest(t, "POST", true, http.StatusNotFound, nil, ""))

	err := client.DeleteOrgToken(context.Background(), "example")
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}

func TestGetOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string%2Ffart", verifyRequest(t, "GET", true, http.StatusOK, nil, "orgtoken/get_success.json"))

	result, err := client.GetOrgToken(context.Background(), "string/fart")
	assert.NoError(t, err, "Unexpected error getting token")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string%2Ffart", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetOrgToken(context.Background(), "string/fart")
	assert.Error(t, err, "Should have gotten an error from a missing token")
	assert.Nil(t, result, "Should have gotten a nil result from a missing token")
}

func TestSearchOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	name := "foo/fart"
	offset := 2
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", url.PathEscape(name))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/token", verifyRequest(t, "GET", true, http.StatusOK, params, "orgtoken/search_success.json"))

	results, err := client.SearchOrgTokens(context.Background(), limit, name, offset)
	assert.NoError(t, err, "Unexpected error search token")
	assert.Equal(t, int32(2), results.Count, "Incorrect number of results")
}

func TestSearchOrgTokenBad(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	name := "foo/fart"
	offset := 2
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", url.PathEscape(name))
	params.Add("offset", strconv.Itoa(offset))

	mux.HandleFunc("/v2/token", verifyRequest(t, "GET", true, http.StatusBadRequest, params, ""))

	_, err := client.SearchOrgTokens(context.Background(), limit, name, offset)
	assert.Error(t, err, "Unexpected error search token")
}

func TestUpdateOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string%2Ffart", verifyRequest(t, "PUT", true, http.StatusOK, nil, "orgtoken/update_success.json"))

	result, err := client.UpdateOrgToken(context.Background(), "string/fart", &orgtoken.CreateUpdateTokenRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error updating token")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateMissingOrgToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/token/string", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	result, err := client.UpdateOrgToken(context.Background(), "string", &orgtoken.CreateUpdateTokenRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should have gotten an error from an update on a missing token")
	assert.Nil(t, result, "Should have gotten a nil result from an update on a missing token")
}
