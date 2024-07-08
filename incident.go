package signalfx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/detector"
)

// IncidentAPIURL is the base URL for interacting with alert muting rules.
const IncidentAPIURL = "/v2/incident"

// Get incident with the given id
func (c *Client) GetIncident(ctx context.Context, id string) (*detector.Incident, error) {
	resp, err := c.doRequest(ctx, "GET", IncidentAPIURL+"/"+id, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Bad status %d: %s", resp.StatusCode, message)
	}

	finalIncident := &detector.Incident{}

	err = json.NewDecoder(resp.Body).Decode(finalIncident)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalIncident, err
}

// Get all incidents
func (c *Client) GetIncidents(ctx context.Context, includeResolved bool, limit int, query string, offset int) ([]*detector.Incident, error) {
	params := url.Values{}
	params.Add("includeResolved", strconv.FormatBool(includeResolved))
	params.Add("limit", strconv.Itoa(limit))
	params.Add("query", query)
	params.Add("offset", strconv.Itoa(offset))
	resp, err := c.doRequest(ctx, "GET", IncidentAPIURL, params, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(string(body))
	}
	if err != nil {
		return nil, err
	}

	var incidents []*detector.Incident
	err = json.NewDecoder(resp.Body).Decode(&incidents)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return incidents, err
}
