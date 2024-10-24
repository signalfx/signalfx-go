package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/team"
)

// TeamAPIURL is the base URL for interacting with teams.
const TeamAPIURL = "/v2/team"

// CreateTeam creates a team.
func (c *Client) CreateTeam(ctx context.Context, t *team.CreateUpdateTeamRequest) (*team.Team, error) {
	payload, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", TeamAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalTeam := &team.Team{}

	err = json.NewDecoder(resp.Body).Decode(finalTeam)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalTeam, err
}

// DeleteTeam deletes a team.
func (c *Client) DeleteTeam(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", TeamAPIURL+"/"+id, nil, nil)
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

// GetTeam gets a team.
func (c *Client) GetTeam(ctx context.Context, id string) (*team.Team, error) {
	resp, err := c.doRequest(ctx, "GET", TeamAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalTeam := &team.Team{}

	err = json.NewDecoder(resp.Body).Decode(finalTeam)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalTeam, err
}

// UpdateTeam updates a team.
func (c *Client) UpdateTeam(ctx context.Context, id string, t *team.CreateUpdateTeamRequest) (*team.Team, error) {
	payload, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "PUT", TeamAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalTeam := &team.Team{}

	err = json.NewDecoder(resp.Body).Decode(finalTeam)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalTeam, err
}

// SearchTeam searches for teams, given a query string in `name`.
func (c *Client) SearchTeam(ctx context.Context, limit int, name string, offset int, tags string) (*team.SearchResults, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	resp, err := c.doRequest(ctx, "GET", TeamAPIURL, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalTeams := &team.SearchResults{}

	err = json.NewDecoder(resp.Body).Decode(finalTeams)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalTeams, err
}

// LinkDetectorToTeam links a detector to a team.
func (c *Client) LinkDetectorToTeam(ctx context.Context, id string, detectorID string) error {
	targetURL := fmt.Sprintf("%s/%s/detector/%s", TeamAPIURL, id, detectorID)
	resp, err := c.doRequest(ctx, "POST", targetURL, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}

	return nil
}

// UnLinkDetectorFromTeam unlinks a detector from a team.
func (c *Client) UnlinkDetectorFromTeam(ctx context.Context, id string, detectorID string) error {
	targetURL := fmt.Sprintf("%s/%s/detector/%s", TeamAPIURL, id, detectorID)
	resp, err := c.doRequest(ctx, "DELETE", targetURL, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}

	return nil
}

// LinkDashboardGroupToTeam links a dashboard group to a team.
func (c *Client) LinkDashboardGroupToTeam(ctx context.Context, id string, dashboardGroupID string) error {
	targetURL := fmt.Sprintf("%s/%s/dashboardgroup/%s", TeamAPIURL, id, dashboardGroupID)
	resp, err := c.doRequest(ctx, "POST", targetURL, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}

	return nil
}

// UnlinkDashboardGroupFromTeam unlinks a dashboard group from a team.
func (c *Client) UnlinkDashboardGroupFromTeam(ctx context.Context, id string, dashboardGroupID string) error {
	targetURL := fmt.Sprintf("%s/%s/dashboardgroup/%s", TeamAPIURL, id, dashboardGroupID)
	resp, err := c.doRequest(ctx, "DELETE", targetURL, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}

	return nil
}
