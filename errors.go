package signalfx

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ResponseError captures the error details
// and allows for it to be inspected by external libraries.
type ResponseError struct {
	code    int
	route   string
	details string
}

var _ error = (*ResponseError)(nil)

func newResponseError(resp *http.Response, target int, targets ...int) error {
	if resp == nil {
		return nil
	}

	// Once upgraded to go 1.22+, replace for slices.Contains
	for _, code := range append([]int{target}, targets...) {
		if resp.StatusCode == code {
			return nil
		}
	}

	details, _ := io.ReadAll(resp.Body)

	return &ResponseError{
		code:    resp.StatusCode,
		route:   resp.Request.URL.Path,
		details: string(details),
	}
}

// AsResponseError is a convenience function to check the error
// to see if it contains an `ResponseError` and returns the value with true.
// If the error was initially joined using [errors.Join], it will check each error
// within the list and return the first matching error.
func AsResponseError(err error) (*ResponseError, bool) {
	// When `errors.Join` is called, it returns an error that
	// matches the provided interface.
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		for _, err := range joined.Unwrap() {
			if re, ok := AsResponseError(err); ok {
				return re, ok
			}
		}
		return nil, false
	}

	for err != nil {
		if re, ok := err.(*ResponseError); ok {
			return re, true
		}
		// In case the error is wrapped using `fmt.Errorf`
		// this will also account for that.
		err = errors.Unwrap(err)
	}
	return nil, false
}

func (re *ResponseError) Error() string {
	return fmt.Sprintf("route %q had issues with status code %d", re.route, re.code)
}

func (re *ResponseError) Code() int {
	return re.code
}

func (re *ResponseError) Route() string {
	return re.route
}

func (re *ResponseError) Details() string {
	return re.details
}
