package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/signalfx/signalfx-go/navigator"
)

// NavigatorAPIURL is the base URL for interacting with navigator.
const NavigatorAPIURL = "/v2/navigator"

// CreateNavigator creates a navigator.
func (c *Client) CreateNavigator(ctx context.Context, navigatorRequest *navigator.CreateNavigatorRequest) (*navigator.Navigator, error) {
	return c.executeNavigatorRequest(ctx, NavigatorAPIURL, http.MethodPost, http.StatusOK, navigatorRequest, nil)
}

// DeleteNavigator deletes a navigator.
func (c *Client) DeleteNavigator(ctx context.Context, id string) error {
	_, err := c.executeNavigatorRequest(ctx, NavigatorAPIURL+"/"+id, http.MethodDelete, http.StatusOK, nil, nil)
	if err == io.EOF {
		// Expected error as delete request returns status 200 instead of 204
		return nil
	}
	return err
}

// GetNavigator gets a navigator.
func (c *Client) GetNavigator(ctx context.Context, id string) (*navigator.Navigator, error) {
	return c.executeNavigatorRequest(ctx, NavigatorAPIURL+"/"+id, http.MethodGet, http.StatusOK, nil, nil)
}

// UpdateNavigator updates a navigator.
func (c *Client) UpdateNavigator(ctx context.Context, id string, navigatorRequest *navigator.UpdateNavigatorRequest) (*navigator.Navigator, error) {
	return c.executeNavigatorUpdateRequest(ctx, NavigatorAPIURL+"/"+id, http.MethodPut, http.StatusOK, navigatorRequest, nil)
}

func (c *Client) executeNavigatorRequest(ctx context.Context, url string, method string, expectedValidStatus int, navigatorRequest *navigator.CreateNavigatorRequest, params url.Values) (*navigator.Navigator, error) {
	var body io.Reader
	if navigatorRequest != nil {
		payload, err := json.Marshal(navigatorRequest)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(payload)
	}

	resp, err := c.doRequest(ctx, method, url, params, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := newResponseError(resp, expectedValidStatus); err != nil {
		return nil, err
	}

	if expectedValidStatus == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, nil
	}

	returnedNavigator := &navigator.Navigator{}
	err = json.NewDecoder(resp.Body).Decode(returnedNavigator)
	_, _ = io.Copy(io.Discard, resp.Body)
	return returnedNavigator, err
}

func (c *Client) executeNavigatorUpdateRequest(ctx context.Context, url string, method string, expectedValidStatus int, navigatorRequest *navigator.UpdateNavigatorRequest, params url.Values) (*navigator.Navigator, error) {
	var body io.Reader
	if navigatorRequest != nil {
		payload, err := json.Marshal(navigatorRequest)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(payload)
	}

	resp, err := c.doRequest(ctx, method, url, params, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := newResponseError(resp, expectedValidStatus); err != nil {
		return nil, err
	}

	if expectedValidStatus == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, nil
	}

	returnedNavigator := &navigator.Navigator{}
	err = json.NewDecoder(resp.Body).Decode(returnedNavigator)
	_, _ = io.Copy(io.Discard, resp.Body)
	return returnedNavigator, err
}
