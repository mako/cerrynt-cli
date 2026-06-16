// Package articleview implements the article reading screen.
// It receives a single article and displays its content.
// Content is scrollable with j/k when it exceeds the terminal height.
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

// statusBarHeight is the number of terminal rows consumed by the status bar.
// styles.StatusBar has MarginTop(1), so it renders as a blank line + the bar itself.
const statusBarHeight = 2

// Model is the Bubble Tea model for the article view screen.
type Model struct {
	article domain.Article
	width   int
	height  int
	offset  int // first visible line index (scroll position)
}

// New constructs an articleview Model. width and height should be the current
// terminal dimensions so that word-wrap and scrolling are correct from the first render.
func New(article domain.Article, width, height int) Model {
	return Model{
		article: article,
		width:   width,
		height:  height,
	}
}

// Init satisfies tea.Model. No startup commands needed.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key input and window resize events.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Clamp offset in case the window shrunk.
		m.offset = m.clampOffset(m.offset)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Default.Down):
			m.offset = m.clampOffset(m.offset + 1)

		case key.Matches(msg, keymap.Default.Up):
			if m.offset > 0 {
				m.offset--
			}

		case key.Matches(msg, keymap.Default.Back):
			return m, func() tea.Msg { return BackMsg{} }
		}
	}

	return m, nil
}

// View renders the visible portion of the article content plus a status bar.
func (m Model) View() string {
	lines := m.buildLines()

	visible := m.height - statusBarHeight
	if visible < 1 {
		visible = 1
	}

	// Clamp again at render time in case dimensions changed between Update calls.
	offset := m.clampOffset(m.offset)

	end := offset + visible
	if end > len(lines) {
		end = len(lines)
	}

	content := strings.Join(lines[offset:end], "\n")

	return content + "\n" + styles.StatusBar.Render(m.statusBarText(lines, visible, offset))
}

// buildLines renders the full article into a flat slice of display lines.
// This is used both by View() for rendering and by Update() for scroll clamping.
//
// Keeping rendering in one place means View() is a simple slice operation and
// scroll position is always relative to the same line layout.
func (m Model) buildLines() []string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	// Leave a small left margin so content is not flush against the edge.
	wrapWidth := w - 2
	if wrapWidth < 20 {
		wrapWidth = 20
	}

	var lines []string

	a := m.article

	// Title — may span multiple lines on narrow terminals.
	for _, l := range wordWrap(a.Title, wrapWidth) {
		lines = append(lines, styles.Title.Render(l))
	}

	// Metadata
	lines = append(lines,
		styles.Faint.Render(fmt.Sprintf("Published  %s", a.Published.Format("2 Jan 2006  15:04"))),
		styles.Faint.Render(fmt.Sprintf("URL        %s", a.URL)),
		"", // blank separator between metadata and body
	)

	// Body
	if a.Summary != "" {
		lines = append(lines, wordWrap(a.Summary, wrapWidth)...)
	} else {
		lines = append(lines, styles.Faint.Render("No summary available."))
	}

	return lines
}

// clampOffset ensures offset is within the valid scroll range.
func (m Model) clampOffset(offset int) int {
	lines := m.buildLines()
	visible := m.height - statusBarHeight
	if visible < 1 {
		visible = 1
	}
	max := len(lines) - visible
	if max < 0 {
		max = 0
	}
	if offset > max {
		return max
	}
	if offset < 0 {
		return 0
	}
	return offset
}

// statusBarText builds the status bar content, including a scroll percentage
// when the content is longer than the visible area.
func (m Model) statusBarText(lines []string, visible, offset int) string {
	hint := "esc back  •  q quit"

	if len(lines) <= visible {
		// All content fits — no scroll indicator needed.
		return hint
	}

	max := len(lines) - visible
	pct := 0
	if max > 0 {
		pct = offset * 100 / max
	}

	return fmt.Sprintf("[%d%%]  j/k scroll  •  %s", pct, hint)
}

// wordWrap breaks text into lines of at most width characters on word boundaries.
// It returns a slice of lines rather than a joined string, so callers can
// treat each line individually (for scrolling, styling, etc.).
func wordWrap(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
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

	return append(lines, current)
}
