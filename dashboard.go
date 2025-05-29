package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/dashboard"
	"github.com/signalfx/signalfx-go/util"
)

// TODO Create simple dashboard

// DashboardAPIURL is the base URL for interacting with dashboard.
const DashboardAPIURL = "/v2/dashboard"

// CreateDashboard creates a dashboard.
func (c *Client) CreateDashboard(ctx context.Context, dashboardRequest *dashboard.CreateUpdateDashboardRequest) (*dashboard.Dashboard, error) {
	return c.executeDashboardRequest(ctx, DashboardAPIURL, http.MethodPost, http.StatusOK, dashboardRequest, nil)
}

// DeleteDashboard deletes a dashboard.
func (c *Client) DeleteDashboard(ctx context.Context, id string) error {
	_, err := c.executeDashboardRequest(ctx, DashboardAPIURL+"/"+id, http.MethodDelete, http.StatusOK, nil, nil)
	if err == io.EOF {
		// Expected error as delete request returns status 200 instead of 204
		return nil
	}
	return err
}

// GetDashboard gets a dashboard.
func (c *Client) GetDashboard(ctx context.Context, id string) (*dashboard.Dashboard, error) {
	return c.executeDashboardRequest(ctx, DashboardAPIURL+"/"+id, http.MethodGet, http.StatusOK, nil, nil)
}

// UpdateDashboard updates a dashboard.
func (c *Client) UpdateDashboard(ctx context.Context, id string, dashboardRequest *dashboard.CreateUpdateDashboardRequest) (*dashboard.Dashboard, error) {
	return c.executeDashboardRequest(ctx, DashboardAPIURL+"/"+id, http.MethodPut, http.StatusOK, dashboardRequest, nil)
}

// ValidateDashboard validates a dashboard with default mode.
func (c *Client) ValidateDashboard(ctx context.Context, dashboardRequest *dashboard.CreateUpdateDashboardRequest) error {
	return c.ValidateDashboardWithMode(ctx, dashboardRequest, util.FULL)
}

// ValidateDashboard validates a dashboard.
func (c *Client) ValidateDashboardWithMode(ctx context.Context, dashboardRequest *dashboard.CreateUpdateDashboardRequest, validationMode util.ValidationMode) error {
	if err := util.ValidateValidationMode(validationMode); err != nil {
		return err
	}
	params := url.Values{}
	params.Add("validationMode", string(validationMode))
	_, err := c.executeDashboardRequest(ctx, DashboardAPIURL+"/validate", http.MethodPost, http.StatusNoContent, dashboardRequest, params)
	return err
}

func (c *Client) executeDashboardRequest(ctx context.Context, url string, method string, expectedValidStatus int, dashboardRequest *dashboard.CreateUpdateDashboardRequest, params url.Values) (*dashboard.Dashboard, error) {
	var body io.Reader
	if dashboardRequest != nil {
		payload, err := json.Marshal(dashboardRequest)
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

	if err := newResponseError(resp, expectedValidStatus); err != nil {
		return nil, err
	}

	if expectedValidStatus == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, nil
	}

	returnedDashboard := &dashboard.Dashboard{}
	err = json.NewDecoder(resp.Body).Decode(returnedDashboard)
	_, _ = io.Copy(io.Discard, resp.Body)
	return returnedDashboard, err
}

// SearchDashboard searches for dashboards, given a query string in `name`.
func (c *Client) SearchDashboard(ctx context.Context, limit int, name string, offset int, tags string) (*dashboard.SearchResult, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	resp, err := c.doRequest(ctx, "GET", DashboardAPIURL, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDashboards := &dashboard.SearchResult{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboards)
	_, _ = io.Copy(io.Discard, resp.Body)

	return finalDashboards, err
}
