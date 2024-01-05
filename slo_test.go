package signalfx

import (
	"context"
	"fmt"
	"github.com/signalfx/signalfx-go/detector"
	"github.com/signalfx/signalfx-go/slo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const id = "12345"

func TestGetSloWithRollingWindowTarget(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/v2/slo/%s", id), verifyRequest(t, http.MethodGet, true, http.StatusOK, nil, "slo/get_success_rolling_window_target.json"))

	result, err := client.GetSlo(context.Background(), id)
	assert.NoError(t, err, "Unexpected error getting SLO")
	assert.Equal(t, id, result.Id, "Id does not match")
	assert.Equal(t, "SLO testing", result.Name, "Name does not match")
	assert.Equal(t, 99.99, result.Targets[0].Slo, "SloObject target does not match")
	assert.Equal(t, "7d", result.Targets[0].CompliancePeriod, "SloObject compliance period does not match")
	assert.Equal(t, detector.MAJOR, result.Targets[0].SloAlertRules[0].BreachSloAlertRule.Rules[0].Severity, "SloObject rule severity does not match")
	assert.Equal(t, "5m", result.Targets[0].SloAlertRules[0].BreachSloAlertRule.Rules[0].Parameters.FireLasting, "SloObject rule fire lasting does not match")
}

func TestGetSloWithCalendarWindowTarget(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/v2/slo/%s", id), verifyRequest(t, "GET", true, http.StatusOK, nil, "slo/get_success_calendar_window_target.json"))

	result, err := client.GetSlo(context.Background(), id)
	assert.NoError(t, err, "Unexpected error getting SLO")
	assert.Equal(t, id, result.Id, "Id does not match")
	assert.Equal(t, "SLO testing", result.Name, "Name does not match")
	assert.Equal(t, 95.0, result.Targets[0].Slo, "SloObject target does not match")
	assert.Equal(t, "month", result.Targets[0].CycleType, "SloObject cycle type does not match")
	assert.Equal(t, detector.CRITICAL, result.Targets[0].SloAlertRules[0].ErrorBudgetLeftSloAlertRule.Rules[0].Severity, "SloObject rule severity does not match")
	assert.Equal(t, "10m", result.Targets[0].SloAlertRules[0].ErrorBudgetLeftSloAlertRule.Rules[0].Parameters.FireLasting, "SloObject rule fire lasting does not match")
}

func TestGetMissingSlo(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/v2/slo/%s", id), verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetSlo(context.Background(), "string")
	assert.Error(t, err, "Expected error getting missing SLO")
	assert.Nil(t, result, "Expected nil result getting SLO")
}

func TestCreateSlo(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/slo", verifyRequest(t, "POST", true, http.StatusOK, nil, "slo/get_success_rolling_window_target.json"))

	result, err := client.CreateSlo(context.Background(), &slo.SloObject{
		BaseSlo: slo.BaseSlo{
			Type: slo.RequestBased,
		},
	})
	assert.NoError(t, err, "Unexpected error creating SLO")
	assert.Equal(t, id, result.Id, "Id does not match")
}

func TestCreateBadSlo(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/slo", verifyRequest(t, "POST", true, http.StatusBadRequest, nil, ""))

	result, err := client.CreateSlo(context.Background(), &slo.SloObject{
		BaseSlo: slo.BaseSlo{
			Type: slo.RequestBased,
		},
	})

	assert.Error(t, err, "Should have gotten an error from a bad create")
	assert.Nil(t, result, "Should have a null SLO on bad create")
}

func TestDeleteSlo(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/v2/slo/%s", id), verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteSlo(context.Background(), id)
	assert.NoError(t, err, "Unexpected error deleting SLO")
}

func TestDeleteMissingSlo(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/v2/slo/%s", id), verifyRequest(t, "DELETE", true, http.StatusNotFound, nil, ""))

	err := client.DeleteSlo(context.Background(), id)
	assert.Error(t, err, "Should have gotten an error from a missing SLO")
}

func TestUpdateSlo(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/v2/slo/%s", id), verifyRequest(t, "PUT", true, http.StatusOK, nil, "slo/get_success_rolling_window_target.json"))

	result, err := client.UpdateSlo(context.Background(), id, &slo.SloObject{
		BaseSlo: slo.BaseSlo{
			Type: slo.RequestBased,
		},
	})
	assert.NoError(t, err, "Unexpected error updating SLO")
	assert.Equal(t, id, result.Id, "Id does not match")
}

func TestUpdateMissingSlo(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/v2/slo/%s", id), verifyRequest(t, "PUT", true, http.StatusNotFound, nil, ""))

	result, err := client.UpdateSlo(context.Background(), id, &slo.SloObject{
		BaseSlo: slo.BaseSlo{
			Type: slo.RequestBased,
		},
	})
	assert.Error(t, err, "Should have gotten an error from an update on a missing SLO")
	assert.Nil(t, result, "Should have gotten a nil result from an update on a missing SLO")
}
