package signalfx

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

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
func verifyRequest(t *testing.T, method string, status int, params url.Values, resultPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, ok := r.Header[AuthHeaderKey]; ok {
			assert.Equal(t, []string{TestToken}, val, "Incorrect auth token in headers")
		} else {
			assert.Fail(t, "Failed to find auth token in headers")
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

func TestCreateChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart", verifyRequest(t, "POST", http.StatusOK, nil, "chart/create_success.json"))

	result, err := client.CreateChart(&Chart{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error creating chart")
	assert.Equal(t, "string", result.Name, "Name does not match")
}

func TestDeleteChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "DELETE", http.StatusOK, nil, ""))

	err := client.DeleteChart("string")
	assert.NoError(t, err, "Unexpected error getting chart")
}

func TestGetChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "GET", http.StatusOK, nil, "chart/get_success.json"))

	result, err := client.GetChart("string")
	assert.NoError(t, err, "Unexpected error getting chart")
	assert.Equal(t, result.Name, "string", "Name does not match")
}

func TestSearchChart(t *testing.T) {
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

	mux.HandleFunc("/v2/chart", verifyRequest(t, "GET", http.StatusOK, params, "chart/get_success.json"))

	results, err := client.SearchChart(limit, name, offset, tags)
	assert.NoError(t, err, "Unexpected error search chart")
	assert.Equal(t, int64(0), results.Count, "Incorrect number of results")
}

func TestUpdateChart(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/chart/string", verifyRequest(t, "PUT", http.StatusOK, nil, "chart/update_success.json"))

	result, err := client.UpdateChart("string", &Chart{
		Name: "string",
	})
	assert.NoError(t, err, "Unexpected error updating chart")
	assert.Equal(t, "string", result.Name, "Name does not match")
}
