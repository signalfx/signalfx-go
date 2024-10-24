package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/signalfx/signalfx-go/sessiontoken"
)

// SessionTokenAPIURL is the base URL for interacting with org tokens.
const SessionTokenAPIURL = "/v2/session"

// CreateOrgToken creates a org token.
func (c *Client) CreateSessionToken(ctx context.Context, tokenRequest *sessiontoken.CreateTokenRequest) (*sessiontoken.Token, error) {
	payload, err := json.Marshal(tokenRequest)
	if err != nil {
		return nil, err
	}

	// we need to explicitly pass an empty token (which means it wont get set in the header)
	// the API accepts either no token or a valid token, but not an empty token.
	resp, err := c.doRequestWithToken(ctx, "POST", SessionTokenAPIURL, nil, bytes.NewReader(payload), "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	sessionToken := &sessiontoken.Token{}

	err = json.NewDecoder(resp.Body).Decode(sessionToken)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return sessionToken, err
}

// DeleteOrgToken deletes a token.
func (c *Client) DeleteSessionToken(ctx context.Context, token string) error {
	resp, err := c.doRequestWithToken(ctx, "DELETE", SessionTokenAPIURL, nil, nil, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return nil
}
