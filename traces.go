/* Traces API
 *
 * https://dev.splunk.com/observability/reference/api/trace_id/latest#endpoint-getlatestsegment
 */

package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/signalfx/signalfx-go/traces"
	"io"
	"io/ioutil"
	"net/http"
)

// TracesAPIURL is the base URL for interacting with traces.
const TracesAPIURL = "/v2/apm/trace"

// GetTrace gets an integration as map.
func (c *Client) GetTrace(ctx context.Context, id string) (traces.Trace, error) {
	out := make(traces.Trace, 0)

	err := c.getTrace(ctx, id, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) getTrace(ctx context.Context, id string, out interface{}) error {
	return c.doTraceRequest(ctx, TracesAPIURL+"/"+id+"/latest", "GET", http.StatusOK, nil, out)
}
func (c *Client) doTraceRequest(ctx context.Context, url string, method string, status int, in interface{}, out interface{}) error {
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
