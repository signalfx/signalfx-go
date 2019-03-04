package signalfx

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// DefaultAPIURL is the default URL for making API requests
const DefaultAPIURL = "https://api.signalfx.com"

// AuthHeaderKey is the HTTP header used to pass along the auth token
// Note that while HTTP headers are case insensitive this header is case
// sensitive on the tests for convenience.
const AuthHeaderKey = "X-Sf-Token"

// Client is a SignalFx API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

// TODO override, reuse?

// NewClient creates a new SignalFx client using the specified token.
func NewClient(token string, apiURL string) (*Client, error) {
	configuredURL := apiURL
	if configuredURL == "" {
		configuredURL = DefaultAPIURL
	}

	return &Client{
		baseURL: configuredURL,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		authToken: token,
	}, nil
}

func (c *Client) doRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error) {
	destURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	destURL.Path = path

	if params != nil {
		destURL.RawQuery = params.Encode()
	}
	req, err := http.NewRequest(method, destURL.String(), body)
	req.Header.Set(AuthHeaderKey, c.authToken)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}
