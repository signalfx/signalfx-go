package signalfx

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/signalfx/signalfx-go/chart"
	"github.com/stretchr/testify/assert"
)

const TestToken = "abc123"

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewClient(TestToken, APIUrl(server.URL))

	return func() {
		server.Close()
	}
}

// TODO: rename this to verifyRequestAndCreateResponse
func verifyRequest(t *testing.T, method string, expectToken bool, status int, params url.Values, resultPath string) func(w http.ResponseWriter, r *http.Request) {
	return createResponse(t, status, resultPath, func(t *testing.T, r *http.Request) {
		verifyHeaders(t, r, expectToken)
		verifyParams(t, r, params)

		assert.Equal(t, method, r.Method, "Incorrect HTTP method")
	})
}

func verifyRequestWithJsonBody(t *testing.T, method string, expectToken bool, status int, params url.Values, jsonBody string, resultPath string) func(w http.ResponseWriter, r *http.Request) {
	return createResponse(t, status, resultPath, func(t *testing.T, r *http.Request) {
		verifyHeaders(t, r, expectToken)
		verifyParams(t, r, params)
		verifyJsonBody(t, r, jsonBody)

		assert.Equal(t, method, r.Method, "Incorrect HTTP method")
	})
}

func verifyHeaders(t *testing.T, r *http.Request, expectToken bool) {
	if val, ok := r.Header[AuthHeaderKey]; ok {
		assert.Equal(t, []string{TestToken}, val, "Incorrect auth token in headers")
	} else {
		if expectToken {
			assert.Fail(t, "Failed to find auth token in headers")
		}
	}

	if val, ok := r.Header["Content-Type"]; ok {
		assert.Equal(t, []string{"application/json"}, val, "Incorrect content-type in headers")
	} else {
		assert.Fail(t, "Failed to find content type in headers")
	}
}

func verifyParams(t *testing.T, r *http.Request, params url.Values) {
	if params != nil {
		incomingParams := r.URL.Query()
		for k := range params {
			assert.Contains(t, incomingParams, k, "Request is missing expected query parameter %q", k)
			assert.Equal(t, params.Get(k), incomingParams.Get(k), "Params do match for parameter '"+k+"': '"+incomingParams.Get(k)+"'")
		}
		for k := range incomingParams {
			assert.Contains(t, params, k, "Request contains unexpected query parameter %q", k)
		}
	}
}

func verifyJsonBody(t *testing.T, r *http.Request, jsonBody string) {
	actualBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		assert.Fail(t, "Error reading request body: %v", err)
	}

	actualJsonBody := string(actualBody)
	assert.JSONEq(t, jsonBody, actualJsonBody, "Expected body: %s, got: %s", jsonBody, actualJsonBody)
}

func createResponse(t *testing.T, status int, resultPath string, requestValidator func(t *testing.T, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if requestValidator != nil {
			requestValidator(t, r)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if resultPath != "" {
			_, _ = fmt.Fprintf(w, fixture(resultPath))
		}
	}
}

func TestPathHandling(t *testing.T) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	defer server.Close()

	client, _ = NewClient(TestToken, APIUrl(server.URL+"/extra/path"))

	mux.HandleFunc("/extra/path/v2/chart", verifyRequest(t, "POST", true, http.StatusOK, nil, "chart/create_success.json"))

	result, err := client.CreateChart(context.Background(), &chart.CreateUpdateChartRequest{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error creating chart")
	assert.Equal(t, "string", result.Name, "Name does not match")
}
