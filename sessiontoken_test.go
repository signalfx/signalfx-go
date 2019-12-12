package signalfx

import (
	"net/http"
	"testing"

	"github.com/signalfx/signalfx-go/sessiontoken"

	"github.com/stretchr/testify/assert"
)

func TestCreateSessionToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/session", verifyRequest(t, "POST", http.StatusOK, nil, "sessiontoken/create_success.json"))

	result, err := client.CreateSessionToken(&sessiontoken.CreateTokenRequest{
		Email: "testemail@test.com",
		Password: "testpassword",
	})
	assert.NoError(t, err, "Unexpected error creating orgtoken")
	assert.Equal(t, "testemail@test.com", result.Email, "Email does not match")
	assert.Equal(t, "mytokenvalue", result.AccessToken, "Access token does not match")
}

func TestCreateBadSessionToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/session", verifyRequest(t, "POST", http.StatusBadRequest, nil, ""))

	result, err := client.CreateSessionToken(&sessiontoken.CreateTokenRequest{
		Email: "email",
	})
	assert.Error(t, err, "Should have gotten an error from a bad create")
	assert.Nil(t, result, "Should have a null token on bad create")
}

func TestDeleteSessionToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/session", verifyRequest(t, "DELETE", http.StatusNoContent, nil, ""))

	err := client.DeleteSessionToken(TestToken)
	assert.NoError(t, err, "Unexpected error deleting token")
}

func TestDeleteMissingSessionToken(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/session", verifyRequest(t, "DELETE", http.StatusNotFound, nil, ""))

	err := client.DeleteSessionToken(TestToken)
	assert.Error(t, err, "Should have gotten an error from a missing delete")
}

