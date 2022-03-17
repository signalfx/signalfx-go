/*
 * Integrations API
 *
 * https://dev.splunk.com/observability/reference/api/integrations/latest
 */

package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// IntegrationAPIURL is the base URL for interacting with integrations.
const IntegrationAPIURL = "/v2/integration"

// GetIntegration gets an integration as map.
func (c *Client) GetIntegration(ctx context.Context, id string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	err := c.getIntegration(ctx, id, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// DeleteIntegration deletes an integration.
func (c *Client) DeleteIntegration(ctx context.Context, id string) error {
	return c.doIntegrationRequest(ctx, IntegrationAPIURL+"/"+id, "DELETE", http.StatusNoContent, nil, nil)
}

func (c *Client) createIntegration(ctx context.Context, in interface{}, out interface{}) error {
	return c.doIntegrationRequest(ctx, IntegrationAPIURL, "POST", http.StatusOK, in, out)
}

func (c *Client) getIntegration(ctx context.Context, id string, out interface{}) error {
	return c.doIntegrationRequest(ctx, IntegrationAPIURL+"/"+id, "GET", http.StatusOK, nil, out)
}

func (c *Client) updateIntegration(ctx context.Context, id string, in interface{}, out interface{}) error {
	return c.doIntegrationRequest(ctx, IntegrationAPIURL+"/"+id, "PUT", http.StatusOK, in, out)
}

func (c *Client) doIntegrationRequest(ctx context.Context, url string, method string, status int, in interface{}, out interface{}) error {
	var body io.Reader

	if in != nil {
		payload, err := json.Marshal(in)
		if err != nil {
			return err
		}

		body = bytes.NewReader(payload)
	}

	resp, err := c.doRequest(ctx, method, url, nil, body)
	if resp != nil {
		//noinspection GoUnhandledErrorResult
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != status {
		message, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d: %s", resp.StatusCode, message)
	}

	if out != nil {
		err = json.NewDecoder(resp.Body).Decode(out)
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return err
}
