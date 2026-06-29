package data

import (
	"context"
	"fmt"
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

// Feeds fetches the user's subscribed feeds from the API.
func (s *APIStore) Feeds(ctx context.Context) ([]domain.Feed, error) {
	reqURL := strings.TrimRight(s.client.BaseURL, "/") + "/feeds"
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("data: create feeds request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.Token)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("data: execute feeds request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("data: unexpected status %d", resp.StatusCode)
	}

	dtos, err := api.DecodeFeeds(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("data: decode feeds: %w", err)
	}

	return api.MapFeeds(dtos), nil
}

// Articles fetches articles for a given feed from the API.
func (s *APIStore) Articles(ctx context.Context, feedID string) ([]domain.Article, error) {
	path := fmt.Sprintf("/feeds/%s/articles", url.PathEscape(feedID))
	reqURL := strings.TrimRight(s.client.BaseURL, "/") + path
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("data: create articles request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.Token)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("data: execute articles request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("data: unexpected status %d", resp.StatusCode)
	}

	dtos, err := api.DecodeArticles(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("data: decode articles: %w", err)
	}

	return api.MapArticles(dtos), nil
}
