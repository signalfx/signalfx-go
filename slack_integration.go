package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/signalfx/signalfx-go/integration"
)

// CreateSlackIntegration creates a Slack integration.
func (c *Client) CreateSlackIntegration(ctx context.Context, si *integration.SlackIntegration) (*integration.SlackIntegration, error) {
	payload, err := json.Marshal(si)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", IntegrationAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalIntegration := integration.SlackIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// GetSlackIntegration retrieves a Slack integration.
func (c *Client) GetSlackIntegration(ctx context.Context, id string) (*integration.SlackIntegration, error) {
	resp, err := c.doRequest(ctx, "GET", IntegrationAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalIntegration := integration.SlackIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// UpdateSlackIntegration updates a Slack integration.
func (c *Client) UpdateSlackIntegration(ctx context.Context, id string, si *integration.SlackIntegration) (*integration.SlackIntegration, error) {
	payload, err := json.Marshal(si)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "PUT", IntegrationAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalIntegration := integration.SlackIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// DeleteSlackIntegration deletes a Slack integration.
func (c *Client) DeleteSlackIntegration(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", IntegrationAPIURL+"/"+id, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return err
}
