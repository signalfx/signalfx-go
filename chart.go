package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/chart"
)

// ChartAPIURL is the base URL for interacting with charts.
const ChartAPIURL = "/v2/chart"
const CreateSloChartAPIURL = ChartAPIURL + "/createSloChart"
const UpdateSloChartAPIURL = ChartAPIURL + "/updateSloChart"

// CreateChart creates a chart.
func (c *Client) CreateChart(ctx context.Context, chartRequest *chart.CreateUpdateChartRequest) (*chart.Chart, error) {
	return c.executeChartRequest(ctx, ChartAPIURL, http.MethodPost, http.StatusOK, chartRequest)
}

// CreateSloChart creates a SLO chart.
func (c *Client) CreateSloChart(ctx context.Context, chartRequest *chart.CreateUpdateSloChartRequest) (*chart.Chart, error) {
	return c.executeChartRequest(ctx, CreateSloChartAPIURL, http.MethodPost, http.StatusOK, chartRequest)
}

// DeleteChart deletes a chart.
func (c *Client) DeleteChart(ctx context.Context, id string) error {
	_, err := c.executeChartRequest(ctx, ChartAPIURL+"/"+id, http.MethodDelete, http.StatusOK, nil)
	if err == io.EOF {
		// Expected error as delete request returns status 200 instead of 204
		return nil
	}
	return err
}

// GetChart gets a chart.
func (c *Client) GetChart(ctx context.Context, id string) (*chart.Chart, error) {
	return c.executeChartRequest(ctx, ChartAPIURL+"/"+id, http.MethodGet, http.StatusOK, nil)
}

// UpdateChart updates a chart.
func (c *Client) UpdateChart(ctx context.Context, id string, chartRequest *chart.CreateUpdateChartRequest) (*chart.Chart, error) {
	return c.executeChartRequest(ctx, ChartAPIURL+"/"+id, http.MethodPut, http.StatusOK, chartRequest)
}

// UpdateSloChart updates an SLO chart.
func (c *Client) UpdateSloChart(ctx context.Context, id string, chartRequest *chart.CreateUpdateSloChartRequest) (*chart.Chart, error) {
	return c.executeChartRequest(ctx, UpdateSloChartAPIURL+"/"+id, http.MethodPut, http.StatusOK, chartRequest)
}

// ValidateChart validates a chart.
func (c *Client) ValidateChart(ctx context.Context, chartRequest *chart.CreateUpdateChartRequest) error {
	_, err := c.executeChartRequest(ctx, ChartAPIURL+"/validate", http.MethodPost, http.StatusNoContent, chartRequest)
	return err
}

func (c *Client) executeChartRequest(ctx context.Context, url string, method string, expectedValidStatus int, chartRequest interface{}) (*chart.Chart, error) {
	var body io.Reader
	if chartRequest != nil {
		payload, err := json.Marshal(chartRequest)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(payload)
	}

	resp, err := c.doRequest(ctx, method, url, nil, body)
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

	returnedChart := &chart.Chart{}
	err = json.NewDecoder(resp.Body).Decode(returnedChart)
	_, _ = io.Copy(io.Discard, resp.Body)
	return returnedChart, err
}

// SearchCharts searches for charts, given a query string in `name`.
func (c *Client) SearchCharts(ctx context.Context, limit int, name string, offset int, tags string) (*chart.SearchResult, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	if name != "" {
		params.Add("name", name)
	}
	params.Add("offset", strconv.Itoa(offset))
	if tags != "" {
		params.Add("tags", tags)
	}

	resp, err := c.doRequest(ctx, "GET", ChartAPIURL, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalCharts := &chart.SearchResult{}

	err = json.NewDecoder(resp.Body).Decode(finalCharts)
	_, _ = io.Copy(io.Discard, resp.Body)

	return finalCharts, err
}
