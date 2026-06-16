// Package data provides feed and article data for the application.
// For MVP this is static mock data. It will be replaced by an API client
// in a future phase without changes to callers.
package data

import (
	"time"

	"github.com/mako/cerrynt-cli/internal/domain"
)

// Feeds returns a static list of mock RSS feeds.
func Feeds() []domain.Feed {
	return []domain.Feed{
		{ID: "1", Title: "Hacker News", URL: "https://news.ycombinator.com/rss", UnreadCount: 12},
		{ID: "2", Title: "Go Blog", URL: "https://go.dev/blog/feed.atom", UnreadCount: 3},
		{ID: "3", Title: "lobste.rs", URL: "https://lobste.rs/rss", UnreadCount: 7},
		{ID: "4", Title: "Drew DeVault's Blog", URL: "https://drewdevault.com/blog/index.xml", UnreadCount: 1},
		{ID: "5", Title: "The Changelog", URL: "https://changelog.com/feed", UnreadCount: 0},
	}
}

// Articles returns a static list of mock articles for the given feed ID.
func Articles(feedID string) []domain.Article {
	now := time.Now()

	all := map[string][]domain.Article{
		"1": {
			{
				ID: "101", FeedID: "1",
				Title:     "Ask HN: What are you building this month?",
				URL:       "https://news.ycombinator.com/item?id=101",
				Summary:   "Monthly thread for sharing projects in progress.",
				Published: now.Add(-2 * time.Hour),
				IsRead:    false,
			},
			{
				ID: "102", FeedID: "1",
				Title:     "Show HN: I built a terminal RSS reader in Go",
				URL:       "https://news.ycombinator.com/item?id=102",
				Summary:   "A Bubble Tea based RSS client with vim keybindings.",
				Published: now.Add(-5 * time.Hour),
				IsRead:    true,
			},
			{
				ID: "103", FeedID: "1",
				Title:     "Why Go is eating the CLI tooling world",
				URL:       "https://news.ycombinator.com/item?id=103",
				Summary:   "An analysis of Go's adoption in developer tooling.",
				Published: now.Add(-8 * time.Hour),
				IsRead:    false,
			},
		},
		"2": {
			{
				ID: "201", FeedID: "2",
				Title:     "Go 1.23 is released",
				URL:       "https://go.dev/blog/go1.23",
				Summary:   "Highlights from the latest Go release including range-over-func iterators.",
				Published: now.Add(-48 * time.Hour),
				IsRead:    false,
			},
			{
				ID: "202", FeedID: "2",
				Title:     "Structured logging with slog",
				URL:       "https://go.dev/blog/slog",
				Summary:   "A deep dive into the new log/slog package introduced in Go 1.21.",
				Published: now.Add(-72 * time.Hour),
				IsRead:    true,
			},
		},
		"3": {
			{
				ID: "301", FeedID: "3",
				Title:     "Why I switched from Vim to Neovim",
				URL:       "https://lobste.rs/s/abc123",
				Summary:   "A developer shares their migration experience and config setup.",
				Published: now.Add(-1 * time.Hour),
				IsRead:    false,
			},
			{
				ID: "302", FeedID: "3",
				Title:     "The return of the terminal",
				URL:       "https://lobste.rs/s/def456",
				Summary:   "How terminal-first tools are making a comeback among developers.",
				Published: now.Add(-3 * time.Hour),
				IsRead:    false,
			},
		},
		"4": {
			{
				ID: "401", FeedID: "4",
				Title:     "My experience with sourcehut after two years",
				URL:       "https://drewdevault.com/2024/01/01/sourcehut.html",
				Summary:   "Reflections on running a small software forge.",
				Published: now.Add(-24 * time.Hour),
				IsRead:    false,
			},
		},
		"5": {},
	}

	return all[feedID]
}
