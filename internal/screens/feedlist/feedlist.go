// Package feedlist implements the feed list screen.
// This is the first screen the user sees. It displays all subscribed feeds
// and allows the user to select one to open.
package feedlist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mako/cerrynt-cli/internal/domain"
	"github.com/mako/cerrynt-cli/internal/keymap"
	"github.com/mako/cerrynt-cli/internal/styles"
)

// FeedSelectedMsg is emitted when the user opens a feed.
// The parent model (app) receives this and navigates to the article list.
type FeedSelectedMsg struct {
	Feed domain.Feed
}

// Model is the Bubble Tea model for the feed list screen.
type Model struct {
	feeds   []domain.Feed
	cursor  int
	loading bool
	err     error
	spinner spinner.Model
}

// New constructs a feedlist Model in the loading state.
// Feeds are populated later via SetFeeds once the async load completes.
func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return Model{
		loading: true,
		spinner: s,
	}
}

// SetFeeds returns a new Model populated with feeds and no longer loading.
func (m Model) SetFeeds(feeds []domain.Feed) Model {
	m.feeds = feeds
	m.loading = false
	m.err = nil
	return m
}

// SetError returns a new Model in the error state.
func (m Model) SetError(err error) Model {
	m.loading = false
	m.err = err
	return m
}

// Init returns the spinner tick command so the spinner animates while loading.
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles key presses and spinner ticks.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Advance the spinner regardless of loading state — it costs nothing
	// and avoids a special case.
	var spinCmd tea.Cmd
	m.spinner, spinCmd = m.spinner.Update(msg)

	if m.loading || m.err != nil {
		// Discard key input while loading or in error state.
		return m, spinCmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, spinCmd
	}

	switch {
	case key.Matches(keyMsg, keymap.Default.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case key.Matches(keyMsg, keymap.Default.Down):
		if m.cursor < len(m.feeds)-1 {
			m.cursor++
		}

	case key.Matches(keyMsg, keymap.Default.Open):
		if len(m.feeds) > 0 {
			selected := m.feeds[m.cursor]
			// Return a command that produces a FeedSelectedMsg.
			// The parent app model inspects this message to trigger navigation.
			return m, func() tea.Msg {
				return FeedSelectedMsg{Feed: selected}
			}
		}
	}

	return m, spinCmd
}

// View renders the feed list as a string.
func (m Model) View() string {
	s := styles.Title.Render("Feeds") + "\n\n"

	switch {
	case m.loading:
		s += "  " + m.spinner.View() + " Loading feeds...\n"

	case m.err != nil:
		s += styles.Faint.Render("  Error: "+m.err.Error()) + "\n"

	default:
		for i, feed := range m.feeds {
			line := fmt.Sprintf("  %s", feed.Title)
			if i == m.cursor {
				line = styles.Selected.Render(fmt.Sprintf("> %s", feed.Title))
			}
			s += line + "\n"
		}
	}

	var pos string
	if !m.loading && m.err == nil && len(m.feeds) > 0 {
		pos = fmt.Sprintf("[%d/%d]  ", m.cursor+1, len(m.feeds))
	}
	s += "\n" + styles.StatusBar.Render(pos+"j/k move  •  enter open  •  q quit")

	return s
}
