package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/dashboard_group"
	"github.com/signalfx/signalfx-go/util"
)

// DashboardGroupAPIURL is the base URL for interacting with dashboard.
const DashboardGroupAPIURL = "/v2/dashboardgroup"

// TODO Clone dashboard to group

// CreateDashboardGroup creates a dashboard.
func (c *Client) CreateDashboardGroup(ctx context.Context, dashboardGroupRequest *dashboard_group.CreateUpdateDashboardGroupRequest, skipImplicitDashboard bool) (*dashboard_group.DashboardGroup, error) {
	params := url.Values{}
	if skipImplicitDashboard {
		params.Add("empty", "true")
	}

	return c.executeDashboardGroupRequest(ctx, DashboardGroupAPIURL, http.MethodPost, http.StatusOK, dashboardGroupRequest, params)
}

// DeleteDashboardGroup deletes a dashboard.
func (c *Client) DeleteDashboardGroup(ctx context.Context, id string) error {
	_, err := c.executeDashboardGroupRequest(ctx, DashboardGroupAPIURL+"/"+id, http.MethodDelete, http.StatusNoContent, nil, nil)
	return err
}

// GetDashboardGroup gets a dashboard group.
func (c *Client) GetDashboardGroup(ctx context.Context, id string) (*dashboard_group.DashboardGroup, error) {
	return c.executeDashboardGroupRequest(ctx, DashboardGroupAPIURL+"/"+id, http.MethodGet, http.StatusOK, nil, nil)
}

// UpdateDashboardGroup updates a dashboard group.
func (c *Client) UpdateDashboardGroup(ctx context.Context, id string, dashboardGroupRequest *dashboard_group.CreateUpdateDashboardGroupRequest) (*dashboard_group.DashboardGroup, error) {
	return c.executeDashboardGroupRequest(ctx, DashboardGroupAPIURL+"/"+id, http.MethodPut, http.StatusOK, dashboardGroupRequest, nil)
}

// ValidateDashboardGroup validates a dashboard grouop with default mode.
func (c *Client) ValidateDashboardGroup(ctx context.Context, dashboardGroupRequest *dashboard_group.CreateUpdateDashboardGroupRequest) error {
	return c.ValidateDashboardGroupWithMode(ctx, dashboardGroupRequest, util.FULL)
}

// ValidateDashboardGroupWithMode validates a dashboard grouop.
func (c *Client) ValidateDashboardGroupWithMode(ctx context.Context, dashboardGroupRequest *dashboard_group.CreateUpdateDashboardGroupRequest, validationMode util.ValidationMode) error {
	if err := util.ValidateValidationMode(validationMode); err != nil {
		return err
	}
	params := url.Values{}
	params.Add("validationMode", string(validationMode))
	_, err := c.executeDashboardGroupRequest(ctx, DashboardGroupAPIURL+"/validate", http.MethodPost, http.StatusNoContent, dashboardGroupRequest, params)
	return err
}

func (c *Client) executeDashboardGroupRequest(ctx context.Context, url string, method string, expectedValidStatus int, dashboardGroupRequest *dashboard_group.CreateUpdateDashboardGroupRequest, params url.Values) (*dashboard_group.DashboardGroup, error) {
	var body io.Reader
	if dashboardGroupRequest != nil {
		payload, err := json.Marshal(dashboardGroupRequest)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(payload)
	}

	resp, err := c.doRequest(ctx, method, url, params, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, expectedValidStatus); err != nil {
		return nil, err
	}

	if expectedValidStatus == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, nil
	}

	returnedDashboardGroup := &dashboard_group.DashboardGroup{}
	err = json.NewDecoder(resp.Body).Decode(returnedDashboardGroup)
	_, _ = io.Copy(io.Discard, resp.Body)
	return returnedDashboardGroup, err
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
	_, _ = io.Copy(io.Discard, resp.Body)

	return finalDashboardGroups, err
}
