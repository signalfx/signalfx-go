package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/signalfx/signalfx-go/graphql"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GraphQLAPIURL is the base URL for interacting with the graphql API which the frontend uses.
// We need to use the frontend api because signalfx has a bunch of things missing from the normal apis.
const GraphQLAPIURL = "/v2/apm/graphql"

// Get incident with the given id
func (c *Client) GraphQLRequest(ctx context.Context, request graphql.Request, response any) error {
	params := url.Values{
		"op": {request.OperationName},
	}

	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(request)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(ctx, "POST", GraphQLAPIURL, params, body)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Bad status %d: %s", resp.StatusCode, message)
	}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(&response)

	return err
}
