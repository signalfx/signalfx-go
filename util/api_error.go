package util

// APIError is an error entry in an API response envelope. APIs in this
// repository return the status code as either a string or an integer.
type APIError struct {
	Code    StringOrInteger `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
}
