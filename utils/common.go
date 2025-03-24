package utils

import (
	"bytes"
	"errors"
	"net/http"
)

// MakeAPICall sends an HTTP request based on the given method, URL, and payload.
func MakeAPICall(method string, url string, payload []byte) (*http.Response, error) {
	client := &http.Client{}

	// Create a new HTTP request
	var req *http.Request
	var err error

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		// For methods requiring a payload
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else if method == http.MethodGet || method == http.MethodDelete {
		// For methods not requiring a payload
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unsupported HTTP method")
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
