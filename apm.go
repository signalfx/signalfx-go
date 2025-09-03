package signalfx

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/signalfx/signalfx-go/apm"
)

const APMTopologyURL = "/v2/apm/topology"

func (c *Client) ListTopology(ctx context.Context, req *apm.RetrieveServiceTopologyRequest) (*apm.RetrieveServiceTopologyResponse, error) {
	content := c.leaseBuffer()
	defer c.releaseBuffer(content)

	if err := json.NewEncoder(content).Encode(req); err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, APMTopologyURL, nil, content)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var topologyResp apm.RetrieveServiceTopologyResponse
	if err := json.NewDecoder(resp.Body).Decode(&topologyResp); err != nil {
		return nil, err
	}

	return &topologyResp, nil
}
