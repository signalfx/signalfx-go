/*
 * Integrations API for BigPanda
 *
 * https://dev.splunk.com/observability/reference/api/integrations/latest
 * https://docs.splunk.com/Observability/admin/notif-services/bigpanda.html
 */

package signalfx

import (
	"context"

	"github.com/signalfx/signalfx-go/integration"
)

// CreateBigPandaIntegration creates a BigPanda integration.
func (c *Client) CreateBigPandaIntegration(ctx context.Context, in *integration.BigPandaIntegration) (*integration.BigPandaIntegration, error) {
	out := integration.BigPandaIntegration{}

	err := c.createIntegration(ctx, in, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// GetBigPandaIntegration retrieves a BigPanda integration.
func (c *Client) GetBigPandaIntegration(ctx context.Context, id string) (*integration.BigPandaIntegration, error) {
	out := integration.BigPandaIntegration{}

	err := c.getIntegration(ctx, id, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// UpdateBigPandaIntegration updates a BigPanda integration.
func (c *Client) UpdateBigPandaIntegration(ctx context.Context, id string, in *integration.BigPandaIntegration) (*integration.BigPandaIntegration, error) {
	out := integration.BigPandaIntegration{}

	err := c.updateIntegration(ctx, id, in, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// DeleteBigPandaIntegration deletes a BigPanda integration.
func (c *Client) DeleteBigPandaIntegration(ctx context.Context, id string) error {
	return c.DeleteIntegration(ctx, id)
}
