package signalfx

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteIntegration(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error deleting integration")
}

func TestDeleteMissingIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "DELETE", true, http.StatusNotFound, nil, ""))

	err := client.DeleteIntegration(context.Background(), "string")
	assert.Error(t, err, "Should get error error deleting missing integration")
}

func TestGetIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "GET", true, http.StatusOK, nil, "integration/get_success.json"))

	result, err := client.GetIntegration(context.Background(), "string")
	assert.NoError(t, err, "Unexpected error getting integration")
	id := result["id"].(string)
	assert.Equal(t, id, "string", "Missing ID")
}

func TestGetMissingIntegration(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/integration/string", verifyRequest(t, "GET", true, http.StatusNotFound, nil, ""))

	result, err := client.GetIntegration(context.Background(), "string")
	assert.Error(t, err, "Should get an error getting missing integration")
	assert.Nil(t, result, "Should get a nil result from a missing integration")
}
