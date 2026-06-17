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
		{ID: "1", Title: "Hacker News", URL: "https://news.ycombinator.com/rss"},
		{ID: "2", Title: "Go Blog", URL: "https://go.dev/blog/feed.atom"},
		{ID: "3", Title: "lobste.rs", URL: "https://lobste.rs/rss"},
		{ID: "4", Title: "Drew DeVault's Blog", URL: "https://drewdevault.com/blog/index.xml"},
		{ID: "5", Title: "The Changelog", URL: "https://changelog.com/feed"},
	}
}

// Articles returns a static list of mock articles for the given feed ID.
func Articles(feedID string) []domain.Article {
	now := time.Now()

	all := map[string][]domain.Article{
		"1": {
			{
				ID:      "101",
				FeedID:  "1",
				Title:   "Ask HN: What are you building this month?",
				URL:     "https://news.ycombinator.com/item?id=101",
				Summary: "It's that time again — the monthly thread where HN members share what they are currently working on. This month's thread has attracted over 600 comments covering everything from indie SaaS tools and developer productivity apps to hardware experiments and open source libraries. A recurring theme this month is terminal-first tooling: several people are building RSS readers, task managers, and database explorers designed to run entirely in the terminal. One commenter noted that tools like Newsboat, ranger, and lazygit have inspired a new generation of TUI applications, many written in Go using the Bubble Tea framework.",
				Published: now.Add(-2 * time.Hour),
				IsRead:    false,
			},
			{
				ID:      "102",
				FeedID:  "1",
				Title:   "Show HN: I built a terminal RSS reader in Go",
				URL:     "https://news.ycombinator.com/item?id=102",
				Summary: "I have been using Newsboat for years but wanted something that could sync read state across multiple machines without relying on a self-hosted server. So I built Cerrynt: a terminal-first RSS reader with a Rails backend for sync. The CLI is written in Go using Bubble Tea for the TUI layer and Lip Gloss for styling. It supports vim keybindings, respects the terminal color scheme by default, and reads config from the XDG config directory. The backend is a standard Rails API that stores feeds, articles, and read state. The CLI talks to it over HTTP with a token stored in the config file. Early days but it already replaces Newsboat for my daily use.",
				Published: now.Add(-5 * time.Hour),
				IsRead:    true,
			},
			{
				ID:      "103",
				FeedID:  "1",
				Title:   "Why Go is eating the CLI tooling world",
				URL:     "https://news.ycombinator.com/item?id=103",
				Summary: "A look at why Go has become the dominant language for command-line tools over the past five years. The article covers several factors: static binaries with no runtime dependencies, fast compilation, a strong standard library, and excellent cross-compilation support. The author also credits the ecosystem — tools like Cobra for argument parsing, Viper for configuration, and Bubble Tea for terminal UIs have dramatically lowered the barrier to writing polished CLI applications. The piece includes a survey of popular Go CLI tools including kubectl, gh, lazygit, and k9s, and examines what they have in common architecturally.",
				Published: now.Add(-8 * time.Hour),
				IsRead:    false,
			},
		},
		"2": {
			{
				ID:      "201",
				FeedID:  "2",
				Title:   "Go 1.23 is released",
				URL:     "https://go.dev/blog/go1.23",
				Summary: "Go 1.23 is now available. The most anticipated addition in this release is range-over-function iterators, which allow custom types to be iterable using the standard for-range syntax. This has been a long-requested feature and landed after extended discussion and a prototype in the x/exp repository. The release also includes improvements to the timer implementation that fix a long-standing goroutine leak in programs that create many short-lived timers, a new unique package for interning comparable values, and several toolchain improvements including faster builds for large modules. The compatibility guarantee continues to hold: all existing Go programs compile and run correctly with 1.23.",
				Published: now.Add(-48 * time.Hour),
				IsRead:    false,
			},
			{
				ID:      "202",
				FeedID:  "2",
				Title:   "Structured logging with slog",
				URL:     "https://go.dev/blog/slog",
				Summary: "The log/slog package, introduced in Go 1.21, provides structured, leveled logging for Go programs. This post explains the design goals behind slog, how it differs from existing logging libraries, and how to use it effectively. The key insight is that slog separates the logging API (the Logger type and its methods) from the output implementation (the Handler interface), making it easy to route logs to different backends — JSON files, human-readable consoles, external aggregators — without changing the call sites. The post also covers the performance characteristics of slog, including the use of the Enabled method to avoid expensive argument evaluation when a log level is disabled.",
				Published: now.Add(-72 * time.Hour),
				IsRead:    true,
			},
			{
				ID:      "203",
				FeedID:  "2",
				Title:   "Evolving the Go standard library with math/rand/v2",
				URL:     "https://go.dev/blog/randv2",
				Summary: "Go 1.22 introduced math/rand/v2, the first v2 package in the Go standard library. This post explains the decisions behind the new package: removing the global source that leaked state between packages, switching the default algorithm to a faster and higher-quality generator, and cleaning up the API by removing methods that had been kept only for historical compatibility. The post is also a case study in how Go plans to evolve the standard library over time without breaking the compatibility guarantee — by using the /v2 module path convention borrowed from the broader Go module ecosystem.",
				Published: now.Add(-96 * time.Hour),
				IsRead:    false,
			},
		},
		"3": {
			{
				ID:      "301",
				FeedID:  "3",
				Title:   "Why I switched from Vim to Neovim",
				URL:     "https://lobste.rs/s/abc123",
				Summary: "After eight years of Vim, I migrated to Neovim last quarter and have not looked back. The move was motivated primarily by the Lua configuration ecosystem: being able to write plugin configs in a real programming language rather than Vimscript has made my setup significantly more maintainable. This post documents the migration: which plugins I replaced, how I handled the muscle memory transition, and what I found unexpectedly better or worse. The LSP integration in Neovim was the biggest practical win — neovim-lspconfig plus Mason means I get proper completion, diagnostics, and go-to-definition in every language without maintaining per-language plugin configurations.",
				Published: now.Add(-1 * time.Hour),
				IsRead:    false,
			},
			{
				ID:      "302",
				FeedID:  "3",
				Title:   "The return of the terminal",
				URL:     "https://lobste.rs/s/def456",
				Summary: "Terminal applications are having a genuine renaissance. After years of Electron apps and web-based developer tools, there is a visible shift back toward TUI applications among a vocal and growing segment of the developer community. Tools like lazygit, k9s, btop, and Newsboat are seeing increased adoption. This article explores the reasons: performance, keyboard-first interaction, SSH compatibility, and a backlash against resource-heavy GUI applications. The author argues that the terminal is not a niche interest but rather the natural interface for developers who spend most of their time in a shell anyway — and that modern TUI frameworks like Bubble Tea and Ratatui have made building these tools significantly more accessible.",
				Published: now.Add(-3 * time.Hour),
				IsRead:    false,
			},
		},
		"4": {
			{
				ID:      "401",
				FeedID:  "4",
				Title:   "My experience with sourcehut after two years",
				URL:     "https://drewdevault.com/2024/01/01/sourcehut.html",
				Summary: "Two years ago I moved all of my personal and professional projects off GitHub and onto sourcehut. This post is an honest retrospective on that decision. The good: email-based patch workflow forces thoughtful contribution, the CI system is flexible and YAML-based without vendor lock-in, the platform is fast and has no JavaScript dependency for basic use, and I have full confidence in the long-term sustainability of the project. The challenging: the contributor friction is real and has reduced the volume of outside contributions noticeably, some tooling that expects GitHub is painful to work around, and the notification model takes adjustment. Overall the move was worth it for my use case, though I would not recommend it unconditionally.",
				Published: now.Add(-24 * time.Hour),
				IsRead:    false,
			},
		},
		"5": {},
	}

	return all[feedID]
}
