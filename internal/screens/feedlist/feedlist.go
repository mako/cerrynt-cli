// Package feedlist implements the feed list screen.
// This is the first screen the user sees. It displays all subscribed feeds
// with their unread counts and allows the user to select one to open.
package feedlist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
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
	feeds  []domain.Feed
	cursor int
}

// New constructs a feedlist Model from a slice of feeds.
func New(feeds []domain.Feed) Model {
	return Model{feeds: feeds}
}

// Init satisfies the tea.Model interface. No I/O commands needed on startup.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key presses and returns an updated model and an optional command.
//
// In Bubble Tea, Update() is a pure function: given a message, return the next
// state. Side effects (network calls, navigation) are expressed as tea.Cmd
// values — functions that run asynchronously and produce the next message.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
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

	return m, nil
}

// View renders the feed list as a string. Bubble Tea calls this after every
// Update() and replaces the terminal output with the result.
func (m Model) View() string {
	s := styles.Title.Render("Feeds") + "\n\n"

	for i, feed := range m.feeds {
		var unread string
		if feed.UnreadCount > 0 {
			unread = fmt.Sprintf(" (%d)", feed.UnreadCount)
		}

		line := fmt.Sprintf("  %s%s", feed.Title, unread)

		if i == m.cursor {
			// Prefix selected row with ">" and apply bold style.
			line = styles.Selected.Render(fmt.Sprintf("> %s%s", feed.Title, unread))
		} else if feed.UnreadCount == 0 {
			line = styles.Faint.Render(line)
		}

		s += line + "\n"
	}

	var pos string
	if len(m.feeds) > 0 {
		pos = fmt.Sprintf("[%d/%d]  ", m.cursor+1, len(m.feeds))
	}
	s += "\n" + styles.StatusBar.Render(pos+"j/k move  •  enter open  •  q quit")

	return s
}
