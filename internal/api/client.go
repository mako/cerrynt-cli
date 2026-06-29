// Package api provides an HTTP client for communicating with the Cerrynt API.
package api

import (
	"net/http"
	"time"
)

// Client is responsible only for holding HTTP communication parameters
// for the Cerrynt API. Request execution logic is handled by the caller.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// New constructs a Client for the given API base URL and auth token.
//
// The client owns its http.Client with a 15-second timeout.
func New(baseURL, token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}
