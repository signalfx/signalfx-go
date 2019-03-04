package signalfx

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// DashboardGroupAPIURL is the base URL for interacting with dashboard.
const DashboardGroupAPIURL = "/v2/dashboardgroup"

// TODO Clone dashboard to group

// DashboardGroup is a Dashboard Group.
type DashboardGroup struct {
	AuthorizedWriters struct {
		Teams []string `json:"teams,omitempty"`
		Users []string `json:"users,omitempty"`
	} `json:"authorizedWriters,omitempty"`
	Dashboards  []string    `json:"dashboards,omitempty"`
	Description string      `json:"description,omitempty"`
	Name        string      `json:"name,omitempty"`
	Teams       interface{} `json:"teams,omitempty"`
}

// DashboardGroupSearch is the result of a query for DashboardGroups
type DashboardGroupSearch struct {
	Count   int64 `json:"count,omitempty"`
	Results []DashboardGroup
}

// CreateDashboardGroup creates a dashboard.
func (c *Client) CreateDashboardGroup(dashboardGroup *DashboardGroup) (*DashboardGroup, error) {
	payload, err := json.Marshal(dashboardGroup)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", DashboardGroupAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDashboardGroup := &DashboardGroup{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboardGroup)

	return finalDashboardGroup, err
}

// DeleteDashboardGroup deletes a dashboard.
func (c *Client) DeleteDashboardGroup(id string) error {
	resp, err := c.doRequest("DELETE", DashboardGroupAPIURL+"/"+id, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Unexpected status code: " + resp.Status)
	}

	return nil
}

// GetDashboardGroup gets a dashboard group.
func (c *Client) GetDashboardGroup(id string) (*DashboardGroup, error) {
	resp, err := c.doRequest("GET", DashboardGroupAPIURL+"/"+id, nil, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDashboardGroup := &DashboardGroup{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboardGroup)

	return finalDashboardGroup, err
}

// UpdateDashboardGroup updates a dashboard group.
func (c *Client) UpdateDashboardGroup(id string, dashboardGroup *DashboardGroup) (*DashboardGroup, error) {
	payload, err := json.Marshal(dashboardGroup)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("PUT", DashboardGroupAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDashboardGroup := &DashboardGroup{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboardGroup)

	return finalDashboardGroup, err
}

// SearchDashboardGroup searches for dashboard groups, given a query string in `name`.
func (c *Client) SearchDashboardGroup(limit int, name string, offset int, tags string) (*DashboardGroupSearch, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	resp, err := c.doRequest("GET", DashboardGroupAPIURL, params, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDashboardGroups := &DashboardGroupSearch{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboardGroups)

	return finalDashboardGroups, err
}
