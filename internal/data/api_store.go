package data

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mako/cerrynt-cli/internal/api"
	"github.com/mako/cerrynt-cli/internal/domain"
)

// APIStore is a Store backed by the Cerrynt API.
type APIStore struct {
	client *api.Client
}

// NewAPIStore constructs a new APIStore.
func NewAPIStore(baseURL, token string) *APIStore {
	return &APIStore{
		client: api.New(baseURL, token),
	}
}

// doRequest centralizes HTTP request creation, authorization, and execution.
// It uses a callback to ensure the response body is always safely closed,
// minimizing caller mistakes and preventing leaks.
func (s *APIStore) doRequest(ctx context.Context, method, reqURL string, decode func(io.Reader) error) error {
	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		return fmt.Errorf("data: create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.Token)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		// Let context cancellation errors propagate naturally, prefixed consistently.
		return fmt.Errorf("data: execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("data: unexpected status %d", resp.StatusCode)
	}

	return decode(resp.Body)
}

// Feeds fetches the user's subscribed feeds from the API.
func (s *APIStore) Feeds(ctx context.Context) ([]domain.Feed, error) {
	reqURL := strings.TrimRight(s.client.BaseURL, "/") + "/feeds"

	var dtos []api.FeedDTO
	err := s.doRequest(ctx, http.MethodGet, reqURL, func(r io.Reader) error {
		var decodeErr error
		dtos, decodeErr = api.DecodeFeeds(r)
		return decodeErr // Preserves underlying api: prefix
	})
	if err != nil {
		return nil, err
	}

	return api.MapFeeds(dtos), nil
}

// Articles fetches articles for a given feed from the API.
func (s *APIStore) Articles(ctx context.Context, feedID string) ([]domain.Article, error) {
	path := fmt.Sprintf("/feeds/%s/articles", url.PathEscape(feedID))
	reqURL := strings.TrimRight(s.client.BaseURL, "/") + path

	var dtos []api.ArticleDTO
	err := s.doRequest(ctx, http.MethodGet, reqURL, func(r io.Reader) error {
		var decodeErr error
		dtos, decodeErr = api.DecodeArticles(r)
		return decodeErr // Preserves underlying api: prefix
	})
	if err != nil {
		return nil, err
	}

	return api.MapArticles(dtos), nil
}
