package signalfx

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "DELETE", http.StatusOK, nil, ""))

	err := client.DeleteIntegration("string")
	assert.NoError(t, err, "Unexpected error deleting integration")
}

func TestGetIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "GET", http.StatusOK, nil, "integration/get_success.json"))

	_, err := client.GetIntegration("string")
	assert.NoError(t, err, "Unexpected error getting integration")
}
