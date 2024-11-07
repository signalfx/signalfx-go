package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/signalfx/signalfx-go/slo"
)

const SloAPIURL = "/v2/slo"

func (c *Client) GetSlo(ctx context.Context, id string) (*slo.SloObject, error) {
	return c.executeSloRequest(ctx, SloAPIURL+"/"+id, http.MethodGet, http.StatusOK, nil)
}

func (c *Client) CreateSlo(ctx context.Context, sloRequest *slo.SloObject) (*slo.SloObject, error) {
	return c.executeSloRequest(ctx, SloAPIURL, http.MethodPost, http.StatusOK, sloRequest)
}

func (c *Client) ValidateSlo(ctx context.Context, sloRequest *slo.SloObject) error {
	_, err := c.executeSloRequest(ctx, SloAPIURL+"/validate", http.MethodPost, http.StatusNoContent, sloRequest)
	return err
}

func (c *Client) UpdateSlo(ctx context.Context, id string, sloRequest *slo.SloObject) (*slo.SloObject, error) {
	return c.executeSloRequest(ctx, SloAPIURL+"/"+id, http.MethodPut, http.StatusOK, sloRequest)
}

func (c *Client) DeleteSlo(ctx context.Context, id string) error {
	_, err := c.executeSloRequest(ctx, SloAPIURL+"/"+id, http.MethodDelete, http.StatusNoContent, nil)
	return err
}

func (c *Client) executeSloRequest(ctx context.Context, url string, method string, expectedValidStatus int, sloRequest *slo.SloObject) (*slo.SloObject, error) {
	var body io.Reader

	if sloRequest != nil {
		payload, err := json.Marshal(sloRequest)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(payload)
	}

	resp, err := c.doRequest(ctx, method, url, nil, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, expectedValidStatus); err != nil {
		return nil, err
	}

	if expectedValidStatus == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, nil
	}

	returnedSlo := &slo.SloObject{}
	err = json.NewDecoder(resp.Body).Decode(returnedSlo)
	_, _ = io.Copy(io.Discard, resp.Body)
	return returnedSlo, err
}
