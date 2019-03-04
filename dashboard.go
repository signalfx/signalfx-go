package signalfx

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TODO Create simple dashboard

// DashboardAPIURL is the base URL for interacting with dashboard.
const DashboardAPIURL = "/v2/dashboard"

// Dashboard is a dashboard.
type Dashboard struct {
	AuthorizedWriters struct {
		Teams []string `json:"teams,omitempty"`
		Users []string `json:"users,omitempty"`
	} `json:"authorizedWriters,omitempty"`
	ChartDensity string `json:"chartDensity,omitempty"`
	Charts       []struct {
		ChartID string `json:"chartId,omitempty"`
		Column  int64  `json:"column,omitempty"`
		Height  int64  `json:"height,omitempty"`
		Row     int64  `json:"row,omitempty"`
		Width   int64  `json:"width,omitempty"`
	} `json:"charts,omitempty"`
	Description   string `json:"description,omitempty"`
	EventOverlays []struct {
		Not      bool     `json:"NOT,omitempty"`
		Property string   `json:"property,omitempty"`
		Value    []string `json:"value,omitempty"`
	} `json:"eventOverlays"`
	Filters struct {
		Sources []struct {
			Not      bool     `json:"NOT,omitempty"`
			Property string   `json:"property,omitempty"`
			Value    []string `json:"value,omitempty"`
		} `json:"sources"`
		Time struct {
			End   string `json:"end,omitempty"`
			Start string `json:"start,omitempty"`
		} `json:"time,omitempty"`
		Variables []struct {
			Alias                string   `json:"alias,omitempty"`
			PreferredSuggestions []string `json:"preferredSuggestions,omitempty"`
			Property             string   `json:"property,omitempty"`
			Required             bool     `json:"required,omitempty"`
			Restricted           bool     `json:"restricted,omitempty"`
			Value                []string `json:"value,omitempty"`
		} `json:"variables,omitempty"`
	} `json:"filters,omitempty"`
	GroupID               string `json:"groupId,omitempty"`
	MaxDelayOverride      int64  `json:"maxDelayOverride,omitempty"`
	Name                  string `json:"name,omitempty"`
	SelectedEventOverlays []struct {
		EventColorIndex int64 `json:"eventColorIndex,omitempty"`
		EventLine       bool  `json:"eventLine,omitempty"`
		EventSignal     struct {
			EventSearchText string `json:"eventSearchText,omitempty"`
			EventType       string `json:"eventType,omitempty"`
		} `json:"eventSignal,omitempty"`
		Label     string `json:"label,omitempty"`
		OverlayID string `json:"overlayId,omitempty"`
		Sources   [][]struct {
			Not      bool     `json:"NOT,omitempty"`
			Property string   `json:"property,omitempty"`
			Value    []string `json:"value,omitempty"`
		} `json:"sources,omitempty"`
	} `json:"selectedEventOverlays,omitempty"`
}

// DashboardSearch is the result of a query for Dashboards
type DashboardSearch struct {
	Count   int64 `json:"count,omitempty"`
	Results []Dashboard
}

// CreateDashboard creates a dashboard.
func (c *Client) CreateDashboard(dashboard *Dashboard) (*Dashboard, error) {
	payload, err := json.Marshal(dashboard)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", DashboardAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDashboard := &Dashboard{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboard)

	return finalDashboard, err
}

// DeleteDashboard deletes a dashboard.
func (c *Client) DeleteDashboard(id string) error {
	resp, err := c.doRequest("DELETE", DashboardAPIURL+"/"+id, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Unexpected status code: " + resp.Status)
	}

	return nil
}

// GetDashboard gets a dashboard.
func (c *Client) GetDashboard(id string) (*Dashboard, error) {
	resp, err := c.doRequest("GET", DashboardAPIURL+"/"+id, nil, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDashboard := &Dashboard{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboard)

	return finalDashboard, err
}

// UpdateDashboard updates a dashboard.
func (c *Client) UpdateDashboard(id string, dashboard *Dashboard) (*Dashboard, error) {
	payload, err := json.Marshal(dashboard)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("PUT", DashboardAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDashboard := &Dashboard{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboard)

	return finalDashboard, err
}

// SearchDashboard searches for dashboards, given a query string in `name`.
func (c *Client) SearchDashboard(limit int, name string, offset int, tags string) (*DashboardSearch, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	resp, err := c.doRequest("GET", DashboardAPIURL, params, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDashboards := &DashboardSearch{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboards)

	return finalDashboards, err
}
