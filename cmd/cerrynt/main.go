package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mako/cerrynt-cli/internal/app"
	"github.com/mako/cerrynt-cli/internal/config"
	"github.com/mako/cerrynt-cli/internal/data"
	"github.com/mako/cerrynt-cli/internal/domain"
	"github.com/mako/cerrynt-cli/internal/state"
)

func main() {
	feeds, err := loadFeeds()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cerrynt: %v\n", err)
		os.Exit(1)
	}

	st, statePath, err := loadState()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cerrynt: %v\n", err)
		os.Exit(1)
	}

	// config is the source of truth for which feeds the user subscribes to.
	// MockStore is given those feeds so the Store layer serves them without
	// independently deciding what the feed list is.
	store := data.NewMockStore(feeds)

	p := tea.NewProgram(
		app.New(store, st, statePath),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cerrynt: %v\n", err)
		os.Exit(1)
	}
}

// loadFeeds resolves the config path, loads the config file, and returns the
// feeds as domain types. All error handling for a missing or malformed config
// is centralised here so the rest of the app stays unaware of the source.
func loadFeeds() ([]domain.Feed, error) {
	path, err := config.DefaultPath()
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config file not found: %s\n\nCreate it with your feeds to get started.", path)
		}
		return nil, fmt.Errorf("load config: %w", err)
	}

	feeds := make([]domain.Feed, len(cfg.Feeds))
	for i, f := range cfg.Feeds {
		feeds[i] = f.ToDomain()
	}
	return feeds, nil
}

// loadState resolves the state path and loads persisted local state.
// A missing state file is normal on first run and results in an empty State
// rather than an error. Only unexpected errors (e.g. corrupt JSON, permission
// denied) are returned as fatal.
func loadState() (*state.State, string, error) {
	path, err := state.DefaultPath()
	if err != nil {
		return nil, "", fmt.Errorf("resolve state path: %w", err)
	}

	st, err := state.Load(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// First run: no state file yet. Start fresh.
			return &state.State{Read: make(map[string]bool)}, path, nil
		}
		return nil, "", fmt.Errorf("load state: %w", err)
	}

	return st, path, nil
}
