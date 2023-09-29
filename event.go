/* Events API
 *
 * https://dev.splunk.com/observability/reference/api/trace_id/latest#endpoint-getlatestsegment
 */

package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/signalfx/signalfx-go/event"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// EventsAPIURL is the base URL for interacting with event.
const EventsAPIURL = "/v2/event/find"

// GetEvents gets an integration as map.
func (c *Client) GetEvents(ctx context.Context, query string, startTime time.Time, endTime time.Time, limit uint32, offset uint32) (event.Events, error) {
	out := make(event.Events, 0)

	err := c.getEvents(ctx, &out, query, startTime, endTime, limit, offset)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) getEvents(ctx context.Context, out interface{}, query string, startTime time.Time, endTime time.Time, limit uint32, offset uint32) error {
	return c.doEventsRequest(ctx, EventsAPIURL, "GET", c.encodeParams(query, startTime, endTime, limit, offset), http.StatusOK, nil, out)
}

func (c *Client) encodeParams(query string, startTime time.Time, endTime time.Time, limit uint32, offset uint32) map[string][]string {
	return map[string][]string{
		"query":      {query},
		"start_time": {strconv.FormatInt(startTime.UnixMilli(), 10)},
		"end_time":   {strconv.FormatInt(endTime.UnixMilli(), 10)},
		"limit":      {strconv.FormatUint(uint64(limit), 10)},
		"offset":     {strconv.FormatUint(uint64(offset), 10)},
	}
}

func (c *Client) doEventsRequest(ctx context.Context, url2 string, method string, params map[string][]string, status int, in interface{}, out interface{}) error {
	var body io.Reader

	if in != nil {
		payload, err := json.Marshal(in)
		if err != nil {
			return err
		}

		body = bytes.NewReader(payload)
	}

	// Typecast from map[strigng]string to url.Values
	paramsToPass := url.Values(params)

	resp, err := c.doRequest(ctx, method, url2, paramsToPass, body)
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
