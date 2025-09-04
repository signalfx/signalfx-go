package signalfx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/signalfx/signalfx-go/apm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPMTopologyListService(t *testing.T) {
	t.Parallel()

	start := time.Unix(100, 0)
	end := start.Add(7 * 24 * 60 * time.Minute)

	formatted := fmt.Sprintf("%d/%d", start.UnixMilli(), end.UnixMilli())

	for _, tc := range []struct {
		name    string
		body    string
		results string
		status  int
		req     *apm.RetrieveServiceTopologyRequest
		errVal  string
	}{
		{
			name:    "successfully retrieves topology",
			status:  http.StatusOK,
			body:    `{"timeRange":"` + formatted + `"}`,
			results: "create_success.json",
			req:     &apm.RetrieveServiceTopologyRequest{TimeRange: formatted},
			errVal:  "",
		},
		{
			name:    "returns error on bad request",
			status:  http.StatusBadRequest,
			body:    `{"timeRange":"` + formatted + `"}`,
			results: "failed_error.json",
			req:     &apm.RetrieveServiceTopologyRequest{TimeRange: formatted},
			errVal:  `route "/v2/apm/topology" had issues with status code 400`,
		},
		{
			name:    "returns error on invalid json response",
			status:  http.StatusBadRequest,
			body:    `{"timeRange":""}`,
			results: "failed_error.json",
			req:     &apm.RetrieveServiceTopologyRequest{},
			errVal:  `route "/v2/apm/topology" had issues with status code 400`,
		},
		{
			name:    "returns error on nil request",
			status:  http.StatusInternalServerError,
			body:    `null`,
			results: "failed_error.json",
			req:     nil,
			errVal:  `route "/v2/apm/topology" had issues with status code 500`,
		},
	} {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			tmux := http.NewServeMux()
			tmux.HandleFunc(APMTopologyURL, func(w http.ResponseWriter, r *http.Request) {
				actual := bytes.NewBuffer(nil)
				_, _ = io.Copy(actual, r.Body)
				_ = r.Body.Close()

				assert.JSONEq(t, tc.body, actual.String(), "Request body must match")

				f, err := os.Open(path.Join("testdata/fixtures/apm", tc.results))
				if err != nil {
					panic(err)
				}
				defer f.Close()

				w.WriteHeader(tc.status)
				_, _ = io.Copy(w, f)
			})
			s := httptest.NewServer(tmux)
			t.Cleanup(s.Close)

			c, err := NewClient("mocked", APIUrl(s.URL))
			require.NoError(t, err, "Must not error when creating client")

			v, err := c.ListTopology(context.Background(), tc.req)
			if tc.errVal != "" {
				assert.EqualError(t, err, tc.errVal, "Must match the expected value")
				assert.Nil(t, v, "Result must be nil")
			} else {
				assert.NoError(t, err, "Must not error")
				assert.NotNil(t, v, "Result must not be nil")
			}
		})
	}
}
