package signalfx

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestError(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		re   *ResponseError
		err  string
	}{
		{
			name: "simple request error",
			re: &ResponseError{
				code:    404,
				route:   "/internal/heathz",
				details: "failed due to overworked",
			},
			err: "route \"/internal/heathz\" had issues with status code 404: failed due to overworked",
		},
		{
			name: "missed details",
			re: &ResponseError{
				code:  404,
				route: "/internal/heathz",
			},
			err: "route \"/internal/heathz\" had issues with status code 404",
		},
		{
			name: "missed code",
			re: &ResponseError{
				route: "/internal/heathz",
			},
			err: "route \"/internal/heathz\" had issues with status code 0",
		},
		{
			name: "missed all info",
			re:   &ResponseError{},
			err:  "route \"\" had issues with status code 0",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.EqualError(t, tc.re, tc.err, "Must match the expected error string")
		})
	}
}

func TestNewRequestError(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		resp   *http.Response
		expect error
	}{
		{
			name:   "no response",
			resp:   nil,
			expect: nil,
		},
		{
			name: "successful response",
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Request: &http.Request{
					URL: &url.URL{Host: "localhost", Path: "/internal/heathz"},
				},
				Body: io.NopCloser(strings.NewReader("successs")),
			},
			expect: nil,
		},
		{
			name: "failed response",
			resp: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					URL: &url.URL{Host: "localhost", Path: "/internal/heathz"},
				},
				Body: io.NopCloser(strings.NewReader("issue trying to connect")),
			},
			expect: &ResponseError{
				code:    http.StatusInternalServerError,
				route:   "/internal/heathz",
				details: "issue trying to connect",
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expect, newResponseError(tc.resp, http.StatusOK), "Must match the exected value")
		})
	}
}

func TestIsRequestError(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "no error provided",
			err:      nil,
			expected: false,
		},
		{
			name:     "not a request error",
			err:      errors.New("failed"),
			expected: false,
		},
		{
			name:     "is a request error",
			err:      &ResponseError{},
			expected: true,
		},
		{
			name:     "joined errors",
			err:      errors.Join(errors.New("boom"), &ResponseError{}),
			expected: true,
		},
		{
			name:     "fmt error",
			err:      fmt.Errorf("check permissions: %w", &ResponseError{}),
			expected: true,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, ok := AsResponseError(tc.err)
			assert.Equal(t, tc.expected, ok, "Must match the expected value")
		})
	}
}

func TestRequestErrorMethods(t *testing.T) {
	t.Parallel()

	re := &ResponseError{code: http.StatusOK, route: "/heathz", details: "service alive"}

	assert.Equal(t, 200, re.Code(), "Must match the expected code")
	assert.Equal(t, "/heathz", re.Route(), "Must match the expected route")
	assert.Equal(t, "service alive", re.Details(), "Must match the expected details")
}
