package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/organization"
)

// OrganizationAPIURL is the base URL for interacting with detectors.
const OrganizationAPIURL = "/v2/organization"
const OrganizationMemberAPIURL = "/v2/organization/member"
const OrganizationMembersAPIURL = "/v2/organization/members"

// GetOrganization gets an organization.
func (c *Client) GetOrganization(ctx context.Context, id string) (*organization.Organization, error) {
	resp, err := c.doRequest(ctx, "GET", OrganizationAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalOrganization := &organization.Organization{}

	err = json.NewDecoder(resp.Body).Decode(finalOrganization)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalOrganization, err
}

// GetMember gets a member.
func (c *Client) GetMember(ctx context.Context, id string) (*organization.Member, error) {
	resp, err := c.doRequest(ctx, "GET", OrganizationMemberAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalMember := &organization.Member{}

	err = json.NewDecoder(resp.Body).Decode(finalMember)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalMember, err
}

// DeleteMember deletes a detector.
func (c *Client) DeleteMember(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", OrganizationMemberAPIURL+"/"+id, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return nil
}

// InviteMember invites a member to the organization.
func (c *Client) InviteMember(ctx context.Context, inviteRequest *organization.CreateUpdateMemberRequest) (*organization.Member, error) {
	payload, err := json.Marshal(inviteRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", OrganizationMemberAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalMember := &organization.Member{}

	err = json.NewDecoder(resp.Body).Decode(finalMember)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalMember, err
}

// Updates admin status of a member.
func (c *Client) UpdateMember(ctx context.Context, id string, updateRequest *organization.UpdateMemberRequest) (*organization.Member, error) {
	payload, err := json.Marshal(updateRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "PUT", OrganizationMemberAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalMember := &organization.Member{}

	err = json.NewDecoder(resp.Body).Decode(finalMember)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalMember, err
}

// InviteMembers invites many members to the organization.
func (c *Client) InviteMembers(ctx context.Context, inviteRequest *organization.InviteMembersRequest) (*organization.InviteMembersRequest, error) {
	payload, err := json.Marshal(inviteRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", OrganizationMembersAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalMembers := &organization.InviteMembersRequest{}

	err = json.NewDecoder(resp.Body).Decode(finalMembers)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalMembers, err
}

// GetOrganizationMembers gets members for an org, with an optional search.
func (c *Client) GetOrganizationMembers(ctx context.Context, limit int, query string, offset int, orderBy string) (*organization.MemberSearchResults, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("query", query)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("orderBy", orderBy)

	resp, err := c.doRequest(ctx, "GET", OrganizationMemberAPIURL, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalMembers := &organization.MemberSearchResults{}

	err = json.NewDecoder(resp.Body).Decode(finalMembers)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalMembers, err
}
