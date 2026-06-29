// Package articlelist implements the article list screen.
// It receives a feed and its articles, lets the user browse them,
// and emits messages that app.go uses to drive navigation.
package articlelist

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
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
	loading  bool
	err      error
	spinner  spinner.Model
}

// New constructs an articlelist Model in the loading state for the given feed.
// Articles are populated later via SetArticles once the async load completes.
func New(feed domain.Feed) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return Model{
		feed:    feed,
		loading: true,
		spinner: s,
	}
}

// SetArticles returns a new Model populated with articles and no longer loading.
func (m Model) SetArticles(articles []domain.Article) Model {
	m.articles = articles
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

// Update handles key input and spinner ticks.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var spinCmd tea.Cmd
	m.spinner, spinCmd = m.spinner.Update(msg)

	if m.loading || m.err != nil {
		// ESC is still handled in error/loading state so the user can go back.
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(keyMsg, keymap.Default.Back) {
				return m, func() tea.Msg { return BackMsg{} }
			}
		}
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

	return m, spinCmd
}

// MarkRead returns a new Model with the article matching id marked as read.
// This is called by app.go immediately when the user opens an article so the
// list reflects the new state when the user navigates back.
//
// Note on slice copying: domain.Article is a plain struct (no reference fields),
// so copy() produces a true independent copy — modifying elements of the new
// slice does not affect the original.
func (m Model) MarkRead(id string) Model {
	articles := make([]domain.Article, len(m.articles))
	copy(articles, m.articles)
	for i, a := range articles {
		if a.ID == id {
			articles[i].IsRead = true
			break
		}
	}
	m.articles = articles
	return m
}

// View renders the article list for the current feed.
func (m Model) View() string {
	s := styles.Title.Render(m.feed.Title) + "\n\n"

	switch {
	case m.loading:
		s += "  " + m.spinner.View() + " Loading articles...\n"

	case m.err != nil:
		s += styles.Faint.Render("  Error: "+m.err.Error()) + "\n"

	default:
		// Count unread for the header subtitle.
		unread := 0
		for _, a := range m.articles {
			if !a.IsRead {
				unread++
			}
		}
		if unread > 0 {
			// Re-render title with unread count now that we have data.
			s = styles.Title.Render(fmt.Sprintf("%s — %d unread", m.feed.Title, unread)) + "\n\n"
		}

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
	}

	var pos string
	if !m.loading && m.err == nil && len(m.articles) > 0 {
		pos = fmt.Sprintf("[%d/%d]  ", m.cursor+1, len(m.articles))
	}
	s += "\n" + styles.StatusBar.Render(pos+"j/k move  •  enter open  •  esc back  •  q quit")

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
