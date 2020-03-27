package infra

import (
	"net/http"
	"os"
)

type roundTripperFn func(req *http.Request) (*http.Response, error)

func (fn roundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// NewHTTPClient creates the default http client for the application
func NewHTTPClient(transport http.RoundTripper) *http.Client {
	client := &http.Client{
		Transport: transport,
	}

	return client
}

// NewFSClient creates a new http.Client that reads from the file system
func NewFSClient() *http.Client {
	fsTripper := func(req *http.Request) (*http.Response, error) {
		file, err := os.OpenFile(req.URL.String(), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}

		resp := &http.Response{StatusCode: 200, Body: file, Header: make(http.Header)}
		return resp, nil
	}

	return NewHTTPClient(roundTripperFn(fsTripper))
}
