package signalfx

import (
	"context"
	"net/http"
	"testing"

	automated_archival "github.com/signalfx/signalfx-go/automated-archival"
	"github.com/stretchr/testify/assert"
)

func TestCreateSettings(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(AutomatedArchivalApiURL+"/settings", verifyRequest(t, http.MethodPost, true, http.StatusOK, nil, "automated_archival/create_settings_success.json"))
	result, err := client.CreateSettings(context.Background(), &automated_archival.AutomatedArchivalSettings{
		Version:        1,
		Enabled:        true,
		LookbackPeriod: "P60D",
		GracePeriod:    "P30D",
	})

	assert.NoError(t, err, "Unexpected error creating Automated Archival settings")
	assert.Equal(t, int64(1), result.Version, "Version doesn't match")
	assert.Equal(t, true, result.Enabled, "Enabled doesn't match")
	assert.Equal(t, "P60D", result.LookbackPeriod, "Lookback period doesn't match")
	assert.Equal(t, "P30D", result.GracePeriod, "Grace period doesn't match")
}
