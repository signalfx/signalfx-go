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

// IsResponseError is convenience function to see
// if it can convert into RequestError.
func IsResponseError(err error) (*ResponseError, bool) {
	var re *ResponseError
	if errors.As(err, &re) {
		return err.(*ResponseError), true
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
