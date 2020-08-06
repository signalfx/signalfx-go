package signalfx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/signalfx/signalfx-go/team"
	"github.com/stretchr/testify/assert"
)

func TestCreateTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team", verifyRequest(t, "POST", true, http.StatusOK, nil, "team/create_success.json"))

	result, err := client.CreateTeam(context.Background(), &team.CreateUpdateTeamRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error creating team")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestCreateBadTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	result, err := client.CreateTeam(context.Background(), &team.CreateUpdateTeamRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should get an error from bad tream")
	assert.Nil(t, result, "Should get nil result from bad team")
}

func TestDeleteTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteTeam(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting Team")
}

func TestDeleteMissingTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string", verifyRequest(t, "DELETE", true, http.StatusNotFound, nil, ""))

	err := client.DeleteTeam(context.Background(), "string")
	assert.Error(t, err, "Should have gotten an error code for a missing team")
}

func TestGetTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "team/get_success.json"))

	result, err := client.GetTeam(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting Team")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetTeam(context.Background(), "string")
	assert.Error(t, err, "Expected error getting missing team")
	assert.Nil(t, result, "Expected nil result getting missing team")
}

func TestSearchTeam(t *testing.T) {
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

	mux.HandleFunc("/v2/team", verifyRequest(t, "GET", true, http.StatusOK, params, "team/search_success.json"))

	results, err := client.SearchTeam(context.Background(), limit, name, offset, tags)
	assert.NoError(t, err, "Unexpected error search Team")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestSearchTeamBad(t *testing.T) {
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

	mux.HandleFunc("/v2/team", verifyRequest(t, "GET", true, http.StatusOK, params, ""))

	_, err := client.SearchTeam(context.Background(), limit, name, offset, tags)
	assert.Error(t, err, "Unexpected error search Team")
}

func TestUpdateTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string", verifyRequest(t, "PUT", true, http.StatusOK, nil, "team/update_success.json"))

	result, err := client.UpdateTeam(context.Background(), "string", &team.CreateUpdateTeamRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error updating Team")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateMissingTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	result, err := client.UpdateTeam(context.Background(), "string", &team.CreateUpdateTeamRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should've gotten an error from a missing team update")
	assert.Nil(t, result, "Should've gotten a nil dashboard from a missing team update")
}

func TestLinkDetectorToTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string/detector/string2", verifyRequest(t, "POST", true, http.StatusNoContent, nil, ""))

	err := client.LinkDetectorToTeam(context.Background(), "string", "string2")
	assert.NoError(t, err, "Unexpected error linking team")
}

func TestUnlinkDetectorFromTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string/detector/string2", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.UnlinkDetectorFromTeam(context.Background(), "string", "string2")
	assert.NoError(t, err, "Unexpected error unlinking team")
}

func TestLinkDashboardGroupToTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string/dashboardgroup/string2", verifyRequest(t, "POST", true, http.StatusNoContent, nil, ""))

	err := client.LinkDashboardGroupToTeam(context.Background(), "string", "string2")
	assert.NoError(t, err, "Unexpected error linking dashboard group")
}

func TestUnlinkDashboardGroupFromTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string/dashboardgroup/string2", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.UnlinkDashboardGroupFromTeam(context.Background(), "string", "string2")
	assert.NoError(t, err, "Unexpected error linking dashboard group")
}
