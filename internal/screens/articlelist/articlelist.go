// Package articlelist implements the article list screen.
// It receives a feed and its articles, lets the user browse them,
// and emits messages that app.go uses to drive navigation.
package articlelist

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mako/cerrynt-cli/internal/domain"
	"github.com/mako/cerrynt-cli/internal/keymap"
	"github.com/mako/cerrynt-cli/internal/styles"
)

// ArticleSelectedMsg is emitted when the user opens an article.
type ArticleSelectedMsg struct {
	Article domain.Article
}

// BackMsg is emitted when the user presses ESC to return to the feed list.
type BackMsg struct{}

// Model is the Bubble Tea model for the article list screen.
type Model struct {
	feed     domain.Feed
	articles []domain.Article
	cursor   int
}

// New constructs an articlelist Model for the given feed and its articles.
func New(feed domain.Feed, articles []domain.Article) Model {
	return Model{feed: feed, articles: articles}
}

// Init satisfies tea.Model. No startup commands needed.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key input and returns navigation commands.
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
		if m.cursor < len(m.articles)-1 {
			m.cursor++
		}

	case key.Matches(keyMsg, keymap.Default.Open):
		if len(m.articles) > 0 {
			selected := m.articles[m.cursor]
			return m, func() tea.Msg {
				return ArticleSelectedMsg{Article: selected}
			}
		}

	case key.Matches(keyMsg, keymap.Default.Back):
		return m, func() tea.Msg { return BackMsg{} }
	}

	return m, nil
}

// View renders the article list for the current feed.
func (m Model) View() string {
	s := styles.Title.Render(m.feed.Title) + "\n"

	if len(m.articles) == 0 {
		s += styles.Faint.Render("  No articles.") + "\n"
	} else {
		for i, article := range m.articles {
			age := formatAge(article.Published)
			title := article.Title

			if i == m.cursor {
				s += styles.Selected.Render(fmt.Sprintf("> %s  %s", title, age)) + "\n"
			} else if article.IsRead {
				s += styles.Faint.Render(fmt.Sprintf("  %s  %s", title, age)) + "\n"
			} else {
				s += fmt.Sprintf("  %s  %s\n", title, age)
			}
		}
	}

	s += styles.StatusBar.Render("j/k move  •  enter open  •  esc back  •  q quit")

	return s
}

// formatAge returns a human-readable age string for a published timestamp.
func formatAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
