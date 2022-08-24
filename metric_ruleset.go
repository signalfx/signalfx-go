package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/signalfx/signalfx-go/metric_ruleset"
)

// MetricRulesetApiURL is the base URL for interacting with metric rulesets.
const MetricRulesetApiURL = "/v2/metricruleset"

// GetMetricRuleset gets a metric ruleset.
func (c *Client) GetMetricRuleset(ctx context.Context, id string) (*metric_ruleset.GetMetricRulesetResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, MetricRulesetApiURL+"/"+id, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status %d: %s", resp.StatusCode, message)
	}

	metricRuleset := &metric_ruleset.GetMetricRulesetResponse{}
	err = json.NewDecoder(resp.Body).Decode(&metricRuleset)
	io.Copy(io.Discard, resp.Body)

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

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status %d: %s", resp.StatusCode, message)
	}

	createdMetricRuleset := &metric_ruleset.CreateMetricRulesetResponse{}
	err = json.NewDecoder(resp.Body).Decode(&createdMetricRuleset)
	io.Copy(io.Discard, resp.Body)

	return createdMetricRuleset, err
}

// UpdateMetricRuleset updates a metric ruleset.
func (c *Client) UpdateMetricRuleset(ctx context.Context, id string, metricRuleset *metric_ruleset.UpdateMetricRulesetRequest) (*metric_ruleset.UpdateMetricRulesetResponse, error) {
	payload, err := json.Marshal(metricRuleset)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, http.MethodPut, MetricRulesetApiURL+"/"+id, nil, bytes.NewReader(payload))
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status %d: %s", resp.StatusCode, message)
	}

	updatedMetricRuleset := &metric_ruleset.UpdateMetricRulesetResponse{}
	err = json.NewDecoder(resp.Body).Decode(&updatedMetricRuleset)
	io.Copy(io.Discard, resp.Body)

	return updatedMetricRuleset, err
}

// DeleteMetricRuleset deletes a metric ruleset.
func (c *Client) DeleteMetricRuleset(ctx context.Context, id string) (error) {
	resp, err := c.doRequest(ctx, http.MethodDelete, MetricRulesetApiURL+"/"+id, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		message, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status %d: %s", resp.StatusCode, message)
	}

	io.Copy(io.Discard, resp.Body)

	return nil
}