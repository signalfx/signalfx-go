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

// TODO: Use HTTPSuccess from testify?
func verifyRequest(t *testing.T, method string, expectToken bool, status int, params url.Values, resultPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
