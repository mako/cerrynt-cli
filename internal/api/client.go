// Package api provides an HTTP client for communicating with the Cerrynt API.
package api

import (
	"net/http"
	"time"
)

// Client is responsible only for HTTP communication with the Cerrynt API.
// It holds the connection parameters and owns the underlying http.Client.
//
// Request methods (added in later steps) will use the unexported httpClient
// to make calls, attaching the auth token and resolving paths against baseURL.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New constructs a Client for the given API base URL and auth token.
//
// The client owns its http.Client with a 15-second timeout. Using a private
// instance rather than http.DefaultClient avoids modifying global shared state
// and allows transport configuration per-client in the future.
//
// Contexts for cancellation belong on individual request methods, not here.
func New(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}
