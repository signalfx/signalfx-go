package signalfx

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/signalfx/signalfx-go/detector"
	"github.com/stretchr/testify/assert"
)

func TestGetIncident(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incident/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "incident/get_incident.json"))

	result, err := client.GetIncident(context.Background(), "string")
	expected := detector.Incident{
		Active:       true,
		AnomalyState: "ANOMALOUS",
		DetectLabel:  "string",
		DetectorId:   "string",
		Events: []*detector.Event{
			{
				AnomalyState:     "ANOMALOUS",
				DetectLabel:      "string",
				DetectorId:       "string",
				DetectorName:     "string",
				EventAnnotations: &map[string]interface{}{},
				Id:               "string",
				IncidentId:       "string",
				Inputs: &map[string]interface{}{
					"A": float64(5),
					"B": float64(6),
				},
				Severity:  "Critical",
				Timestamp: 1557484230000,
			},
		},
		IncidentId:                "string",
		IsMuted:                   true,
		Severity:                  "Critical",
		TriggeredNotificationSent: true,
		TriggeredWhileMuted:       true,
	}
	assert.NoError(t, err, "Unexpected error getting incident")
	assert.Equal(t, expected, *result, "Extected incident does not match the actual incident reponse")
}

func TestGetIncidents(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incident", verifyRequest(t, "GET", true, http.StatusOK, nil, "incident/get_incidents.json"))

	result, err := client.GetIncidents(context.Background(), false, 10, "", 0)
	assert.NoError(t, err, "Unexpected error getting all incidents")
	assert.Equal(t, len(result), 2, "Incorrect number of incidents returned")
	assert.Equal(t, result[0].IncidentId, "string", "Name does not match")
	assert.Equal(t, result[0].Active, true, "Active field does not match")
	assert.Equal(t, result[1].IncidentId, "string1", "Name does not match")
}

func TestGetIncidentsClosesResponseBody(t *testing.T) {
	body := &closeTrackingBody{Reader: strings.NewReader("[]")}
	client, err := NewClient(TestToken, HTTPClient(&http.Client{Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: body, Request: request}, nil
	})}))
	assert.NoError(t, err)

	_, err = client.GetIncidents(context.Background(), false, 10, "", 0)
	assert.NoError(t, err)
	assert.True(t, body.closed)
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

type closeTrackingBody struct {
	io.Reader
	closed bool
}

func (b *closeTrackingBody) Close() error {
	b.closed = true
	return nil
}
