package data

import (
	"context"

	"github.com/mako/cerrynt-cli/internal/domain"
)

// Store is the interface for all data access in the application.
// The mock implementation wraps the static functions in this package.
// Future implementations will call the Cerrynt Rails API.
//
// Methods accept a context so callers can cancel in-flight requests when
// the user navigates away. For the mock store the context is ignored.
type Store interface {
	Feeds(ctx context.Context) ([]domain.Feed, error)
	Articles(ctx context.Context, feedID string) ([]domain.Article, error)
}

// MockStore is a Store backed by feeds provided at construction time and
// static mock article data keyed by feed ID.
//
// The feed list comes from the caller (config.yaml via main.go) so that the
// Store is not responsible for deciding which feeds the user has subscribed to.
// Article data remains static mock content for the MVP.
type MockStore struct {
	feeds []domain.Feed
}

// NewMockStore constructs a MockStore that serves the given feeds.
func NewMockStore(feeds []domain.Feed) MockStore {
	return MockStore{feeds: feeds}
}

// Feeds returns the configured feed list.
func (s MockStore) Feeds(_ context.Context) ([]domain.Feed, error) {
	return s.feeds, nil
}

// Articles returns the static mock articles for the given feed ID.
func (s MockStore) Articles(_ context.Context, feedID string) ([]domain.Article, error) {
	return Articles(feedID), nil
}
