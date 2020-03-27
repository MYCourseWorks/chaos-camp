package infra

import (
	"encoding/json"
	"fmt"
)

// ClientError comment
type ClientError interface {
	error // inherits the error interface
	ResponseBody() ([]byte, error)
	ResponseHeaders() (int, map[string]string)
}

// HTTPError comment
type HTTPError struct {
	Cause  error  `json:"-"`
	Detail string `json:"detail"`
	Status int    `json:"-"`
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(cause error, status int, detail string) *HTTPError {
	return &HTTPError{
		Cause:  cause,
		Detail: detail,
		Status: status,
	}
}

func (e *HTTPError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + " : " + e.Cause.Error()
}

// ResponseBody returns JSON response body.
func (e *HTTPError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}
	return body, nil
}

// ResponseHeaders returns http status code and headers.
func (e *HTTPError) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}
