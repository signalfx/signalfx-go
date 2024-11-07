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

	"github.com/signalfx/signalfx-go/dashboard_group"
)

// DashboardGroupAPIURL is the base URL for interacting with dashboard.
const DashboardGroupAPIURL = "/v2/dashboardgroup"

// TODO Clone dashboard to group

// CreateDashboardGroup creates a dashboard.
func (c *Client) CreateDashboardGroup(ctx context.Context, dashboardGroupRequest *dashboard_group.CreateUpdateDashboardGroupRequest, skipImplicitDashboard bool) (*dashboard_group.DashboardGroup, error) {
	payload, err := json.Marshal(dashboardGroupRequest)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	if skipImplicitDashboard {
		params.Add("empty", "true")
	}

	resp, err := c.doRequest(ctx, "POST", DashboardGroupAPIURL, params, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDashboardGroup := &dashboard_group.DashboardGroup{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboardGroup)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDashboardGroup, err
}

// DeleteDashboardGroup deletes a dashboard.
func (c *Client) DeleteDashboardGroup(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", DashboardGroupAPIURL+"/"+id, nil, nil)
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

// GetDashboardGroup gets a dashboard group.
func (c *Client) GetDashboardGroup(ctx context.Context, id string) (*dashboard_group.DashboardGroup, error) {
	resp, err := c.doRequest(ctx, "GET", DashboardGroupAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDashboardGroup := &dashboard_group.DashboardGroup{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboardGroup)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDashboardGroup, err
}

// UpdateDashboardGroup updates a dashboard group.
func (c *Client) UpdateDashboardGroup(ctx context.Context, id string, dashboardGroupRequest *dashboard_group.CreateUpdateDashboardGroupRequest) (*dashboard_group.DashboardGroup, error) {
	payload, err := json.Marshal(dashboardGroupRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "PUT", DashboardGroupAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDashboardGroup := &dashboard_group.DashboardGroup{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboardGroup)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDashboardGroup, err
}

// ValidateDashboardGroup validates a dashboard grouop.
func (c *Client) ValidateDashboardGroup(ctx context.Context, dashboardGroupRequest *dashboard_group.CreateUpdateDashboardGroupRequest) error {
	payload, err := json.Marshal(dashboardGroupRequest)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(ctx, "POST", DashboardGroupAPIURL+"/validate", nil, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}

// SearchDashboardGroup searches for dashboard groups, given a query string in `name`.
func (c *Client) SearchDashboardGroups(ctx context.Context, limit int, name string, offset int) (*dashboard_group.SearchResult, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))

	resp, err := c.doRequest(ctx, "GET", DashboardGroupAPIURL, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDashboardGroups := &dashboard_group.SearchResult{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboardGroups)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDashboardGroups, err
}
