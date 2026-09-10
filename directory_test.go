package signalfx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/signalfx/signalfx-go/directory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectoryRequestPath(t *testing.T) {
	valid := []struct {
		name      string
		path      string
		allowRoot bool
		want      string
	}{
		{
			name:      "root",
			allowRoot: true,
			want:      "/v2/directory",
		},
		{
			name: "nested path",
			path: "team/dashboards",
			want: "/v2/directory/team/dashboards",
		},
		{
			name: "spaces and at sign",
			path: "~users/user@example.com/My Dashboards",
			want: "/v2/directory/~users/user%40example.com/My+Dashboards",
		},
		{
			name: "literal plus differs from space",
			path: "A+B/A B",
			want: "/v2/directory/A%2BB/A+B",
		},
		{
			name: "unicode",
			path: "zażółć/東京",
			want: "/v2/directory/za%C5%BC%C3%B3%C5%82%C4%87/%E6%9D%B1%E4%BA%AC",
		},
		{
			name: "server supported punctuation",
			path: "test'path/test(path)/test,path/test!path/test*path/a#b&c=d",
			want: "/v2/directory/test%27path/test%28path%29/test%2Cpath/test%21path/test*path/a%23b%26c%3Dd",
		},
	}

	for _, tc := range valid {
		t.Run(tc.name, func(t *testing.T) {
			got, err := directoryRequestPath(tc.path, tc.allowRoot)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}

	t.Run("server limits", func(t *testing.T) {
		_, err := directoryRequestPath(strings.Repeat("a", 249), false)
		require.NoError(t, err)

		_, err = directoryRequestPath("a/b/c/d/e/f/g/h/i/j", false)
		require.NoError(t, err)
	})

	invalid := []struct {
		name string
		path string
	}{
		{name: "root mutation", path: ""},
		{name: "leading slash", path: "/team"},
		{name: "trailing slash", path: "team/"},
		{name: "empty segment", path: "team//dashboards"},
		{name: "leading whitespace", path: "team/ dashboards"},
		{name: "trailing whitespace", path: "team/dashboards "},
		{name: "dot segment", path: "team/./dashboards"},
		{name: "dot dot segment", path: "team/../dashboards"},
		{name: "percent", path: "team/100%"},
		{name: "filesystem punctuation", path: "team/dashboards?"},
		{name: "control character", path: "team/dash\tboards"},
		{name: "misplaced tilde", path: "team/dash~boards"},
		{name: "extra tilde", path: "team/~dash~boards"},
		{name: "too many segments", path: "a/b/c/d/e/f/g/h/i/j/k"},
		{name: "too long", path: strings.Repeat("a", 250)},
		{name: "invalid UTF-8", path: string([]byte{0xff})},
	}

	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			_, err := directoryRequestPath(tc.path, false)
			assert.Error(t, err)
		})
	}
}

func TestGetDirectoryEntry(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		mux.HandleFunc("/v2/directory", verifyRequest(t, http.MethodGet, true, http.StatusOK, nil, "directory/root_success.json"))

		result, err := client.GetDirectoryEntry(context.Background(), "")
		require.NoError(t, err)
		require.NotNil(t, result.Data)
		assert.Equal(t, directory.EntryType, result.Data.Type)
		assert.Equal(t, "", result.Data.Path)
		assert.Equal(t, []string{"/v2/directory/~organization", "/v2/directory/~users"}, result.Data.Children)
	})

	t.Run("nested path with base URL prefix", func(t *testing.T) {
		mux = http.NewServeMux()
		server = httptest.NewServer(mux)
		defer server.Close()

		client, _ = NewClient(TestToken, APIUrl(server.URL+"/extra/path/"))
		mux.HandleFunc("/", createResponse(t, http.StatusOK, "directory/nested_success.json", func(t *testing.T, request *http.Request) {
			verifyHeaders(t, request, true)
			assert.Equal(t, http.MethodGet, request.Method)
			assert.Equal(t, "/extra/path/v2/directory/~users/user%40example.com/My+Dashboards", request.RequestURI)
		}))

		result, err := client.GetDirectoryEntry(context.Background(), "~users/user@example.com/My Dashboards")
		require.NoError(t, err)
		require.NotNil(t, result.Data)
		assert.Equal(t, "~users/user@example.com/My Dashboards", result.Data.Path)
		assert.Equal(t, []string{"/v2/template/test-dashboard", "/v2/template/test-chart"}, result.Data.Templates)
		require.Len(t, result.Includes, 1)
		assert.True(t, result.Includes[0].Identity)
	})

	t.Run("unoccupied", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		mux.HandleFunc("/v2/directory/test", verifyRequest(t, http.MethodGet, true, http.StatusOK, nil, "directory/unoccupied_success.json"))

		result, err := client.GetDirectoryEntry(context.Background(), "test")
		require.NoError(t, err)
		require.NotNil(t, result.Data)
		assert.Equal(t, "test", result.Data.Path)
		assert.False(t, result.Data.Pinned)
		assert.Empty(t, result.Data.Templates)
		assert.Empty(t, result.Data.Children)
		assert.Empty(t, result.Errors)
	})

	t.Run("unexpected status preserves response details", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		mux.HandleFunc("/v2/directory/test", verifyRequest(t, http.MethodGet, true, http.StatusForbidden, nil, "directory/error.json"))

		result, err := client.GetDirectoryEntry(context.Background(), "test")
		assert.Nil(t, result)
		require.Error(t, err)

		responseError, ok := AsResponseError(err)
		require.True(t, ok)
		assert.Equal(t, http.StatusForbidden, responseError.Code())
		assert.Equal(t, "/v2/directory/test", responseError.Route())
		assert.JSONEq(t, fixture("directory/error.json"), responseError.Details())
	})

	t.Run("malformed response", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		mux.HandleFunc("/v2/directory/test", func(response http.ResponseWriter, _ *http.Request) {
			response.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(response, "{")
		})

		result, err := client.GetDirectoryEntry(context.Background(), "test")
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "decode directory response")
	})

	t.Run("response body is closed", func(t *testing.T) {
		body := &trackingDirectoryBody{Reader: strings.NewReader(fixture("directory/root_success.json"))}
		testClient, err := NewClient(TestToken,
			APIUrl("https://example.test"),
			HTTPClient(&http.Client{Transport: directoryRoundTripFunc(func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
					Header:     make(http.Header),
					Request:    request,
				}, nil
			})}),
		)
		require.NoError(t, err)

		result, err := testClient.GetDirectoryEntry(context.Background(), "")
		require.NoError(t, err)
		require.NotNil(t, result.Data)
		assert.True(t, body.closed)
	})

	t.Run("transport error", func(t *testing.T) {
		transportErr := errors.New("test transport failure")
		testClient, err := NewClient(TestToken,
			APIUrl("https://example.test"),
			HTTPClient(&http.Client{Transport: directoryRoundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, transportErr
			})}),
		)
		require.NoError(t, err)

		result, err := testClient.GetDirectoryEntry(context.Background(), "")
		assert.Nil(t, result)
		assert.ErrorIs(t, err, transportErr)
	})
}

func TestPatchDirectoryEntry(t *testing.T) {
	tests := []struct {
		name       string
		patch      *directory.Patch
		wantBody   string
		statusCode int
	}{
		{
			name: "complete ordered membership and pinned",
			patch: &directory.Patch{
				Templates: directoryStrings("/v2/template/dashboard", "/v2/template/chart"),
				Pinned:    directoryBool(true),
			},
			wantBody:   `{"templates":["/v2/template/dashboard","/v2/template/chart"],"pinned":true}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "explicit empty membership",
			patch:      &directory.Patch{Templates: directoryStrings()},
			wantBody:   `{"templates":[]}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "explicit null membership",
			patch:      &directory.Patch{Templates: directoryNullStrings()},
			wantBody:   `{"templates":null}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "false is not omitted",
			patch:      &directory.Patch{Pinned: directoryBool(false)},
			wantBody:   `{"pinned":false}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "omitted fields",
			patch:      &directory.Patch{},
			wantBody:   `{}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "legacy created status",
			patch:      &directory.Patch{Pinned: directoryBool(true)},
			wantBody:   `{"pinned":true}`,
			statusCode: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc("/v2/directory/test", verifyRequestWithJsonBody(t, http.MethodPatch, true, tc.statusCode, nil, tc.wantBody, "directory/nested_success.json"))

			result, err := client.PatchDirectoryEntry(context.Background(), "test", tc.patch)
			require.NoError(t, err)
			require.NotNil(t, result.Data)
			assert.Equal(t, "~users/user@example.com/My Dashboards", result.Data.Path)
		})
	}

	t.Run("nil patch is rejected without a request", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		requested := false
		mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {
			requested = true
		})

		result, err := client.PatchDirectoryEntry(context.Background(), "test", nil)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "must not be nil")
		assert.False(t, requested)
	})

	t.Run("invalid path is rejected without a request", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		requested := false
		mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {
			requested = true
		})

		result, err := client.PatchDirectoryEntry(context.Background(), "other//user", &directory.Patch{})
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "must not be empty")
		assert.False(t, requested)
	})

	t.Run("unexpected status", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		mux.HandleFunc("/v2/directory/test", verifyRequestWithJsonBody(t, http.MethodPatch, true, http.StatusBadRequest, nil, `{}`, "directory/error.json"))

		result, err := client.PatchDirectoryEntry(context.Background(), "test", &directory.Patch{})
		assert.Nil(t, result)
		require.Error(t, err)

		responseError, ok := AsResponseError(err)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, responseError.Code())
	})
}

func TestDeleteDirectoryEntry(t *testing.T) {
	t.Run("exact path", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		mux.HandleFunc("/", createResponse(t, http.StatusNoContent, "", func(t *testing.T, request *http.Request) {
			verifyHeaders(t, request, true)
			assert.Equal(t, http.MethodDelete, request.Method)
			assert.Equal(t, "/v2/directory/~users/user%40example.com/DASHM-1966+test%2Bfolder", request.RequestURI)
		}))

		err := client.DeleteDirectoryEntry(context.Background(), "~users/user@example.com/DASHM-1966 test+folder")
		assert.NoError(t, err)
	})

	t.Run("root is rejected without a request", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		requested := false
		mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {
			requested = true
		})

		err := client.DeleteDirectoryEntry(context.Background(), "")
		assert.ErrorContains(t, err, "root cannot be modified")
		assert.False(t, requested)
	})

	t.Run("unexpected status", func(t *testing.T) {
		teardown := setup()
		defer teardown()

		mux.HandleFunc("/v2/directory/test", verifyRequest(t, http.MethodDelete, true, http.StatusForbidden, nil, "directory/error.json"))

		err := client.DeleteDirectoryEntry(context.Background(), "test")
		require.Error(t, err)

		responseError, ok := AsResponseError(err)
		require.True(t, ok)
		assert.Equal(t, http.StatusForbidden, responseError.Code())
		assert.JSONEq(t, fixture("directory/error.json"), responseError.Details())
	})
}

func directoryBool(value bool) *bool {
	return &value
}

func directoryStrings(values ...string) *[]string {
	if values == nil {
		values = []string{}
	}
	return &values
}

func directoryNullStrings() *[]string {
	var values []string
	return &values
}

type directoryRoundTripFunc func(*http.Request) (*http.Response, error)

func (function directoryRoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

type trackingDirectoryBody struct {
	io.Reader
	closed bool
}

func (body *trackingDirectoryBody) Close() error {
	body.closed = true
	return nil
}
