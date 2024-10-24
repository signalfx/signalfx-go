package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/datalink"
)

// DataLinkAPIURL is the base URL for interacting with data link.
const DataLinkAPIURL = "/v2/crosslink"

// CreateDataLink creates a data link.
func (c *Client) CreateDataLink(ctx context.Context, dataLinkRequest *datalink.CreateUpdateDataLinkRequest) (*datalink.DataLink, error) {
	payload, err := json.Marshal(dataLinkRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", DataLinkAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDataLink := &datalink.DataLink{}

	err = json.NewDecoder(resp.Body).Decode(finalDataLink)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDataLink, err
}

// DeleteDataLink deletes a data link.
func (c *Client) DeleteDataLink(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", DataLinkAPIURL+"/"+id, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// The API returns a 200 here, which I think is a mistake so covering for
	// future changes.
	if err = newResponseError(resp, http.StatusNoContent, http.StatusOK); err != nil {
		return err
	}

	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return nil
}

// GetDataLink gets a data link.
func (c *Client) GetDataLink(ctx context.Context, id string) (*datalink.DataLink, error) {
	resp, err := c.doRequest(ctx, "GET", DataLinkAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDataLink := &datalink.DataLink{}

	err = json.NewDecoder(resp.Body).Decode(finalDataLink)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDataLink, err
}

// UpdateDataLink updates a data link.
func (c *Client) UpdateDataLink(ctx context.Context, id string, dataLinkRequest *datalink.CreateUpdateDataLinkRequest) (*datalink.DataLink, error) {
	payload, err := json.Marshal(dataLinkRequest)
	if err != nil {
		return nil, err
	}

	encodedName := url.PathEscape(id)
	resp, err := c.doRequest(ctx, "PUT", DataLinkAPIURL+"/"+encodedName, nil, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDataLink := &datalink.DataLink{}

	err = json.NewDecoder(resp.Body).Decode(finalDataLink)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDataLink, err
}

// SearchDataLinks searches for data links given a query string in `name`.
func (c *Client) SearchDataLinks(ctx context.Context, limit int, context string, offset int) (*datalink.SearchResults, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("context", context)
	params.Add("offset", strconv.Itoa(offset))

	resp, err := c.doRequest(ctx, "GET", DataLinkAPIURL, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDataLinks := &datalink.SearchResults{}

	err = json.NewDecoder(resp.Body).Decode(finalDataLinks)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDataLinks, err
}
