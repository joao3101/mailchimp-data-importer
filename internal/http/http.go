// Package http implements the http interface
package http

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPClientWrapper defines the interface to interact with the HTTP client
type HTTPClientWrapper interface {
	MakeHTTPRequest(req *http.Request) ([]byte, error)
}

func NewHTTPClientWrapper() HTTPClientWrapper {
	return &httpClientWrapper{
		httpClient: &http.Client{},
	}
}

type httpClientWrapper struct {
	httpClient
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// makeHTTPRequest makes an HTTP request, handles the response and returns the body
func (c *httpClientWrapper) MakeHTTPRequest(req *http.Request) ([]byte, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("nil response from %v", req.URL.String())
	}
	if res.StatusCode > 299 { //nolint:gomnd
		return nil, fmt.Errorf("error code %d making request", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error while reading the response bytes:%v", err)
	}
	return body, nil
}
