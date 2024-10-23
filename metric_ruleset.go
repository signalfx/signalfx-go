package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/signalfx/signalfx-go/metric_ruleset"
)

// MetricRulesetApiURL is the base URL for interacting with metric rulesets.
const MetricRulesetApiURL = "/v2/metricruleset"

// GetMetricRuleset gets a metric ruleset.
func (c *Client) GetMetricRuleset(ctx context.Context, id string) (*metric_ruleset.GetMetricRulesetResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, MetricRulesetApiURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	metricRuleset := &metric_ruleset.GetMetricRulesetResponse{}
	err = json.NewDecoder(resp.Body).Decode(&metricRuleset)
	io.Copy(ioutil.Discard, resp.Body)

	return metricRuleset, err
}

// CreateMetricRuleset creates a metric ruleset.
func (c *Client) CreateMetricRuleset(ctx context.Context, metricRuleset *metric_ruleset.CreateMetricRulesetRequest) (*metric_ruleset.CreateMetricRulesetResponse, error) {
	payload, err := json.Marshal(metricRuleset)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, MetricRulesetApiURL, nil, bytes.NewReader(payload))
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	createdMetricRuleset := &metric_ruleset.CreateMetricRulesetResponse{}
	err = json.NewDecoder(resp.Body).Decode(&createdMetricRuleset)
	io.Copy(ioutil.Discard, resp.Body)

	return createdMetricRuleset, err
}

// UpdateMetricRuleset updates a metric ruleset.
func (c *Client) UpdateMetricRuleset(ctx context.Context, id string, metricRuleset *metric_ruleset.UpdateMetricRulesetRequest) (*metric_ruleset.UpdateMetricRulesetResponse, error) {
	payload, err := json.Marshal(metricRuleset)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, http.MethodPut, MetricRulesetApiURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	updatedMetricRuleset := &metric_ruleset.UpdateMetricRulesetResponse{}
	err = json.NewDecoder(resp.Body).Decode(&updatedMetricRuleset)
	io.Copy(ioutil.Discard, resp.Body)

	return updatedMetricRuleset, err
}

// DeleteMetricRuleset deletes a metric ruleset.
func (c *Client) DeleteMetricRuleset(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, MetricRulesetApiURL+"/"+id, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}

	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

func (c *Client) GenerateAggregationMetricName(ctx context.Context, generateAggregationNameRequest metric_ruleset.GenerateAggregationNameRequest) (string, error) {
	payload, err := json.Marshal(generateAggregationNameRequest)
	if err != nil {
		return "", err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, MetricRulesetApiURL+"/generateAggregationMetricName", nil, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	aggregationMetricName := string(respBody)

	return aggregationMetricName, err
}
