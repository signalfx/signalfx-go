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

// CreateAWSCloudWatchIntegration creates an AWS CloudWatch integration.
func (c *Client) CreateAWSCloudWatchIntegration(ctx context.Context, acwi *integration.AwsCloudWatchIntegration) (*integration.AwsCloudWatchIntegration, error) {
	payload, err := json.Marshal(acwi)
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

	finalIntegration := integration.AwsCloudWatchIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// GetAWSCloudWatchIntegration retrieves an AWS CloudWatch integration.
func (c *Client) GetAWSCloudWatchIntegration(ctx context.Context, id string) (*integration.AwsCloudWatchIntegration, error) {
	resp, err := c.doRequest(ctx, "GET", IntegrationAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalIntegration := integration.AwsCloudWatchIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// UpdateAWSCloudWatchIntegration updates an AWS CloudWatch integration.
func (c *Client) UpdateAWSCloudWatchIntegration(ctx context.Context, id string, acwi *integration.AwsCloudWatchIntegration) (*integration.AwsCloudWatchIntegration, error) {
	payload, err := json.Marshal(acwi)
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

	finalIntegration := integration.AwsCloudWatchIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return &finalIntegration, err
}

// DeleteAWSCloudWatchIntegration deletes an AWS CloudWatch integration.
func (c *Client) DeleteAWSCloudWatchIntegration(ctx context.Context, id string) error {
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
