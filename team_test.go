package signalfx

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team", verifyRequest(t, "POST", http.StatusOK, nil, "team/create_success.json"))

	result, err := client.CreateTeam(&Team{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error creating team")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeleteTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string", verifyRequest(t, "DELETE", http.StatusOK, nil, ""))

	err := client.DeleteTeam("string")
	assert.NoError(t, err, "Unexpected error getting Team")
}

func TestGetTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string", verifyRequest(t, "GET", http.StatusOK, nil, "team/get_success.json"))

	result, err := client.GetTeam("string")
	assert.NoError(t, err, "Unexpected error getting Team")
	assert.Equal(t, result.Name, "string", "Name does not match")
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

	mux.HandleFunc("/v2/team", verifyRequest(t, "GET", http.StatusOK, params, "team/get_success.json"))

	results, err := client.SearchTeam(limit, name, offset, tags)
	assert.NoError(t, err, "Unexpected error search Team")
	assert.Equal(t, int64(0), results.Count, "Incorrect number of results")
}

func TestUpdateTeam(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/team/string", verifyRequest(t, "PUT", http.StatusOK, nil, "team/update_success.json"))

	result, err := client.UpdateTeam("string", &Team{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error updating Team")
	assert.Equal(t, "string", result.Name, "Name does not match")
}
