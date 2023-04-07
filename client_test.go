package signalfx

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/signalfx/signalfx-go/chart"
	"github.com/signalfx/signalfx-go/signalflow"
	"github.com/stretchr/testify/assert"
)

const TestToken = "abc123"

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func fixture(path string) string {
	b, err := os.ReadFile("testdata/fixtures/" + path)
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
				assert.Contains(t, incomingParams, k, "Request is missing expected query parameter %q", k)
				assert.Equal(t, params.Get(k), incomingParams.Get(k), "Params do match for parameter '"+k+"': '"+incomingParams.Get(k)+"'")
			}
			for k := range incomingParams {
				assert.Contains(t, params, k, "Request contains unexpected query parameter %q", k)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		// Allow empty bodies
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

func Test_Close_DoesNotLeakGoroutines(t *testing.T) {
	// it can happen that there are still some dangling goroutines from other tests.
	// This test compares goroutines before and after creating the client - let's wait for other goroutines to close first.
	_ = waitForAllGoroutinesToDie(2, time.Second*2)
	goroutinesBefore := runtime.NumGoroutine()
	signalflow.ReconnectDelay = time.Millisecond * 200

	clientsCount := 10
	clients := make([]*signalflow.Client, clientsCount)
	for i := 0; i < clientsCount; i++ {
		clients[i] = newClient(t)
		go executeRequest(t, clients[i])
	}

	time.Sleep(time.Second) // I don't see a better way to wait for client to hang on sendMessage..

	for i := 0; i < clientsCount; i++ {
		clients[i].Close()
	}

	goroutinesAfter := waitForAllGoroutinesToDie(goroutinesBefore, time.Second*10)

	if goroutinesAfter != goroutinesBefore {
		t.Logf("goroutinesAfter[%d] != goroutinesBefore[%d]\n", goroutinesAfter, goroutinesBefore)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGQUIT) // send SIGQUIT to dump goroutines for easier debug
		t.Errorf("goroutinesAfter[%d] != goroutinesBefore[%d]\n", goroutinesAfter, goroutinesBefore)
	}
}

func newClient(t *testing.T) *signalflow.Client {
	t.Helper()

	client, err := signalflow.NewClient(
		signalflow.StreamURLForRealm("invalidRealm"),
		signalflow.StreamURL("wss://stream.invalidRealm.fakeSignalFX.com"),
		signalflow.AccessToken("invalidToken"),
		signalflow.WriteTimeout(time.Millisecond*100),
		signalflow.MetadataTimeout(time.Millisecond*100),
	)
	if err != nil {
		t.Fatalf("error when creating client: %s", err.Error())
	}

	return client
}

func executeRequest(t *testing.T, c *signalflow.Client) {
	t.Helper()

	_, err := c.Execute(&signalflow.ExecuteRequest{
		Program:      "invalid | program",
		Start:        time.Now(),
		ResolutionMs: (time.Second * 15).Milliseconds(),
		Immediate:    false,
		MaxDelay:     time.Minute,
	})

	if err != nil {
		t.Fatalf("error when executing request: %s", err.Error())
	}
}

func waitForAllGoroutinesToDie(goroutinesBefore int, maxWaitTime time.Duration) int {
	i := 0
	sleepTime := time.Millisecond * 200
	goroutinesAfter := runtime.NumGoroutine()
	for goroutinesAfter != goroutinesBefore {
		time.Sleep(sleepTime)
		if sleepTime*time.Duration(i) == maxWaitTime {
			break
		}
		i++
		goroutinesAfter = runtime.NumGoroutine()
	}
	return goroutinesAfter
}
