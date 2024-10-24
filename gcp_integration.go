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

// CreateGCPIntegration creates a GCP integration.
func (c *Client) CreateGCPIntegration(ctx context.Context, gcpi *integration.GCPIntegration) (*integration.GCPIntegration, error) {
	payload, err := json.Marshal(gcpi)
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

	finalIntegration := integration.GCPIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// GetGCPIntegration retrieves a GCP integration.
func (c *Client) GetGCPIntegration(ctx context.Context, id string) (*integration.GCPIntegration, error) {
	resp, err := c.doRequest(ctx, "GET", IntegrationAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalIntegration := integration.GCPIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// UpdateGCPIntegration updates a GCP integration.
func (c *Client) UpdateGCPIntegration(ctx context.Context, id string, gcpi *integration.GCPIntegration) (*integration.GCPIntegration, error) {
	payload, err := json.Marshal(gcpi)
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

	finalIntegration := integration.GCPIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// DeleteGCPIntegration deletes a GCP integration.
func (c *Client) DeleteGCPIntegration(ctx context.Context, id string) error {
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
