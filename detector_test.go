package signalfx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/signalfx/signalfx-go/detector"
	"github.com/stretchr/testify/assert"
)

func TestCreateDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector", verifyRequest(t, "POST", true, http.StatusOK, nil, "detector/create_success.json"))

	result, err := client.CreateDetector(context.Background(), &detector.CreateUpdateDetectorRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error creating detector")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestCreateBadDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	result, err := client.CreateDetector(context.Background(), &detector.CreateUpdateDetectorRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should have gotten an error from a bad create")
	assert.Nil(t, result, "Should have a null detector on bad create")
}

func TestDeleteDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteDetector(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error deleting detector")
}

func TestDeleteMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector", verifyRequest(t, "POST", true, http.StatusNotFound, nil, ""))

	err := client.DeleteDetector(context.Background(), "example")
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}

func TestDisableDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/disable", verifyRequest(t, "PUT", true, http.StatusNoContent, nil, ""))

	err := client.DisableDetector(context.Background(), "string", []string{"example"})
	assert.NoError(t, err, "Unexpected error disabling detector")
}

func TestDisableMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/disable", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	err := client.DisableDetector(context.Background(), "string", []string{"example"})
	assert.Error(t, err, "Should have gotten an error from a missing disable")
}

func TestEnableDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/enable", verifyRequest(t, "PUT", true, http.StatusNoContent, nil, ""))

	err := client.EnableDetector(context.Background(), "string", []string{"example"})
	assert.NoError(t, err, "Unexpected error disabling detector")
}

func TestEnableMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string/enable", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	err := client.EnableDetector(context.Background(), "string", []string{"example"})
	assert.Error(t, err, "Should have gotten an error from a missing enable")
}

func TestGetDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "detector/get_success.json"))

	result, err := client.GetDetector(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting detector")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestGetMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetDetector(context.Background(), "string")
	assert.Error(t, err, "Should have gotten an error from a missing detector")
	assert.Nil(t, result, "Should have gotten a nil result from a missing detector")
}

func TestSearchDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	name := "foo"
	offset := 2
	tags := "bar"
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	mux.HandleFunc("/v2/detector", verifyRequest(t, "GET", true, http.StatusOK, params, "detector/search_success.json"))

	results, err := client.SearchDetectors(context.Background(), limit, name, offset, tags)
	assert.NoError(t, err, "Unexpected error search detector")
	assert.Equal(t, int32(1), results.Count, "Incorrect number of results")
}

func TestSearchDetectorBad(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	name := "foo"
	offset := 2
	tags := "bar"
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	mux.HandleFunc("/v2/detector", verifyRequest(t, "GET", true, http.StatusBadRequest, params, ""))

	_, err := client.SearchDetectors(context.Background(), limit, name, offset, tags)
	assert.Error(t, err, "Unexpected error search detector")
}

func TestUpdateDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "PUT", true, http.StatusOK, nil, "detector/update_success.json"))

	result, err := client.UpdateDetector(context.Background(), "string", &detector.CreateUpdateDetectorRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error updating detector")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestUpdateMissingDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/string", verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	result, err := client.UpdateDetector(context.Background(), "string", &detector.CreateUpdateDetectorRequest{
		Name: "string",
	})
	assert.Error(t, err, "Should have gotten an error from an update on a missing detector")
	assert.Nil(t, result, "Should have gotten a nil result from an update on a missing detector")
}

func TestGetDetectorEvents(t *testing.T) {
	teardown := setup()
	defer teardown()

	from := 1557534630000
	to := 1557534640000
	offset := 12
	limit := 2

	params := url.Values{}
	params.Add("from", strconv.Itoa(from))
	params.Add("to", strconv.Itoa(to))
	params.Add("offset", strconv.Itoa(offset))
	params.Add("limit", strconv.Itoa(limit))

	mux.HandleFunc("/v2/detector/string/events", verifyRequest(t, "GET", true, http.StatusOK, params, "detector/get_events.json"))

	result, err := client.GetDetectorEvents(context.Background(), "string", from, to, offset, limit)
	assert.NoError(t, err, "Unexpected error getting detector")
	assert.Equal(t, result[0].AnomalyState, "ANOMALOUS", "AnomalyState does not match")
}

func TestGetDetectorIncidents(t *testing.T) {
	teardown := setup()
	defer teardown()

	offset := 12
	limit := 2

	params := url.Values{}
	params.Add("offset", strconv.Itoa(offset))
	params.Add("limit", strconv.Itoa(limit))

	mux.HandleFunc("/v2/detector/string/incidents", verifyRequest(t, "GET", true, http.StatusOK, params, "detector/get_incidents.json"))

	result, err := client.GetDetectorIncidents(context.Background(), "string", offset, limit)
	assert.NoError(t, err, "Unexpected error getting detector")
	assert.Equal(t, result[0].AnomalyState, "ANOMALOUS", "AnomalyState does not match")
}

func TestValidateDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/validate", verifyRequest(t, "POST", true, http.StatusNoContent, nil, ""))

	err := client.ValidateDetector(context.Background(), &detector.ValidateDetectorRequestModel{
		Name:        "string",
		ProgramText: "signal = data('cpu.utilization').mean(by=['sf_metric', 'sfx_realm']).publish(label='A'); detect(when(A > threshold(10), lasting='2m'), auto_resolve_after='3d').publish('Test detector validation')",
		Rules: []*detector.Rule{
			&detector.Rule{
				Description: "Maximum > 10 for 2m",
				DetectLabel: "Test detector validation",
				Severity:    detector.INFO,
			},
		},
	})
	assert.NoError(t, err, "Unexpected error validating detector programText")
}

func TestValidateBadDetector(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/detector/validate", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	err := client.ValidateDetector(context.Background(), &detector.ValidateDetectorRequestModel{
		Name:        "string",
		ProgramText: "signal = ('cpu.utilization').mean(by=['sf_metric', 'sfx_realm']).publish(label='A'); detect(when(A > threshold(10), lasting='2m'), auto_resolve_after='3d').publish('Test detector validation')",
		Rules: []*detector.Rule{
			&detector.Rule{
				Description: "Maximum > 10 for 2m",
				DetectLabel: "Test detector validation",
				Severity:    detector.INFO,
			},
		},
	})
	assert.Error(t, err, "Should have gotten an error from invalid detector")
}
