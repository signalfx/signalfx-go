package signalfx

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIncident(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incident/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "incident/get_incident.json"))

	result, err := client.GetIncident(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting incident")
	assert.Equal(t, result.IncidentId, "string", "Name does not match")
	assert.Equal(t, true, result.Active, "Active field does not match")
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
