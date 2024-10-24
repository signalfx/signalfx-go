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

// CreateJiraIntegration creates an Jira integration.
func (c *Client) CreateJiraIntegration(ctx context.Context, ji *integration.JiraIntegration) (*integration.JiraIntegration, error) {
	payload, err := json.Marshal(ji)
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

	finalIntegration := integration.JiraIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// GetJiraIntegration retrieves an Jira integration.
func (c *Client) GetJiraIntegration(ctx context.Context, id string) (*integration.JiraIntegration, error) {
	resp, err := c.doRequest(ctx, "GET", IntegrationAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalIntegration := integration.JiraIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// UpdateJiraIntegration updates an Jira integration.
func (c *Client) UpdateJiraIntegration(ctx context.Context, id string, ji *integration.JiraIntegration) (*integration.JiraIntegration, error) {
	payload, err := json.Marshal(ji)
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

	finalIntegration := integration.JiraIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// DeleteJiraIntegration deletes an Jira integration.
func (c *Client) DeleteJiraIntegration(ctx context.Context, id string) error {
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
