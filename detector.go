package signalfx

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TODO retrieve events
// TODO retrieve incidents
// TODO validate detector definition

// DetectorAPIURL is the base URL for interacting with detectors.
const DetectorAPIURL = "/v2/detector"

// Detector is a detector.
type Detector struct {
	AuthorizedWriters struct {
		Teams []string `json:"teams"`
		Users []string `json:"users"`
	} `json:"authorizedWriters"`
	CustomProperties string `json:"customProperties"`
	Description      string `json:"description"`
	MaxDelay         int64  `json:"maxDelay"`
	Name             string `json:"name"`
	ProgramText      string `json:"programText"`
	Rules            []struct {
		Description   string `json:"description"`
		DetectLabel   string `json:"detectLabel"`
		Disabled      bool   `json:"disabled"`
		Notifications []struct {
			Channel      string `json:"channel"`
			CredentialID string `json:"credentialId"`
			Type         string `json:"type"`
		} `json:"notifications"`
		ParameterizedBody    string `json:"parameterizedBody"`
		ParameterizedSubject string `json:"parameterizedSubject"`
		RunbookURL           string `json:"runbookUrl"`
		Severity             string `json:"severity"`
	} `json:"rules"`
	Tags                 []string `json:"tags"`
	Teams                []string `json:"teams"`
	Timezone             string   `json:"timezone"`
	VisualizationOptions []struct {
		DisableSampling bool `json:"disableSampling"`
		ShowDataMarkers bool `json:"showDataMarkers"`
		ShowEventLines  bool `json:"showEventLines"`
		Time            struct {
			End   int64  `json:"end"`
			Range int64  `json:"range"`
			Start int64  `json:"start"`
			Type  string `json:"type"`
		} `json:"time"`
	} `json:"visualizationOptions"`
}

// DetectorSearch is the result of a query for Detectors
type DetectorSearch struct {
	Count   int64 `json:"count,omitempty"`
	Results []Detector
}

// CreateDetector creates a detector.
func (c *Client) CreateDetector(detector *Detector) (*Detector, error) {
	payload, err := json.Marshal(detector)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", DetectorAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDetector := &Detector{}

	err = json.NewDecoder(resp.Body).Decode(finalDetector)

	return finalDetector, err
}

// DeleteDetector deletes a detector.
func (c *Client) DeleteDetector(id string) error {
	resp, err := c.doRequest("DELETE", DetectorAPIURL+"/"+id, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Unexpected status code: " + resp.Status)
	}

	return nil
}

// DisableDetector disables a detector.
func (c *Client) DisableDetector(id string) error {
	resp, err := c.doRequest("PUT", DetectorAPIURL+"/"+id+"/disable", nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// EnableDetector enables a detector.
func (c *Client) EnableDetector(id string) error {
	resp, err := c.doRequest("PUT", DetectorAPIURL+"/"+id+"/enable", nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetDetector gets a detector.
func (c *Client) GetDetector(id string) (*Detector, error) {
	resp, err := c.doRequest("GET", DetectorAPIURL+"/"+id, nil, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDetector := &Detector{}

	err = json.NewDecoder(resp.Body).Decode(finalDetector)

	return finalDetector, err
}

// UpdateDetector updates a detector.
func (c *Client) UpdateDetector(id string, detector *Detector) (*Detector, error) {
	payload, err := json.Marshal(detector)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("PUT", DetectorAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDetector := &Detector{}

	err = json.NewDecoder(resp.Body).Decode(finalDetector)

	return finalDetector, err
}

// SearchDetector searches for detectors, given a query string in `name`.
func (c *Client) SearchDetector(limit int, name string, offset int, tags string) (*DetectorSearch, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	resp, err := c.doRequest("GET", DetectorAPIURL, params, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDetectors := &DetectorSearch{}

	err = json.NewDecoder(resp.Body).Decode(finalDetectors)

	return finalDetectors, err
}
