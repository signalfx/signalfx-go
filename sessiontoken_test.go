package signalfx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/signalfx/signalfx-go/sessiontoken"

	"github.com/stretchr/testify/assert"
)

func verifyNoTokenRequest(t *testing.T, method string, status int, params url.Values, resultPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header[AuthHeaderKey]; ok {
			assert.Fail(t, "Didn't expect to find token in headers")
		}

		if val, ok := r.Header["Content-Type"]; ok {
			assert.Equal(t, []string{"application/json"}, val, "Incorrect content-type in headers")
		} else {
			assert.Fail(t, "Failed to find content type in headers")
		}

		assert.Equal(t, method, r.Method, "Incorrect HTTP method")

		if params != nil {
			incomingParams := r.URL.Query()
			for k := range params {
				assert.Equal(t, params.Get(k), incomingParams.Get(k), "Params do match for parameter '"+k+"': '"+incomingParams.Get(k)+"'")
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		// Allow empty bodies
		if resultPath != "" {
			fmt.Fprintf(w, fixture(resultPath))
		}
	}
}

func TestCreateSessionToken(t *testing.T) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewClient("", APIUrl(server.URL))
	defer server.Close()

	mux.HandleFunc("/v2/session", verifyNoTokenRequest(t, "POST", http.StatusOK, nil, "sessiontoken/create_success.json"))

	result, err := client.CreateSessionToken(&sessiontoken.CreateTokenRequest{
		Email: "testemail@test.com",
		Password: "testpassword",
	})
	assert.NoError(t, err, "Unexpected error creating orgtoken")
	assert.Equal(t, "testemail@test.com", result.Email, "Email does not match")
	assert.Equal(t, "mytokenvalue", result.AccessToken, "Access token does not match")
}

func TestCreateBadCredentials(t *testing.T) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient("", APIUrl(server.URL))
	defer server.Close()

	mux.HandleFunc("/v2/session", verifyNoTokenRequest(t, "POST", http.StatusBadRequest, nil, ""))

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

