// Package articleview implements the article reading screen.
// It receives a single article and displays its content.
// The only navigation action is ESC to go back.
package articleview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mako/cerrynt-cli/internal/domain"
	"github.com/mako/cerrynt-cli/internal/keymap"
	"github.com/mako/cerrynt-cli/internal/styles"
)

// BackMsg is emitted when the user presses ESC to return to the article list.
type BackMsg struct{}

// Model is the Bubble Tea model for the article view screen.
type Model struct {
	article domain.Article
}

// New constructs an articleview Model for the given article.
func New(article domain.Article) Model {
	return Model{article: article}
}

// Init satisfies tea.Model. No startup commands needed.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key input. Only ESC is meaningful here.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if key.Matches(keyMsg, keymap.Default.Back) {
		return m, func() tea.Msg { return BackMsg{} }
	}

	return m, nil
}

// View renders the article content.
func (m Model) View() string {
	a := m.article

	s := styles.Title.Render(a.Title) + "\n"
	s += styles.Faint.Render(fmt.Sprintf("Published  %s", a.Published.Format("2 Jan 2006  15:04"))) + "\n"
	s += styles.Faint.Render(fmt.Sprintf("URL        %s", a.URL)) + "\n"
	s += "\n"

	if a.Summary != "" {
		s += wordWrap(a.Summary, 72) + "\n"
	} else {
		s += styles.Faint.Render("No summary available.") + "\n"
	}

	s += styles.StatusBar.Render("esc back  •  q quit")

	return s
}

// wordWrap breaks text into lines of at most width characters, breaking on
// word boundaries. This is a simple implementation sufficient for the MVP.
func wordWrap(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	current := words[0]

	for _, word := range words[1:] {
		if len(current)+1+len(word) > width {
			lines = append(lines, current)
			current = word
		} else {
			current += " " + word
		}
	}
	lines = append(lines, current)

	return strings.Join(lines, "\n")
}
