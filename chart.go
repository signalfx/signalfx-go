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

	"github.com/signalfx/signalfx-go/chart"
)

// ChartAPIURL is the base URL for interacting with charts.
const ChartAPIURL = "/v2/chart"

// CreateChart creates a chart.
func (c *Client) CreateChart(ctx context.Context, chartRequest *chart.CreateUpdateChartRequest) (*chart.Chart, error) {
	payload, err := json.Marshal(chartRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", ChartAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalChart := &chart.Chart{}

	err = json.NewDecoder(resp.Body).Decode(finalChart)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalChart, err
}

// DeleteChart deletes a chart.
func (c *Client) DeleteChart(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", ChartAPIURL+"/"+id, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return nil
}

// GetChart gets a chart.
func (c *Client) GetChart(ctx context.Context, id string) (*chart.Chart, error) {
	resp, err := c.doRequest(ctx, "GET", ChartAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalChart := &chart.Chart{}

	err = json.NewDecoder(resp.Body).Decode(finalChart)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalChart, err
}

// UpdateChart updates a chart.
func (c *Client) UpdateChart(ctx context.Context, id string, chartRequest *chart.CreateUpdateChartRequest) (*chart.Chart, error) {
	payload, err := json.Marshal(chartRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "PUT", ChartAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalChart := &chart.Chart{}

	err = json.NewDecoder(resp.Body).Decode(finalChart)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalChart, err
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
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalCharts, err
}
