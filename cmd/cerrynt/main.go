package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mako/cerrynt-cli/internal/app"
	"github.com/mako/cerrynt-cli/internal/config"
	"github.com/mako/cerrynt-cli/internal/domain"
)

func main() {
	feeds, err := loadFeeds()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cerrynt: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		app.New(feeds),
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

