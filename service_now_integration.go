/*
 * Integrations API for ServiceNow
 *
 * https://dev.splunk.com/observability/reference/api/integrations/latest
 * https://docs.splunk.com/Observability/admin/notif-services/servicenow.html
 */

package signalfx

import (
	"context"
	"github.com/signalfx/signalfx-go/integration"
)

// CreateServiceNowIntegration creates SNOW integration.
func (c *Client) CreateServiceNowIntegration(ctx context.Context, in *integration.ServiceNowIntegration) (*integration.ServiceNowIntegration, error) {
	out := integration.ServiceNowIntegration{}

	err := c.createIntegration(ctx, in, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// GetServiceNowIntegration retrieves SNOW integration.
func (c *Client) GetServiceNowIntegration(ctx context.Context, id string) (*integration.ServiceNowIntegration, error) {
	out := integration.ServiceNowIntegration{}

	err := c.getIntegration(ctx, id, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// UpdateServiceNowIntegration updates SNOW integration.
func (c *Client) UpdateServiceNowIntegration(ctx context.Context, id string, in *integration.ServiceNowIntegration) (*integration.ServiceNowIntegration, error) {
	out := integration.ServiceNowIntegration{}

	err := c.updateIntegration(ctx, id, in, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// DeleteServiceNowIntegration deletes SNOW integration.
func (c *Client) DeleteServiceNowIntegration(ctx context.Context, id string) error {
	return c.DeleteIntegration(ctx, id)
}
