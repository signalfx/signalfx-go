package signalfx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/signalfx/signalfx-go/organization"
)

func TestGetOrganization(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/organization", verifyRequest(t, "GET", true, http.StatusOK, nil, "organization/get_success.json"))

	_, err := client.GetOrganization(context.Background(), "")
	assert.NoError(t, err, "Unexpected error getting organization")
}

func TestGetMissingOrganization(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/organization/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetDetector(context.Background(), "string")
	assert.Error(t, err, "Should have gotten an error from a missing organization")
	assert.Nil(t, result, "Should have gotten a nil result from a missing organization")
}

func TestGetMember(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/organization/member/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "organization/get_member_success.json"))

	result, err := client.GetMember(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting member")
	assert.Equal(t, result.Id, "string", "Id does not match")
}

func TestGetMissingMember(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/organization/member/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetMember(context.Background(), "string")
	assert.Error(t, err, "Should have gotten an error from a missing member")
	assert.Nil(t, result, "Should have gotten a nil result from a missing member")
}

func TestInviteMember(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/organization/member", verifyRequest(t, "POST", true, http.StatusOK, nil, "organization/invite_member_success.json"))

	results, err := client.InviteMember(context.Background(), &organization.CreateUpdateMemberRequest{
		Email: "string",
	})
	assert.NoError(t, err, "Unexpected error inviting member")
	assert.Equal(t, "string", results.Email, "Incorrect email")
}

func TestUpdateMember(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/organization/member/string", verifyRequest(t, "PUT", true, http.StatusOK, nil, "organization/update_member_success.json"))

	results, err := client.UpdateMember(context.Background(), "string", &organization.UpdateMemberRequest{
		Admin: true,
	})
	assert.NoError(t, err, "Unexpected error updating member")
	assert.Equal(t, true, results.Admin, "Incorrect admin config")
}

func TestGetInviteMembers(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/organization/members", verifyRequest(t, "POST", true, http.StatusOK, nil, "organization/invite_members_success.json"))

	members := make([]*organization.Member, 1)
	members[0] = &organization.Member{
		Email: "string",
	}
	results, err := client.InviteMembers(context.Background(), &organization.InviteMembersRequest{
		Members: members,
	})
	assert.NoError(t, err, "Unexpected error inviting members")
	assert.Equal(t, 1, len(results.Members), "Incorrect email")
}

func TestGetOrganizationMembers(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	query := "foo"
	offset := 2
	orderBy := "bar"
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("query", query)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("orderBy", orderBy)

	mux.HandleFunc("/v2/organization/member", verifyRequest(t, "GET", true, http.StatusOK, params, "organization/get_organization_members_success.json"))

	results, err := client.GetOrganizationMembers(context.Background(), limit, query, offset, orderBy)
	assert.NoError(t, err, "Unexpected error getting members")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestGetOrganizationMembersBad(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	query := "foo"
	offset := 2
	orderBy := "bar"
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("query", query)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("orderBy", orderBy)

	mux.HandleFunc("/v2/organization/member", verifyRequest(t, "GET", true, http.StatusBadRequest, params, ""))

	_, err := client.GetOrganizationMembers(context.Background(), limit, query, offset, orderBy)
	assert.Error(t, err, "Unexpected error getting members")
}

func TestDeleteMember(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/organization/member/string", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteMember(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error deleting member")
}

func TestDeleteMissingMember(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/organization/member", verifyRequest(t, "POST", true, http.StatusNotFound, nil, ""))

	err := client.DeleteMember(context.Background(), "example")
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}
